package main

import (
	"io"
	"log"
	"os"
	"time"
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
			streamFile(f, ac)
			time.Sleep(1 * time.Second)
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
	bufferSize := 128000 // High water mark

	defer f.Close()

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
