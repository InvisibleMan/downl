package main

import (
	"log"

	"github.com/gosuri/uiprogress"
)

type ProgresBarPresenter struct {
	Bars map[*Task]*uiprogress.Bar
}

type TextLogPresenter struct {
}

func (p *TextLogPresenter) Show(events map[*Task]TaskEvent) {
	for k, v := range events {
		if v.IsStart {
			log.Printf("Start dowload file: '%v'. With size: %v bytes\n", k.Url, v.TotalSize)
		}
		if v.IsFinish {
			log.Printf("Finish dowload file: '%v'\n", k.Url)
		}

		if !(v.IsStart || v.IsFinish) {
			var percent int64 = v.CurrentSize * 100 / v.TotalSize
			log.Printf("Process download File: '%v'. (%v %%)\n", k.Url, percent)
		}
	}
}

func (p *TextLogPresenter) Stop() {

}

func NewProgresBarPresenter() *ProgresBarPresenter {
	p := &ProgresBarPresenter{}
	p.Bars = make(map[*Task]*uiprogress.Bar)
	uiprogress.Start()
	return p
}

func (p *ProgresBarPresenter) Show(events map[*Task]TaskEvent) {
	for task, event := range events {
		if event.IsFinish {
			if bar, ok := p.Bars[task]; ok {
				_ = bar.Set(event.TotalKB())
			} else {
				log.Println("Something Wrong. Bar for the current Task doesn't exist")
			}

		} else if event.IsStart {
			if _, ok := p.Bars[task]; !ok {
				total := event.TotalKB()
				url := task.Url

				bar := uiprogress.AddBar(total).AppendCompleted()
				bar.PrependFunc(func(b *uiprogress.Bar) string {
					return url + ": "
				})
				p.Bars[task] = bar
			} else {
				log.Println("Somethin Wrong. Bar for the current Task already exists!")
			}
		} else {
			if bar, ok := p.Bars[task]; ok {
				_ = bar.Set(event.CurrentKB())
			} else {
				log.Println("Something Wrong. Bar for the current Task doesn't exist")
			}
		}
	}
}

func (p *ProgresBarPresenter) Stop() {
	uiprogress.Stop()
}
