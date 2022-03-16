package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hashwavelab/rainbowmist/oracle"

	"github.com/hashwavelab/rainbowmist/pb"

	"google.golang.org/grpc"
)

const (
	port = ":8889"
)

type server struct {
	pb.UnimplementedRainbowmistServer
	Oracle *oracle.Oracle
}

func (s *server) GetPrice(ctx context.Context, in *pb.GetPriceRequest) (*pb.GetPriceReply, error) {
	p, err := s.Oracle.GetPrice(in.BaseAsset, in.QuoteAsset)
	if err != nil {
		return &pb.GetPriceReply{Status: false}, nil
	}
	pString := fmt.Sprintf("%."+in.Decimals+"f", p)
	return &pb.GetPriceReply{Price: pString, Status: true}, nil
}

func (s *server) GetUSDPrice(ctx context.Context, in *pb.GetUSDPriceRequest) (*pb.GetPriceReply, error) {
	p, err := s.Oracle.GetUSDPrice(in.Asset)
	if err != nil {
		return &pb.GetPriceReply{Status: false}, nil
	}
	pString := fmt.Sprintf("%."+in.Decimals+"f", p)
	return &pb.GetPriceReply{Price: pString, Status: true}, nil
}

func InitGrpcServer(oracle *oracle.Oracle) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRainbowmistServer(s, &server{Oracle: oracle})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
