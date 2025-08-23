package main

import (
	"flag"
	"grpc-learn/echo"
	"grpc-learn/echo-client/client"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := echo.NewEchoClient(conn)
	//client.CallUnary(c)
	//client.CallServerStream(c)
	//client.CallClientStream(c)
	client.CallDoubleStream(c)
}
