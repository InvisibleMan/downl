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

	// totalCount int
	Events chan TaskEvent
	Chunks chan int

	mux           sync.Mutex
	LastEvents    map[*Task]TaskEvent
	presenterStop chan int
	presenterQuit chan int
	Presenter     Presenter
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
	outTasks := make(chan *Task)
	inTasks := make([]*Task, 0, len(targets))

	for _, t := range targets {
		inTasks = append(inTasks, &Task{Url: t})
	}

	tp.startChunksQueue(stopWorld)
	tp.runWorkers(outTasks)
	log.Println("Start download files:")

	go tp.startEventsListen()
	go tp.putTasks(inTasks)

	for i := 0; i < len(targets); i++ {
		<-outTasks
	}
	close(outTasks)
	stopWorld <- 1
	tp.waitPresenter()
	// time.Sleep(3 * time.Second)
}

func (tp *TaskProcessor) runWorkers(outTasks chan *Task) {
	for i := 0; i < tp.Threads; i++ {
		go func() {
			NewWorker(tp.Chunks, tp.Events).Run(tp.Tasks, outTasks)
		}()
	}
}

func (tp *TaskProcessor) putTasks(tasks []*Task) {
	for _, t := range tasks {
		tp.Tasks <- t
	}

	close(tp.Tasks)
}

func (tp *TaskProcessor) startChunksQueue(stop chan int) {
	var countdown = int(time.Second) / (tp.Speed * KB / CHUNK_SIZE)

	// log.Printf("Duration ONE SECOND: '%d'\n", time.Duration(1*time.Second))
	// log.Printf("Speed: '%d'\n", tp.Speed)
	// log.Printf("CHUNK_SIZE: '%d'\n", CHUNK_SIZE/KB)
	// log.Printf("Countdown: '%d'\n", countdown)

	ticker := time.NewTicker(time.Duration(countdown))

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				tp.Chunks <- 1

			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (tp *TaskProcessor) SetPresenter(presenter Presenter, refreshMilliseconds int) {
	ticker := time.NewTicker(time.Duration(refreshMilliseconds) * time.Millisecond)
	tp.presenterStop = make(chan int)
	tp.presenterQuit = make(chan int)

	updater := func() {
		tp.mux.Lock()
		if len(tp.LastEvents) > 0 {
			events := tp.LastEvents
			presenter.Show(events)
			tp.LastEvents = make(map[*Task]TaskEvent)
		}
		tp.mux.Unlock()
	}

	go func() {
		for {
			select {
			case <-tp.presenterStop:
				ticker.Stop()
				updater()

				tp.presenterQuit <- 1
				return
			case <-ticker.C:
				updater()
			}
		}
	}()

	tp.Presenter = presenter
}

func (tp *TaskProcessor) waitPresenter() {
	close(tp.Events)
	<-tp.presenterQuit
	tp.Presenter.Stop()
}

func (tp *TaskProcessor) stopPresenter() {
	tp.presenterStop <- 1
}

func (tp *TaskProcessor) startEventsListen() {
	for event := range tp.Events {
		tp.mux.Lock()

		if lastE, ok := tp.LastEvents[event.Task]; ok {
			event.IsStart = lastE.IsStart
		}
		tp.LastEvents[event.Task] = event

		tp.mux.Unlock()
	}
	tp.stopPresenter()
}

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
