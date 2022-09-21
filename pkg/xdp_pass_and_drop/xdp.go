package xdp_pass_and_drop

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/joyme123/ebpf-demo/pkg/utils"
)

type Program string

const (
	ProgramXDPPass    Program = "xdp_pass"
	ProgramXDPDrop    Program = "xdp_drop"
	ProgramXDPAborted Program = "xdp_aborted"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf xdp ../../ebpf/xdp_pass_and_drop.c -- -I./include -O2 -Wall

type App struct {
	objs *xdpObjects
}

func NewXdpPassAndDropApp() (*App, error) {
	c := &App{
		objs: &xdpObjects{},
	}
	err := loadXdpObjects(c.objs, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Launch load xdp program
// prog should be program name: xdp_pass, xdp_drop, xdp_aborted
func (c *App) Launch(ctx context.Context, intf string, prog string) error {
	log.Info("launch app")
	var progFd int
	switch Program(prog) {
	case ProgramXDPPass:
		progFd = c.objs.XdpPassFunc.FD()
	case ProgramXDPDrop:
		progFd = c.objs.XdpDropFunc.FD()
	case ProgramXDPAborted:
		progFd = c.objs.XdpAbortedFunc.FD()
	}
	err := utils.AddXdpToLink(intf, progFd)
	if err != nil {
		return err
	}
	<-ctx.Done()
	log.Info("receive context done, cleanup")
	return c.cleanup(intf)
}

func (c *App) cleanup(intf string) error {
	return utils.DelXdpFromLink(intf)
}
