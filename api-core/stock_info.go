package main

import (
	"strconv"
	"time"

	"github.com/jinzhu/copier"
)

type stockInfo interface {
	Kind() string
	SetPrice(string, priceKind)
	SetName(string)
	Perform()
}

type priceKind string

const (
	_LAST_PRICE      priceKind = "last_price"
	_BEST_BUY_OFFER  priceKind = "best_buy_offer"
	_BEST_SELL_OFFER priceKind = "best_sell_offer"
)

type option struct {
	Price          float64   `json:"price"`
	BestSellOffer  float64   `json:"best_sell_offer"`
	BestBuyOffer   float64   `json:"best_buy_offer"`
	Name           string    `json:"name"`
	Stock          *stock    `json:"stock"`
	Strike         float64   `json:"strike"`
	Expiration     float64   `json:"expiration"`
	Volume         float64   `json:"volume"`
	ExpirationDate time.Time `json:"expiration_date"`
	Updated        time.Time `json:"updated"`
}

func (opt option) Kind() string {
	return "option"
}

func (opt *option) SetPrice(p string, kind priceKind) {
	price, _ := strconv.ParseFloat(p, 32)
	switch kind {
	case _LAST_PRICE:
		opt.Price = price
	case _BEST_BUY_OFFER:
		opt.BestBuyOffer = price
	case _BEST_SELL_OFFER:
		opt.BestSellOffer = price
	}
}

func (opt *option) IsMarketOpen() bool {
	if opt.Updated.Hour() >= 10 && opt.Updated.Hour() <= 18 {
		return true
	}
	return false
}

func (opt *option) SetName(p string) {
	opt.Name = p
}

func (opt *option) Perform() {
	//Perform VDX
	if !opt.IsMarketOpen() {
		return
	}
	if opt.Stock == nil || opt.Stock.Price == 0 || opt.Price == 0 || opt.BestBuyOffer == 0 || opt.Stock.BestSellOffer == 0 {
		return
	}

	if _, ok := snapshotMap[opt.Name]; !ok {
		snapshotMap[opt.Name] = make(map[filter]*snapshot)
	}

	// PERFORM MAX PROFIT
	{
		if _, ok := snapshotMap[opt.Name][_MAX_PROFIT]; !ok {
			snapshotMap[opt.Name][_MAX_PROFIT] = &snapshot{}
		}
		maxProfit := (opt.Strike + opt.Price - opt.Stock.Price) / opt.Stock.Price
		if snapshotMap[opt.Name][_MAX_PROFIT].Value < maxProfit {
			snapshotMap[opt.Name][_MAX_PROFIT].Value = maxProfit
			newOpt := &option{}
			newStock := &stock{}
			copier.Copy(newOpt, stockMap[opt.Name])
			copier.Copy(newStock, stockMap[opt.Name].(*option).Stock)
			newOpt.Stock = newStock
			snapshotMap[opt.Name][_MAX_PROFIT].Option = *newOpt
		}
	}

	// PERFORM MIN PROFIT
	{
		if _, ok := snapshotMap[opt.Name][_MIN_PROFIT]; !ok {
			snapshotMap[opt.Name][_MIN_PROFIT] = &snapshot{}
		}
		minProfit := opt.Price / opt.Stock.Price
		if snapshotMap[opt.Name][_MIN_PROFIT].Value < minProfit {
			snapshotMap[opt.Name][_MIN_PROFIT].Value = minProfit
			newOpt := &option{}
			newStock := &stock{}
			copier.Copy(newOpt, stockMap[opt.Name])
			copier.Copy(newStock, stockMap[opt.Name].(*option).Stock)
			newOpt.Stock = newStock
			snapshotMap[opt.Name][_MIN_PROFIT].Option = *newOpt
		}
	}

	// PERFORM ON MARKET MIN PROFIT
	{
		f := _ONMARKET__MIN_PROFIT
		if _, ok := snapshotMap[opt.Name][f]; !ok {
			snapshotMap[opt.Name][f] = &snapshot{Filter: f}
		}
		minProfit := opt.BestBuyOffer / opt.Stock.BestSellOffer
		if snapshotMap[opt.Name][f].Value < minProfit {
			snapshotMap[opt.Name][f].Value = minProfit
			newOpt := &option{}
			newStock := &stock{}
			copier.Copy(newOpt, stockMap[opt.Name])
			copier.Copy(newStock, stockMap[opt.Name].(*option).Stock)
			newOpt.Stock = newStock
			snapshotMap[opt.Name][f].Option = *newOpt
		}
		snapshotMap[opt.Name][f].trigger()
	}
}

type stock struct {
	Price         float64 `json:"price"`
	BestSellOffer float64 `json:"best_sell_offer"`
	BestBuyOffer  float64 `json:"best_buy_offer"`
	Name          string  `json:"name"`
	Volume        float64 `json:"volume"`
}

func (st stock) Kind() string {
	return "stock"
}

func (st *stock) SetPrice(p string, kind priceKind) {
	price, _ := strconv.ParseFloat(p, 32)
	switch kind {
	case _LAST_PRICE:
		st.Price = price
	case _BEST_BUY_OFFER:
		st.BestBuyOffer = price
	case _BEST_SELL_OFFER:
		st.BestSellOffer = price
	}
}

func (st stock) getChildren() *option {
	for _, i := range stockMap {
		if i.Kind() == (option{}).Kind() {
			if i.(*option).Stock == nil {
				return nil
			}
			if i.(*option).Stock.Name == st.Name {
				return i.(*option)
			}
		}
	}
	return nil
}

func (st *stock) SetMarketStatus(p string) {
}

func (st *stock) SetName(p string) {
	st.Name = p
}

func (st *stock) Perform() {
}
