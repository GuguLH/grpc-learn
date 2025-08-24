package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func UnaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("Client UnaryInterceptor")

	var credsConfigured bool
	for _, opt := range opts {
		_, ok := opt.(*grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(GetPerRPCCredentials(FetchToken())))
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func StreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	fmt.Println("Client StreamInterceptor")

	var credsConfigured bool
	for _, opt := range opts {
		_, ok := opt.(*grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(GetPerRPCCredentials(FetchToken())))
	}

	return streamer(ctx, desc, cc, method, opts...)
}
