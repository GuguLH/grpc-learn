package main

import (
	"flag"
	"fmt"
	"grpc-learn/name"
	"grpc-learn/name-server/server"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 60051, "server port")
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	name.RegisterNameServer(s, &server.NameServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	//testData()
}

func testData() {
	server.Register("echo", "localhost:50051")
	server.Register("echo", "localhost:50052")
	time.Sleep(2 * time.Second)
	server.Register("echo", "localhost:50053")
	server.Register("echo", "localhost:50054")
	time.Sleep(2 * time.Second)
	server.Register("echo", "localhost:50055")
	server.Register("echo", "localhost:50056")
	time.Sleep(2 * time.Second)
	server.Register("echo", "localhost:50051")
	server.Register("echo", "localhost:50052")
	time.Sleep(2 * time.Second)
	server.Register("echo", "localhost:50053")
	server.Register("echo", "localhost:50054")
	time.Sleep(2 * time.Second)
	server.Register("echo", "localhost:50055")
	server.Register("echo", "localhost:50056")

	allData := server.GetAllData()
	fmt.Println(allData)

	server.Delete("echo", "localhost:50056")
	fmt.Println(server.GetByServiceName("echo"))
	allData = server.GetAllData()
	fmt.Println(allData)
}
