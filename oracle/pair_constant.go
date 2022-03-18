package oracle

import (
	"errors"
	"sync"
)

type PairConstant struct {
	sync.RWMutex
	Weight     float64
	BaseAsset  string
	QuoteAsset string
	Price      float64
}

func (_p *PairConstant) getPriceOf(a string) (float64, error) {
	_p.RLock()
	defer _p.RUnlock()
	if _p.Price == 0 {
		return 0, errors.New("price is 0")
	} else if a == _p.BaseAsset {
		return _p.Price, nil
	} else if a == _p.QuoteAsset {
		return 1 / _p.Price, nil
	} else {
		return 0, errors.New("asset not found in PairConstant")
	}
}

func (_p *PairConstant) getWeight() (float64, error) {
	_p.RLock()
	defer _p.RUnlock()
	return _p.Weight, nil
}

func (_p *PairConstant) update(p float64) {
	_p.Lock()
	defer _p.Unlock()
	_p.Price = p
}
