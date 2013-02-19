package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.Handle("/ws/audio", websocket.Handler(handleWebsocket))
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleWebsocket(s *websocket.Conn) {
	var data [8192]uint8
	log.Printf("Opened WebSocket")
	startTime := time.Now()

	f, err := os.OpenFile(fmt.Sprintf("%s.mp3", time.Now()), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("os.OpenFile failed: %v", err)
		return
	}
	defer f.Close()

	sum := 0

	for {
		n, err := s.Read(data[:])
		if err != nil {
			log.Printf("s.Read failed: %v", err)
			break
		}
		sum += n
		if _, err := f.Write(data[:n]); err != nil {
			log.Printf("f.Write failed: %v", err)
		}
	}

	duration := time.Since(startTime)

	log.Printf("Closed WebSocket, received %d bytes, took %s (%.3f kb/s)", sum, duration, (float64(sum) / duration.Seconds()) / float64(1024))
}
