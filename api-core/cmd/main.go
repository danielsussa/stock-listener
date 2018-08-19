package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type resource struct {
	Stocks map[string]stock
}

type stock struct {
	// input data
	Kind          string
	Parent        string
	Strike        float64
	LastNotifSend time.Time
	Expiration    time.Time

	// resource data
	Price float64 `json:"-"`
	Ask   float64 `json:"-"`
	Bid   float64 `json:"-"`

	//stats data
	_vdx    float64
	_maxVdx float64
}

func (res *resource) updateStock(msgSpl []string) {
	stock := res.Stocks[msgSpl[1]]
	for i, msg := range msgSpl {
		if i%2 != 0 {
			//Last Price
			if msg == "2" {
				f, _ := strconv.ParseFloat(msgSpl[i+1], 32)
				stock.Price = f
			}
			//ASK
			if msg == "3" {
				f, _ := strconv.ParseFloat(msgSpl[i+1], 32)
				stock.Ask = f
			}
			//BID
			if msg == "4" {
				f, _ := strconv.ParseFloat(msgSpl[i+1], 32)
				stock.Bid = f
			}
		}
	}
	res.Stocks[msgSpl[1]] = stock
}

func (res *resource) calculateVDX(stockName string) {
	option := res.Stocks[stockName]

	if option.Kind == "option" {
		stock := res.Stocks[option.Parent]
		timeNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		expDif := option.Expiration.Sub(timeNow).Hours() / 24
		fmt.Println(expDif)
		vdx := (option.Price / stock.Price) * (120 - expDif) * (option.Strike - stock.Price)
		maxVdx := (option.Bid / stock.Ask) * (120 - expDif) * (option.Strike - stock.Ask)
		option._vdx = vdx
		option._maxVdx = maxVdx
	}
	res.Stocks[stockName] = option

	//=(D2/D3)*(120-D4)*(D5-D3)
	//2=PR. OPT
	//3=PR. ACAO
	//4=DIAS VENC
	//5=STRIKE

}

func convertFile() (res resource) {
	b, err := ioutil.ReadFile("/build/assets/input.json") // just pass the file name
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	res := convertFile()

	// connect to this socket
	conn, err := net.Dial("tcp", os.Getenv("TCP_URL"))

	if err != nil {
		panic(err)
	}

	// send to socket
	fmt.Fprintf(conn, "olaaa\n")

	for {
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
		msgSpl := strings.Split(message, ":")
		if len(msgSpl) < 3 {
			continue
		}
		res.updateStock(msgSpl)
		res.calculateVDX(msgSpl[1])

	}
}
