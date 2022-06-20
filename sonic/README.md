# SONiC

SONiC 関連情報のまとめ

- Running SONiC
  - [on KVM (sonic-vs)](running-sonic-kvm.md)
  - [on Docker (docker-sonic-vs)](running-sonic-docker.md)
  - on Fixed Function ASIC(TBD: running-sonic.md) ⇒ 機材入手したら作成
  - on Tofino ASIC (TBD: running-sonic-tofino.md) ⇒ 公開できない？
- [SONiC Switch Image](sonic-image.md)
  - [Getting pre-built image](sonic-image.md#getting-pre-built-image)
  - [Building SONiC image from source code](sonic-image.md)
  - [SONiC Build Image Memo](sonic-buildimage-memo.md) ビルド関連情報のメモ集約
- SONiC Internal
  - Architecture and modules
  - Source Code Analysis
  - sonic-vs dataplane
- SONiC Release Process
- SAI ... Switch Abstraction Interface
  - [SAI overview](sai.md)

## TODO

### 開発効率化
- Running SONiC on Docker： clab を利用した複数台Fabricのサンプル手順＆スクリプト作成
- sonic-build: 各モジュールだけをビルドする方法（変更加えた部分だけをビルドし、現在の3時間⇒10分程度で開発を回せるように）
- sonic-builder の `examples/sonic-cloud-init.yaml` で PLATFORM 毎のビルドスクリプトを生成しているが、機能の有効無効を含めてもっと便利な方法を検討
  - e.g. 機能毎の yaml ファイルを用意しておく、Target名は `./build-sonic.sh <target>` のように引数で指定する、など。
  - 但し、スクリプトを作り込むと sonic-buildimage が変わった時に追従できないため、手作業の方が分かりやすい場合は手作業のままにしておく。


### 追加したい機能（Contributionネタ）

- SONiC on KVM での `show interface counters` サポート

## SONiC設定サンプル

SONiC設定サンプル集。仮想環境の場合は仮想マシンやコンテナ関連の構築スクリプトを載せている場合があります。
[Edgecore SONiC のサポートサイト](https://support.edge-core.com/hc/en-us/categories/360002134713-Edgecore-SONiC) が充実しているので、本家の前にこちらを参照するのも良い。（OSS版ではサポートされていない機能や動作が異なる場合もあるので注意）


- [demo01: Layer 2/3 with VLAN (type=bridge) on KVM](running-sonic-kvm.md#demo01-layer-23-with-vlan-typebridge)
  - libvirt domain 設定（sonic.xml） で `<interface type='bridge'>` を利用したサンプル
  - ホスト側に bridge & veth pair を作成する必要があるため、netns をホストと見立てたテストには煩雑
  - 逆に、スイッチやルーターの仮想インスタンスを接続する場合には有用な方式
- [demo02: Layer 2/3 with VLAN (type=network) on KVM](running-sonic-kvm.md#demo02-layer-23-with-vlan-typenetwork)
  - libvirt domain 設定（sonic.xml） で `<interface type='network'>` を利用したサンプル
  - ホスト側に bridge & veth pair 設定が不要なため、netns をホストと見立てたテストを簡単に実施可能


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