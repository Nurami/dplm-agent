package main

import (
	"fmt"
	"time"
)

var events = map[string]int{
	"genEvent1": 1,
	"genEvent2": 2,
}

func genEvent1() {
	for {
		time.Sleep(5 * time.Second)
		nameOfFunc := getNameOfCurrentFunction()
		fmt.Println(nameOfFunc)
		mainChannel <- events[nameOfFunc]
	}
}

func genEvent2() {
	for {
		time.Sleep(2 * time.Second)
		nameOfFunc := getNameOfCurrentFunction()
		fmt.Println(nameOfFunc)
		mainChannel <- events[nameOfFunc]
	}
}
