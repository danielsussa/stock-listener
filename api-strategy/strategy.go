package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

type stock struct {
	Price          float64
	StockPrice     float64
	Name           string
	Strike         float64
	Expiration     float64
	ExpirationDate time.Time
	Vdx            float64
}

func getAllSeries() []string {
	return []string{"j", "k", "l"}
}

func getInputStock() map[string]string {
	return map[string]string{
		//"petr4": "petr",
		//"vale3": "vale",
		//"bbas3": "bbas",
		"bbdc3": "bbdc",
	}
}

var stockMap map[string]*stock

func serveWeb() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		newMap := make(map[string]*stock, 0)
		for k, val := range stockMap {
			f, _ := strconv.ParseFloat(c.QueryParam("min"), 32)
			if val.Vdx > f {
				newMap[k] = val
			}
		}
		return c.JSON(http.StatusOK, newMap)
	})
	e.GET("/exist", func(c echo.Context) error {
		newMap := make([]stock, 0)
		for _, val := range stockMap {
			f, _ := strconv.ParseFloat(c.QueryParam("min"), 32)
			if val.Vdx > f {
				newMap = append(newMap, *val)
			}
		}
		return c.JSON(http.StatusOK, newMap)
	})
	e.Logger.Fatal(e.Start(":8099"))
}

func main() {
	fmt.Println("Start CORE")
	go serveWeb()

	// connect to this socket
	conn, err := net.Dial("tcp", "datafeeddl1.cedrofinances.com.br:81")

	if err != nil {
		panic(err)
	}

	// send to socket
	fmt.Fprintf(conn, "\n")
	fmt.Fprintf(conn, "kanczuk\n")
	fmt.Fprintf(conn, "102030\n")

	go sendAllMessages(&conn)

	stockMap = make(map[string]*stock, 0)

	for {
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')

		msgSpl := strings.Split(message, ":")
		if len(msgSpl) < 10 {
			continue
		}

		fmt.Print("Message from server: " + message)

		msgMap := transformMsgIntoMap(msgSpl)

		//if is a stock
		if msgMap["45"] == "1" {
			stock := &stock{}
			stockMap[msgMap["name"]] = stock
			f, _ := strconv.ParseFloat(msgMap["2"], 32)
			stock.Price = f
		}
		// If is a option
		if msgMap["45"] == "2" {
			if _, ok := stockMap[msgMap["81"]]; !ok {
				panic("Cannot continue")
			}
			option := &stock{}
			stock := stockMap[msgMap["81"]]

			//get stock

			//price
			price, _ := strconv.ParseFloat(msgMap["2"], 32)

			//price
			strike, _ := strconv.ParseFloat(msgMap["121"], 32)

			//expiration Date
			tExp, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", msgMap["125"][0:4], msgMap["125"][4:6], msgMap["125"][6:8]))
			tNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

			option.Expiration = tExp.Sub(tNow).Hours() / 24
			option.ExpirationDate = tExp
			option.Name = msgMap["name"]
			option.Price = price
			option.StockPrice = stock.Price
			option.Strike = strike
			option.Vdx = (option.Price / stock.Price) * (120 - option.Expiration) * (option.Strike - stock.Price)
			stockMap[msgMap["name"]] = option

		}

	}

}

func transformMsgIntoMap(msgSpl []string) (t map[string]string) {
	t = make(map[string]string)
	t["name"] = msgSpl[1]
	for i, msg := range msgSpl {
		if i%2 != 0 {
			if len(msgSpl) > i+1 {
				t[msg] = msgSpl[i+1]
			}
		}
	}
	return
}

func sendAllMessages(conn *net.Conn) {
	inputStock := getInputStock()
	allSeries := getAllSeries()
	for stockName, optPrefix := range inputStock {
		fmt.Fprintf(*conn, fmt.Sprintf("sqt %s\n", stockName))

		for _, serie := range allSeries {
			fmt.Println(stockName, serie)
			for i := 1; i <= 300; i++ {
				fmt.Fprintf(*conn, fmt.Sprintf("sqt %s%s%d\n", optPrefix, serie, i))
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
	fmt.Println("FINISHED")
}
