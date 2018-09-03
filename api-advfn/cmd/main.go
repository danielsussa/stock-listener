package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/danielsussa/stock-listener/api-advfn/read_page"
	"github.com/gocarina/gocsv"
)

func getAllURLS() []string {
	return []string{
		"https://br.advfn.com/bolsa-de-valores/bovespa/petrobras-PETR4/opcoes",
		"https://br.advfn.com/bolsa-de-valores/bovespa/kroton-KROT3/opcoes",
	}
}

func main() {
	optsSlc := read_page.NewOptions()

	for _, url := range getAllURLS() {
		options := read_page.ReadMainPage(url)

		// Step 4 Iterate over options

		for k, opt := range options {
			//time.Sleep(200 * time.Millisecond)
			read_page.ReadDetailPage(&opt)
			opt.Name = k
			if opt.Price != 0 {
				//fmt.Println(fmt.Sprintf("%s -> cannot setup option %s", time.Now(), k))
				min := opt.MinProfit(true)
				max := opt.MaxProfit(true)

				if min > 4 && max > 0 && opt.Expiration < 120 && opt.Kind == "C" && opt.Price > 1 {
					optsSlc = append(optsSlc, opt)
					fmt.Println(fmt.Sprintf("%8s (%5.2f)-> (profit: %5.2f <min|max> %6.2f | spr: %5.2f) - Price: %.2f | Stk.Price: %5.2f | Exp: %3.0f => Kind: %s | Style: %s",
						k, opt.Strike, min, max, (opt.Stock.Price - opt.Price),
						opt.Price, opt.Stock.Price, opt.Expiration,
						opt.Kind, opt.Style,
					))
				}
			}
		}
	}

	gocsv.TagSeparator = ";"
	csvContent, err := gocsv.MarshalString(&optsSlc) // Get all clients as CSV string
	//err = gocsv.MarshalFile(&clients, clientsFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
	csvContent = strings.Replace(csvContent, ".", ",", -1)
	ioutil.WriteFile("output.csv", []byte(csvContent), 0644)

	//spew.Dump(options)

}
