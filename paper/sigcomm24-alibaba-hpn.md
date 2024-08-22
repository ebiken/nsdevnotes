# SIGCOMM'22 Alibaba HPN: A Data Center Network for Large Language Model Training

- link: https://dl.acm.org/doi/10.1145/3651890.3672265

## Summary

TODO

## SONiC AI Working Group - 2024/08/20

- Topology
  - leaf spine (2 layer)
  - 8GPUs per server
  - 1024 servers per segment
  - 15 segments per Pod
  - 15K GPUs per pod. This was from avilable power in the building.
- ACCL: Alibaba xCCL
  - Add entolopy per queue pair. Can do with customized NCCL.
- Using BF3
  - CX7 does not support Adaptive Routing (AR)
  - Price difference was not much so picked BF3
  - (not sure if this is regular BF3 or SuperNIC)
- 2 NIC per GPU.
  - Custom Linux Kernel to send ARP req/responce to both ToRs.
- Alternative Marking DSCP (A.M.D)

