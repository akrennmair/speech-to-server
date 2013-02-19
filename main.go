package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/hex"
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

	f, err := os.OpenFile(fmt.Sprintf("%s.mp3", time.Now()), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("os.OpenFile failed: %v", err)
		return
	}
	defer f.Close()

	for {
		n, err := s.Read(data[:])
		if err != nil {
			log.Printf("s.Read failed: %v", err)
			return
		}
		log.Printf("received %d bytes:\n%s", n, hex.Dump(data[:n]))
		if _, err := f.Write(data[:n]); err != nil {
			log.Printf("f.Write failed: %v", err)
		}
	}
}
