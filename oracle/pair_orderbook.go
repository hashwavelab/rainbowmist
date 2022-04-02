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

func (p *PairAtObex) getPriceOf(a string) (float64, error) {
	p.RLock()
	defer p.RUnlock()
	if p.Price == 0 {
		return 0, errors.New("price is 0")
	} else if a == p.BaseAsset {
		return p.Price, nil
	} else if a == p.QuoteAsset {
		return 1 / p.Price, nil
	} else {
		return 0, errors.New("asset not found in PairAtBinance")
	}
}

func (p *PairAtObex) getWeight() (float64, error) {
	p.RLock()
	defer p.RUnlock()
	return p.Weight, nil
}

func (p *PairAtObex) update(price float64) {
	p.Lock()
	defer p.Unlock()
	// make sure only a sensible price can be stored
	ok := priceSenseCheck(price)
	if !ok {
		return
	}
	p.Price = price
}

func priceSenseCheck(price float64) bool {
	if price < 1.0/1000000000.0 || price > 1000000000.0 {
		return false
	}
	return true
}
