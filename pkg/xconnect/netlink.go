package xconnect

import (
	"github.com/hashicorp/go-multierror"
	"github.com/vishvananda/netlink"
)

func lookupLink(intf string) (*netlink.Link, error) {
	link, err := netlink.LinkByName(intf)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func xdpFlags(linkType string) int {
	if linkType == "veth" || linkType == "tuntap" {
		return 2
	}
	return 0
}

func (c *App) updateLinkMap(intfs []string) error {
	var errs error

	for _, intf := range intfs {
		link, err := lookupLink(intf)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		c.linkMap[intf] = link
	}

	return errs
}

func (c *App) cleanupLinkMap(intfs []string) {
	for _, intf := range intfs {
		delete(c.linkMap, intf)
	}
}

func (c *App) addXdpToLink(intfs []string) error {
	var errs error
	for _, intf := range intfs {
		link := c.linkMap[intf]
		err := netlink.LinkSetXdpFdWithFlags(*link, c.objs.XdpXconnect.FD(), xdpFlags((*link).Type()))
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (c *App) delXdpFromLink(intfs []string) error {
	var errs error
	for _, intf := range intfs {
		link := c.linkMap[intf]
		err := netlink.LinkSetXdpFdWithFlags(*link, -1, xdpFlags((*link).Type()))
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}
