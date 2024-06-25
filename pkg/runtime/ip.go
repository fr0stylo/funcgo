package runtime

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/vishvananda/netlink"
)

type Ipmanager struct {
	lock     sync.RWMutex
	template string
	iptable  []bool
}

func (r *Ipmanager) acquire(i int) *netlink.Addr {
	r.lock.RLock()
	defer r.lock.RUnlock()
	r.iptable[i] = true
	ip, _ := netlink.ParseAddr(fmt.Sprintf(r.template, i))
	return ip
}

func (r *Ipmanager) Acquire() *netlink.Addr {
	for {
		for i, v := range r.iptable {
			if !v {
				return r.acquire(i)
			}
		}
	}
}

func (r *Ipmanager) Release(ip *netlink.Addr) error {
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

func (r *Ipmanager) Base() *netlink.Addr {
	ip, _ := netlink.ParseAddr(fmt.Sprintf(r.template, 1))
	return ip
}

func (r *Ipmanager) Gateway() *netlink.Addr {
	ip, _ := netlink.ParseAddr(fmt.Sprintf(r.template, 1))
	return ip
}

type IPManager interface {
	Acquire() *netlink.Addr
	Release(*netlink.Addr) error
	Base() *netlink.Addr
	Gateway() *netlink.Addr
}

func NewIPManager(template string) *Ipmanager {
	mngr := &Ipmanager{
		lock:     sync.RWMutex{},
		template: template,
		iptable:  make([]bool, 256),
	}

	mngr.acquire(0)
	mngr.acquire(1)
	mngr.acquire(255)

	return mngr
}

var DefaultIpManager = NewIPManager("10.10.10.%d/24")
