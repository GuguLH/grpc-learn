package server

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func getMetadataByMap(m map[string]string) metadata.MD {
	// 通过Map初始化metadata
	md := metadata.New(m)
	return md
}

func getMetadataByKV(kv ...string) metadata.MD {
	// 通过键值对的方式初始化metadata
	md := metadata.Pairs(kv...)
	return md
}

func getOutgoingContext(ctx context.Context, md metadata.MD) context.Context {
	// OutgoingContext 用于请求发送方,包装数据出去
	// IncomingContext 用于请求接收方,用于获取发送方的数据
	// Context 通过序列化成 http2 header 的方式传输
	// new 方法会覆盖ctx原有的数据
	return metadata.NewOutgoingContext(ctx, md)
}

// 将数据附加到OutgoingContext
func appendToOutgoingContext(ctx context.Context, kv ...string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, kv...)
}
