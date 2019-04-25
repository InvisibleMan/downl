package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello Download Master")

	var files = []string{
		"http://localhost:8083/file1.dat",
		"http://localhost:8083/file2.dat"}

	var dest = "./"
	var maxThreads = 2
	var maxSpeed = 500 * 1024

	startDownloads(files, dest, maxThreads, maxSpeed)
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
