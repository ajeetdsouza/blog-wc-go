package main

import (
	"io"
	"os"
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

	totalCount := Count{}
	lastCharIsSpace := true

	const bufferSize = 16 * 1024
	buffer := make([]byte, bufferSize)

	for {
		bytes, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		count := GetCount(Chunk{lastCharIsSpace, buffer[:bytes]})
		lastCharIsSpace = IsSpace(buffer[bytes-1])

		totalCount.LineCount += count.LineCount
		totalCount.WordCount += count.WordCount
	}

	fileStat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	byteCount := fileStat.Size()

	println(totalCount.LineCount, totalCount.WordCount, byteCount, file.Name())
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
