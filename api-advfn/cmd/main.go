package main

import (
	"io/ioutil"
	"strings"
	"sync"
	"time"

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

		var wg sync.WaitGroup
		wg.Add(len(options))
		for k, opt := range options {
			opt.Name = k
			go opt.ReadAndPrint()
			time.Sleep(50 * time.Millisecond)
			wg.Done()
		}
		wg.Wait()
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
