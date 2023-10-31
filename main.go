package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	audioChannel := newAudioChannel()
	hub := newHub(*audioChannel)
	go audioChannel.processAudio()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", hub.wsServer)
	server := &http.Server{
		Addr:              "localhost:8080",
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
