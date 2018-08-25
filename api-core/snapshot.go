package main

import (
	"fmt"
	"time"
)

type filter string

const (
	_MAX_PROFIT           filter = "max_profit"
	_MIN_PROFIT           filter = "min_profit"
	_ONMARKET__MIN_PROFIT filter = "onmarket_min_profit"
)

type snapshot struct {
	Value           float64
	Filter          filter
	LastTriggerTime time.Time
	Option          option
}

var snapshotMap map[string]map[filter]*snapshot

func (s *snapshot) trigger() {
	if s.Option.Updated.Sub(s.LastTriggerTime).Minutes() < 15 {
		return
	}
	switch s.Filter {
	case _ONMARKET__MIN_PROFIT:
		if s.Value > 0.05 {
			s.LastTriggerTime = s.Option.Updated
			fmt.Println("Trigger")
			//sendMessage(fmt.Sprintf("Verificar ação: %s , opção: %s a preço %f - horario: %s", s.Option.Stock.Name, s.Option.Name, s.Option.BestBuyOffer, s.Option.Updated.String()))
		}
	}
}
