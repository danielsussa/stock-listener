package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

type Tr struct {
	Class string   `xml:"class,attr"`
	Td    []string `xml:"td"`
	Tg    tg       `xml:"tg"`
	Timg  timg     `xml:"timg"`
}

type tg struct {
	Content string `xml:"a"`
}

type timg struct {
	Kind string `xml:"kind,attr"`
}

type XMLTable struct {
	Class   string   `xml:"class,attr"`
	XMLName xml.Name `xml:"table"`
	Tr      []Tr     `xml:"tr"`
}

type option struct {
	Kind   string // Call or Put
	Style  string //American or European
	IPrice float64
	Url    string
	EPrice float64
	Price  float64
	Strike float64
	Stock  stock
}

func main() {

	// Step 1 Main call to get main HTML
	res, _ := resty.R().Get("https://br.advfn.com/bolsa-de-valores/bovespa/petrobras-PETR4/opcoes")

	s := string(res.Body())

	// Step 2 Get STOCK
	stock := getStock(s)

	//Step 3: Get all stocks
	options := getOptions(s, stock)

	// Step 4 Iterate over options

	for k, opt := range options {
		time.Sleep(2 * time.Second)
		fmt.Println(fmt.Sprintf("%s -> get detail from %s", time.Now(), k))
		getOptionDetail(&opt)
		if opt.Price == 0 {
			fmt.Println(fmt.Sprintf("%s -> cannot setup option %s", time.Now(), k))
		}
	}

	//spew.Dump(options)

}

type stock struct {
	Price float64
}

func getOptionDetail(opt *option) {
	res, _ := resty.R().Get(opt.Url)
	s := string(res.Body())

	spl := strings.Split(s, "quoteElementPiece")
	for _, v := range spl {
		if len(v) == 0 {
			continue
		}
		if string(v[0]) == "10" {
			pStr := v[strings.Index(v, ">")+1 : strings.Index(v, "<")]
			price, _ := strconv.ParseFloat(pStr, 32)
			opt.Price = price
		}
	}
}

func getStock(s string) stock {
	stk := stock{}
	spl := strings.Split(s, "quoteElementPiece")
	for _, v := range spl {
		if string(v[0]) == "10" {
			pStr := v[strings.Index(v, ">")+1 : strings.Index(v, "<")]
			price, _ := strconv.ParseFloat(pStr, 32)
			stk.Price = price
		}
	}
	return stk
}

func getOptions(s string, stk stock) map[string]option {
	s = strings.SplitN(s, "id_options", 2)[1]
	s = strings.SplitN(s, ">", 2)[1]
	s = strings.SplitN(s, "</table>", 2)[0]
	s = fmt.Sprintf("%s%s%s", "<table>", s, "</table>")

	{
		spl := strings.Split(s, "<im")

		s = ""
		for _, v := range spl {
			if string(v[0]) == "g" {
				k := strings.SplitN(v, ">", 2)
				s += "<im" + k[0] + "/>" + k[1]
			} else {
				s += v
			}
		}
	}

	{
		spl := strings.Split(s, "</td>")
		s = ""
		for _, v := range spl {
			if strings.Contains(v, "href") {
				s += "<tg " + v[4:len(v)] + "</tg>"
			} else if strings.Contains(v, "img") {
				kind := ""
				if strings.Contains(v, "E.gif") {
					kind = "kind=\"european\""
				}
				if strings.Contains(v, "A.gif") {
					kind = "kind=\"american\""
				}
				s += "<timg " + kind + " " + v[4:len(v)] + "</timg>"
			} else {
				s += v[0:len(v)] + "</td>"
			}
		}
	}

	//fmt.Println(s)

	table := XMLTable{}

	err := xml.Unmarshal([]byte(s), &table)
	if err != nil {
		panic(err)
	}

	optMap := make(map[string]option, 0)

	for _, v := range table.Tr {

		strike, _ := strconv.ParseFloat(strings.Replace(v.Td[1], ",", ".", -1), 32)
		optMap[v.Tg.Content] = option{
			Style:  v.Timg.Kind,
			Kind:   v.Td[2],
			Strike: strike,
			Stock:  stk,
		}
	}

	return optMap
}
