package client

import (
	"github.com/goex-top/market_center_client"
	goex "github.com/nntaoli-project/GoEx"
)

type Client struct {
	ExchangeName string
	CurrencyPair string
	api          goex.API
	c            *market_center_client.Client
}

func NewClient(exchangeName, currencyPair string, api goex.API, c *market_center_client.Client) *Client {
	if api == nil && c == nil {
		return nil
	}
	return &Client{
		ExchangeName: exchangeName,
		CurrencyPair: currencyPair,
		api:          api,
		c:            c,
	}
}
func (c *Client) GetTicker() *goex.Ticker {
	var tick *goex.Ticker
	var err error
	if c.api != nil {
		tick, err = c.api.GetTicker(goex.NewCurrencyPair2(c.CurrencyPair))
	} else if c.c != nil {
		tick, err = c.c.GetTicker(c.ExchangeName, c.CurrencyPair)
	}
	if err != nil {
		return nil
	}
	return tick
}

func (c *Client) GetDepth() *goex.Depth {
	var depth *goex.Depth
	var err error
	if c.api != nil {
		depth, err = c.api.GetDepth(20, goex.NewCurrencyPair2(c.CurrencyPair))
	} else if c.c != nil {
		depth, err = c.c.GetDepth(c.ExchangeName, c.CurrencyPair)
	}
	if err != nil {
		return nil
	}
	return depth
}

func (c *Client) GetKline() *goex.Kline {
	panic("not support yet")
}
