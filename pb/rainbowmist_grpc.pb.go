// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package rainbowmist

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RainbowmistClient is the client API for Rainbowmist service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RainbowmistClient interface {
	GetPrice(ctx context.Context, in *GetPriceRequest, opts ...grpc.CallOption) (*GetPriceReply, error)
	GetUSDPrice(ctx context.Context, in *GetUSDPriceRequest, opts ...grpc.CallOption) (*GetPriceReply, error)
}

type rainbowmistClient struct {
	cc grpc.ClientConnInterface
}

func NewRainbowmistClient(cc grpc.ClientConnInterface) RainbowmistClient {
	return &rainbowmistClient{cc}
}

func (c *rainbowmistClient) GetPrice(ctx context.Context, in *GetPriceRequest, opts ...grpc.CallOption) (*GetPriceReply, error) {
	out := new(GetPriceReply)
	err := c.cc.Invoke(ctx, "/Rainbowmist/GetPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rainbowmistClient) GetUSDPrice(ctx context.Context, in *GetUSDPriceRequest, opts ...grpc.CallOption) (*GetPriceReply, error) {
	out := new(GetPriceReply)
	err := c.cc.Invoke(ctx, "/Rainbowmist/GetUSDPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RainbowmistServer is the server API for Rainbowmist service.
// All implementations must embed UnimplementedRainbowmistServer
// for forward compatibility
type RainbowmistServer interface {
	GetPrice(context.Context, *GetPriceRequest) (*GetPriceReply, error)
	GetUSDPrice(context.Context, *GetUSDPriceRequest) (*GetPriceReply, error)
	mustEmbedUnimplementedRainbowmistServer()
}

// UnimplementedRainbowmistServer must be embedded to have forward compatible implementations.
type UnimplementedRainbowmistServer struct {
}

func (UnimplementedRainbowmistServer) GetPrice(context.Context, *GetPriceRequest) (*GetPriceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrice not implemented")
}
func (UnimplementedRainbowmistServer) GetUSDPrice(context.Context, *GetUSDPriceRequest) (*GetPriceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUSDPrice not implemented")
}
func (UnimplementedRainbowmistServer) mustEmbedUnimplementedRainbowmistServer() {}

// UnsafeRainbowmistServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RainbowmistServer will
// result in compilation errors.
type UnsafeRainbowmistServer interface {
	mustEmbedUnimplementedRainbowmistServer()
}

func RegisterRainbowmistServer(s grpc.ServiceRegistrar, srv RainbowmistServer) {
	s.RegisterService(&Rainbowmist_ServiceDesc, srv)
}

func _Rainbowmist_GetPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RainbowmistServer).GetPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Rainbowmist/GetPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RainbowmistServer).GetPrice(ctx, req.(*GetPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rainbowmist_GetUSDPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUSDPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RainbowmistServer).GetUSDPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Rainbowmist/GetUSDPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RainbowmistServer).GetUSDPrice(ctx, req.(*GetUSDPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Rainbowmist_ServiceDesc is the grpc.ServiceDesc for Rainbowmist service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rainbowmist_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Rainbowmist",
	HandlerType: (*RainbowmistServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPrice",
			Handler:    _Rainbowmist_GetPrice_Handler,
		},
		{
			MethodName: "GetUSDPrice",
			Handler:    _Rainbowmist_GetUSDPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/rainbowmist.proto",
}
