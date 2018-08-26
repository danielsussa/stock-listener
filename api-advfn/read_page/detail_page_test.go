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

	if fmt.Sprintf("%.2f", opt.Price) != "0.69" {
		t.Fail()
	}

	if opt.Expiration != 114 {
		t.Log("Error on expiration")
		t.Fail()
	}

}
