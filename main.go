package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var htdocsDir = flag.String("htdocs", ".", "htdocs directory")

func main() {
	flag.Parse()

	http.Handle("/ws/audio", websocket.Handler(handleWebsocket))
	http.Handle("/", http.FileServer(http.Dir(*htdocsDir)))

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleWebsocket(s *websocket.Conn) {
	var data [8192]uint8
	log.Printf("Opened WebSocket")
	startTime := time.Now()

	f, err := os.OpenFile(fmt.Sprintf("%s/uploads/%s.mp3", *htdocsDir, time.Now().Format(time.RFC3339)), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("os.OpenFile failed: %v", err)
		return
	}
	defer f.Close()

	sum := 0
	count := 0

	for {
		n, err := s.Read(data[:])
		if err != nil {
			log.Printf("s.Read failed: %v", err)
			break
		}
		log.Printf("Received WebSocket frame: %d bytes", n)
		count++
		sum += n
		if _, err := f.Write(data[:n]); err != nil {
			log.Printf("f.Write failed: %v", err)
		}
	}

	duration := time.Since(startTime)

	log.Printf("Closed WebSocket, received %d frames (%d bytes), took %s (%.3f kb/s)", count, sum, duration, (float64(sum) / duration.Seconds()) / float64(1024))
}
