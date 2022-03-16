package client

import (
	"context"
	pb "rainbowmist/pb"
	"strconv"
	sync "sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	RPCTimeout               = time.Second * 2
	SyncInterval             = time.Second * 5
	TimeAllowedSinceLastSync = time.Second * 30
	Decimals                 = "10"
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
			go _o.syncWatchList()
		}
	}()
}

func (_o *Oracle) syncWatchList() {
	_o.USDQuoteWatchList.Range(func(k, v interface{}) bool {
		go _o.syncPrice(v.(*USDQuote))
		return true
	})
}

func (_o *Oracle) connect() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, _o.address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	return conn, err
}

func (_o *Oracle) syncPrice(q *USDQuote) {
	conn, err := _o.connect()
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewRainbowmistClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	r, err := c.GetUSDPrice(ctx, &pb.GetUSDPriceRequest{
		Asset:    q.Asset,
		Decimals: Decimals,
	})
	if err != nil {
		return
	}
	if !r.Status {
		return
	}
	price, _ := strconv.ParseFloat(r.Price, 64)
	q.SetPrice(price)
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
	_q.LastSync = time.Now() //maybe get from rainbowmist instead // TODO, Implement this please!
	_q.Price = price
}
