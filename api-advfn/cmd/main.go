package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/danielsussa/stock-listener/api-advfn/read_page"
)

func main() {

	options := read_page.ReadMainPage("https://br.advfn.com/bolsa-de-valores/bovespa/petrobras-PETR4/opcoes")

	// Step 4 Iterate over options

	for k, opt := range options {
		//time.Sleep(200 * time.Millisecond)
		read_page.ReadDetailPage(&opt)
		if opt.Price == 0 {
			fmt.Println(fmt.Sprintf("%s -> cannot setup option %s", time.Now(), k))
		} else {
			fmt.Println(fmt.Sprintf("%s -> get detail from %s", time.Now(), k))
			spew.Dump(opt)
		}
	}

	//spew.Dump(options)

}
