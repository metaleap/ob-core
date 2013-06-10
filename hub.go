package obcore

import (
	"sync"
)

var (
	hubsCache struct {
		sync.Mutex
		m map[string]*Hub
	}
)

func RootHub() *Hub {
	return GetHub("/")
}

func GetHub(path string) *Hub {
	hubsCache.Lock()
	defer hubsCache.Unlock()
	hub, ok := hubsCache.m[path]
	if !ok {
		hubsCache.m[path] = newHub()
	}
	return hub
}

type Hub struct {
	parent *Hub
}

func newHub() (me *Hub) {
	me = &Hub{}
	return
}

func (me *Hub) Parent() (parent *Hub) {
	return nil
}
