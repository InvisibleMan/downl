package main

import (
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var KB int = 1024
var CHUNK_SIZE int = 4 * KB // 4KB

type TaskProcessor struct {
	Speed   int // max speed limit in KB
	Threads int // max concurent downloads

	sync.WaitGroup

	// Destination
	tasks             chan string
	used              chan int
	lastIntervalCount int
	totalCount        int
}

func (tp *TaskProcessor) Download(targets []string) {
	tp.Add(len(targets))
	tp.tasks = make(chan string, tp.Threads)

	go tp.runWorkers()
	go tp.putTasks(targets)
	go printStatus(tp)

	tp.Wait()
}

func (tp *TaskProcessor) putTasks(targets []string) {
	defer close(tp.tasks)

	for _, target := range targets {
		tp.tasks <- target
	}
}

func (tp *TaskProcessor) runWorkers() {
	chunks := tp.startChunksQueue()
	usedChunks := tp.processUsedChunks()

	log.Println("Start download files:")
	for task := range tp.tasks {
		go tp.downloadTarget(task, chunks, usedChunks)
	}
}

func (tp *TaskProcessor) startChunksQueue() <-chan int {
	var countdown = int(time.Second) / (tp.Speed * KB / CHUNK_SIZE)

	log.Printf("Duration ONE SECOND: '%d'\n", time.Duration(1*time.Second))
	log.Printf("Speed: '%d'\n", tp.Speed)
	log.Printf("CHUNK_SIZE: '%d'\n", CHUNK_SIZE/KB)
	log.Printf("Countdown: '%d'\n", countdown)

	ticker := time.NewTicker(time.Duration(countdown))
	queueChunks := make(chan int, 1)

	go func() {
		for _ = range ticker.C {
			queueChunks <- 1
		}
	}()

	return queueChunks
}

func (tp *TaskProcessor) downloadTarget(file string, chunks <-chan int, chunksUsed chan int) {
	resp, err := http.Get(file)
	if err != nil {
		tp.Done()
		return
	}

	defer resp.Body.Close()

	log.Printf("Start dowload file: '%v'. With size: %v bytes\n", file, resp.ContentLength)
	p := make([]byte, CHUNK_SIZE)
	reader := resp.Body
	var read int64

	var fullSize int64 = resp.ContentLength
	var lastPercent int64 = 0
	var diffPercent int64 = 0

	for {
		<-chunks
		_, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				log.Printf("Final read chunk. File: '%v' \n", file)
				tp.Done()
				return
			}
			log.Println(err)
		}

		chunksUsed <- CHUNK_SIZE
		read += int64(CHUNK_SIZE)
		percent := read * 100 / fullSize
		if percent-lastPercent > diffPercent {
			// log.Printf("Read chunk. File: '%v'. (%v %%)\n", file, percent)
			lastPercent = percent
		}
		runtime.Gosched()
	}
}

func (tp *TaskProcessor) processUsedChunks() chan int {
	tp.used = make(chan int)

	go func() {
		for chunk := range tp.used {
			tp.addDownloadCount(chunk)
		}
	}()

	return tp.used
}

func (tp *TaskProcessor) addDownloadCount(chunk int) {
	// sync.Mutex
	tp.lastIntervalCount = tp.lastIntervalCount + chunk
	tp.totalCount = tp.totalCount + chunk
}

func printStatus(tp *TaskProcessor) {

	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		prevCount := 0
		for _ = range ticker.C {
			log.Printf("Current download speed: %d KB\n", (tp.totalCount-prevCount)/KB)
			prevCount = tp.totalCount
		}
	}()

	// total download size
	// current speed = count for interval / interval
}
