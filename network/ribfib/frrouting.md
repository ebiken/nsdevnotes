# Rib/Fib notes: FRRouting


```
[[Routingd]] -> ((ROUTEs)) -> [RIB](ZEBRA)[FIB] -> [dataplane]
```

> In FRR, the Routing Information Base (RIB) resides inside zebra. Routing protocols communicate their best routes to zebra, and zebra computes the best route across protocols for each prefix. This latter information makes up the Forwarding Information Base (FIB). Zebra feeds the FIB to the kernel.


## TODO

- Rib/Fib それぞれの構造体を抜き出す
- Rib -> Fib の流れを調べる
- 

## Reference

- [An abstract workflow for BGP implementations October 7, 2020](https://pluginized-protocols.org/xbgp/2020/10/07/xbgp.html)

## Related Source Code

- https://github.com/FRRouting/frr/blob/master/zebra/rib.h

## memo

### Cisco RIB, CEF and FIB

[RIBs and FIBs (aka IP Routing Table and CEF Table)](https://blog.ipspace.net/2010/09/ribs-and-fibs.html), updated on Dec 26, 2020 (original posted on Sep 2010)

- Ideally, we would use RIB to forward IP packets, but we can’t as some entries in it (static routes and BGP routes) could have next hops that are not directly connected.
- BGP route has no outgoing interface and its next hop is not directly connected; the router has to perform recursive lookups to find the outgoing interface
- Forwarding Information Base (FIB) and Cisco Express Forwarding (CEF) were introduced to make layer-3 switching performance consistent. When IP routes are copied from RIB to FIB, their next hops are resolved, outgoing interfaces are computed and multiple entries are created when the next-hop resolution results in multiple paths to the same destination.
- For example, when the BGP route from the previous printout is inserted into FIB, its next-hop is changed to point to the actual next-hop router. The information about the recursive next-hop is retained, as it allows the router to update the FIB (CEF table) without rescanning and recomputing the whole RIB if the path toward the BGP next-hop changes.

```
RR#show ip cef 10.0.11.11 detail
10.0.11.11/32, epoch 0, flags rib only nolabel, rib defined all labels
  recursive via 10.0.1.1
    nexthop 10.0.2.1 FastEthernet0/0 label 19
```

### Code reading


```c
/* Nexthop structure. */
struct rnh {
  uint8_t flags;
#define ZEBRA_NHT_CONNECTED 0x1
#define ZEBRA_NHT_DELETED 0x2
#define ZEBRA_NHT_RESOLVE_VIA_DEFAULT 0x4
	/* VRF identifier. */
	vrf_id_t vrf_id;
	afi_t afi;
	safi_t safi;
	uint32_t seqno;
	struct route_entry *state;
	struct prefix resolved_route;
	struct list *client_list;
	/* pseudowires dependent on this nh */
	struct list *zebra_pseudowire_list;
	struct route_node *node;
	/*
	 * if this has been filtered for the client
	 */
	int filtered[ZEBRA_ROUTE_MAX];
	struct rnh_list_item rnh_list_item;
}

struct route_entry {}

/*
 * Structure that represents a single destination (prefix).
 */
typedef struct rib_dest_t_ { }

/*
 * rib_table_info_t
 *
 * Structure that is hung off of a route_table that holds information about
 * the table.
 */
struct rib_table_info {
	struct zebra_vrf *zvrf;
	afi_t afi;
	safi_t safi;
	uint32_t table_id;
}



```