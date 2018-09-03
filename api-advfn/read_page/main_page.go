package read_page

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty"
)

type TrTmp struct {
	Class string   `xml:"class,attr"`
	Td    []string `xml:"td"`
	Tg    tg       `xml:"tg"`
	Timg  timg     `xml:"timg"`
}

type tg struct {
	Url     string `xml:"url,attr"`
	Content string `xml:"a"`
}

type timg struct {
	Kind string `xml:"kind,attr"`
}

type XMLTable struct {
	Class   string   `xml:"class,attr"`
	XMLName xml.Name `xml:"table"`
	Tr      []TrTmp  `xml:"tr"`
}

func ReadMainPage(url string) map[string]option {
	// Step 1 Main call to get main HTML
	res, _ := resty.R().Get(url)

	s := string(res.Body())

	// Step 2 Get STOCK
	stock := getStock(s)

	//Step 3: Get all stocks
	options := getOptions(s, stock)
	return options
}

func getStock(s string) stock {
	stk := stock{}
	spl := strings.Split(s, "quoteElementPiece")
	for _, v := range spl {
		if string(v[0:1]) == "6" {
			pStr := v[strings.Index(v, ">")+1 : strings.Index(v, "<")]
			pStr = strings.TrimSpace(pStr)
			pStr = strings.Replace(pStr, ",", ".", -1)
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
				url := strings.Split(v, "href=")[1]
				url = strings.Split(url[1:], "\"")[0]
				url = "\"" + url + "\""
				s += "<tg " + "url=" + url + " " + v[4:len(v)] + "</tg>"
			} else if strings.Contains(v, "img") {
				kind := ""
				if strings.Contains(v, "E.gif") {
					kind = "kind=\"E\""
				}
				if strings.Contains(v, "A.gif") {
					kind = "kind=\"A\""
				}
				s += "<timg " + kind + " " + v[4:len(v)] + "</timg>"
			} else {
				s += v[0:len(v)] + "</td>"
			}
		}
	}

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
			Url:    v.Tg.Url,
			Kind:   v.Td[2],
			Strike: strike,
			Stock:  stk,
		}
	}

	return optMap
}
