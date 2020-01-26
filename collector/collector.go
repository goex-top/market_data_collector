package collector

import (
	"context"
	"fmt"
	"github.com/goex-top/market_center"
	"github.com/goex-top/market_data_collector/client"
	"github.com/goex-top/market_data_collector/storage"
	"time"
)

func NewCollector(ctx context.Context, c *client.Client, period int64, flag market_center.DataFlag, csvStore *storage.CsvStorage) {
	fmt.Printf("(%s) %s new collector with flag[%d]\n", c.ExchangeName, c.CurrencyPair, flag)
	go func() {
		tick := time.NewTicker(time.Millisecond * time.Duration(period))
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("(%s) %s collector exit\n", c.ExchangeName, c.CurrencyPair)
				return
			case <-tick.C:
				if flag&market_center.DataFlag_Depth != 0 {
					depth := c.GetDepth()
					if depth != nil {
						if depth.UTime.UnixNano() == -6795364578871345152 {
							depth.UTime = time.Now()
						}
						csvStore.SaveDepth(depth)
					}
				}

				if flag&market_center.DataFlag_Ticker != 0 {
					ticker := c.GetTicker()
					if ticker != nil {
						if ticker.Date == 0 {
							ticker.Date = uint64(time.Now().UnixNano() / int64(time.Millisecond))
						}
						csvStore.SaveTicker(ticker)
					}
				}
			}
		}
	}()
}
