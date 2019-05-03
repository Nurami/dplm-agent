package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
)

var (
	currentState int
)

type FSM struct {
	StartingState     int                  `json:"startingState"`
	StatesWithActions [][]stateWithActions `json:"statesWithActions"`
}

type stateWithActions struct {
	State   int      `json:"state"`
	Actions []action `json:"actions"`
}

type action struct {
	Name   string `json:"name"`
	Params []int  `json:"params"`
}

func main() {
	//считывание json
	file, err := ioutil.ReadFile("example.json")
	check(err)
	fsm := FSM{}
	err = json.Unmarshal(file, &fsm)
	check(err)

	k := make(chan int)
	go generateEvent(k)
	fsm.startFSM(k)
}

func (data *FSM) startFSM(ch chan int) {
	for {
		event := <-ch
		currentNode := data.StatesWithActions[currentState][event]
		currentState = currentNode.State
		for _, v := range currentNode.Actions {
			_, err := call(functions, v.Name, v.Params)
			check(err)
		}
	}
}

func call(m map[string]interface{}, name string, params []int) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

//мапа "имя функции" - функция
//далее slice функций строится либо типа interface{} (реализация будет через рефлексию),
//либо конкретный интерфейс (для каждой функции должен быть создан тип func)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// есть два вида функций: совершающие считывание данных с реального мира, вызывающие определенные
// события, и контролирующие устройства(поднять, опустить..).
