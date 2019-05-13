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

var (
	currentState int
	mainChannel  chan int
	log          = logging.MustGetLogger("logger")
	//TODO: впихнуть необходим номер агента
	logsFormat = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{shortfunc} %{level:s} %{id:d} %{message}`,
	)
	mutex = &sync.Mutex{}
	url   = ""
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
	go logByPeriod(10)
	//TODO: обдумать и реализовать ожидание(waitgroup?), если это возможно
	//Главному потоку Необходимо поспать, чтобы переменная логирования успела инициализироваться
	time.Sleep(time.Second)

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

//TODO: удаление отправленных файлов
//TODO: отправка данных на сервер

func logByPeriod(duration int) {
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
		time.Sleep(time.Duration(duration) * time.Second)
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

func sendLogsToServerByPeriod(ip string, period int) {
	files, err := ioutil.ReadDir("logs/")
	//TODO: обработка ошибки
	if err != nil {
	}
	for i := 0; i < len(files)-1; i++ {
		filename := files[i].Name()
		err = upload(filename)
		//TODO: обработка ошибки
		if err != nil {
		}
	}
}

func upload(filename string) error {
	file, err := os.Open(filename)
	//TODO: обработка ошибки
	if err != nil {
		panic(err)
	}
	defer file.Close()
	resp, err := http.Post(url, "binary/octet-stream", file)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	message, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(message))
	//TODO: возврат ошибки
	return nil
}

//мапа "имя функции" - функция
//далее slice функций строится либо типа interface{} (реализация будет через рефлексию),
//либо конкретный интерфейс (для каждой функции должен быть создан тип func)

// есть два вида функций: совершающие считывание данных с реального мира, вызывающие определенные
// события, и контролирующие устройства(поднять, опустить..).
