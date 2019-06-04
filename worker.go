package main

import (
	"io"
	"log"
	"net/http"
	"runtime"
)

type Worker struct {
	Chunks <-chan int
	Events chan<- TaskEvent
}

func NewWorker(chunks <-chan int, events chan TaskEvent) *Worker {
	return &Worker{Chunks: chunks, Events: events}
}

func (w *Worker) Run(tasks chan *Task) {
	for t := range tasks {
		w.StartDownload(t)
	}
}

func (w *Worker) StartDownload(t *Task) {
	url := t.Url

	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
	} else {
		w.Events <- TaskEvent{Task: t, IsFinish: true}
		return
	}

	p := make([]byte, CHUNK_SIZE)
	reader := resp.Body
	var read int64

	var fullSize int64 = resp.ContentLength

	w.Events <- TaskEvent{Task: t, IsStart: true, TotalSize: fullSize}

	for {
		<-w.Chunks
		_, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				read += int64(len(p))
				w.Events <- TaskEvent{Task: t, IsFinish: true, TotalSize: fullSize, CurrentSize: read}
				return
			}
			log.Println(err)
		}

		read += int64(CHUNK_SIZE)
		w.Events <- TaskEvent{Task: t, TotalSize: fullSize, CurrentSize: read}

		runtime.Gosched()
	}
}
