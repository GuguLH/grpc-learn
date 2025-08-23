package main

import (
	"flag"
	"grpc-learn/echo"
	"grpc-learn/echo-client/client"
	"log"

	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "")
)

func getOptions() (opts []grpc.DialOption) {
	opts = make([]grpc.DialOption, 0)
	opts = append(opts, client.GetMTlsOpt())
	return opts
}

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, getOptions()...)
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
