package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

type stockInfo interface {
	Kind() string
	SetPrice(string)
	SetMarketStatus(string)
	SetName(string)
	Perform()
}

type option struct {
	Price          float64
	Name           string
	Stock          *stock
	MinProfit      float64
	MaxProfit      float64
	Strike         float64
	Expiration     float64
	Volume         float64
	IsMarketOpen   bool
	ExpirationDate time.Time
	Vdx            float64
}

func (opt option) Kind() string {
	return "option"
}

func (opt *option) SetPrice(p string) {
	price, _ := strconv.ParseFloat(p, 32)
	opt.Price = price
}

func (opt *option) SetMarketStatus(p string) {
	if p == "0" {
		opt.IsMarketOpen = true
	}
	opt.IsMarketOpen = false
}

func (opt *option) SetName(p string) {
	opt.Name = p
}

func (opt *option) Perform() {
	//Perform VDX
	if opt.Stock == nil {
		return
	}
	opt.Vdx = (opt.Price / opt.Stock.Price) * (120 - opt.Expiration) * (opt.Strike - opt.Stock.Price)

	// 	option.MinProfit = option.Price / stock.Price
	// 	option.MaxProfit = (option.Strike + option.Price - stock.Price) / stock.Price
	opt.MinProfit = opt.Price / opt.Stock.Price
	opt.MaxProfit = (opt.Strike + opt.Price - opt.Stock.Price) / opt.Stock.Price
}

type stock struct {
	Price  float64
	Name   string
	Volume float64
}

func (st stock) Kind() string {
	return "stock"
}

func (st *stock) SetPrice(p string) {
	price, _ := strconv.ParseFloat(p, 32)
	st.Price = price
}

func (st *stock) SetMarketStatus(p string) {

}

func (st *stock) SetName(p string) {
	st.Name = p
}

func (st *stock) Perform() {
}

var stockMap map[string]stockInfo
var msgList []string

func main() {
	stocks, options := convertFile()

	//Add to stockMap
	stockMap = make(map[string]stockInfo)
	for _, st := range stocks {
		stockMap[strings.ToUpper(st)] = &stock{}
	}
	for _, opt := range options {
		stockMap[strings.ToUpper(opt)] = &option{}
	}

	// connect to this socket
	conn, err := net.Dial("tcp", "datafeeddl1.cedrofinances.com.br:81")

	if err != nil {
		panic(err)
	}

	// send to socket
	fmt.Fprintf(conn, "\n")
	fmt.Fprintf(conn, "kanczuk\n")
	fmt.Fprintf(conn, "102030\n")

	go sendAllMessages(&conn, stocks, options)
	go serveWeb()

	for {
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')

		msgSpl := strings.Split(message, ":")
		saveMsgToFile()

		if len(msgSpl) < 3 || msgSpl[0] == "E" {
			continue
		}
		if msgSpl[0] != "T" {
			continue
		}

		fmt.Print("Message from server: " + message)
		msgList = append(msgList, message)

		msgMap := transformMsgIntoMap(msgSpl)

		//Setup name
		stockMap[msgMap["name"]].SetName(msgMap["name"])

		if _, ok := msgMap["45"]; ok {
			//If is a option
			if msgMap["45"] == "2" {
				if _, ok := msgMap["81"]; ok {
					st := stockMap[msgMap["81"]].(*stock)
					stockMap[msgMap["name"]].(*option).Stock = st
				}
			}
		}

		//Setup price
		if _, ok := msgMap["2"]; ok {
			stockMap[msgMap["name"]].SetPrice(msgMap["2"])
		}

		//Setup Market Status
		if _, ok := msgMap["84"]; ok {
			stockMap[msgMap["name"]].SetMarketStatus(msgMap["84"])
		}

		//Setup strike
		if _, ok := msgMap["121"]; ok {
			strike, _ := strconv.ParseFloat(msgMap["121"], 32)
			if stockMap[msgMap["name"]].Kind() == "option" {
				stockMap[msgMap["name"]].(*option).Strike = strike
			}
		}

		//Setup volume
		if _, ok := msgMap["9"]; ok {
			vol, _ := strconv.ParseFloat(msgMap["9"], 32)
			if stockMap[msgMap["name"]].Kind() == "option" {
				stockMap[msgMap["name"]].(*option).Volume = vol
			}
		}

		//Setup Expiration date
		if _, ok := msgMap["125"]; ok {
			tExp, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", msgMap["125"][0:4], msgMap["125"][4:6], msgMap["125"][6:8]))
			tNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

			if stockMap[msgMap["name"]].Kind() == "option" {
				stockMap[msgMap["name"]].(*option).Expiration = tExp.Sub(tNow).Hours() / 24
				stockMap[msgMap["name"]].(*option).ExpirationDate = tExp
			}
		}
		stockMap[msgMap["name"]].Perform()

	}
}

func serveWeb() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, stockMap)
	})
	e.Logger.Fatal(e.Start(":8099"))
}

func saveMsgToFile() {
	if len(msgList) >= 10 {
		f, _ := os.OpenFile("api-core/output/out.txt", os.O_APPEND|os.O_WRONLY, 0644)

		allTxt := ""
		for _, txt := range msgList {
			allTxt += txt
		}

		f.WriteString(allTxt)

		f.Close()
		msgList = make([]string, 0)
	}

}

func transformMsgIntoMap(msgSpl []string) (t map[string]string) {
	t = make(map[string]string)
	t["name"] = msgSpl[1]
	for i, msg := range msgSpl {
		if strings.Contains(msg, "!") {
			break
		}
		if i%2 != 0 {
			t[msg] = msgSpl[i+1]
		}
	}
	return
}

func convertFile() (stocks []string, options []string) {

	{
		b, err := ioutil.ReadFile("api-core/assets/options.json") // just pass the file name
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, &options)
		if err != nil {
			panic(err)
		}
	}
	{
		b, err := ioutil.ReadFile("api-core/assets/stocks.json") // just pass the file name
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, &stocks)
		if err != nil {
			panic(err)
		}
	}

	return
}

func sendAllMessages(conn *net.Conn, stocks []string, options []string) {
	for _, stock := range stocks {
		stock = strings.ToLower(stock)
		fmt.Println(fmt.Sprintf("sqt %s\n", stock))
		fmt.Fprintf(*conn, fmt.Sprintf("sqt %s\n", stock))
		time.Sleep(200 * time.Millisecond)
	}

	time.Sleep(1000 * time.Millisecond)

	for _, option := range options {
		option = strings.ToLower(option)
		fmt.Println(fmt.Sprintf("sqt %s\n", option))
		fmt.Fprintf(*conn, fmt.Sprintf("sqt %s\n", option))
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("FINISHED")
}
