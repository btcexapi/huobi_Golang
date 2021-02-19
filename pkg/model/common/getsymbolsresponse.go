package common

import "github.com/shopspring/decimal"

type GetSymbolsResponse struct {
	Status string   `json:"status"`
	Data   []Symbol `json:"data"`
}

type Symbol struct {
	BaseCurrency           string          `json:"base-currency"`//交易对中的基础币种
	QuoteCurrency          string          `json:"quote-currency"`//交易对中的报价币种:一般是USDT
	PricePrecision         int             `json:"price-precision"`//交易对报价的精度
	AmountPrecision        int             `json:"amount-precision"`//交易对基础币种计数精度
	SymbolPartition        string          `json:"symbol-partition"`
	Symbol                 string          `json:"symbol"`
	State                  string          `json:"state"`
	ValuePrecision         int             `json:"value-precision"`
	LimitOrderMinOrderAmt  decimal.Decimal `json:"limit-order-min-order-amt"`
	LimitOrderMaxOrderAmt  decimal.Decimal `json:"limit-order-max-order-amt"`
	SellMarketMinOrderAmt  decimal.Decimal `json:"sell-market-min-order-amt"`
	SellMarketMaxOrderAmt  decimal.Decimal `json:"sell-market-max-order-amt"`
	BuyMarketMaxOrderValue decimal.Decimal `json:"buy-market-max-order-value"`
	MinOrderValue          decimal.Decimal `json:"min-order-value"`
	MaxOrderValue          decimal.Decimal `json:"max-order-value"`
	LeverageRatio          decimal.Decimal `json:"leverage-ratio"`
}
