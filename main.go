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

	endEvents := make(chan int)
	startDownloads(params.Targets, params.MaxSpeed, endEvents)

	<-endEvents
}
