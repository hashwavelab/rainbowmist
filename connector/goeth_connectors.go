package connector

import (
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

type GoEthConnectors struct {
	sync.RWMutex
	connectors map[string]*ethclient.Client
}

func NewGoEthConnectors() *GoEthConnectors {
	return &GoEthConnectors{
		connectors: make(map[string]*ethclient.Client),
	}
}

func (c *GoEthConnectors) GetClinet(url string) *ethclient.Client {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.connectors[url]; !ok {
		cli, err := ethclient.Dial(url)
		if err != nil {
			return nil
		}
		c.connectors[url] = cli
	}
	return c.connectors[url]
}
