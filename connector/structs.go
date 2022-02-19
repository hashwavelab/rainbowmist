package connector

type OrderBookInfo struct {
	Symbol     string
	BaseAsset  string
	QuoteAsset string
}

type OrderBookTicker struct {
	BaseAsset  string
	QuoteAsset string
	Price      float64
}
