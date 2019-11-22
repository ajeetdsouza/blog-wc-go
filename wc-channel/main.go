package main

import (
	"io"
	"os"
	"runtime"
)

func main() {
	if len(os.Args) < 2 {
		panic("no file path specified")
	}
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	chunks := make(chan Chunk)
	counts := make(chan Count)

	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		go ChunkCounter(chunks, counts)
	}

	const bufferSize = 16 * 1024
	lastCharIsSpace := true

	for {
		buffer := make([]byte, bufferSize)
		bytes, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		chunks <- Chunk{lastCharIsSpace, buffer[:bytes]}
		lastCharIsSpace = IsSpace(buffer[bytes-1])
	}
	close(chunks)

	totalCount := Count{}
	for i := 0; i < numWorkers; i++ {
		count := <-counts
		totalCount.LineCount += count.LineCount
		totalCount.WordCount += count.WordCount
	}
	close(counts)

	fileStat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	byteCount := fileStat.Size()

	println("%d %d %d %s\n", totalCount.LineCount, totalCount.WordCount, byteCount, file.Name())
}

func ChunkCounter(chunks <-chan Chunk, counts chan<- Count) {
	totalCount := Count{}
	for {
		chunk, ok := <-chunks
		if !ok {
			break
		}
		count := GetCount(chunk)
		totalCount.LineCount += count.LineCount
		totalCount.WordCount += count.WordCount
	}
	counts <- totalCount
}

type Chunk struct {
	PrevCharIsSpace bool
	Buffer          []byte
}

type Count struct {
	LineCount int
	WordCount int
}

func GetCount(chunk Chunk) Count {
	count := Count{}

	prevCharIsSpace := chunk.PrevCharIsSpace
	for _, b := range chunk.Buffer {
		switch b {
		case '\n':
			count.LineCount++
			prevCharIsSpace = true
		case ' ', '\t', '\r', '\v', '\f':
			prevCharIsSpace = true
		default:
			if prevCharIsSpace {
				prevCharIsSpace = false
				count.WordCount++
			}
		}
	}

	return count
}

func IsSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v' || b == '\f'
}
