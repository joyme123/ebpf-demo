.PHONY: build,generate

generate:
	go generate pkg/xconnect/xdp.go
	go generate pkg/xdp_pass_and_drop/xdp.go
	go generate pkg/map_counter/xdp.go

build: generate
	go build -o xdp-demo main.go
