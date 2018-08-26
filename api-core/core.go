package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

var mg mailgun.Mailgun

type stockInfoMap map[string]stockInfo

var stockMap stockInfoMap
var msgList []string

func main() {
	stocks, options := convertFile()

	snapshotMap = make(map[string]map[filter]*snapshot)

	mg = mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_PRIVATE"), os.Getenv("MAILGUN_PUBLIC"))

	//Add to stockMap
	stockMap = make(map[string]stockInfo)
	for _, st := range stocks {
		stockMap[strings.ToUpper(st)] = &stock{}
	}
	for _, opt := range options {
		stockMap[strings.ToUpper(opt)] = &option{}
	}

	var c connector

	c = &fileConnector{}
	if os.Getenv("CONNECTOR_KIND") == "tcp" {
		c = &tcpConnector{}
	}

	c.connect()

	// send to socket
	c.sendMessage("")
	c.sendMessage("kanczuk")
	c.sendMessage("102030")

	go sendAllMessages(c, stocks, options)

	for {
		// listen for reply
		message := c.getMessage()

		msgSpl := strings.Split(message, ":")
		//saveMsgToFile()

		if len(msgSpl) < 3 || msgSpl[0] == "E" {
			continue
		}
		if msgSpl[0] != "T" {
			continue
		}

		//fmt.Print("Message from server: " + message)
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
			stockMap[msgMap["name"]].SetPrice(msgMap["2"], _LAST_PRICE)
		}
		if _, ok := msgMap["3"]; ok {
			stockMap[msgMap["name"]].SetPrice(msgMap["3"], _BEST_BUY_OFFER)
		}
		if _, ok := msgMap["4"]; ok {
			stockMap[msgMap["name"]].SetPrice(msgMap["4"], _BEST_SELL_OFFER)
		}
		//Setup Market Status
		// if _, ok := msgMap["84"]; ok {
		// 	stockMap[msgMap["name"]].SetMarketStatus(msgMap["84"])
		// }

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
		//Setup Last TIME changes
		t, err := time.Parse("2006-01-02T15:04:05", fmt.Sprintf("2006-01-02T%s:%s:%s", msgMap["time"][0:2], msgMap["time"][2:4], msgMap["time"][4:6]))
		if err != nil {
			panic(err)
		}
		if stockMap[msgMap["name"]].Kind() == "option" {
			l := stockMap[msgMap["name"]].(*option).Updated
			stockMap[msgMap["name"]].(*option).Updated = time.Date(l.Year(), l.Month(), l.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
		}
		//Setup Last DATE changes
		if _, ok := msgMap["1"]; ok {
			t, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", msgMap["1"][0:4], msgMap["1"][4:6], msgMap["1"][6:8]))
			if err != nil {
				panic(err)
			}
			if stockMap[msgMap["name"]].Kind() == "option" {
				l := stockMap[msgMap["name"]].(*option).Updated
				stockMap[msgMap["name"]].(*option).Updated = time.Date(t.Year(), t.Month(), t.Day(), l.Hour(), l.Minute(), l.Second(), 0, time.UTC)
			}
		}
		stockMap[msgMap["name"]].Perform()

		fmt.Println("\033[2J")
		fmt.Printf("\n\033[0;0H")
		var keys []string
		for k := range stockMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if stockMap[k].Kind() == (option{}).Kind() {
				opt := stockMap[k].(*option)
				if opt.Stock == nil {
					break
				}
				if len(k) < 8 {
					k += "_"
				}
				minProfit := opt.BestBuyOffer / opt.Stock.BestSellOffer * 100
				if opt.BestBuyOffer > 0 && opt.Price > 0 {
					h, m, s := opt.Updated.Clock()
					fmt.Printf("[%s] -> B: %.2f | P: %.2f | A: %.2f | STR: %.2f | [%s] | P: %.2f | PROFIT: %.2f - %d:%d:%d \n",
						k,
						opt.BestBuyOffer,
						opt.Price,
						opt.BestSellOffer,
						opt.Strike,
						opt.Stock.Name,
						opt.Stock.Price,
						minProfit,
						h, m, s)
				}
			}
		}

	}
}

func saveMsgToFile() {
	if len(msgList) >= 30 {
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
	t["time"] = msgSpl[2]
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

func sendAllMessages(c connector, stocks []string, options []string) {
	for _, stock := range stocks {
		stock = strings.ToLower(stock)
		c.sendMessage(fmt.Sprintf("sqt %s\n", stock))
		time.Sleep(200 * time.Millisecond)
	}

	time.Sleep(1000 * time.Millisecond)

	for _, option := range options {
		option = strings.ToLower(option)
		//fmt.Println(fmt.Sprintf("sqt %s\n", option))
		c.sendMessage(fmt.Sprintf("sqt %s\n", option))
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("FINISHED")
}

func sendMessage(body string) {
	message := mg.NewMessage("danielsussa@gmail.com", "Stock", body, "danielsussa@gmail.com")
	resp, id, err := mg.Send(message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
