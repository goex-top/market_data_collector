# Market Data Collector
![HitCount](http://hits.dwyl.io/goex-top/market_data_collector.svg)

Collect market data for quant analysis

## Storage
Store daily data in different `csv` files in `csv` folder, compress it to `tar` folder

`csv` folder(older file was removed automatically)
```
├── depth_binance.com_BTC_USDT_2020-01-26.csv
├── depth_fcoin.com_BTC_USDT_2020-01-26.csv
├── depth_huobi.pro_BTC_USDT_2020-01-26.csv
└── depth_okex.com_BTC_USDT_2020-01-26.csv
```

`tar` folder
```
.
├── binance.com_BTC_USDT_2020-01-24.tar.gz
├── binance.com_BTC_USDT_2020-01-25.tar.gz
├── fcoin.com_BTC_USDT_2020-01-24.tar.gz
├── fcoin.com_BTC_USDT_2020-01-25.tar.gz
├── huobi.pro_BTC_USDT_2020-01-24.tar.gz
├── huobi.pro_BTC_USDT_2020-01-25.tar.gz
├── okex.com_BTC_USDT_2020-01-24.tar.gz
└── okex.com_BTC_USDT_2020-01-25.tar.gz
```
## Support Data
* Ticker 
* Depth
* ~~Kline~~

## TODO
* InfluxDB

### 观星者

![观星者](https://starchart.cc/goex-top/market_data_collector.svg)