package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/go-audio/wav"
)

type AudioChannel struct {
	channel chan []byte
}

func newAudioChannel() *AudioChannel {
	return &AudioChannel{
		channel: make(chan []byte),
	}
}

func (ac *AudioChannel) processAudio() {
	files := getFiles("./audio_files/")
	for {
		for _, f := range files {
			// WAV files supported only
			streamFile(f, ac)
		}
	}
}

func getFiles(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal("No files in directory: ", dir, err)
	}

	var r []string
	for _, file := range files {
		path := dir + file.Name()
		r = append(r, path)
	}
	return r
}

func streamFile(file string, ac *AudioChannel) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Unable to read file", err)
		return
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		log.Println("Skipping non wav file: ", f.Name())
		return
	}
	bufferSize := 128000 // High water mark
	wavHeader := parseWAVHeader(f)
	select {
	case ac.channel <- wavHeader:
	}

	buffer := make([]byte, bufferSize)
	for {
		n, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("File %s parsed\n", file)
				time.Sleep(22 * time.Second)
				break
			}
			log.Fatal(err)
		}
		select {
		case ac.channel <- buffer[:n]:
		}
	}

}

func parseWAVHeader(file *os.File) []byte {
	header := make([]byte, 44)
	_, err := file.Read(header)
	if err != nil {
		log.Fatal(err)
	}
	return header
}
