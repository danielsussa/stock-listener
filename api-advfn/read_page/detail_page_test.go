package read_page

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestTimeConsuming(t *testing.T) {
	b, err := ioutil.ReadFile("detail.html") // just pass the file name
	if err != nil {
		panic(err)
	}

	opt := option{}
	getOptionDetail(string(b), &opt)

	if opt.Name != "PETRJ49" {
		t.Log("Error on name: ", opt.Name)
		t.Fail()
	}

	if fmt.Sprintf("%.2f", opt.Strike) != "19.42" {
		t.Log("Error on Strike: ", opt.Strike)
		t.Fail()
	}

	if fmt.Sprintf("%.2f", opt.Price) != "1.29" {
		t.Log("Error on Price: ", opt.Price)
		t.Fail()
	}

	if opt.Expiration != 48 {
		t.Log("Error on expiration")
		t.Fail()
	}

	if opt.Kind != "C" {
		t.Log("Error on Kind")
		t.Fail()
	}

	if opt.Style != "A" {
		t.Log("Error on Kind")
		t.Fail()
	}

	opt.Stock.Price = 16.55
	if opt.MinProfitPerMonth() != 4.87 {
		t.Log("Error processing profit", opt.MinProfitPerMonth())
		t.Fail()
	}

	opt.Strike = 17.12
	opt.Kind = "C"
	if opt.MaxProfitPerMonth() != 7.02 {
		t.Log("Error processing max profit", opt.MaxProfitPerMonth())
		t.Fail()
	}

}
