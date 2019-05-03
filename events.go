package main

import (
	"fmt"
	"math/rand"
	"time"
)

//первый вид, абстрактный пример
func generateEvent(ch chan int) {
	for {
		//3 - число максимально возможных событий
		rand := rand.Intn(3)
		ch <- rand
		fmt.Printf("%d was generated\n", rand)
		time.Sleep(5 * time.Second)
	}
}

func genEvent1() {

}

func genEvent2() {

}
