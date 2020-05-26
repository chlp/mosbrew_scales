package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"go.bug.st/serial.v1"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GoMosbrewConfig struct {
	name     string
	baudRate int
	httpPort int
}

func getConfig(fileName string) GoMosbrewConfig {
	fmt.Printf("Загружаем конфиг %v.\r\n", fileName)
	p := GoMosbrewConfig{name: "", baudRate: 0, httpPort: 0}
	comPortName := ""
	comPortBaudRateStr := ""
	httpPortStr := ""
	dat, err := ioutil.ReadFile(fileName)
	needWriteConfig := false
	if err != nil {
		fmt.Printf("Не получается прочесть конфиг %v. Создаем новый.\r\n", fileName)
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			fmt.Printf("Не найдены COM порты\r\n")
		}
		for _, port := range ports {
			fmt.Printf("Найден порт: %v\r\n", port)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Введите COM PORT (пример: COM7): ")
		comPortName, _ = reader.ReadString('\n')
		fmt.Print("Введите BAUD RATE (пример: 38400, 9600): ")
		comPortBaudRateStr, _ = reader.ReadString('\n')
		httpPortStr = "8844"
		needWriteConfig = true
	} else {
		configData := strings.Split(string(dat), "\n")
		if len(configData) < 3 {
			fmt.Printf("Конфиг %v должен содержать 3 строки\r\n", fileName)
			time.Sleep(time.Second * 5)
			os.Exit(1)
		}
		comPortName = configData[0]
		comPortBaudRateStr = configData[1]
		httpPortStr = configData[2]
	}
	reg, err := regexp.Compile("[^/._a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	comPortName = reg.ReplaceAllString(comPortName, "")
	comPortBaudRateStr = reg.ReplaceAllString(comPortBaudRateStr, "")
	httpPortStr = reg.ReplaceAllString(httpPortStr, "")
	if needWriteConfig {
		_ = ioutil.WriteFile(fileName, []byte(fmt.Sprintf("%v\r\n%v\r\n%v", comPortName, comPortBaudRateStr, httpPortStr)), 0644)
	}
	comPortBaudRate, err := strconv.Atoi(comPortBaudRateStr)
	if err != nil {
		fmt.Printf("Недопустимый BAUD RATE: %v.\r\n", comPortBaudRateStr)
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		fmt.Printf("Недопустимый HTTP PORT: %v.\r\n", httpPortStr)
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	p.name = comPortName
	p.baudRate = comPortBaudRate
	p.httpPort = httpPort
	return p
}

func weightFromBug(buff []byte) int {
	weight1, _ := strconv.Atoi(fmt.Sprintf("%x", buff[0]))
	weight2, _ := strconv.Atoi(fmt.Sprintf("%x", uint(buff[1])))
	weight2 *= 100
	weight3, _ := strconv.Atoi(fmt.Sprintf("%x", uint(buff[2])))
	weight3 *= 10000
	return weight1 + weight2 + weight3
}

func main() {
	log.Println("Старт программы")

	ex, err := os.Executable()
	if err != nil {
		log.Println("Проблемы с ex")
		log.Fatal(err)
	}
	path := filepath.Dir(ex)
	if err != nil {
		log.Println("Проблемы с path")
		log.Fatal(err)
	}
	var config GoMosbrewConfig
	configPath := path + string(os.PathSeparator) + "config.txt"
	config = getConfig(configPath)

	weight := 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%v", weight)
	})
	var httpListener func()
	httpListener = func() {
		log.Printf("Запускаем HTTP-сервер %v\r\n", config.httpPort)
		err = http.ListenAndServe(":"+strconv.Itoa(config.httpPort), nil)
		if err != nil {
			log.Printf("Проблемы с HTTP %v\r\n", config.httpPort)
			log.Println(err)
			time.Sleep(time.Second * 3)
			go httpListener()
			return
		}
	}
	go httpListener()

	log.Printf("Используем COM PORT: %v; BAUD RATE: %v; HTTP PORT: %v\r\n", config.name, config.baudRate, config.httpPort)
	mode := &serial.Mode{
		BaudRate: config.baudRate,
	}

	var comPortListener func()
	comPortListener = func() {
		comPort, err := serial.Open(config.name, mode)
		if err != nil {
			log.Printf("Не получается работать с COM PORT %v; BAUD RATE: %v.\r\n", config.name, config.baudRate)
			time.Sleep(time.Second * 3)
			go comPortListener()
			return
		}

		log.Printf("Слушаем COM PORT %v : %v\r\n", config.name, config.baudRate)
		buff := make([]byte, 6)
		for {
			n, err := comPort.Read(buff)
			if err != nil {
				log.Println("проблемы с данными")
				log.Println(err)
				time.Sleep(time.Second * 3)
				go comPortListener()
				return
			}
			if n == 0 {
				log.Println("окончание данных")
				time.Sleep(time.Second * 3)
				go comPortListener()
				return
			}
			weightTmp := weightFromBug(buff)
			log.Printf("Данные получены: %v :", weightTmp)
			for _, n := range buff {
				fmt.Printf("%08b ", n)
			}
			if weightTmp >= 0 {
				fmt.Printf(" : корретно")
				if math.Abs(float64(weight-weightTmp)) > 10 {
					weight = weightTmp
				} else {
					fmt.Printf(", но значение близко к предыдущему")
				}
				fmt.Printf("\r\n")
			} else {
				fmt.Printf(" : некорректно: %v\r\n", err)
			}
		}
	}
	go comPortListener()

	for {
		log.Println("продолжаем выполнение")
		time.Sleep(time.Second * 10)
	}
}