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

	defer f.Close()

	for {
		buffer := make([]byte, 256)
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
