package read_page

import (
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
	QtdNegs    float64
	VolNegs    float64
	Stock      stock
}

func NewOptions() []option {
	return []option{}
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
