package server

import (
	"sync"
	"time"
)

type Address struct {
	serviceName string
	addr        string
	expireAt    int64
}

type nameStore struct {
	data         map[string]map[string]*Address
	dataLocker   sync.RWMutex
	expireAtData map[int64]*Address
	expireLocker sync.RWMutex
}

var serviceNameData *nameStore
var expireTs = time.Second * 10

func init() {
	serviceNameData = &nameStore{
		data:         map[string]map[string]*Address{},
		expireAtData: map[int64]*Address{},
	}
	// TODO 异步删除过期的注册信息
}
