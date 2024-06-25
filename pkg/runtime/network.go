package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/vishvananda/netlink"
)

func waitForIface() (netlink.Link, error) {
	log.Info("Starting to wait for network interface")
	start := time.Now()
	for {
		log.Info(".")
		if time.Since(start) > 5*time.Second {
			log.Info("\n")
			return nil, fmt.Errorf("failed to find veth interface in 5 seconds")
		}
		// get list of all interfaces
		lst, err := netlink.LinkList()
		if err != nil {
			log.Info("\n")
			return nil, err
		}
		for _, l := range lst {
			// if we found "veth" interface - it's time to continue setup
			if l.Type() == "veth" {
				log.Info(l.Attrs().Name)
				return l, nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

const suidNet = "unet"

func putIface(pid int) error {
	log.Info("Putting veth interface into container")
	//
	//cmd := exec.Command(suidNet, strconv.Itoa(pid))
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//must("putIface", cmd.Run())
	//-bridgeAddress string
	//Address to assign to bridge device (CIDR notation) (default "10.10.10.1/24")
	//-bridgeName string
	//Name to assign to bridge device (default "brg0")
	//-containerAddress string
	//Address to assign to the container (CIDR notation) (default "10.10.10.2/24")
	//-pid int
	//pid of a process in the container's network namespace
	netsetgoCmd := exec.Command("netsetgo", "-pid", strconv.Itoa(pid))
	if err := netsetgoCmd.Run(); err != nil {
		fmt.Printf("Error running netsetgo - %s\n", err)
		os.Exit(1)
	}

	return nil
}

type Cfg struct {
	IP string
}

func setupIface(link netlink.Link, cfg Cfg) error {
	// up loopback
	lo, err := netlink.LinkByName("lo")
	if err != nil {
		return fmt.Errorf("lo interface: %v", err)
	}
	if err := netlink.LinkSetUp(lo); err != nil {
		return fmt.Errorf("up veth: %v", err)
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("up veth: %v", err)
	}
	addr, err := netlink.ParseAddr(cfg.IP)
	if err != nil {
		return fmt.Errorf("parse IP: %v", err)
	}
	log.Infof("IP:  %s", cfg.IP)
	log.Info(netlink.RouteList(link, netlink.FAMILY_V4))

	return netlink.AddrAdd(link, addr)
}

func SetupNet(ip string) error {
	_, err := waitForIface()
	if err != nil {
		return err
	}
	//if err := setupIface(lnk, Cfg{
	//	IP: ip,
	//}); err != nil {
	//	return err
	//}

	return nil
}
