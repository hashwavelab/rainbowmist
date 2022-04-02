package oracle

import (
	"context"
	"errors"
	"log"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashwavelab/rainbowmist/connector"
	pair "github.com/hashwavelab/rainbowmist/evm/contracts/UniV2Pair"
)

var (
	GoEthConnectors = connector.NewGoEthConnectors()
)

type PairAtUniV2 struct {
	sync.RWMutex
	pairInstance   *pair.Pair
	lastSynced     time.Time
	Weight         float64
	Asset0         string
	Asset1         string
	Asset0Decimals float64
	Asset1Decimals float64
	Asset0Reserve  float64
	Asset1Reserve  float64
}

func NewPairAtUniV2(url, address string, a0, a1 string, a0d, a1d, weight float64) *PairAtUniV2 {
	client := GoEthConnectors.GetClinet(url)
	if client == nil {
		log.Fatal("cannot get client")
	}
	pairInstance, err := pair.NewPair(common.HexToAddress(address), client)
	if err != nil {
		log.Fatal(err)
	}
	pair := &PairAtUniV2{
		Weight:         weight,
		pairInstance:   pairInstance,
		Asset0:         a0,
		Asset1:         a1,
		Asset0Decimals: a0d,
		Asset1Decimals: a1d,
	}
	return pair
}

func (p *PairAtUniV2) getPriceOf(a string) (float64, error) {
	p.RLock()
	defer p.RUnlock()
	if p.Asset0Reserve == 0 || p.Asset1Reserve == 0 {
		return 0, errors.New("reserve is 0")
	} else if time.Since(p.lastSynced) > 3*time.Minute {
		return 0, errors.New("outsynced")
	} else if a == p.Asset0 {
		return p.Asset1Reserve / p.Asset0Reserve * math.Pow(10, p.Asset0Decimals-p.Asset1Decimals), nil
	} else if a == p.Asset1 {
		return p.Asset0Reserve / p.Asset1Reserve * math.Pow(10, p.Asset1Decimals-p.Asset0Decimals), nil
	} else {
		return 0, errors.New("asset not found in PairAtBinance")
	}
}

func (p *PairAtUniV2) getWeight() (float64, error) {
	p.RLock()
	defer p.RUnlock()
	return p.Weight, nil
}

func (p *PairAtUniV2) update(price float64) {
	p.Lock()
	defer p.Unlock()
	CallOpts, cancel := getStandard3sCallOpts()
	defer cancel()
	res, err := p.pairInstance.GetReserves(CallOpts)
	if err != nil {
		log.Println("sync error", p.Asset0, p.Asset1)
	}
	p.lastSynced = time.Now()
	p.Asset0Reserve, _ = new(big.Float).SetInt(res.Reserve0).Float64()
	p.Asset1Reserve, _ = new(big.Float).SetInt(res.Reserve1).Float64()
}

func getStandard3sCallOpts() (*bind.CallOpts, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	return &bind.CallOpts{Context: ctx}, cancel
}
