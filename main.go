package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/op/go-logging"
)

const agentID = "0"

var (
	currentState int
	mainChannel  chan int
	log          = logging.MustGetLogger("logger")
	logsFormat   = logging.MustStringFormatter(
		agentID + ` %{time:2006-01-02 15:04:05} %{shortfunc} %{level:s} %{id:d} %{message}`,
	)
	mutex = &sync.Mutex{}
	url   = "http://localhost:8080/logs"
)

//FSM - это структура, описывающая конечный автомат
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
	go logToANewFileByPeriod(10)
	//TODO: обдумать и реализовать ожидание(waitgroup?), если это возможно
	//Главному потоку Необходимо поспать, чтобы переменная логирования успела инициализироваться
	time.Sleep(time.Second)
	go sendLogsToServerByPeriod(20)

	fsm := FSM{}
	err := fsm.createFromJSONFile("example.json")
	if err != nil {
		log.Panic(err)
	}

	mainChannel = make(chan int)

	go genEvent1()
	go genEvent2()
	fsm.startFSM()
}

func logToANewFileByPeriod(period int) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	for {
		logsName := "logs/" + time.Now().Format("2006.01.02-15.04.05") + ".log"
		file, err := os.Create(logsName)
		mutex.Lock()
		log = logging.MustGetLogger(logsName)
		var backend *logging.LogBackend
		if err != nil {
			fmt.Println(time.Now(), " ", err)
			backend = logging.NewLogBackend(os.Stdout, "", 0)
		} else {
			backend = logging.NewLogBackend(file, "", 0)
		}
		backendFormatter := logging.NewBackendFormatter(backend, logsFormat)
		logging.SetBackend(backendFormatter)
		mutex.Unlock()
		time.Sleep(time.Duration(period) * time.Second)
		file.Close()
	}
}

func (fsm *FSM) startFSM() {
	for {
		event := <-mainChannel
		currentNode := fsm.StatesWithActions[currentState][event]
		currentState = currentNode.State
		for _, v := range currentNode.Actions {
			_, err := call(functions, v.Name, v.Params)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

func call(m map[string]interface{}, name string, params []int) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("Число параметров не верно в функции " + name)
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
		return
	}
	err = json.Unmarshal(file, &fsm)
	if err != nil {
		return
	}
	return nil
}

func sendLogsToServerByPeriod(period int) {
	for {
		time.Sleep(time.Duration(period) * time.Second)
		files, err := ioutil.ReadDir("logs/")
		if err != nil {
			log.Critical(err)
		}
		for i := 0; i < len(files)-1; i++ {
			filename := files[i].Name()
			err = upload("logs/" + filename)
			if err != nil {
				log.Critical(err)
			} else {
				err = os.Remove("logs/" + filename)
				if err != nil {
					log.Critical(err)
				}
			}
		}
	}
}

func upload(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	resp, err := http.Post(url, "binary/octet-stream", file)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	message, _ := ioutil.ReadAll(resp.Body)
	//TODO: убрать фмт и в лог
	fmt.Println(string(message))
	return nil
}

//мапа "имя функции" - функция
//далее slice функций строится либо типа interface{} (реализация будет через рефлексию),
//либо конкретный интерфейс (для каждой функции должен быть создан тип func)

// есть два вида функций: совершающие считывание данных с реального мира, вызывающие определенные
// события, и контролирующие устройства(поднять, опустить..).
