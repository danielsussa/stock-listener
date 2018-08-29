package read_page

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestStockInfo(t *testing.T) {
	b, err := ioutil.ReadFile("home.html") // just pass the file name
	if err != nil {
		panic(err)
	}

	stock := getStock(string(b))

	if fmt.Sprintf("%.2f", stock.Price) != "18.30" {
		spew.Dump(stock)
		t.Fail()
	}

}
