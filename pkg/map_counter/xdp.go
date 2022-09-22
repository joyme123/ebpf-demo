package map_counter

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/joyme123/ebpf-demo/constants"
	"github.com/joyme123/ebpf-demo/pkg/utils"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf xdp ../../ebpf/map_counter.c -- -I./include -O2 -Wall

type App struct {
	objs *xdpObjects
}

func NewMapCounterApp() (*App, error) {
	c := &App{
		objs: &xdpObjects{},
	}

	err := loadXdpObjects(c.objs, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Launch loads xdp program
func (c *App) Launch(ctx context.Context, intf string, printPeriodSeconds int) error {
	log.Info("launch app")

	progFd := c.objs.XdpStats1Func.FD()
	err := utils.AddXdpToLink(intf, progFd)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(printPeriodSeconds) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Info("receive context done, cleanup")
			return c.cleanup(intf)
		case <-ticker.C:
			c.printStats()
		}
	}
}

func (c *App) cleanup(intf string) error {
	return utils.DelXdpFromLink(intf)
}

type Record struct {
	RxPackets uint64
	RxBytes   uint64
}

func (c *App) printStats() {
	records := []*Record{}
	if err := c.objs.xdpMaps.XdpStatsMap.Lookup(uint32(constants.XDP_PASS), &records); err != nil {
		log.Errorf("lookup from bpf map failed: %v", err)
	}

	total := Record{}
	for i, record := range records {
		log.Infof("cpu %d received packets: %d", i, record.RxPackets)
		log.Infof("cpu %d received bytes: %d", i, record.RxBytes)
		total.RxBytes = total.RxBytes + record.RxBytes
		total.RxPackets = total.RxPackets + record.RxPackets
	}
	log.Infof("total received packets: %d", total.RxPackets)
	log.Infof("total received bytes: %d", total.RxBytes)
}
