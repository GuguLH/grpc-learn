package client

import (
	"log"

	"google.golang.org/grpc/resolver"
)

const (
	MyScheme      = "myscheme"
	MyServiceName = "myecho"
)

var addrs = []string{"localhost:50051", "localhost:50052", "localhost:50053"}

type MyResolverBuilder struct {
}

func (*MyResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &MyResolver{
		target:     target,
		cc:         cc,
		addrsStore: map[string][]string{MyServiceName: addrs},
	}
	return r, nil
}

func (r *MyResolver) Start() {
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

}

func (r *MyResolver) Close() {}
