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

type ProcessChunk struct {
	TimeRemaining float64
	NoChunks      int
}

const CHUNKTIME = 5 // Process in 5 second chunks

func newAudioChannel() *AudioChannel {
	return &AudioChannel{
		channel: make(chan AudioChunk),
	}
}

func (ac *AudioChannel) processAudio() {
	files := getFiles("./audio_files/")
	for {
		for _, file := range files {
			ac.processAudioFile(file)
		}
	}
}

func (ac *AudioChannel) processAudioFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Read file", file)
		return
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	prepared := prepareAudio(decoder)
	processed := false
	for !processed {
		buffer, isProcessed := prepared.decodeChunk(CHUNKTIME, decoder)
		processed = isProcessed
		chunk := createAudioChunk(buffer, decoder)
		select {
		case ac.channel <- chunk:
		}
	}

}

func prepareAudio(d *wav.Decoder) *ProcessChunk {
	d.ReadInfo()
	audioTime, err := d.Duration()
	at := audioTime.Seconds()
	if err != nil {
		log.Fatal("Unable to determine audio duration", err)
	}
	numChunks := int(at) / CHUNKTIME
	return &ProcessChunk{
		TimeRemaining: at,
		NoChunks:      numChunks,
	}
}

func createAudioChunk(buffer *audio.IntBuffer, decoder *wav.Decoder) AudioChunk {
	return AudioChunk{
		Data:        buffer.Data,
		Format:      buffer.Format,
		BitDepth:    int(decoder.BitDepth),
		AudioFormat: int(decoder.WavAudioFormat),
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

func (p *ProcessChunk) decodeChunk(chunkTime uint32, d *wav.Decoder) (buf *audio.IntBuffer, done bool) {
	if p.TimeRemaining < float64(chunkTime) {
		chunkTime = uint32(p.TimeRemaining)
	}

	p.TimeRemaining = p.TimeRemaining - float64(chunkTime)

	numSamples := int(chunkTime * d.SampleRate)
	numChannels := int(d.NumChans)
	intBufferSize := numChannels * numSamples

	buf = createIntBuffer(d, intBufferSize)
	d.PCMBuffer(buf)

	done = d.EOF() || isAudioBufferEnd(buf.Data)
	return buf, done
}

func createIntBuffer(d *wav.Decoder, bufferSize int) *audio.IntBuffer {
	return &audio.IntBuffer{
		Format: d.Format(),
		Data:   make([]int, bufferSize),
	}
}

func isAudioBufferEnd(slice []int) bool {
	for _, v := range slice {
		if v != 0 {
			return false
		}
	}
	return true
}
