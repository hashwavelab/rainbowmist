package oracle

import (
	"errors"
	"sync"
)

type Pair struct {
	sync.RWMutex
	Map map[string]PairAtOneSourceSession
}

type PairAtOneSourceSession interface {
	getPriceOf(string) (float64, error)
	getWeight() (float64, error)
	update(float64)
}

func NewPair() *Pair {
	pair := &Pair{
		Map: make(map[string]PairAtOneSourceSession),
	}
	return pair
}

func (_p *Pair) SetupPairAtOneSource(sourceName, sourceType, baseAsset, quoteAsset string, weight float64) {
	_p.Lock()
	defer _p.Unlock()
	switch sourceType {
	case "obex_spot":
		ps := &PairAtObex{
			Weight:     weight,
			BaseAsset:  baseAsset,
			QuoteAsset: quoteAsset,
		}
		_p.Map[sourceName] = ps
	case "constant":
		ps := &PairConstant{
			Weight:     weight,
			BaseAsset:  baseAsset,
			QuoteAsset: quoteAsset,
		}
		_p.Map[sourceName] = ps
	}
}

func (_p *Pair) Update(sourceName string, price float64) {
	_p.Lock()
	defer _p.Unlock()
	ps, ok := _p.Map[sourceName]
	if !ok {
		return
	}
	ps.update(price)
}

func (_p *Pair) GetPriceOf(a string) (float64, error) {
	_p.RLock()
	defer _p.RUnlock()
	prices := make([]float64, 0)
	weights := make([]float64, 0)
	var weightedPrice float64 = 0
	var weightSum float64 = 0
	for _, v := range _p.Map {
		p, err := v.getPriceOf(a)
		if err != nil {
			return 0, errors.New("asset not included in at least one source")
		}
		w, err := v.getWeight()
		if err != nil {
			return 0, errors.New("weight not set in at least one source")
		}
		prices = append(prices, p)
		weights = append(weights, w)
		weightSum += w
	}
	for i := 0; i < len(prices); i++ {
		weightedPrice += prices[i] * weights[i] / weightSum
	}
	return weightedPrice, nil
}
