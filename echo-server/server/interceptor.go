package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	fmt.Println("Server UnaryInterceptor")
	fmt.Println(info)

	err = oauth2Valid(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func StreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	fmt.Println("Server StreamInterceptor")
	fmt.Println(info)

	err := oauth2Valid(ss.Context())
	if err != nil {
		return err
	}

	return handler(srv, ss)
}

func oauth2Valid(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("元数据获取失败,身份验证失败")
	}
	authorization := md.Get("authorization")

	if !valid(authorization) {
		return errors.New("身份令牌校验失败,身份验证失败")
	}
	return nil
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == fetchToken()
}

func fetchToken() string {
	return "some-secret-token"
}
