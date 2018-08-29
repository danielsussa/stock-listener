package main

import (
	"fmt"

	"github.com/danielsussa/stock-listener/api-advfn/read_page"
)

func main() {

	options := read_page.ReadMainPage("https://br.advfn.com/bolsa-de-valores/bovespa/petrobras-PETR4/opcoes")

	// Step 4 Iterate over options

	for k, opt := range options {
		//time.Sleep(200 * time.Millisecond)
		read_page.ReadDetailPage(&opt)
		if opt.Price == 0 {
			//fmt.Println(fmt.Sprintf("%s -> cannot setup option %s", time.Now(), k))
		} else {
			//fmt.Println(fmt.Sprintf("%s -> get detail from %s", time.Now(), k))
			min := opt.MinProfitPerMonth()
			max := opt.MaxProfitPerMonth()
			if min < max && min > 3 && opt.Expiration < 120 {
				fmt.Println(fmt.Sprintf("%8s (%5.2f)-> (profit: %5.2f <min|max> %6.2f | spr: %5.2f) - Price: %.2f | Stk.Price: %5.2f | Exp: %3.0f => Kind: %s | Style: %s",
					k, opt.Strike, min, max, (opt.Stock.Price - opt.Price),
					opt.Price, opt.Stock.Price, opt.Expiration,
					opt.Kind, opt.Style,
				))
			}
		}
	}

	//spew.Dump(options)

}
