# MISC NOTES

> This is temporary notes which should be sorted out in future.

- [Switch ASIC Buffer Size !! Un-official information !!](https://people.ucsc.edu/~warner/buffer.html)
- "Azure is using Arista Switch with Tofino running VxLAN" mentioned on public website: [https://mnkcg.com/products/p4-ansible/](https://mnkcg.com/products/p4-ansible/)
  - An Arista network switch has shipped to Microsoft Azure cloud with support for 256K VTEP. This switch has a Tofino asic using all stages. There is no room left on the asic for any incremental P4 program merging. This switch requires a forklift upgrade to remove some existing features and add new ones. P4-Ansible has automated the forklift upgrade once a user provides what feature(s) to remove and add.
