package main

import (
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type AudioChannel struct {
	channel chan AudioChunk
}

type AudioChunk struct {
	Data        []int
	Format      *audio.Format
	BitDepth    int
	AudioFormat int
}

func newAudioChannel() *AudioChannel {
	return &AudioChannel{
		channel: make(chan AudioChunk),
	}
}

func (ac *AudioChannel) processAudio() {
	files := getFiles("./audio_files/")
	for {
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				log.Fatal("Read file", err)
				return
			}
			defer f.Close()

			decoder := wav.NewDecoder(f)
			/*parsed = streamFile(decoder, ac)
			if parsed {
				f.Close()
			}*/
			for {
				buffer, done := decodeChunk(5, decoder)
				if done {
					break
				}
				chunk := AudioChunk{
					Data:        buffer.Data,
					Format:      buffer.Format,
					BitDepth:    int(decoder.BitDepth),
					AudioFormat: int(decoder.WavAudioFormat),
				}
				select {
				case ac.channel <- chunk:
				}
			}
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

func decodeChunk(chunkTime uint32, d *wav.Decoder) (*audio.IntBuffer, bool) {
	d.ReadInfo()
	done := false
	var buf *audio.IntBuffer
	numSamples := int(chunkTime * d.SampleRate)
	numChannels := int(d.NumChans)
	intBufferSize := numChannels * numSamples
	buf = &audio.IntBuffer{
		Format: d.Format(),
		Data:   make([]int, intBufferSize),
	}
	d.PCMBuffer(buf)
	if d.EOF() || isAudioBufferEnd(buf.Data) {
		log.Println("EOF")
		done = true
		return buf, done
	}
	return buf, done
}

func isAudioBufferEnd(slice []int) bool {
	for _, v := range slice {
		if v != 0 {
			return false
		}
	}
	return true
}
