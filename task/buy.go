package task

import (
	"github.com/huobirdcenter/huobi_golang/config"
	"github.com/huobirdcenter/huobi_golang/logging/applogger"
	"github.com/huobirdcenter/huobi_golang/pkg/client"
	"github.com/huobirdcenter/huobi_golang/pkg/model/common"
	"github.com/huobirdcenter/huobi_golang/pkg/model/order"
	"github.com/jinzhu/now"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

////购买几次不重复的币
//var count = 2
////最大最小涨幅
//var minIncrease = decimal.NewFromFloat(0.15)
//var maxIncrease = decimal.NewFromFloat(0.3)
////购买多少usdt的币
//var utotal = 3000.0

func Buy() {
	//幂等，保证某种币只买一次
	idempotence := make(map[string]struct{})
	symbolMap := make(map[string]common.Symbol)
	commonClient := new(client.CommonClient).Init(config.Host)
	//加载交易对
	symbols, err := commonClient.GetSymbols()
	if err != nil {
		applogger.Error(err.Error())
		return
	}
	for _, symbol := range symbols {
		symbolMap[symbol.Symbol] = symbol
	}
	marketClient := new(client.MarketClient).Init(config.Host)
	orderClient := new(client.OrderClient).Init(config.AccessKey, config.SecretKey, config.Host)
	ticker := time.NewTicker(time.Millisecond * 100) //一秒十次
	for range ticker.C {
		//判断是否是凌晨
		if time.Now().Sub(now.BeginningOfDay()).Seconds() > 3600 {
			continue
		}
		resp, err := marketClient.GetAllSymbolsLast24hCandlesticksAskBid()
		if err != nil {
			applogger.Error(err.Error())
		}
		for _, tick := range resp {
			//结束
			if config.Count == 0 {
				return
			}
			//首先要是usdt交易对
			if !strings.Contains(tick.Symbol, "usdt") || strings.Contains(tick.Symbol, "3") {
				continue
			}
			//applogger.Info("Symbol: %s, High: %v, Low: %v, Ask[%v, %v], Bid[%v, %v], Open: %v Close: %v",
			//	tick.Symbol, tick.High, tick.Low, tick.Ask, tick.AskSize, tick.Bid, tick.BidSize,tick.Open,tick.Close)
			//涨幅计算
			increase := tick.Close.Sub(tick.Open).Div(tick.Open)
			//涨幅小于maxIncrease，大于minIncrease 下单
			if increase.GreaterThanOrEqual(config.MinIncrease) && increase.LessThanOrEqual(config.MaxIncrease) {
				//防止某个交易对重复下单
				_, ok := idempotence[tick.Symbol]
				if ok {
					continue
				}
				symbol, ok := symbolMap[tick.Symbol]
				if !ok {
					continue
				}
				applogger.Info("Symbol: %s Open: %v Close: %v 涨幅:%v", tick.Symbol, tick.Open, tick.Close, increase)
				//todo市价购买3000刀
				applogger.Info("=====购买=====")
				//精度
				pricePrecision := symbol.PricePrecision
				amountPrecision := symbol.AmountPrecision
				//当前价格的1.1倍
				price := tick.Close.Mul(decimal.NewFromFloat(1.08))
				price = price.Round(int32(pricePrecision))
				//数量
				amount := decimal.NewFromFloat(config.Utotal).Div(price).Mul(decimal.NewFromFloat(0.99))
				amount = amount.Round(int32(amountPrecision))
				if amount.GreaterThan(symbol.LimitOrderMaxOrderAmt) {
					amount = symbol.LimitOrderMaxOrderAmt
				}
				request := order.PlaceOrderRequest{
					AccountId: config.AccountId,
					Type:      "buy-ioc",
					Source:    "spot-api",
					Symbol:    tick.Symbol,
					Price:     price.String(),
					Amount:    amount.String(),
				}
				resp, err := orderClient.PlaceOrder(&request)
				if err != nil {
					applogger.Error(err.Error())
					continue
				}
				switch resp.Status {
				case "ok":
					applogger.Info("Place order successfully, order id: %s", resp.Data)
				case "error":
					applogger.Error("Place order error: %s", resp.ErrorMessage)
				}
				//幂等
				idempotence[tick.Symbol] = struct{}{}
				config.Count--
			}
		}
	}

}
