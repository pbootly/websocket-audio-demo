package main

import (
	"encoding/binary"
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

func writeFile(chunkNo int, buf *audio.IntBuffer) string {
	dir := "./fileout/"
	file := fmt.Sprintf("%s%d", dir, chunkNo)
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	out := convertToInt32(buf.Data)
	err = binary.Write(f, binary.LittleEndian, out)
	if err != nil {
		log.Fatal(err)
	}
	f.Sync()
	return f.Name()
}

func readAndWrite(bufFile string, format *audio.Format, bitDepth int, audioFormat int) {
	file, err := os.Open(bufFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()

	numInt32 := fileSize / 4

	int32Slice := make([]int32, numInt32)
	err = binary.Read(file, binary.LittleEndian, &int32Slice)
	if err != nil {
		log.Fatal(err)
	}

	// in32slice is now our buffer
	iBuf := convertToIntSlice(int32Slice)
	var audioBuf *audio.IntBuffer
	audioBuf = &audio.IntBuffer{
		Data:   iBuf,
		Format: format,
	}
	fout := fmt.Sprintf("%s.wav", file.Name())
	out, err := os.Create(fout)
	if err != nil {
		log.Fatal("os create", err)
	}

	log.Printf("Audiobuf: %v\n", audioBuf.Format)
	e := wav.NewEncoder(out,
		audioBuf.Format.SampleRate,
		bitDepth,
		audioBuf.Format.NumChannels,
		audioFormat,
	)

	log.Println(e)

	if err = e.Write(audioBuf); err != nil {
		log.Fatal("ERROR WRITING NEW STUFF", err)
	}

	if err = e.Close(); err != nil {
		log.Fatal("ERROR CLOSING NEW STUFF", err)
	}
	out.Close()

}

func convertToInt32(audioData []int) []int32 {
	i32 := make([]int32, len(audioData))
	for i, v := range audioData {
		i32[i] = int32(v)
	}
	return i32
}

func convertToIntSlice(int32Slice []int32) []int {
	intSlice := make([]int, len(int32Slice))
	for i, v := range int32Slice {
		intSlice[i] = int(v)
	}
	return intSlice
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
		RunServer()
	} else {
		log.Println("client")
		RunClient()
	}

	/*file, err := os.Open("./we-the-people.wav")
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
		fileName := writeFile(i, buf)
		log.Println("Raw file:", fileName)
		//writeOut(i, buf, d)
		readAndWrite(fileName, buf.Format, int(d.BitDepth), int(d.WavAudioFormat))
		i++
	}*/

}
