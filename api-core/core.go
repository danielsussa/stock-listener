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

type stockInfo struct {
	Price          float64
	Name           string
	Parent         string
	MinProfit      float64
	MaxProfit      float64
	Strike         float64
	Expiration     float64
	Volume         float64
	ExpirationDate time.Time
	Vdx            float64
}

var stockMap map[string]*stockInfo
var msgList []string

func main() {
	stocks, options := convertFile()

	//Add to stockMap
	stockMap = make(map[string]*stockInfo)
	for _, stock := range stocks {
		stockMap[strings.ToUpper(stock)] = &stockInfo{}
	}
	for _, option := range options {
		stockMap[strings.ToUpper(option)] = &stockInfo{}
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
	go saveMsgToFile()
	go serveWeb()

	for {
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')

		msgSpl := strings.Split(message, ":")
		if len(msgSpl) < 3 || msgSpl[0] == "E" {
			continue
		}
		if msgSpl[0] != "T" {
			continue
		}

		fmt.Print("Message from server: " + message)
		msgList = append(msgList, message)

		msgMap := transformMsgIntoMap(msgSpl)

		if _, ok := msgMap["45"]; ok {
			//If is a option
			if msgMap["45"] == "2" {
				if _, ok := msgMap["81"]; ok {
					stockMap[msgMap["name"]].Parent = msgMap["81"]
				}
			}
		}

		//Setup price
		if _, ok := msgMap["2"]; ok {
			price, _ := strconv.ParseFloat(msgMap["2"], 32)
			stockMap[msgMap["name"]].Price = price
		}

		//Setup strike
		if _, ok := msgMap["121"]; ok {
			strike, _ := strconv.ParseFloat(msgMap["121"], 32)
			stockMap[msgMap["name"]].Strike = strike
		}

		//Setup volume
		if _, ok := msgMap["9"]; ok {
			vol, _ := strconv.ParseFloat(msgMap["9"], 32)
			stockMap[msgMap["name"]].Volume = vol
		}

		//Setup Expiration date
		if _, ok := msgMap["125"]; ok {
			tExp, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", msgMap["125"][0:4], msgMap["125"][4:6], msgMap["125"][6:8]))
			tNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

			stockMap[msgMap["name"]].Expiration = tExp.Sub(tNow).Hours() / 24
			stockMap[msgMap["name"]].ExpirationDate = tExp
		}

		// 	//expiration Date

		// 	option.Name = msgMap["name"]
		// 	option.Price = price
		// 	option.StockPrice = stock.Price
		// 	option.Strike = strike
		// 	option.Volume = volume
		// 	option.Vdx = (option.Price / stock.Price) * (120 - option.Expiration) * (option.Strike - stock.Price)
		// 	option.MinProfit = option.Price / stock.Price
		// 	option.MaxProfit = (option.Strike + option.Price - stock.Price) / stock.Price
		// 	stockMap[msgMap["name"]] = option

		// }

	}
}

func serveWeb() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, stockMap)
	})
	e.GET("/exist", func(c echo.Context) error {
		newMap := make([]string, 0)
		for _, val := range stockMap {
			f, _ := strconv.ParseFloat(c.QueryParam("min"), 32)
			if val.Vdx > f {
				newMap = append(newMap, val.Name)
			}
		}
		return c.JSON(http.StatusOK, newMap)
	})
	e.Logger.Fatal(e.Start(":8099"))
}

func saveMsgToFile() {
	for {
		if len(msgList) > 0 {
			f, _ := os.OpenFile("api-core/output/out.txt", os.O_APPEND|os.O_WRONLY, 0644)

			allTxt := ""
			for _, txt := range msgList {
				allTxt += txt
			}

			f.WriteString(allTxt)

			f.Close()
			msgList = make([]string, 0)
		}

		time.Sleep(1 * time.Minute)
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
