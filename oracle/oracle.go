package oracle

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/hashwavelab/rainbowmist/connector"
)

type Oracle struct {
	sync.RWMutex
	MapLock sync.RWMutex
	Map     map[string]*Pair
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
		Map: make(map[string]*Pair),
	}
	return o
}

func (_o *Oracle) Init() {
	_o.initBinance()
	_o.initKucoin()
	_o.initConstant()
}

func (_o *Oracle) initBinance() {
	sourceName := "binance"
	bc := connector.NewBinanceConnector()
	infos, _ := bc.GetExchangeInfo()
	for _, info := range infos {
		pair := _o.GetNewPair(info.BaseAsset, info.QuoteAsset)
		pair.SetupPairAtOneSource(sourceName, "obex_spot", info.BaseAsset, info.QuoteAsset, 8)
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

func (_o *Oracle) initKucoin() {
	sourceName := "kucoin"
	bc := connector.NewKucoinConnector()
	infos, _ := bc.GetExchangeInfo()
	for _, info := range infos {
		pair := _o.GetNewPair(info.BaseAsset, info.QuoteAsset)
		pair.SetupPairAtOneSource(sourceName, "obex_spot", info.BaseAsset, info.QuoteAsset, 2)
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

func (_o *Oracle) initConstant() {
	sourceName := "constant"
	pair := _o.GetNewPair("BUSD", "USD")
	pair.SetupPairAtOneSource(sourceName, "constant", "BUSD", "USD", 100)
	pair.Update(sourceName, 1)
	pair1 := _o.GetNewPair("USDC", "USD")
	pair1.SetupPairAtOneSource(sourceName, "constant", "USDC", "USD", 100)
	pair1.Update(sourceName, 1)
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
	//p, _ := pair.GetPriceOf(a0)
	//log.Println(a0, a1, pairKey, p)
	return pair.GetPriceOf(a0)
}

func (_o *Oracle) GetUSDPrice(a string) (float64, error) {
	price0, err0 := _o.GetPrice(a, "USD")
	if err0 != nil {
		if a == "USDT" {
			price1, err := _o.GetPrice(a, "BUSD")
			if err != nil {
				return 0, errors.New("pair not found")
			}
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
		return price1 * price2 * price3, nil
	} else {
		return price0, nil
	}
}
