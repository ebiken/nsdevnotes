# SONiC

SONiC 関連情報のまとめ

- Running SONiC
  - [on KVM (sonic-vs)](running-sonic-kvm.md)
  - [on Docker (docker-sonic-vs)](running-sonic-docker.md)
  - [on Fixed Function ASIC](running-sonic.md)
  - on Tofino ASIC (TBD: running-sonic-tofino.md)
- [SONiC Switch Image](sonic-image.md)
  - [Getting pre-built image](sonic-image.md#getting-pre-built-image)
  - [Building SONiC image from source code](sonice-image.md)
- SONiC Internal
  - Architecture and modules
  - Source Code Analysis
  - sonic-vs dataplane
- SONiC Release Process
- SAI ... Switch Abstraction Interface
  - [SAI overview](sai.md)

## TODO

- Running SONiC on KVM
- Building SONiC image
- Running SONiC on KVM
  - Build Docker image and upload to docker.io
- 関連技術のメモ
  - Linux ebtables
  - Linux KVM/Virsh
  - container lab

### 追加したい機能（Contributionネタ）

- SONiC on KVM での `show interface counters` サポート

## References

### Documents

- Official Pages
  - Main repo : https://github.com/sonic-net/SONiC
    - `This repository contains documentation, Wiki, master project management, and website for the Software for Open Networking in the Cloud (SONiC).`
  - Wiki: https://github.com/sonic-net/SONiC/wiki
  - Source Code (build Repo) : https://github.com/Azure/sonic-buildimage
  - Roadmap : https://github.com/sonic-net/SONiC/wiki/Sonic-Roadmap-Planning
  - Feature Docs (HDL) : https://github.com/sonic-net/SONiC/tree/master/doc

### Slides / Videos

- [SONiCをはじめてみよう（2019）](https://speakerdeck.com/imasaruoki/sonicwohazimetemiyou)
  - 起動～コマンド利用方法の解説。
- [(PDF) JANOG44: OSSなWhitebox用NOSのSONiCが商用で使われている理由を考える （2019）](https://www.janog.gr.jp/meeting/janog44/application/files/1415/6396/6082/janog44_sonic_kuwata-00.pdf)
  - SONiC仮想マシン４台構成の設定例あり
- [(PDF) JANOG49: SONiCの開発状況アップデート](https://www.janog.gr.jp/meeting/janog49/wp-content/uploads/2022/01/JANOGWeeeeeK%C3%AE%C3%B7%C3%A8J%C3%84%C3%A6%C3%B9%E2%94%90_APRESIA_v.0.1.pdf)
  - https://www.janog.gr.jp/meeting/janog49/sonic/
  - SONiC ロードマップや機能リスト

### Other informative links

- µSONiC (micro-SONiC) used in Optical White Box (TIP Goldstone)
  - "Open whitebox architecture for smart integration of optical networking and data center technology"
    - https://ieeexplore.ieee.org/document/9275288)
    - Page 6 `µSONiC is used in cases where it is necessary to control an Ethernet ASIC. µSONiC is a lightweight package of SONiC, the Microsoft NOS, in which only the components that control SONiC’s Ethernet ASIC are extracted and containerized for easy deployment on Kubernetes.`
  - Two repos under GitHub `oopt-goldstone` organization 
    - https://github.com/oopt-goldstone/usonic_new
    - https://github.com/oopt-goldstone/usonic