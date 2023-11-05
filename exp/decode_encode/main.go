package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// Server will read file as chunks
// Client will read chunks and encode them into seperate wav files

// First test - chunk files and re-assemble into their own .wav files

func decodeChunk(file *os.File, chunkTime uint32) (*wav.Decoder, *audio.IntBuffer) {
	d := wav.NewDecoder(file)
	d.ReadInfo()
	var buf *audio.IntBuffer
	numSamples := int(chunkTime * d.SampleRate)
	numChannels := int(d.NumChans)
	intBufferSize := numChannels * numSamples
	buf = &audio.IntBuffer{
		Format: d.Format(),
		Data:   make([]int, intBufferSize),
	}
	d.PCMBuffer(buf)
	if d.EOF() {
		return d, buf
	}
	return d, buf
}

func writeOut() {

}

func main() {
	server := flag.Bool("server", false, "Start server")
	flag.Parse()
	if *server {
		log.Println("server")
	} else {
		log.Println("client")
	}

	// Testing
	file, err := os.Open("./we-the-people.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	d, buf := decodeChunk(file, 5)
	out, err := os.Create("./kick.wav")
	if err != nil {
		panic(fmt.Sprintf("couldn't create output file - %v", err))
	}

	e := wav.NewEncoder(out,
		buf.Format.SampleRate,
		int(d.BitDepth),
		buf.Format.NumChannels,
		int(d.WavAudioFormat),
	)
	if err = e.Write(buf); err != nil {
		panic(err)
	}

	if err = e.Close(); err != nil {
		panic(err)
	}
	out.Close()

	out, err = os.Open("./kick.wav")
	if err != nil {
		panic(err)
	}
	d2 := wav.NewDecoder(out)
	d2.ReadInfo()
	fmt.Println("new:", d2)
	out.Close()
	//os.Remove(out.Name())
}
