// You should run this binary with suid set.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"

	"github.com/fr0stylo/funcgo/pkg/runtime"
)

const (
	bridgeName = "unc0"
	vethPrefix = "uv"
)

var (
	gateway = runtime.DefaultIpManager.Gateway()
)

func createBridge() error {
	// try to get bridge by name, if it already exists then just exit
	_, err := net.InterfaceByName(bridgeName)
	if err == nil {
		return nil
	}
	if !strings.Contains(err.Error(), "no such network interface") {
		return err
	}
	la := netlink.NewLinkAttrs()
	la.Name = bridgeName
	la.TxQLen = -1
	br := &netlink.Bridge{
		LinkAttrs: la,
	}
	if err := netlink.LinkAdd(br); err != nil {
		return fmt.Errorf("bridge creation: %v", err)
	}
	// set up ip addres for bridge
	if err := netlink.AddrAdd(br, gateway); err != nil {
		return fmt.Errorf("add address %v to bridge: %v", gateway, err)
	}
	// sets up bridge ( ip link set dev unc0 up )
	if err := netlink.LinkSetUp(br); err != nil {
		return err
	}
	return nil
}

func createVethPair(pid int) error {
	// get bridge to set as master for one side of veth-pair
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}
	// generate names for interfaces
	x1, x2 := rand.Intn(10000), rand.Intn(10000)
	parentName := fmt.Sprintf("%s%d", vethPrefix, x1)
	peerName := fmt.Sprintf("%s%d", vethPrefix, x2)
	// create *netlink.Veth
	la := netlink.NewLinkAttrs()
	la.Name = parentName
	la.MasterIndex = br.Attrs().Index
	vp := &netlink.Veth{LinkAttrs: la, PeerName: peerName}
	if err := netlink.LinkAdd(vp); err != nil {
		return fmt.Errorf("veth pair creation %s <-> %s: %v", parentName, peerName, err)
	}
	// get peer by name to put it to namespace
	peer, err := netlink.LinkByName(peerName)
	if err != nil {
		return fmt.Errorf("get peer interface: %v", err)
	}

	// put peer side to network namespace of specified PID
	if err := netlink.LinkSetNsPid(peer, pid); err != nil {
		return fmt.Errorf("move peer to ns of %d: %v", pid, err)
	}
	pp, _ := netlink.LinkByName(parentName)

	netlink.LinkSetMaster(pp, br)
	if err := netlink.LinkSetUp(pp); err != nil {
		return err
	}
	if err := netlink.LinkSetUp(vp); err != nil {
		return err
	}
	return netlink.RouteAdd(&netlink.Route{
		LinkIndex: peer.Attrs().Index,
		Scope:     netlink.SCOPE_UNIVERSE,
		Gw:        gateway.IP,
	})
}

func main() {
	pid := 1
	if len(os.Args) > 1 {
		p, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		pid = p
	}
	if err := createBridge(); err != nil {
		log.Fatal(err)
	}
	if err := createVethPair(pid); err != nil {
		log.Fatal(err)
	}
}
