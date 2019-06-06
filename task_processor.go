package main

import (
	"log"
	"sync"
	"time"
)

var KB int = 1024
var CHUNK_SIZE int = 4 * KB // 4KB

type Task struct {
	Url string
}

type TaskEvent struct {
	IsStart     bool
	IsFinish    bool
	TotalSize   int64
	CurrentSize int64
	Task        *Task
}

type TaskProcessor struct {
	Speed   int // max speed limit in KB
	Threads int // max concurent downloads
	// Destination
	Tasks chan *Task
	// wg         sync.WaitGroup
	// totalCount int
	Events chan TaskEvent
	Chunks chan int

	mux        sync.Mutex
	LastEvents map[*Task]TaskEvent
}

func (e *TaskEvent) TotalKB() int {
	return int(e.TotalSize / int64(KB))
}

func (e *TaskEvent) CurrentKB() int {
	return int(e.CurrentSize / int64(KB))
}

func TaskProcessorNew(maxSpeed int, maxThreads int) *TaskProcessor {
	tp := TaskProcessor{Speed: maxSpeed, Threads: maxThreads}

	tp.Tasks = make(chan *Task)
	tp.Chunks = make(chan int)
	tp.Events = make(chan TaskEvent)
	tp.LastEvents = make(map[*Task]TaskEvent)

	return &tp
}

func (tp *TaskProcessor) Download(targets []string) {
	stopWorld := make(chan int)

	tp.startChunksQueue(stopWorld)
	wg := tp.runWorkers(stopWorld)
	log.Println("Start download files:")

	go tp.startEventsListen()
	go tp.putTasks(targets)

	wg.Wait()
	close(stopWorld)
}

func (tp *TaskProcessor) runWorkers(stopWorld chan int) *sync.WaitGroup {
	wg := sync.WaitGroup{}

	for i := 0; i < tp.Threads; i++ {
		wg.Add(1)
		go func() {
			NewWorker(tp.Chunks, tp.Events).Run(tp.Tasks)
			wg.Done()
		}()
	}
	return &wg
}

func (tp *TaskProcessor) putTasks(targets []string) {
	for _, target := range targets {
		tp.Tasks <- &Task{Url: target}
	}

	close(tp.Tasks)
}

func (tp *TaskProcessor) startChunksQueue(stopWorld chan int) {
	var countdown = int(time.Second) / (tp.Speed * KB / CHUNK_SIZE)

	log.Printf("Duration ONE SECOND: '%d'\n", time.Duration(1*time.Second))
	log.Printf("Speed: '%d'\n", tp.Speed)
	log.Printf("CHUNK_SIZE: '%d'\n", CHUNK_SIZE/KB)
	log.Printf("Countdown: '%d'\n", countdown)

	ticker := time.NewTicker(time.Duration(countdown))

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				tp.Chunks <- 1

			case <-stopWorld:
				return
			}
		}
	}()
}

func (tp *TaskProcessor) SetPresenter(presenter Presenter, refreshMilliseconds int) {
	ticker := time.NewTicker(time.Duration(refreshMilliseconds) * time.Millisecond)

	go func() {
		for {
			<-ticker.C
			tp.mux.Lock()
			if len(tp.LastEvents) > 0 {
				presenter.Show(tp.LastEvents)
				tp.LastEvents = make(map[*Task]TaskEvent)
			}
			tp.mux.Unlock()
		}
	}()
}

func (tp *TaskProcessor) startEventsListen() {
	for event := range tp.Events {
		tp.mux.Lock()
		if lastE, ok := tp.LastEvents[event.Task]; ok {
			lastE.IsFinish = event.IsFinish
			lastE.CurrentSize = event.CurrentSize
		} else {
			tp.LastEvents[event.Task] = event
		}

		tp.mux.Unlock()
	}
}

// func (tp *TaskProcessor) processUsedChunks() chan int {
// 	tp.used = make(chan int)

// 	go func() {
// 		for chunk := range tp.used {
// 			tp.addDownloadCount(chunk)
// 		}
// 	}()

// 	return tp.used
// }

// func (tp *TaskProcessor) addDownloadCount(chunk int) {
// 	tp.lastIntervalCount = tp.lastIntervalCount + chunk
// 	tp.totalCount = tp.totalCount + chunk
// }

// func printStatus(tp *TaskProcessor) {
// 	ticker := time.NewTicker(1000 * time.Millisecond)
// 	go func() {
// 		prevCount := 0
// 		for _ = range ticker.C {
// 			log.Printf("Current download speed: %d KB\n", (tp.totalCount-prevCount)/KB)
// 			prevCount = tp.totalCount
// 		}
// 	}()

// 	// total download size
// 	// current speed = count for interval / interval
// }
