.PHONY: build,generate

generate:
	go generate pkg/xconnect/xdp.go

build: generate
	go build -o xdp-demo main.go
