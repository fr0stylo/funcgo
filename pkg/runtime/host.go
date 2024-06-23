package runtime

import (
	"fmt"
	"sync"
)

type Host struct {
	functions map[string]*Function
	m         sync.Mutex
}

func NewHost() *Host {
	return &Host{
		functions: make(map[string]*Function),
		m:         sync.Mutex{},
	}
}

type HostRequest struct {
	FunctionName string
	Request      any
}

func (r *Host) Execute(obj any) ([]byte, error) {
	b, ok := obj.(HostRequest)
	if !ok {
		return nil, fmt.Errorf("host execute requires type of HostRequest but was %t", obj)
	}

	f := r.functions[b.FunctionName]
	if f == nil {
		return nil, fmt.Errorf("function not found")
	}

	return f.Execute(b.Request)
}

func (r *Host) InsertFunction(name string, f *Function) error {
	r.m.Lock()
	defer r.m.Unlock()

	r.functions[name] = f

	return nil
}
