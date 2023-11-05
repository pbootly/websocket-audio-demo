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

// decodeChunk returns a decoder and an audio int buffer for a given .wav file and desired time
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
	if d.EOF() || isEnd(buf.Data) {
		log.Println("EOF")
		done = true
		return buf, done
	}
	return buf, done
}

func writeOut(chunkNo int, buf *audio.IntBuffer, d *wav.Decoder) {
	fout := fmt.Sprintf("./output/%d.wav", chunkNo)
	out, err := os.Create(fout)
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

}

func isEnd(slice []int) bool {
	for _, v := range slice {
		if v != 0 {
			return false
		}
	}
	return true
}

func main() {
	server := flag.Bool("server", false, "Start server")
	flag.Parse()
	if *server {
		log.Println("server")
	} else {
		log.Println("client")
	}

	file, err := os.Open("./we-the-people.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	d := wav.NewDecoder(file)
	i := 0
	for {
		buf, done := decodeChunk(5, d)
		if done {
			break
		}
		writeOut(i, buf, d)
		i++
	}
}
