package main

import (
	"errors"
	"log"
	"os"
	"strconv"
)

type InParams struct {
	MaxSpeed int      // in KB/s
	Targets  []string // urls for files downloads
}

func max(l, r int) int {
	if l > r {
		return l
	}
	return r
}

func GetInParams() (*InParams, error) {
	if len(os.Args) < 3 {
		log.Println("Noting download. Bye")
		return nil, errors.New("Noting download")
	}

	var maxSpeed int

	maxSpeed, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Println("Wrong value of MaxSpeed. Bye")
		return nil, errors.New("Wrong value of MaxSpeed")
	}

	maxSpeed = max(1, maxSpeed)

	return &InParams{
		MaxSpeed: maxSpeed,
		Targets:  os.Args[1:],
	}, nil
}
