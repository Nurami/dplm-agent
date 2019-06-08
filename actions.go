package main

import "fmt"

//мапа функций
var functions = map[string]interface{}{
	"pickUp":      pickUp,
	"pickDown":    pickDown,
	"turnOnLight": turnOnLight,
}

//второй вид функций
func pickUp() {
	log.Info("pickUp starts 0")
	fmt.Println("action: pickUp")
	log.Info("pickUp ends 0")
}

func pickDown() {
	log.Info("pickDown starts 0")
	fmt.Println("action: pickDown")
	log.Info("pickDown ends 0")
}

func turnOnLight() {
	log.Info("turnOnLight starts 0")
	fmt.Println("action: turnOnLight")
	log.Info("turnOnLight ends 0")
}
