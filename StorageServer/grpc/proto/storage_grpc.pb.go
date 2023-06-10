// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: storage.proto

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

const (
	StorageService_StoreFile_FullMethodName    = "/storage.StorageService/StoreFile"
	StorageService_RetrieveFile_FullMethodName = "/storage.StorageService/RetrieveFile"
)

// StorageServiceClient is the client API for StorageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageServiceClient interface {
	StoreFile(ctx context.Context, in *FileChunk, opts ...grpc.CallOption) (*StoreReply, error)
	RetrieveFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (StorageService_RetrieveFileClient, error)
}

type storageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageServiceClient(cc grpc.ClientConnInterface) StorageServiceClient {
	return &storageServiceClient{cc}
}

func (c *storageServiceClient) StoreFile(ctx context.Context, in *FileChunk, opts ...grpc.CallOption) (*StoreReply, error) {
	out := new(StoreReply)
	err := c.cc.Invoke(ctx, StorageService_StoreFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) RetrieveFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (StorageService_RetrieveFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[0], StorageService_RetrieveFile_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceRetrieveFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StorageService_RetrieveFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type storageServiceRetrieveFileClient struct {
	grpc.ClientStream
}

func (x *storageServiceRetrieveFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StorageServiceServer is the server API for StorageService service.
// All implementations must embed UnimplementedStorageServiceServer
// for forward compatibility
type StorageServiceServer interface {
	StoreFile(context.Context, *FileChunk) (*StoreReply, error)
	RetrieveFile(*FileRequest, StorageService_RetrieveFileServer) error
	mustEmbedUnimplementedStorageServiceServer()
}

// UnimplementedStorageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServiceServer struct {
}

func (UnimplementedStorageServiceServer) StoreFile(context.Context, *FileChunk) (*StoreReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreFile not implemented")
}
func (UnimplementedStorageServiceServer) RetrieveFile(*FileRequest, StorageService_RetrieveFileServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveFile not implemented")
}
func (UnimplementedStorageServiceServer) mustEmbedUnimplementedStorageServiceServer() {}

// UnsafeStorageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServiceServer will
// result in compilation errors.
type UnsafeStorageServiceServer interface {
	mustEmbedUnimplementedStorageServiceServer()
}

func RegisterStorageServiceServer(s grpc.ServiceRegistrar, srv StorageServiceServer) {
	s.RegisterService(&StorageService_ServiceDesc, srv)
}

func _StorageService_StoreFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileChunk)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).StoreFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StorageService_StoreFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).StoreFile(ctx, req.(*FileChunk))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_RetrieveFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServiceServer).RetrieveFile(m, &storageServiceRetrieveFileServer{stream})
}

type StorageService_RetrieveFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type storageServiceRetrieveFileServer struct {
	grpc.ServerStream
}

func (x *storageServiceRetrieveFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

// StorageService_ServiceDesc is the grpc.ServiceDesc for StorageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StorageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storage.StorageService",
	HandlerType: (*StorageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreFile",
			Handler:    _StorageService_StoreFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RetrieveFile",
			Handler:       _StorageService_RetrieveFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "storage.proto",
}
