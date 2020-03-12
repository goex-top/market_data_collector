package data

import "testing"

func TestDepthDataFromCSVeData_Load(t *testing.T) {
	dd := NewDepthDataParser("../csv/")
	t.Log(dd.Load("depth_binance.com_BTC_USDT_2020-01-30.csv"))
}
