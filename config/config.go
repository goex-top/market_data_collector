package config

import "github.com/goex-top/market_center"

type Subscribe struct {
	ExchangeName string                 `json:"exchange_name"`
	CurrencyPair string                 `json:"currency_pair"`
	Period       int64                  `json:"period"`
	Flag         market_center.DataFlag `json:"flag"`
	Direct       bool                   `json:"direct"`
}

type Storage struct {
	Csv bool `json:"csv"`
	// TBD
}
type Config struct {
	Subs             []Subscribe `json:"subs"`
	Store            Storage     `json:"store"`
	MarketCenterPath string      `json:"market_center_path"`
}
