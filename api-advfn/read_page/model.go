package read_page

import "math"

type stock struct {
	Price float64
}

type option struct {
	Kind       string // Call or Put
	Style      string //American or European
	Expiration float64
	Url        string
	Price      float64
	AvgPrice   float64
	Strike     float64
	Stock      stock
}

func (opt option) MinProfitPerMonth() float64 {
	d := opt.Expiration / 30.0
	p := opt.Price / opt.Stock.Price / d * 100
	return math.Floor(p*100) / 100
}

func (opt option) MaxProfitPerMonth() float64 {
	d := opt.Expiration / 30.0

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
