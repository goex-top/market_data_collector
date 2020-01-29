# Client for market

Get market data directly
```go
    c := client.NewClient(ExchangeName, CurrencyPair, "", nil)
    c.GetTicker()
```

Get market data with market center
```go

    import (
        mcc "github.com/goex-top/market_center_client"
    )

    mccc := mcc.NewClient()
    isSpot := market_center.IsFutureExchange(ExchangeName)
    if Flag&market_center.DataFlag_Depth != 0 {
        if isSpot {
            mccc.SubscribeSpotDepth(ExchangeName, CurrencyPair, Period)
        } else {
            mccc.SubscribeFutureDepth(ExchangeName, ContractType, CurrencyPair, Period)
        }
    }
    if Flag&market_center.DataFlag_Ticker != 0 {
        if isSpot {
            mccc.SubscribeSpotTicker(ExchangeName, CurrencyPair, Period)
        } else {
            mccc.SubscribeFutureTicker(ExchangeName, ContractType, CurrencyPair, Period)
        }
    }
    c := client.NewClient(ExchangeName, CurrencyPair, "", mccc)
    c.GetTicker()
```