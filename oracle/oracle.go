package oracle

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/hashwavelab/rainbowmist/connector"
	"github.com/hashwavelab/rainbowmist/pix"
	"go.uber.org/zap"
)

type Oracle struct {
	sync.RWMutex
	MapLock sync.RWMutex
	Map     map[string]*Pair
	logger  *zap.Logger
}

type NewPairRecipe struct {
	Asset0  string
	Asset1  string
	Sources []NewPairSourceRecipe
}

type NewPairSourceRecipe struct {
	SourceName string
	SourceType string
	Weight     float64
	BaseAsset  string
	QuoteAsset string
	FixedPrice float64
}

func NewOracle() *Oracle {
	o := &Oracle{
		Map:    make(map[string]*Pair),
		logger: pix.NewLogger("/oracle.log"),
	}
	return o
}

func (_o *Oracle) Init() {
	_o.initBinance(80)
	_o.initKucoin(20)
	_o.initChainPairs()
	_o.initConstant()
}

func (_o *Oracle) initBinance(weight float64) {
	sourceName := "binance"
	bc := connector.NewBinanceConnector()
	infos, _ := bc.GetExchangeInfo()
	for _, info := range infos {
		pair := _o.GetNewPair(info.BaseAsset, info.QuoteAsset)
		pair.SetupPair(sourceName, "obex_spot", info.BaseAsset, info.QuoteAsset, weight)
	}
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			r, err := bc.GetAllTicker()
			if err != nil {
				log.Println("Binance get all ticker ERROR", err)
			} else {
				for _, t := range r {
					pair, ok := _o.ReadNewPair(t.BaseAsset, t.QuoteAsset)
					if ok {
						pair.Update(sourceName, t.Price)
					}
				}
			}
		}
	}()
}

func (_o *Oracle) initKucoin(weight float64) {
	sourceName := "kucoin"
	bc := connector.NewKucoinConnector()
	infos, _ := bc.GetExchangeInfo()
	for _, info := range infos {
		pair := _o.GetNewPair(info.BaseAsset, info.QuoteAsset)
		pair.SetupPair(sourceName, "obex_spot", info.BaseAsset, info.QuoteAsset, weight)
	}
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			r, err := bc.GetAllTicker()
			if err != nil {
				log.Println("Kucoin get all ticker ERROR", err)
			} else {
				for _, t := range r {
					pair, ok := _o.ReadNewPair(t.BaseAsset, t.QuoteAsset)
					if ok {
						pair.Update(sourceName, t.Price)
					}
				}
			}
		}
	}()
}

func (_o *Oracle) initChainPairs() {
	// DefaultChainPairSyncPeriod := 30 * time.Second
	// pair := _o.GetNewPair("JEWEL", "USD")
	// pair.SetupPoolAtChain("harmony", "URL_PH", "chain-evm", "univ2", "0xA1221A5BBEa699f507CC00bDedeA05b5d2e32Eba", "JEWEL", "USDC", 18, 6, 1)
	// go func() {
	// 	pair.Update("harmony", 0)
	// 	ticker := time.NewTicker(DefaultChainPairSyncPeriod)
	// 	for range ticker.C {
	// 		pair.Update("harmony", 0)
	// 	}
	// }()
}

func (_o *Oracle) initConstant() {
	sourceName := "constant"
	constantPriceOne := 1.0
	// Set all USDC, BUSD equivalent to USD (This relies on FTX)
	pair := _o.GetNewPair("BUSD", "USD")
	pair.SetupPair(sourceName, "constant", "BUSD", "USD", 100)
	pair.Update(sourceName, constantPriceOne)
	pair1 := _o.GetNewPair("USDC", "USD")
	pair1.SetupPair(sourceName, "constant", "USDC", "USD", 100)
	pair1.Update(sourceName, constantPriceOne)
	// for avalanche native USDCe and USDTe (Bridged Coins from ETH)
	pair2 := _o.GetNewPair("USDCe", "USD")
	pair2.SetupPair(sourceName, "constant", "USDCe", "USD", 100)
	pair2.Update(sourceName, constantPriceOne)
	pair3 := _o.GetNewPair("USDTe", "USDT")
	pair3.SetupPair(sourceName, "constant", "USDTe", "USDT", 100)
	pair3.Update(sourceName, constantPriceOne)
}

func (_o *Oracle) GetNewPair(a0, a1 string) *Pair {
	pairKey := getPairKey(a0, a1)
	_o.MapLock.Lock()
	defer _o.MapLock.Unlock()
	pair, ok := _o.Map[pairKey]
	if !ok {
		pair = NewPair()
		_o.Map[pairKey] = pair
	}
	return pair
}

func (_o *Oracle) ReadNewPair(a0, a1 string) (*Pair, bool) {
	pairKey := getPairKey(a0, a1)
	_o.MapLock.RLock()
	defer _o.MapLock.RUnlock()
	pair, ok := _o.Map[pairKey]
	return pair, ok
}

func getPairKey(a0, a1 string) string {
	if a0 >= a1 {
		return a0 + "RAINBOW&MIST" + a1
	} else {
		return a1 + "RAINBOW&MIST" + a0
	}
}

// GRPC

func (_o *Oracle) GetPrice(a0, a1 string) (float64, error) {
	// get price of a0 in the unit of a1
	pairKey := getPairKey(a0, a1)
	pair, ok := _o.Map[pairKey]
	if !ok {
		return 0, errors.New("pair not found")
	}
	price, err := pair.GetPriceOf(a0)
	_o.logger.Info("get price", zap.String("baseAsset", a0), zap.String("quoteAsset", a1), zap.Float64("price", price))
	return price, err
}

func (_o *Oracle) GetUSDPrice(a string) (float64, error) {
	price0, err0 := _o.GetPrice(a, "USD")
	if err0 != nil {
		if a == "USDT" {
			price1, err := _o.GetPrice(a, "BUSD")
			if err != nil {
				return 0, errors.New("pair not found")
			}
			_o.logger.Info("get USD price", zap.String("asset", a), zap.Float64("price", price1))
			return price1, nil
		}
		price1, err1 := _o.GetPrice(a, "USDT")
		if err1 != nil {
			return 0, errors.New("pair not found")
		}
		price2, err := _o.GetPrice("USDT", "BUSD")
		if err != nil {
			return 0, errors.New("pair not found")
		}
		price3, err := _o.GetPrice("BUSD", "USD")
		if err != nil {
			return 0, errors.New("pair not found")
		}
		_o.logger.Info("get USD price", zap.String("asset", a), zap.Float64("price", price1*price2*price3))
		return price1 * price2 * price3, nil
	} else {
		_o.logger.Info("get USD price", zap.String("asset", a), zap.Float64("price", price0))
		return price0, nil
	}
}
