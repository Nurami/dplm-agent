package jo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

const (
	state1 = iota
	state2
	state3
)
const (
	sign1 = iota
	sign2
	sign3
)

var workers = [9]func(int, int){}

type fsm struct {
	FsmTableFromJSON [3][3]int `json:"fsmTable"`
	StartState       int       `json:"startState"`
}
type workingState struct {
	work  func(int, int)
	state int
}

var currentSignal int

func main() {
	// tableFSM := [3][3]int{}
	// tableFSM[state1][sign1], tableFSM[state1][sign2], tableFSM[state1][sign3] = state2, state3, state1
	// tableFSM[state2][sign1], tableFSM[state2][sign2], tableFSM[state3][sign3] = -1, state2, state3
	// tableFSM[state3][sign1], tableFSM[state3][sign2], tableFSM[state3][sign3] = state1, -1, state2
	for i := range workers {
		workers[i] = func(state, sign int) {
			fmt.Printf("[%d,%d] work! \n", state, sign)
			time.Sleep(time.Second)
		}
	}

	FSM := fsm{}
	example, err := ioutil.ReadFile("example.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(example, &FSM)

	count := 0
	fsmTable := [3][3]workingState{}
	for i, v := range FSM.FsmTableFromJSON {
		for j, k := range v {
			fsmTable[i][j] = workingState{workers[count], k}
			count++
		}
	}
	ch := make(chan int)
	go doSignal(ch)
	doTableFSM(fsmTable, FSM.StartState, ch)
}

func doTableFSM(fsmTable [3][3]workingState, currentState int, ch chan int) {
	for {
		currentSignal = <-ch
		currentState = fsmTable[currentState][currentSignal].state
		fsmTable[currentState][currentSignal].work(currentState, currentSignal)
	}
}

func doSignal(ch chan int) {
	for {
		currentSignal = rand.Intn(3)
		ch <- currentSignal
		fmt.Printf("%d was generated\n", currentSignal)
		time.Sleep(5 * time.Second)
	}
}
