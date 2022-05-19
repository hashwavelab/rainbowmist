package client

import (
	"context"
	"log"
	sync "sync"
	"time"

	"github.com/hashwavelab/rainbowmist/pb"
	"github.com/hashwavelab/rainbowmist/pix"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	RPCTimeout               = time.Second * 1
	SyncInterval             = time.Second * 2
	TimeAllowedSinceLastSync = time.Second * 5
)

type Oracle struct {
	cli               pb.RainbowmistClient
	USDQuoteWatchList sync.Map //map[string]*USDQuote
}

type USDQuote struct {
	sync.RWMutex
	Asset    string
	LastSync time.Time
	Price    float64
}

func NewOracle(address string) *Oracle {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("grpc dial failed:", err)
	}
	oracle := &Oracle{
		cli: pb.NewRainbowmistClient(conn),
	}
	return oracle
}

func (_o *Oracle) AddUSDWatchPair(a string) {
	_, ok := _o.USDQuoteWatchList.Load(a)
	if ok {
		return
	}
	_o.USDQuoteWatchList.Store(a, &USDQuote{Asset: a})
}

func (_o *Oracle) StartSyncing() {
	go func() {
		ticker := time.NewTicker(SyncInterval)
		for range ticker.C {
			_o.syncWatchList()
		}
	}()
}

func (_o *Oracle) syncWatchList() {
	wg := &sync.WaitGroup{}
	_o.USDQuoteWatchList.Range(func(k, v interface{}) bool {
		wg.Add(1)
		go _o.syncPrice(v.(*USDQuote), wg)
		return true
	})
	wg.Wait()
}

func (_o *Oracle) syncPrice(q *USDQuote, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	r, err := _o.cli.GetUSDPrice(ctx, &pb.GetUSDPriceRequest{
		Asset: q.Asset,
	})
	if err != nil {
		q.setError()
		return
	}
	q.setPrice(r.Price)
}

func (_o *Oracle) GetUSDPrice(a string) (float64, bool) {
	q, ok := _o.USDQuoteWatchList.Load(a)
	if !ok {
		return 0, false
	}
	quote := q.(*USDQuote)
	return quote.getPrice()
}

func (_o *Oracle) GetUSDPriceFunc(a string) (func() (float64, bool), bool) {
	q, ok := _o.USDQuoteWatchList.Load(a)
	if !ok {
		return nil, false
	}
	quote := q.(*USDQuote)
	return quote.getPrice, true
}

func (_q *USDQuote) getPrice() (float64, bool) {
	_q.RLock()
	defer _q.RUnlock()
	if time.Since(_q.LastSync) <= TimeAllowedSinceLastSync {
		return _q.Price, true
	} else {
		return 0, false
	}
}

func (_q *USDQuote) setError() {
	_q.Lock()
	defer _q.Unlock()
	_q.LastSync = time.Time{}
}

func (_q *USDQuote) setPrice(price float64) {
	_q.Lock()
	defer _q.Unlock()
	ok := pix.PriceSenseCheck(price)
	if !ok {
		return
	}
	_q.LastSync = time.Now() //maybe get from rainbowmist instead // TODO, Implement this please!
	_q.Price = price
}
