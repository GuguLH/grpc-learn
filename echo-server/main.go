package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc-learn/echo"
	"grpc-learn/echo-server/server"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &server.EchoServer{})
	log.Printf("server listening at: %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
