package client

import (
	"context"
	sync "sync"
	"time"

	"github.com/hashwavelab/rainbowmist/pb"
	"github.com/hashwavelab/rainbowmist/pix"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	RPCTimeout               = time.Second * 2
	SyncInterval             = time.Second * 2
	TimeAllowedSinceLastSync = time.Second * 15
)

type Oracle struct {
	address           string
	USDQuoteWatchList sync.Map //map[string]*USDQuote
}

type USDQuote struct {
	sync.RWMutex
	Asset    string
	LastSync time.Time
	Price    float64
}

func NewOracle(address string) *Oracle {
	oracle := &Oracle{
		address: address,
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

func (_o *Oracle) connect() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, _o.address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	return conn, err
}

func (_o *Oracle) syncPrice(q *USDQuote, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := _o.connect()
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewRainbowmistClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	r, err := c.GetUSDPrice(ctx, &pb.GetUSDPriceRequest{
		Asset: q.Asset,
	})
	if err != nil {
		return
	}
	q.SetPrice(r.Price)
}

func (_o *Oracle) GetUSDPrice(a string) (float64, bool) {
	q, ok := _o.USDQuoteWatchList.Load(a)
	if !ok {
		return 0, false
	}
	quote := q.(*USDQuote)
	return quote.GetPrice()
}

func (_o *Oracle) GetUSDPriceFunc(a string) (func() (float64, bool), bool) {
	q, ok := _o.USDQuoteWatchList.Load(a)
	if !ok {
		return nil, false
	}
	quote := q.(*USDQuote)
	return quote.GetPrice, true
}

func (_q *USDQuote) GetPrice() (float64, bool) {
	_q.RLock()
	defer _q.RUnlock()
	if time.Since(_q.LastSync) <= TimeAllowedSinceLastSync {
		return _q.Price, true
	} else {
		return 0, false
	}
}

func (_q *USDQuote) SetPrice(price float64) {
	_q.Lock()
	defer _q.Unlock()
	ok := pix.PriceSenseCheck(price)
	if !ok {
		return
	}
	_q.LastSync = time.Now() //maybe get from rainbowmist instead // TODO, Implement this please!
	_q.Price = price
}
