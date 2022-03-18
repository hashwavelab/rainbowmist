package oracle

import (
	"errors"
	"sync"
)

type PairAtObex struct {
	sync.RWMutex
	Weight     float64
	BaseAsset  string
	QuoteAsset string
	Price      float64
}

func (_p *PairAtObex) getPriceOf(a string) (float64, error) {
	_p.RLock()
	defer _p.RUnlock()
	if _p.Price == 0 {
		return 0, errors.New("price is 0")
	} else if a == _p.BaseAsset {
		return _p.Price, nil
	} else if a == _p.QuoteAsset {
		return 1 / _p.Price, nil
	} else {
		return 0, errors.New("asset not found in PairAtBinance")
	}
}

func (_p *PairAtObex) getWeight() (float64, error) {
	_p.RLock()
	defer _p.RUnlock()
	return _p.Weight, nil
}

func (_p *PairAtObex) update(p float64) {
	_p.Lock()
	defer _p.Unlock()
	// make sure only a sensible price can be stored
	ok := priceSenseCheck(p)
	if !ok {
		return
	}
	_p.Price = p
}

func priceSenseCheck(price float64) bool {
	if price < 1.0/1000000000.0 || price > 1000000000.0 {
		return false
	}
	return true
}
