package main

import (
	"flag"
	"fmt"
	"grpc-learn/echo"
	"grpc-learn/echo-server/server"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "")
)

func getOptions() (opts []grpc.ServerOption) {
	opts = make([]grpc.ServerOption, 0)
	opts = append(opts, server.GetMTlsOpt())

	opts = append(opts, grpc.UnaryInterceptor(server.UnaryInterceptor))
	opts = append(opts, grpc.StreamInterceptor(server.StreamInterceptor))
	opts = append(opts, server.GetKeepAliveOpt()...)

	return opts
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(getOptions()...)
	echo.RegisterEchoServer(s, &server.EchoServer{})
	log.Printf("server listening at: %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
