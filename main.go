package main

import (
	"log"

	"github.com/hashwavelab/rainbowmist/oracle"
)

func main() {
	log.Println("Shall rainbow mist be your gateway")
	//connector.SubscribeBinanceBookTicker("BTCUSDT")
	o := oracle.NewOracle()
	// for _, recipe := range defaultRecipes {
	// 	o.AddNewPair(recipe)
	// }
	o.Init()
	InitGrpcServer(o)
}
