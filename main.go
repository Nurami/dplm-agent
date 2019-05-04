package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"
)

var (
	currentState int
	mainChannel  chan int
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
	fsm := FSM{}
	err := fsm.createFromJSONFile("example.json")
	if err != nil {
		panic(err)
	}

	mainChannel = make(chan int)

	go genEvent1()
	go genEvent2()
	fsm.startFSM()
}

func (fsm *FSM) startFSM() {
	for {
		event := <-mainChannel
		currentNode := fsm.StatesWithActions[currentState][event]
		currentState = currentNode.State
		for _, v := range currentNode.Actions {
			_, err := call(functions, v.Name, v.Params)
			if err != nil {
				panic(err)
			}
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

func getNameOfCurrentFunction() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	values := strings.Split(f.Name(), ".")
	return values[len(values)-1]
}

func (fsm *FSM) createFromJSONFile(filename string) (err error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &fsm)
	if err != nil {
		return err
	}
	return nil
}

//мапа "имя функции" - функция
//далее slice функций строится либо типа interface{} (реализация будет через рефлексию),
//либо конкретный интерфейс (для каждой функции должен быть создан тип func)

// есть два вида функций: совершающие считывание данных с реального мира, вызывающие определенные
// события, и контролирующие устройства(поднять, опустить..).
