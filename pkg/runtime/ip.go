package runtime

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/vishvananda/netlink"
)

type ipmanager struct {
	lock     sync.RWMutex
	template string
	iptable  []bool
}

func (r *ipmanager) acquire(i int) *netlink.Addr {
	r.lock.RLock()
	defer r.lock.RUnlock()
	r.iptable[i] = true
	ip, _ := netlink.ParseAddr(fmt.Sprintf(r.template, i))
	return ip
}

func (r *ipmanager) Acquire() *netlink.Addr {
	for {
		for i, v := range r.iptable {
			if !v {
				return r.acquire(i)
			}
		}
	}
}

func (r *ipmanager) Release(ip *netlink.Addr) error {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if ip == nil {
		return fmt.Errorf("no ip provided")
	}
	list := strings.Split(ip.String(), ".")
	l := strings.Split(list[len(list)-1], "/")[0]

	i, err := strconv.ParseInt(l, 10, 64)
	if err != nil {
		log.Fatal(err)
		return err
	}

	r.iptable[int(i)] = false

	return nil
}

func (r *ipmanager) Base() *netlink.Addr {
	ip, _ := netlink.ParseAddr(fmt.Sprintf(r.template, 1))
	return ip
}

func (r *ipmanager) Gateway() *netlink.Addr {
	ip, _ := netlink.ParseAddr(fmt.Sprintf(r.template, 1))
	return ip
}

type IPManager interface {
	Acquire() *netlink.Addr
	Release(*netlink.Addr) error
	Base() *netlink.Addr
	Gateway() *netlink.Addr
}

func NewIPManager(template string) IPManager {
	mngr := &ipmanager{
		lock:     sync.RWMutex{},
		template: template,
		iptable:  make([]bool, 256),
	}

	mngr.acquire(0)
	mngr.acquire(1)
	mngr.acquire(255)

	return mngr
}

var defaultIPManager = NewIPManager("168.0.0.%d/24")
