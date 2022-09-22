#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

#include "map_counter_common.h"

struct bpf_map_def SEC("maps") xdp_stats_map = {
    .type        = BPF_MAP_TYPE_PERCPU_ARRAY,
    .key_size    = sizeof(__u32),
    .value_size  = sizeof(struct datarec),
    .max_entries = XDP_ACTION_MAX,
};

/* LLVM maps __sync_fetch_and_add() as a built-in function tothe BPF atomic add
 * instruction (that is BPF_STX | BPF_XADD | BPF_W for word sizes)
 */
// #ifndef lock_xadd
// #define lock_xadd(ptr, val)    ((void) __sync_fetch_and_add(ptr, val))
// #endif

SEC("xdp_stats1")
int xdp_stats1_func(struct xdp_md *ctx)
{
    // void *data_end = (void *)(long)ctx->data_end;
    // void *data     = (void *)(long)ctx->data;
    struct datarec *rec;
    __u32 key = XDP_PASS;   /* XDP_PASS = 2 */

    /* Lookup in kernel BPF-side return pointer to actual data record */
    rec = bpf_map_lookup_elem(&xdp_stats_map, &key);
    /* BPF kernel-side verifier will reject program if the NULL pointer 
     * check isn't performed here. Even-though this is a static array where
     * we know key lookup XDP_PASS always will succeed.
     */
    if (!rec)
        return XDP_ABORTED;

    /* Multiple CPUs can access data record. Thus, the accounting needs to
     * use an atomic operation.
     */
    // lock_xadd(&rec->rx_packets, 1);
    // lock_xadd(&rec->rx_bytes, ctx->data_end - ctx->data);
    rec->rx_packets++;
    rec->rx_bytes = rec->rx_bytes + (ctx->data_end - ctx->data);

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";

/*
 * user accessible metadata for XDP packet hook
 * new fields must be added to the end of this structure
 * struct xdp_md {
 *	__u32 data;
 * 	__u32 data_end;
 * 	__u32 data_meta;
 * 	// Below access go through struct xdp_rxq_info
 * 	__u32 ingress_ifindex; // rxq->dev->ifindex
 * 	__u32 rx_queue_index;  // rxq->queue_index
 * 
 * 	__u32 egress_ifindex;  // txq->dev->ifindex
 * };
 */
