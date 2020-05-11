package client

import (
	"github.com/goex-top/market_center"
	mcc "github.com/goex-top/market_center/client"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
	"os"
	"sort"
)

type Client struct {
	ExchangeName string
	CurrencyPair string
	ContractType string
	spotApi      goex.API
	futureApi    goex.FutureRestAPI
	c            *mcc.Client
	isDirect     bool
	isSpot       bool
}

func NewClient(exchangeName, currencyPair, contractType string, c *mcc.Client) *Client {
	proxy := os.Getenv("HTTP_PROXY")
	var spotApi goex.API
	var futureApi goex.FutureRestAPI
	var direct = false
	var isSpot = !market_center.IsFutureExchange(exchangeName)
	if c == nil {
		if isSpot {
			spotApi = builder.NewAPIBuilder().HttpProxy(proxy).Build(market_center.SupportAdapter[exchangeName])
		} else {
			futureApi = builder.NewAPIBuilder().HttpProxy(proxy).BuildFuture(market_center.SupportAdapter[exchangeName])
		}
		direct = true
	}
	return &Client{
		ExchangeName: exchangeName,
		CurrencyPair: currencyPair,
		spotApi:      spotApi,
		futureApi:    futureApi,
		ContractType: contractType,
		c:            c,
		isDirect:     direct,
		isSpot:       isSpot,
	}
}

func (c *Client) Close() {
	if c.c != nil {
		c.c.Close()
	}
}

func (c *Client) Name() string {
	return c.ExchangeName
}

func (c *Client) GetTicker() *goex.Ticker {
	var tick *goex.Ticker
	var err error
	if c.isDirect {
		if c.isSpot {
			tick, err = c.spotApi.GetTicker(goex.NewCurrencyPair2(c.CurrencyPair))
		} else {
			tick, err = c.futureApi.GetFutureTicker(goex.NewCurrencyPair2(c.CurrencyPair), c.ContractType)
		}
	} else {
		if c.isSpot {
			tick, err = c.c.GetSpotTicker(c.ExchangeName, c.CurrencyPair)
		} else {
			tick, err = c.c.GetFutureTicker(c.ExchangeName, c.ContractType, c.CurrencyPair)
		}
	}
	if err != nil {
		return nil
	}
	return tick
}

func (c *Client) GetDepth() *goex.Depth {
	var depth *goex.Depth
	var err error
	if c.isDirect {
		if c.isSpot {
			depth, err = c.spotApi.GetDepth(20, goex.NewCurrencyPair2(c.CurrencyPair))
		} else {
			depth, err = c.futureApi.GetFutureDepth(goex.NewCurrencyPair2(c.CurrencyPair), c.ContractType, 20)
		}
	} else {
		if c.isSpot {
			depth, err = c.c.GetSpotDepth(c.ExchangeName, c.CurrencyPair)
		} else {
			depth, err = c.c.GetFutureDepth(c.ExchangeName, c.ContractType, c.CurrencyPair)
		}
	}
	if err != nil || depth.AskList.Len() == 0 || depth.BidList.Len() == 0 {
		return nil
	}

	if depth.AskList[0].Price > depth.AskList[1].Price {
		sort.Slice(depth.AskList, func(i, j int) bool {
			return depth.AskList[i].Price < depth.AskList[j].Price
		})
	}
	if depth.BidList[0].Price < depth.BidList[1].Price {
		sort.Slice(depth.BidList, func(i, j int) bool {
			return depth.BidList[i].Price > depth.BidList[j].Price
		})
	}

	if depth.AskList.Len() > 20 {
		depth.AskList = depth.AskList[:20]
	}
	if depth.BidList.Len() > 20 {
		depth.BidList = depth.BidList[:20]
	}

	return depth
}

func (c *Client) GetKline() *goex.Kline {
	panic("not support yet")
}
