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
	log.Info("pickUp starts")
	fmt.Println("action: pickUp")
	log.Info("pickUp ends")
}

func pickDown() {
	log.Info("pickDown starts")
	fmt.Println("action: pickDown")
	log.Info("pickDown ends")
}

func turnOnLight() {
	log.Info("turnOnLight starts")
	fmt.Println("action: turnOnLight")
	log.Info("turnOnLight ends")
}
