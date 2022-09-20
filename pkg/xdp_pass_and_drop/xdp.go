package xdp_pass_and_drop

import "context"

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf xdp ../../ebpf/xdp_pass_and_drop.c -- -I./include -O2 -Wall

type App struct {
	objs *xdpObjects
}

// Launch ...
func (c *App) Launch(ctx context.Context) {

}
