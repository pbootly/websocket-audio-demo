package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gorilla/websocket"
)

type AudioChunk struct {
	Data        []int
	Format      *audio.Format
	BitDepth    int
	AudioFormat int
}

func RunClient() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		i := 0
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("received message\n")
			audio, _ := deserializeAudioChunk(message)
			if len(audio.Data) != 0 {
				log.Println(audio.Format, audio.BitDepth, audio.AudioFormat, len(audio.Data))
				readIn(i, audio)
				i++
			}
			if i > 20 {
				break
			}
		}
	}()

	select {
	case <-done:
		return
	case <-interrupt:
		log.Println("interrupt")
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
		select {
		case <-done:
		}
		return
	}
}

func deserializeAudioChunk(data []byte) (AudioChunk, error) {
	var chunk AudioChunk
	// Deserialize JSON data into AudioChunk
	err := json.Unmarshal(data, &chunk)
	if err != nil {
		return AudioChunk{}, err
	}
	return chunk, nil
}

func readIn(idx int, msg AudioChunk) {
	rawBuff := msg.Data
	format := msg.Format
	audioBuf := audio.IntBuffer{
		Data:   rawBuff,
		Format: format,
	}
	// Test write to filesystem
	fout := fmt.Sprintf("./fileout/%d.wav", idx)
	log.Println(fout)
	out, err := os.Create(fout)
	if err != nil {
		log.Fatal("os create", err)
	}

	e := wav.NewEncoder(out,
		msg.Format.SampleRate,
		msg.BitDepth,
		msg.Format.NumChannels,
		msg.AudioFormat,
	)
	if err = e.Write(&audioBuf); err != nil {
		log.Fatal("ERROR WRITING NEW STUFF", err)
	}

	if err = e.Close(); err != nil {
		log.Fatal("ERROR CLOSING NEW STUFF", err)
	}

}
