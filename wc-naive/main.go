package main

import (
	"bufio"
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

	const bufferSize = 16 * 1024
	reader := bufio.NewReaderSize(file, bufferSize)

	lineCount := 0
	wordCount := 0
	byteCount := 0

	prevByteIsSpace := true
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		byteCount++

		switch b {
		case '\n':
			lineCount++
			prevByteIsSpace = true
		case ' ', '\t', '\r', '\v', '\f':
			prevByteIsSpace = true
		default:
			if prevByteIsSpace {
				wordCount++
				prevByteIsSpace = false
			}
		}
	}

	println(lineCount, wordCount, byteCount, file.Name())
}
