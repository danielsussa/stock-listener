package read_page

import (
	"strconv"
	"strings"

	"github.com/go-resty/resty"
)

func ReadDetailPage(opt *option) {
	res, err := resty.R().Get("https:" + opt.Url)
	if err != nil {
		panic(err)
	}
	getOptionDetail(string(res.Body()), opt)

}

func getOptionDetail(s string, opt *option) {
	{
		spl := strings.Split(s, "quoteElementPiece")
		for _, v := range spl {
			if len(v) == 0 {
				continue
			}

			if string(v[0:2]) == "10" {
				pStr := v[strings.Index(v, ">")+1 : strings.Index(v, "<")]
				pStr = strings.Replace(pStr, ",", ".", -1)
				price, _ := strconv.ParseFloat(pStr, 32)
				opt.Price = price
			}
		}
	}

	// Get expiration date
	{
		spl := strings.Split(s, "dias</td>")
		txt := spl[0][len(spl[0])-5:]
		txt = strings.TrimSpace(strings.Split(txt, ">")[1])
		expDate, _ := strconv.Atoi(txt)
		opt.Expiration = expDate
	}
}
