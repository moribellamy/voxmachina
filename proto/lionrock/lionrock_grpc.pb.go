// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.5.1
// source: lionrock.proto

package lionrock

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

// TransactionalKeyValueStoreClient is the client API for TransactionalKeyValueStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransactionalKeyValueStoreClient interface {
	// Execute a transaction
	ExecuteTransaction(ctx context.Context, opts ...grpc.CallOption) (TransactionalKeyValueStore_ExecuteTransactionClient, error)
	// Execute a single database operation (with a transaction)
	Execute(ctx context.Context, in *DatabaseRequest, opts ...grpc.CallOption) (*DatabaseResponse, error)
}

type transactionalKeyValueStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionalKeyValueStoreClient(cc grpc.ClientConnInterface) TransactionalKeyValueStoreClient {
	return &transactionalKeyValueStoreClient{cc}
}

func (c *transactionalKeyValueStoreClient) ExecuteTransaction(ctx context.Context, opts ...grpc.CallOption) (TransactionalKeyValueStore_ExecuteTransactionClient, error) {
	stream, err := c.cc.NewStream(ctx, &TransactionalKeyValueStore_ServiceDesc.Streams[0], "/TransactionalKeyValueStore/executeTransaction", opts...)
	if err != nil {
		return nil, err
	}
	x := &transactionalKeyValueStoreExecuteTransactionClient{stream}
	return x, nil
}

type TransactionalKeyValueStore_ExecuteTransactionClient interface {
	Send(*StreamingDatabaseRequest) error
	Recv() (*StreamingDatabaseResponse, error)
	grpc.ClientStream
}

type transactionalKeyValueStoreExecuteTransactionClient struct {
	grpc.ClientStream
}

func (x *transactionalKeyValueStoreExecuteTransactionClient) Send(m *StreamingDatabaseRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *transactionalKeyValueStoreExecuteTransactionClient) Recv() (*StreamingDatabaseResponse, error) {
	m := new(StreamingDatabaseResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *transactionalKeyValueStoreClient) Execute(ctx context.Context, in *DatabaseRequest, opts ...grpc.CallOption) (*DatabaseResponse, error) {
	out := new(DatabaseResponse)
	err := c.cc.Invoke(ctx, "/TransactionalKeyValueStore/execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionalKeyValueStoreServer is the server API for TransactionalKeyValueStore service.
// All implementations must embed UnimplementedTransactionalKeyValueStoreServer
// for forward compatibility
type TransactionalKeyValueStoreServer interface {
	// Execute a transaction
	ExecuteTransaction(TransactionalKeyValueStore_ExecuteTransactionServer) error
	// Execute a single database operation (with a transaction)
	Execute(context.Context, *DatabaseRequest) (*DatabaseResponse, error)
	mustEmbedUnimplementedTransactionalKeyValueStoreServer()
}

// UnimplementedTransactionalKeyValueStoreServer must be embedded to have forward compatible implementations.
type UnimplementedTransactionalKeyValueStoreServer struct {
}

func (UnimplementedTransactionalKeyValueStoreServer) ExecuteTransaction(TransactionalKeyValueStore_ExecuteTransactionServer) error {
	return status.Errorf(codes.Unimplemented, "method ExecuteTransaction not implemented")
}
func (UnimplementedTransactionalKeyValueStoreServer) Execute(context.Context, *DatabaseRequest) (*DatabaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedTransactionalKeyValueStoreServer) mustEmbedUnimplementedTransactionalKeyValueStoreServer() {
}

// UnsafeTransactionalKeyValueStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionalKeyValueStoreServer will
// result in compilation errors.
type UnsafeTransactionalKeyValueStoreServer interface {
	mustEmbedUnimplementedTransactionalKeyValueStoreServer()
}

func RegisterTransactionalKeyValueStoreServer(s grpc.ServiceRegistrar, srv TransactionalKeyValueStoreServer) {
	s.RegisterService(&TransactionalKeyValueStore_ServiceDesc, srv)
}

func _TransactionalKeyValueStore_ExecuteTransaction_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TransactionalKeyValueStoreServer).ExecuteTransaction(&transactionalKeyValueStoreExecuteTransactionServer{stream})
}

type TransactionalKeyValueStore_ExecuteTransactionServer interface {
	Send(*StreamingDatabaseResponse) error
	Recv() (*StreamingDatabaseRequest, error)
	grpc.ServerStream
}

type transactionalKeyValueStoreExecuteTransactionServer struct {
	grpc.ServerStream
}

func (x *transactionalKeyValueStoreExecuteTransactionServer) Send(m *StreamingDatabaseResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *transactionalKeyValueStoreExecuteTransactionServer) Recv() (*StreamingDatabaseRequest, error) {
	m := new(StreamingDatabaseRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _TransactionalKeyValueStore_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionalKeyValueStoreServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/TransactionalKeyValueStore/execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionalKeyValueStoreServer).Execute(ctx, req.(*DatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransactionalKeyValueStore_ServiceDesc is the grpc.ServiceDesc for TransactionalKeyValueStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransactionalKeyValueStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "TransactionalKeyValueStore",
	HandlerType: (*TransactionalKeyValueStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "execute",
			Handler:    _TransactionalKeyValueStore_Execute_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "executeTransaction",
			Handler:       _TransactionalKeyValueStore_ExecuteTransaction_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "lionrock.proto",
}
