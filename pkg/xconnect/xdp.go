package xconnect

import (
	"context"
	"log"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf xdp ../../ebpf/xconnect.c -- -I./include -O2 -Wall

// App stores ebpf programs and maps together with the desired state
type App struct {
	objs    *xdpObjects
	input   map[string]string
	linkMap map[string]*netlink.Link
}

func NewXconnectApp(input map[string]string) (*App, error) {
	c := &App{
		objs:    &xdpObjects{},
		input:   make(map[string]string),
		linkMap: make(map[string]*netlink.Link),
	}

	c.input = makeSymm(input)

	if err := increaseResourceLimits(); err != nil {
		return nil, err
	}
	/*

		specs, err := newXdpSpecs()
		if err != nil {
			return nil, err
		}

		objs, err := specs.Load(nil)
		if err != nil {
			return nil, fmt.Errorf("Can't load objects:%s", err)
		}
	*/
	err := loadXdpObjects(c.objs, nil)
	if err != nil {
		return nil, err
	}

	if err := c.init(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *App) init() error {
	var added []string
	for intf := range c.input {
		added = append(added, intf)
	}

	err := c.updateLinkMap(added)
	if err != nil {
		return err
	}

	return c.updateBpfMap(added, []string{}, []string{})
}

// cleanup clears netlink XDP configuration and closes eBPF objects
func (c *App) cleanup() error {
	var errs error

	var removed []string
	for intf := range c.linkMap {
		removed = append(removed, intf)
	}

	if err := c.delXdpFromLink(removed); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := c.objs.Close(); err != nil {
		errs = multierror.Append(errs, err)
	}

	c.cleanupLinkMap(removed)

	return errs
}

// update ensures running state matches the candidate
func (c *App) update(candidates map[string]string) error {
	candidates = makeSymm(candidates)

	added, changed, orphaned := confDiff(c.input, candidates)

	// Dealing with added interfaces
	err := c.updateLinkMap(added)
	if err != nil {
		return err
	}
	if err := c.addXdpToLink(added); err != nil {
		return err
	}

	// Updating eBPF map
	c.input = candidates
	err = c.updateBpfMap(added, changed, orphaned)
	if err != nil {
		return err
	}

	// Dealing with removed interfaces
	if err := c.delXdpFromLink(orphaned); err != nil {
		return err
	}

	c.cleanupLinkMap(orphaned)
	return nil
}

// updateBpfMap adjusts Bpf Map based on detected changes
func (c *App) updateBpfMap(added, changed, removed []string) error {
	var errs error

	for _, intf := range added {
		link1 := c.linkMap[intf]
		link2 := c.linkMap[c.input[intf]]
		if err := c.objs.XconnectMap.Put(uint32((*link1).Attrs().Index), uint32((*link2).Attrs().Index)); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	for _, intf := range changed {
		link1 := c.linkMap[intf]
		link2 := c.linkMap[c.input[intf]]
		if err := c.objs.XconnectMap.Put(uint32((*link1).Attrs().Index), uint32((*link2).Attrs().Index)); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	for _, intf := range removed {
		link1 := c.linkMap[intf]
		if err := c.objs.XconnectMap.Delete(uint32((*link1).Attrs().Index)); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

// Launch app, watch for changes and perform "warn" reloads.
// This funtion blocks forever and context can be used to gracefully stop it.
// updateCh expects a map between interfaces, similar to input of NewXconnectApp.
func (c *App) Launch(ctx context.Context, updateCh chan map[string]string) {
	var links []string
	for link := range c.linkMap {
		links = append(links, link)
	}
	if err := c.addXdpToLink(links); err != nil {
		log.Fatalf("Failed to set up XDP on links: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("ctx.Done")
			if err := c.cleanup(); err != nil {
				log.Fatal("Cleanup Failed: %s", err)
			}
			return
		case config := <-updateCh:
			if err := c.update(config); err != nil {
				log.Printf("Error updating eBPF: %s", err)
			}
		}
	}
}

func increaseResourceLimits() error {
	return unix.Setrlimit(unix.RLIMIT_MEMLOCK, &unix.Rlimit{
		Cur: unix.RLIM_INFINITY,
		Max: unix.RLIM_INFINITY,
	})
}

func makeSymm(inMap map[string]string) map[string]string {
	res := make(map[string]string)

	var keys []string
	for k := range inMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := inMap[k]
		_, keyFound := res[k]
		_, valueFound := res[v]
		if !keyFound && !valueFound {
			res[k] = v
			res[v] = k
		}
	}
	return res
}

// confDiff compares the running and candidate configurations
// and returns any new, changed or removed interface names
func confDiff(running, candidates map[string]string) ([]string, []string, []string) {
	var new, changed, orphaned []string
	for c1, c2 := range candidates {
		p2, ok := running[c1]
		if !ok {
			new = append(new, c1)
		} else if p2 != c2 {
			changed = append(changed, c1)
		}
	}

	for p1 := range running {
		_, ok := candidates[p1]
		if !ok {
			orphaned = append(orphaned, p1)
		}
	}

	return new, changed, orphaned
}
