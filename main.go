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
	http.HandleFunc("/livestream", streamLiveData)
	http.Handle("/", http.FileServer(http.Dir(*htdocsDir)))

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleWebsocket(s *websocket.Conn) {
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
		data := make([]byte, 8192)
		n, err := s.Read(data)
		if err != nil {
			log.Printf("s.Read failed: %v", err)
			break
		}
		broadcastData(data[:n])
		log.Printf("Received WebSocket frame: %d bytes", n)
		count++
		sum += n
		if _, err := f.Write(data[:n]); err != nil {
			log.Printf("f.Write failed: %v", err)
		}
	}

	endBroadcast()

	duration := time.Since(startTime)

	log.Printf("Closed WebSocket, received %d frames (%d bytes), took %s (%.3f kb/s)", count, sum, duration, (float64(sum) / duration.Seconds()) / float64(1024))
}

func streamLiveData(w http.ResponseWriter, r *http.Request) {
	ch := make(chan []byte)
	registerClient(ch)
	defer unregisterClient(ch)

	f := w.(http.Flusher)
	connClosed := w.(http.CloseNotifier).CloseNotify()
	w.Header().Set("Content-Type", "audio/mpeg")
	w.WriteHeader(http.StatusOK)

	for {
		select {
		case data, ok := <-ch:
			if !ok {
				log.Printf("End of transmission.")
				return
			}
			if _, err := w.Write(data); err != nil {
				log.Printf("Writing data to client failed: %v", err)
				return
			}
			f.Flush()
		case <-connClosed:
			log.Printf("Connection closed, stopping transmission.")
			return
		}
	}
}

var clients = make(map[chan []byte]struct{})

func registerClient(ch chan []byte) {
	clients[ch] = struct{}{}
	log.Printf("Registered client %p", ch)
}

func unregisterClient(ch chan []byte) {
	delete(clients, ch)
	log.Printf("Unregistered client %p", ch)
}

func broadcastData(data []byte) {
	for ch, _ := range clients {
		select {
		case ch <- data:
		}
	}
}

func endBroadcast() {
	for ch, _ := range clients {
		close(ch)
	}
}
