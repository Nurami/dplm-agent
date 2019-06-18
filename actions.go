package main

import (
	"github.com/stianeikeland/go-rpio"
)

var (
	functions = map[string]interface{}{
		"pickUp":      pickUp,
		"pickDown":    pickDown,
		"turnOnLight": turnOnLight,
	}
)

func pickUp() {
	log.Info("pickUp starts 0")
	rpio.Pin(7).High()
	log.Info("pickUp ends 0")
}

func pickDown() {
	log.Info("pickDown starts 0")
	rpio.Pin(7).Low()
	log.Info("pickDown ends 0")
}

func turnOnLight() {
	log.Info("turnOnLight starts 0")
	rpio.Pin(11).High()
	log.Info("turnOnLight ends 0")
}
