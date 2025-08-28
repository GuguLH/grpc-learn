package client

import (
	"context"
	"grpc-learn/name"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const (
	MyScheme      = "myscheme"
	MyServiceName = "myecho"
)

func GetNameResolver(ns *NameServer) grpc.DialOption {
	nameServer = ns
	return grpc.WithResolvers(&MyResolverBuilder{})
}

type MyResolverBuilder struct {
}

func (*MyResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &MyResolver{
		target:     target,
		cc:         cc,
		addrsStore: map[string][]string{MyServiceName: nameServer.getAddressByServiceName(MyServiceName)},
	}
	r.start()
	return r, nil
}

func (r *MyResolver) start() {
	log.Println("Resolver starting...")
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*MyResolverBuilder) Scheme() string {
	return MyScheme
}

type MyResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *MyResolver) ResolveNow(o resolver.ResolveNowOptions) {
	log.Println("Resolver Now")
	log.Println(r.cc)
	r.addrsStore = map[string][]string{MyServiceName: nameServer.getAddressByServiceName(MyServiceName)}
	r.start()
	log.Println(r.cc)
}

func (r *MyResolver) Close() {
	nameServer.Close()
}

var nameServer *NameServer

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

func (ns *NameServer) getAddressByServiceName(serviceName string) []string {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("grpc server close: %v", err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serviceName,
	}
	ret, err := client.GetAddress(context.Background(), in)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	log.Println(ret)
	return ret.Address
}

func (ns *NameServer) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("grpc server close: %v", err)
		}
	}()
	ns.conn.Close()
}
