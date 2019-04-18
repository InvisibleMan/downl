package main

import (
	"fmt"
	"net/http"
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

		// io.

		// resp.Body
	}
}
