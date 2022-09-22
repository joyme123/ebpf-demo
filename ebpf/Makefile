.PHONY: dep,build,all

all: dep build

dep:
	clang -MD -MF headers.d -target bpf -I /usr/src/linux-source-5.15.0/linux-source-5.15.0/tools/lib/bpf -c xconnect.c xdp_pass_and_drop.c map_counter.c

build:
	clang -target bpf -Wall -O2 -emit-llvm -g -Iinclude -c xconnect.c -o - | llc -march=bpf -mcpu=probe -filetype=obj -o xconnect.o
	clang -target bpf -Wall -O2 -emit-llvm -g -Iinclude -c xdp_pass_and_drop.c -o - | llc -march=bpf -mcpu=probe -filetype=obj -o xdp_pass_and_drop.o
	clang -target bpf -Wall -O2 -emit-llvm -g -Iinclude -c map_counter.c -o - | llc -march=bpf -mcpu=probe -filetype=obj -o map_counter.o

