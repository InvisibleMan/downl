package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello Download Master")

	var files = []string{
		"http://192.168.1.100/test_file1.mp4",
		"http://192.168.1.100/test_file2.mp4"}

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
		// Здесь блоками скачиваем файл
		chunkSize := 1024

		p := make([]byte, chunkSize)

		for {
			n, err := reader.Read(p)
			if err != nil {
				if err == io.EOF {
					fmt.Println(string(p[:n])) //should handle any remainding bytes.
					break
				}
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(p[:n]))
		}

	}
}
