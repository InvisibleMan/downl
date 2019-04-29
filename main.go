package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

func max(l, r int) int {
	if l > r {
		return l
	}
	return r
}

func main() {
	fmt.Println("Hello Download Master")

	if len(os.Args) < 3 {
		fmt.Println("Noting download. Bye")
		os.Exit(0)
		return
	}

	var maxSpeed = 0

	maxSpeed, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Wrong value of MaxSpeed. Bye")
		os.Exit(0)
		return
	}

	maxSpeed = max(200, maxSpeed) * 1024
	// var dest = "./"
	// var maxThreads = 2

	var COUNTDOWN = 20                // millisecons
	var defaultChunkSize int64 = 4096 // 4KB

	// https: //gobyexample.com/tickers
	ticker := time.NewTicker(time.Duration(COUNTDOWN) * time.Millisecond)
	chunks := make(chan int, 1)

	fmt.Println("Start download files:")
	for _, file := range os.Args[2:] {
		fmt.Printf(" * %v\n", file)
		go startSingleDownload(file, chunks, defaultChunkSize)
	}

	for _ = range ticker.C {
		chunks <- 1
		// fmt.Println("Tick at", t)
	}

	return

	// for i := 0; i < maxThreads; i++ {
	// }

	// startDownloads(files, dest, maxThreads, maxSpeed)
}

func startSingleDownload(file string, chunks <-chan int, chunkSize int64) {
	resp, err := http.Get(file)
	if err == nil {
		defer resp.Body.Close()
	} else {
		return
	}

	fmt.Printf("Start dowload file: '%v'. With size: %v bytes\n", file, resp.ContentLength)
	p := make([]byte, chunkSize)
	reader := resp.Body
	var readed int64
	var fullSize int64 = resp.ContentLength

	for {
		<-chunks
		_, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Final read chunk. File: '%v' \n", file)
				return
			}
			fmt.Println(err)
		}
		readed += chunkSize
		persents := readed * 100 / fullSize
		fmt.Printf("Read chunk. File: '%v'. (%v %%)\n", file, persents)
		runtime.Gosched()
	}
}

func startDownloads(targets []string, dest string, maxThreads int, maxSpeedBytes int) {
	resp, err := http.Get(targets[0])
	if err == nil {
		defer resp.Body.Close()

		fmt.Printf("SIZE: %v bytes\n", resp.ContentLength)

		// https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185
		// https://github.com/cavaliercoder/grab
		// https://stackoverflow.com/questions/30532886/golang-dynamic-progressbar
		// https://stackoverflow.com/questions/44318345/can-i-increase-golangs-http-stream-chunk-size
		// https://www.reddit.com/r/golang/comments/4xtsbn/help_how_to_read_files_in_blocks/
		// https://gobyexample.com/reading-files
		// Здесь блоками скачиваем файл
		chunkSize := 1024 * 1024 // 1 Mbytes
		p := make([]byte, chunkSize)
		reader := resp.Body

		for {
			_, err := reader.Read(p)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Read chunk...") //should handle any remainding bytes.
					break
				}
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Read chunk...")
			time.Sleep(100 * time.Millisecond)
		}

	}
}
