package main

import (
	"flag"
	"fmt"
	"grpc-learn/echo"
	"grpc-learn/echo-client/client"
	"log"
	"time"

	"google.golang.org/grpc"
)

func getOptions() (opts []grpc.DialOption) {
	opts = make([]grpc.DialOption, 0)
	opts = append(opts, client.GetMTlsOpt())
	opts = append(opts, grpc.WithUnaryInterceptor(client.UnaryInterceptor))
	opts = append(opts, grpc.WithStreamInterceptor(client.StreamInterceptor))
	opts = append(opts, client.GetAuth(client.FetchToken()))
	opts = append(opts, client.GetKeepAliveOpt()...)
	opts = append(opts, client.GetNameResolver(client.NewNameServer("localhost:60051")))
	return opts
}

func main() {
	flag.Parse()
	// 根据地址访问
	// conn, err := grpc.NewClient(*addr, getOptions()...)
	// 根据 协议 + 服务名 通过名称解析器,访问服务器
	conn, err := grpc.NewClient(fmt.Sprintf("%s:///%s", client.MyScheme, client.MyServiceName), getOptions()...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := echo.NewEchoClient(conn)
	client.CallUnary(c)
	time.Sleep(5 * time.Second)
	client.CallServerStream(c)
	time.Sleep(5 * time.Second)
	client.CallClientStream(c)
	time.Sleep(5 * time.Second)
	client.CallDoubleStream(c)
}
