package connector

import (
	"context"
	"log"
	"strconv"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceBookTickerUpdate struct {
	BestAskPrice    string
	BestAskQuantity string
	BestBidPrice    string
	BestBidQuantity string
}

type BinanceConnector struct {
	client  *binance.Client
	infoMap map[string]*OrderBookInfo
}

func NewBinanceConnector() *BinanceConnector {
	c := &BinanceConnector{
		client:  binance.NewClient("", ""),
		infoMap: make(map[string]*OrderBookInfo),
	}
	c.GetExchangeInfo()
	return c
}

func (_c *BinanceConnector) GetExchangeInfo() ([]*OrderBookInfo, error) {
	info, err := _c.client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal("NewBinanceConnector Error", err)
	}
	r := make([]*OrderBookInfo, 0)
	for _, ob := range info.Symbols {
		obi := &OrderBookInfo{
			Symbol:     ob.Symbol,
			BaseAsset:  ob.BaseAsset,
			QuoteAsset: ob.QuoteAsset,
		}
		_c.infoMap[ob.Symbol] = obi
		r = append(r, obi)
	}
	return r, nil
}
func (_c *BinanceConnector) GetAllTicker() ([]*OrderBookTicker, error) {
	bt, err := _c.client.NewListBookTickersService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	r := make([]*OrderBookTicker, 0)
	for _, t := range bt {
		info, ok := _c.infoMap[t.Symbol]
		if !ok {
			continue
		}
		ask, _ := strconv.ParseFloat(t.AskPrice, 64)
		bid, _ := strconv.ParseFloat(t.BidPrice, 64)
		obt := &OrderBookTicker{
			BaseAsset:  info.BaseAsset,
			QuoteAsset: info.QuoteAsset,
			Price:      (ask + bid) / 2,
		}
		r = append(r, obt)
	}
	return r, nil
}
