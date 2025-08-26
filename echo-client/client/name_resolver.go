package client

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const (
	MyScheme      = "myscheme"
	MyServiceName = "myecho"
)

var addrs = []string{"localhost:50052", "localhost:50053", "localhost:50051"}

func GetNameResolver() grpc.DialOption {
	return grpc.WithResolvers(&MyResolverBuilder{})
}

type MyResolverBuilder struct {
}

func (*MyResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &MyResolver{
		target:     target,
		cc:         cc,
		addrsStore: map[string][]string{MyServiceName: addrs},
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
	r.addrsStore = map[string][]string{MyServiceName: []string{"localhost:50053", "localhost:50051"}}
	r.start()
	log.Println(r.cc)
}

func (r *MyResolver) Close() {}
