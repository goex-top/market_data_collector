# Client for market

Get market data directly
```go
		c := client.NewClient(v.ExchangeName, v.CurrencyPair, "", nil)
		c.GetTicker()
```

Get market data with market center
```go

    import (
        mcc "github.com/goex-top/market_center_client"
    )

    mccc := mcc.NewClient()
    isSpot := market_center.IsFutureExchange(v.ExchangeName)
    if v.Flag&market_center.DataFlag_Depth != 0 {
        if isSpot {
            mccc.SubscribeSpotDepth(v.ExchangeName, v.CurrencyPair, v.Period)
        } else {
            mccc.SubscribeFutureDepth(v.ExchangeName, v.ContractType, v.CurrencyPair, v.Period)
        }
    }
    if v.Flag&market_center.DataFlag_Ticker != 0 {
        if isSpot {
            mccc.SubscribeSpotTicker(v.ExchangeName, v.CurrencyPair, v.Period)
        } else {
            mccc.SubscribeFutureTicker(v.ExchangeName, v.ContractType, v.CurrencyPair, v.Period)
        }
    }
    c := client.NewClient(v.ExchangeName, v.CurrencyPair, "", mccc)
    c.GetTicker()
```