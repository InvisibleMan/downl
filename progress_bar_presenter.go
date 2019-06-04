package main

import "log"

type ProgresBarPresenter struct {
}

func (notifer *ProgresBarPresenter) Show(events map[*Task]TaskEvent) {
	for k, v := range events {
		if v.IsFinish {
			log.Printf("Finish dowload file: '%v'\n", k.Url)
		} else if v.IsStart {
			log.Printf("Start dowload file: '%v'. With size: %v bytes\n", k.Url, v.TotalSize)

		} else {
			var percent int64 = v.CurrentSize * 100 / v.TotalSize
			log.Printf("Process download File: '%v'. (%v %%)\n", k.Url, percent)
		}
	}
}
