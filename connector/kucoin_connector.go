package connector

import (
	"log"
	"strconv"

	kucoin "github.com/Kucoin/kucoin-go-sdk"
)

type KucoinConnector struct {
	client  *kucoin.ApiService
	infoMap map[string]*OrderBookInfo
}

func NewKucoinConnector() *KucoinConnector {
	s := kucoin.NewApiService(
		kucoin.ApiKeyOption("key"),
		kucoin.ApiSecretOption("secret"),
		kucoin.ApiPassPhraseOption("passphrase"),
		kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2),
	)
	c := &KucoinConnector{
		client:  s,
		infoMap: make(map[string]*OrderBookInfo),
	}
	c.GetExchangeInfo()
	return c
}

func (_c *KucoinConnector) GetExchangeInfo() ([]*OrderBookInfo, error) {
	rsp, err := _c.client.Symbols("")
	if err != nil {
		log.Fatal("GetExchangeInfo Error", err)
	}
	symbols := kucoin.SymbolsModel{}
	err1 := rsp.ReadData(&symbols)
	if err1 != nil {
		log.Fatal("GetExchangeInfo ReadData Error", err)
	}
	r := make([]*OrderBookInfo, 0)
	for _, symbol := range symbols {
		obi := &OrderBookInfo{
			Symbol:     symbol.Symbol,
			BaseAsset:  symbol.BaseCurrency,
			QuoteAsset: symbol.QuoteCurrency,
		}
		_c.infoMap[symbol.Symbol] = obi
		r = append(r, obi)
	}
	return r, nil
}

func (_c *KucoinConnector) GetAllTicker() ([]*OrderBookTicker, error) {
	rsp, err := _c.client.Tickers()
	if err != nil {
		return nil, err
	}
	tickers := kucoin.TickersResponseModel{}
	err1 := rsp.ReadData(&tickers)
	if err1 != nil {
		return nil, err
	}
	r := make([]*OrderBookTicker, 0)
	for _, t := range tickers.Tickers {
		info, ok := _c.infoMap[t.Symbol]
		if !ok {
			continue
		}
		ask, _ := strconv.ParseFloat(t.Sell, 64)
		bid, _ := strconv.ParseFloat(t.Buy, 64)
		obt := &OrderBookTicker{
			BaseAsset:  info.BaseAsset,
			QuoteAsset: info.QuoteAsset,
			Price:      (ask + bid) / 2,
		}
		r = append(r, obt)
	}
	return r, nil
}
