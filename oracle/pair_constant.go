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

func (p *PairConstant) getPriceOf(a string) (float64, error) {
	p.RLock()
	defer p.RUnlock()
	if p.Price == 0 {
		return 0, errors.New("price is 0")
	} else if a == p.BaseAsset {
		return p.Price, nil
	} else if a == p.QuoteAsset {
		return 1 / p.Price, nil
	} else {
		return 0, errors.New("asset not found in PairConstant")
	}
}

func (p *PairConstant) getWeight() (float64, error) {
	p.RLock()
	defer p.RUnlock()
	return p.Weight, nil
}

func (p *PairConstant) update(price float64) {
	p.Lock()
	defer p.Unlock()
	p.Price = price
}
