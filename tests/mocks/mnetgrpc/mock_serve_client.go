package mnetgrpc

import (
	"context"
	pb "github.com/cpapidas/pegasus/network/netgrpc/proto"
	"google.golang.org/grpc"
)

// MockServeClient mock object for proto.ClientConnection
type MockServeClient struct {
	HandlerSyncMock func(
		ctx context.Context,
		in *pb.HandlerRequest,
		opts ...grpc.CallOption,
	) (*pb.HandlerReply, error)
	HandlerMock func(ctx context.Context, opts ...grpc.CallOption) (pb.Serve_HandlerClient, error)
}

// HandlerSync mock object for proto.ClientConnection HandlerSync function
func (m MockServeClient) HandlerSync(
	ctx context.Context,
	in *pb.HandlerRequest,
	opts ...grpc.CallOption,
) (*pb.HandlerReply, error) {
	if m.HandlerSyncMock != nil {
		return m.HandlerSyncMock(ctx, in, opts...)
	}
	return nil, nil
}

// Handler mock object for proto.ClientConnection Handler function
func (m MockServeClient) Handler(ctx context.Context, opts ...grpc.CallOption) (pb.Serve_HandlerClient, error) {
	if m.HandlerMock != nil {
		return m.HandlerMock(ctx, opts...)
	}
	return nil, nil
}