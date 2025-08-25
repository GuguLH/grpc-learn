package client

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func GetKeepAliveOpt() []grpc.DialOption {
	var kacp = keepalive.ClientParameters{
		// 如果没有活动流,则每10s发送一次ping
		Time: 10 * time.Second,
		// ping超时时长
		Timeout: 1 * time.Second,
		// 当没有任何活动的流的情况下,是否允许被ping
		PermitWithoutStream: true,
	}
	return []grpc.DialOption{grpc.WithKeepaliveParams(kacp)}
}
