// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: emitter/emitter.proto

package emitter

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

const (
	Emitter_GetDNS_FullMethodName   = "/emitter.Emitter/GetDNS"
	Emitter_GetWhois_FullMethodName = "/emitter.Emitter/GetWhois"
)

// EmitterClient is the client API for Emitter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmitterClient interface {
	GetDNS(ctx context.Context, in *GetDNSRequest, opts ...grpc.CallOption) (*ResourceRecords, error)
	GetWhois(ctx context.Context, in *GetWhoisRequest, opts ...grpc.CallOption) (*WhoisRecord, error)
}

type emitterClient struct {
	cc grpc.ClientConnInterface
}

func NewEmitterClient(cc grpc.ClientConnInterface) EmitterClient {
	return &emitterClient{cc}
}

func (c *emitterClient) GetDNS(ctx context.Context, in *GetDNSRequest, opts ...grpc.CallOption) (*ResourceRecords, error) {
	out := new(ResourceRecords)
	err := c.cc.Invoke(ctx, Emitter_GetDNS_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emitterClient) GetWhois(ctx context.Context, in *GetWhoisRequest, opts ...grpc.CallOption) (*WhoisRecord, error) {
	out := new(WhoisRecord)
	err := c.cc.Invoke(ctx, Emitter_GetWhois_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmitterServer is the server API for Emitter service.
// All implementations must embed UnimplementedEmitterServer
// for forward compatibility
type EmitterServer interface {
	GetDNS(context.Context, *GetDNSRequest) (*ResourceRecords, error)
	GetWhois(context.Context, *GetWhoisRequest) (*WhoisRecord, error)
	mustEmbedUnimplementedEmitterServer()
}

// UnimplementedEmitterServer must be embedded to have forward compatible implementations.
type UnimplementedEmitterServer struct {
}

func (UnimplementedEmitterServer) GetDNS(context.Context, *GetDNSRequest) (*ResourceRecords, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDNS not implemented")
}
func (UnimplementedEmitterServer) GetWhois(context.Context, *GetWhoisRequest) (*WhoisRecord, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWhois not implemented")
}
func (UnimplementedEmitterServer) mustEmbedUnimplementedEmitterServer() {}

// UnsafeEmitterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmitterServer will
// result in compilation errors.
type UnsafeEmitterServer interface {
	mustEmbedUnimplementedEmitterServer()
}

func RegisterEmitterServer(s grpc.ServiceRegistrar, srv EmitterServer) {
	s.RegisterService(&Emitter_ServiceDesc, srv)
}

func _Emitter_GetDNS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDNSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmitterServer).GetDNS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Emitter_GetDNS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmitterServer).GetDNS(ctx, req.(*GetDNSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Emitter_GetWhois_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWhoisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmitterServer).GetWhois(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Emitter_GetWhois_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmitterServer).GetWhois(ctx, req.(*GetWhoisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Emitter_ServiceDesc is the grpc.ServiceDesc for Emitter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Emitter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "emitter.Emitter",
	HandlerType: (*EmitterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDNS",
			Handler:    _Emitter_GetDNS_Handler,
		},
		{
			MethodName: "GetWhois",
			Handler:    _Emitter_GetWhois_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "emitter/emitter.proto",
}
