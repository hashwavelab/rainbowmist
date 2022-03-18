/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"time"

	"github.com/hashwavelab/rainbowmist/pb"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:8889"
	defaultName = "piglet"
)

func main() {
	f2()
}

func f1() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRainbowmistClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	r, err := c.GetPrice(ctx, &pb.GetPriceRequest{
		BaseAsset:  "ETH",
		QuoteAsset: "BTC",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Received Reserve: %s", r)
}

func f2() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRainbowmistClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	r, err := c.GetUSDPrice(ctx, &pb.GetUSDPriceRequest{
		Asset: "ONE",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Received Reserve: %s", r)
}
