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

## Format
### ticker
![ticker](ticker.png)

|  symbol | type | description |
|  ----  | ----  | ----  |
| t  | int | timestamp |
| b  | float | best bid |
| s  | float | best ask |
| h  | float | high price |
| l  | float | low price |
| v  | float | volume |

### orderbook
![orderbook](orderbook.png)

|  symbol | type | description |
|  ----  | ----  | ----  |
| t  | int | timestamp |
| a  | array | asks list with size 20, each element is [p,q], p:price, q:qty |
| b  | array | bids list with size 20, each element is [p,q], p:price, q:qty |

## Support Data
* Ticker 
* Depth(Orderbook)
* ~~Kline~~

## TODO
* InfluxDB

### 观星者

![观星者](https://starchart.cc/goex-top/market_data_collector.svg)