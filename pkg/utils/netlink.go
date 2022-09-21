package utils

import (
	"github.com/vishvananda/netlink"
)

func LookupLink(intf string) (*netlink.Link, error) {
	link, err := netlink.LinkByName(intf)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func XdpFlags(linkType string) int {
	if linkType == "veth" || linkType == "tuntap" {
		return 2
	}
	return 0
}

func AddXdpToLink(intf string, xdpObjFd int) error {
	link, err := LookupLink(intf)
	if err != nil {
		return err
	}
	err = netlink.LinkSetXdpFdWithFlags(*link, xdpObjFd, XdpFlags((*link).Type()))
	if err != nil {
		return err
	}
	return nil
}

func DelXdpFromLink(intf string) error {
	link, err := LookupLink(intf)
	if err != nil {
		return err
	}
	return netlink.LinkSetXdpFdWithFlags(*link, -1, XdpFlags((*link).Type()))
}
