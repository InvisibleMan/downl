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

func NewProgresBarPresenter() *ProgresBarPresenter {
	p := &ProgresBarPresenter{}
	p.Bars = make(map[*Task]*uiprogress.Bar)
	uiprogress.Start()
	return p
}

func (p *ProgresBarPresenter) Show(events map[*Task]TaskEvent) {
	for k, v := range events {
		if v.IsFinish {
			if bar, ok := p.Bars[v.Task]; ok {
				_ = bar.Set(v.TotalKB())
			} else {
				log.Println("Something Wrong. Bar for the current Task doesn't exist")
			}

		} else if v.IsStart {
			if _, ok := p.Bars[v.Task]; !ok {
				bar := uiprogress.AddBar(v.TotalKB())
				bar.PrependFunc(func(b *uiprogress.Bar) string {
					return k.Url + ": "
				})
				p.Bars[v.Task] = bar
			} else {
				log.Println("Somethin Wrong. Bar for the current Task already exists!")
			}

		} else {
			if bar, ok := p.Bars[v.Task]; ok {
				bar.PrependFunc(func(b *uiprogress.Bar) string {
					return k.Url + ": "
				})
				_ = bar.Set(v.CurrentKB())
			} else {
				log.Println("Something Wrong. Bar for the current Task doesn't exist")
			}
		}
	}
}
