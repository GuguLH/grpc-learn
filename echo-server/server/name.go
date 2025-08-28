package server

import (
	"context"
	"grpc-learn/name"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NameServer struct {
	conn *grpc.ClientConn
}

func NewNameServer(addr string) *NameServer {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("grpc server close: %v", err)
		}
	}()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	return &NameServer{conn: conn}
}

func (ns *NameServer) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("grpc server close: %v", err)
		}
	}()
	ns.conn.Close()
}

func (ns *NameServer) RegisterName(serverName, addr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("grpc server close: %v", err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serverName,
		Address:     []string{addr},
	}
	_, err := client.Register(context.Background(), in)
	if err != nil {
		log.Println(err)
	}
}

func (ns *NameServer) Delete(serverName, addr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("grpc server close: %v", err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serverName,
		Address:     []string{addr},
	}
	_, err := client.Delete(context.Background(), in)
	if err != nil {
		log.Println(err)
	}
}

func (ns *NameServer) Keepalive(serverName, addr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("grpc server close: %v", err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serverName,
		Address:     []string{addr},
	}
	stream, err := client.Keepalive(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	for {
		stream.Send(in)
		time.Sleep(time.Second)
	}
}
