// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/soda.proto

package proto

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

// SodaServiceClient is the client API for SodaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SodaServiceClient interface {
	RandomNumber(ctx context.Context, in *RandomNumberRequest, opts ...grpc.CallOption) (*RandomNumberResponse, error)
}

type sodaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSodaServiceClient(cc grpc.ClientConnInterface) SodaServiceClient {
	return &sodaServiceClient{cc}
}

func (c *sodaServiceClient) RandomNumber(ctx context.Context, in *RandomNumberRequest, opts ...grpc.CallOption) (*RandomNumberResponse, error) {
	out := new(RandomNumberResponse)
	err := c.cc.Invoke(ctx, "/soda.SodaService/RandomNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SodaServiceServer is the server API for SodaService service.
// All implementations must embed UnimplementedSodaServiceServer
// for forward compatibility
type SodaServiceServer interface {
	RandomNumber(context.Context, *RandomNumberRequest) (*RandomNumberResponse, error)
	mustEmbedUnimplementedSodaServiceServer()
}

// UnimplementedSodaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSodaServiceServer struct {
}

func (UnimplementedSodaServiceServer) RandomNumber(context.Context, *RandomNumberRequest) (*RandomNumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RandomNumber not implemented")
}
func (UnimplementedSodaServiceServer) mustEmbedUnimplementedSodaServiceServer() {}

// UnsafeSodaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SodaServiceServer will
// result in compilation errors.
type UnsafeSodaServiceServer interface {
	mustEmbedUnimplementedSodaServiceServer()
}

func RegisterSodaServiceServer(s grpc.ServiceRegistrar, srv SodaServiceServer) {
	s.RegisterService(&SodaService_ServiceDesc, srv)
}

func _SodaService_RandomNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RandomNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SodaServiceServer).RandomNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/soda.SodaService/RandomNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SodaServiceServer).RandomNumber(ctx, req.(*RandomNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SodaService_ServiceDesc is the grpc.ServiceDesc for SodaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SodaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "soda.SodaService",
	HandlerType: (*SodaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RandomNumber",
			Handler:    _SodaService_RandomNumber_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/soda.proto",
}
