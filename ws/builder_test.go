package ws

import (
	"github.com/goex-top/market_center"
	"testing"
)

func TestBuilderSpot(t *testing.T) {
	t.Log(BuilderSpot(market_center.BINANCE))
}
