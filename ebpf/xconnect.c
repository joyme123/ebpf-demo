// apt install linux-source-5.15.0, use tar -jxvf to unpack source code
// linux headers is locate /usr/src/linux-source-5.15.0/linux-source-5.15.0/tools/lib/bpf
// sudo apt install -y gcc-multilib
// sudo apt install libbpf-dev
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

#define XCONNECT_MAP_SIZE 1024

struct bpf_map_def SEC("maps") xconnect_map = {
    .type = BPF_MAP_TYPE_DEVMAP,
    .key_size = sizeof(int),
    .value_size = sizeof(int),
    .max_entries = XCONNECT_MAP_SIZE,
};

// https://stackoverflow.com/questions/67553794/what-is-variable-attribute-sec-means
SEC("xdp")
int xdp_xconnect(struct xdp_md *ctx)
{
    return bpf_redirect_map(&xconnect_map, ctx->ingress_ifindex, 0);
}

char _license[] SEC("license") = "GPL";
