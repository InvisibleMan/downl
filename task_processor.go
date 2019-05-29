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
	// Destination
	tasks chan string
	wg    sync.WaitGroup
}

func (tp *TaskProcessor) Download(targets []string) {
	tp.wg = sync.WaitGroup{}
	tp.wg.Add(len(targets))
	tp.tasks = make(chan string, tp.Threads)

	go tp.runWorkers()
	go tp.putTasks(targets)

	tp.wg.Wait()
}

func (tp *TaskProcessor) putTasks(targets []string) {
	defer close(tp.tasks)

	for _, target := range targets {
		tp.tasks <- target
	}
}

func (tp *TaskProcessor) runWorkers() {
	chanks := tp.startChanksQueue()

	log.Println("Start download files:")
	for task := range tp.tasks {
		go tp.downloadTarget(task, chanks)
	}
}

func (tp *TaskProcessor) startChanksQueue() <-chan int {
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

func (tp *TaskProcessor) doneTarget() {
	tp.wg.Done()
}

func (tp *TaskProcessor) downloadTarget(file string, chunks <-chan int) {
	resp, err := http.Get(file)
	if err == nil {
		defer resp.Body.Close()
	} else {
		tp.doneTarget()
		return
	}

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
				tp.doneTarget()
				return
			}
			log.Println(err)
		}
		read += int64(CHUNK_SIZE)
		percent := read * 100 / fullSize
		if percent-lastPercent > diffPercent {
			log.Printf("Read chunk. File: '%v'. (%v %%)\n", file, percent)
			lastPercent = percent
		}
		runtime.Gosched()
	}
}
