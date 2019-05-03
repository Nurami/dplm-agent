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
	fmt.Println("action: pickUp")
}

func pickDown() {
	fmt.Println("action: pickDown")
}

func turnOnLight() {
	fmt.Println("action: turnOnLight")
}
