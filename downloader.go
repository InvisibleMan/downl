package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

// Useful links

// https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185
// https://github.com/cavaliercoder/grab
// https://stackoverflow.com/questions/30532886/golang-dynamic-progressbar
// https://stackoverflow.com/questions/44318345/can-i-increase-golangs-http-stream-chunk-size
// https://www.reddit.com/r/golang/comments/4xtsbn/help_how_to_read_files_in_blocks/
// https://gobyexample.com/reading-files

func startDownloads(targets []string, maxSpeedBytes int, endEvent chan int) { // dest string, maxThreads int,
	// var dest = "./"
	// var maxThreads = 2

	var COUNTDOWN = 20                    // millisecond
	var defaultChunkSize int64 = 4 * 1024 // 4KB

	// https: //gobyexample.com/tickers
	ticker := time.NewTicker(time.Duration(COUNTDOWN) * time.Millisecond)
	queueChunks := make(chan int, 1)

	log.Println("Start download files:")
	for _, file := range os.Args[2:] {
		log.Printf(" * %v\n", file)
		go startSingleDownload(file, queueChunks, defaultChunkSize)
	}

	for _ = range ticker.C {
		queueChunks <- 1
		// log.Println("Tick at", t)
	}
}

func startSingleDownload(file string, chunks <-chan int, chunkSize int64) {
	resp, err := http.Get(file)
	if err == nil {
		defer resp.Body.Close()
	} else {
		return
	}

	log.Printf("Start dowload file: '%v'. With size: %v bytes\n", file, resp.ContentLength)
	p := make([]byte, chunkSize)
	reader := resp.Body
	var read int64

	var fullSize int64 = resp.ContentLength

	for {
		<-chunks
		_, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				log.Printf("Final read chunk. File: '%v' \n", file)
				return
			}
			log.Println(err)
		}
		read += chunkSize
		percent := read * 100 / fullSize
		log.Printf("Read chunk. File: '%v'. (%v %%)\n", file, percent)
		runtime.Gosched()
	}
}
