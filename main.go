package main

import (
	"log"
	"os"
)

func main() {
	log.Printf(
		"Starting the downloader. Version: %s (commit: %s, build time: %s).",
		Version, Commit, BuildTime,
	)

	params, err := GetInParams()
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
		return
	}

	tp := TaskProcessor{Speed: params.MaxSpeed, Threads: 2}
	tp.Download(params.Targets)

	// startDownloads(params.Targets, params.MaxSpeed, endEvents)

	// endEvents := make(chan int)
	// <-endEvents
}

// Useful links

// https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185
// https://github.com/cavaliercoder/grab

// https://stackoverflow.com/questions/30532886/golang-dynamic-progressbar
// https://github.com/gosuri/uiprogress

// https://stackoverflow.com/questions/44318345/can-i-increase-golangs-http-stream-chunk-size
// https://www.reddit.com/r/golang/comments/4xtsbn/help_how_to_read_files_in_blocks/
// https://gobyexample.com/reading-files
