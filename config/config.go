package config

import "github.com/goex-top/market_center"

type Subscribe struct {
	ExchangeName string                 `json:"exchange_name" yaml:"exchange_name" default:"binance.com"`
	CurrencyPair string                 `json:"currency_pair" yaml:"currency_pair" default:"BTC_USDT"`
	ContractType string                 `json:"contract_type,omitempty" yaml:"contract_type" default:""`
	Period       int64                  `json:"period" yaml:"period" default:"100"`
	Flag         market_center.DataFlag `json:"flag" yaml:"flag" default:"1"`
}

type Storage struct {
	Csv bool `json:"csv" yaml:"csv" default:"true"`
	// TBD
}
type Config struct {
	Subs             []Subscribe `json:"subs" yaml:"subs" default:"subs"`
	Store            Storage     `json:"store" yaml:"store" default:""`
	WithMarketCenter bool        `json:"with_market_center" yaml:"with_market_center" `
	MarketCenterPath string      `json:"market_center_path" yaml:"market_center_path" default:"/tmp/goex.market.center"`
}
