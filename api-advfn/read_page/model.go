package read_page

import (
	"fmt"
	"math"
)

type stock struct {
	Price float64 `csv:"stock_price"`
}

type option struct {
	Name       string `csv:"name"`
	Kind       string // Call or Put
	Style      string //American or European
	Expiration float64
	Url        string `csv:"-"`
	Price      float64
	AvgPrice   float64
	Strike     float64
	VolNegs    float64
	Stock      stock
}

func NewOptions() []option {
	return []option{}
}

func (opt option) ReadAndPrint() {
	ReadDetailPage(&opt)
	if opt.Price != 0 {
		//fmt.Println(fmt.Sprintf("%s -> cannot setup option %s", time.Now(), k))
		profit := opt.Profit(true)
		protection := opt.Protection()

		if profit > 2 && opt.Expiration < 100 && opt.Kind == "C" && opt.VolNegs > 0 && protection > 5 {
			fmt.Println(fmt.Sprintf("VC | %8s (%5.2f)-> ( %5.2f <prof(%s - %6.2f)marg> %6.2f |  spr: %5.2f | Price: %5.2f | Stk.Price: %5.2f | Vol: %7.0f | Exp: %3.0f => Kind: %s | Style: %s",
				opt.Name, opt.Strike, profit, opt.Modality(), (profit + protection), protection,
				(opt.Stock.Price - opt.Price), //spread
				opt.Price, opt.Stock.Price, opt.VolNegs, opt.Expiration,
				opt.Kind, opt.Style,
			))
		}

		if opt.Kind == "P" && opt.Price < 0.2 && opt.VolNegs > 0 {
			fmt.Println(fmt.Sprintf("PU | %8s (%5.2f)-> Price: %5.2f | Stk.Price: %5.2f | Vol: %7.0f | Exp: %3.0f => Prof per Price: %5.4f",
				opt.Name, opt.Strike, opt.Price, opt.Stock.Price, opt.VolNegs, opt.Expiration,
				(opt.Price/(opt.Stock.Price-opt.Strike))*100,
			))
		}
		// if opt.Price < 0.3 && opt.Expiration < 100 && opt.Kind == "P" && opt.VolNegs > 0 {
		// 	fmt.Println(fmt.Sprintf("PR: %8s (%5.2f)->  Price: %5.2f | Stk.Price: %5.2f",
		// 		opt.Name, opt.Strike, opt.Price, opt.Stock.Price,
		// 	))
		// }
	}
}

func (opt option) Modality() string {
	//ex: strike 19 and stock 20
	if opt.Strike < opt.Stock.Price {
		return "I"
	}
	//ex: strike 19.19 and stock 20
	if opt.Strike-(opt.Strike*0.05) > opt.Stock.Price && opt.Strike+(opt.Strike*0.05) < opt.Stock.Price {
		return "A"
	}
	return "O"
}

func (opt option) Protection() float64 {
	if opt.Stock.Price-opt.Price < opt.Strike {
		return opt.Price / opt.Strike * 100
	}
	return (opt.Stock.Price - opt.Strike) / opt.Stock.Price * 100
}

func (opt option) Profit(perMonth bool) float64 {
	d := opt.Expiration / 30.0
	if !perMonth {
		d = 1
	}

	if opt.Kind == "C" {
		p := (opt.Price + opt.Strike - opt.Stock.Price) / opt.Stock.Price / d * 100
		return math.Floor(p*100) / 100
	}
	if opt.Kind == "P" {
		p := (opt.Price - opt.Strike + opt.Stock.Price) / opt.Stock.Price / d * 100
		return math.Floor(p*100) / 100
	}
	return 0
}
