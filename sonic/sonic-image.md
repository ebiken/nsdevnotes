# SONiC Switch Image

SONiC Image ビルド手順のまとめ。
ナレッジメモは [SONiC Build Image Memo](sonic-buildimage-memo.md) を参照。

- [Cheat Sheet](#cheat-sheet)
- [Getting pre-built image](#getting-pre-built-image)
- [Building SONiC image](#building-sonic-image)
- [仮想環境でのビルド手順](#仮想環境でのビルド手順)
  - [Multipass のセットアップ](#multipass-のセットアップ)
  - [sonic-cloud-init.yaml の編集](#sonic-cloud-inityaml-の編集)
  - [仮想マシンの作成＆ビルドスクリプトの実行](#仮想マシンの作成ビルドスクリプトの実行)

## Cheat Sheet

後述 [Building SONiC image](#仮想環境でのビルド手順) を利用した手順のコピペ作業用メモ

```
sudo snap install multipass
screen

cd ~/nsdevnotes/sonic

time multipass launch 20.04 -n sonic-builder -c32 -m32G -d250G --cloud-init examples/sonic-cloud-init.yaml
> real    2m7.088s

time multipass launch 20.04 -n sonic-builder -c8 -m10G -d250G --cloud-init examples/sonic-cloud-init.yaml
> real    1m48.691s

multipass shell sonic-builder
time ./build-broadcom.sh

```

参考：仮想マシンのリソース割り当てによるビルド時間比較
> Host: Intel(R) Xeon(R) Gold 5218 CPU @ 2.30GHz, 2 sockets, 128GB RAM

| SONIC_BUILD_JOBS |  8core, 10GB | 32core, 32GB |                      |
|-----------------:|:------------:|-------------:|:---------------------|
|                4 |  204m52.503s |  210m37.267s | build-broadcom.sh    |
|                8 |  203m31.549s |  177m55.821s | build-broadcom-8.sh  |
|               16 |      n/a     |  176m17.153s | build-broadcom-16.sh |


## Getting pre-built image

> 自分で改造しない場合は、ビルド済みの image を利用可能

https://sonic-build.azurewebsites.net/ui/sonic/pipelines

- Find `Platform`, `BranchName` and click link on `Builds` (e.g. `vs`, `master`)
- Find the build you want to use and click `Artifacts`
  - make sure the `Result` is `succeeded`
- Click `Name` (e.g. `sonic-buildimage.vs1`)
- Download `*.img.gz` or `*.bin` (e.g. `target/sonic-vs.img.gz`, `target/sonic-broadcom.bin`)

## Building SONiC image

SONiC image は [sonic-buildimage](https://github.com/Azure/sonic-buildimage/) を `clone` し [README.md](https://github.com/Azure/sonic-buildimage/blob/master/README.md) に記載された手順に従うとビルド可能。（`GNU make` を利用）

- 必要な環境（最新情報は README.md を確認）
  - CPU Core：多い方が早い
  - RAM（メモリ）：8GB以上
  - Disk：300GB（200GBでも成功例あり）
  - OS：Ubuntu 20.04（2022/06/17現在）
- 手順
  - 事前準備：pip/jinjaのインストール、Docker環境のセットアップ
  - sonic-buildimage レポジトリの clone
  - コマンド投入

```
sudo apt install -y python3-pip
sudo pip3 install j2cli
> install Docker
sudo gpasswd -a ${USER} docker
git clone https://github.com/Azure/sonic-buildimage.git
sudo modprobe overlay
cd sonic-buildimage
git checkout [branch_name]
make init
make configure PLATFORM=[ASIC_VENDOR]
make SONIC_BUILD_JOBS=4 all
```

## 仮想環境でのビルド手順

SONiCで利用される多数のモジュールを順番にビルドする複雑なステップを踏むため、ビルド環境によって失敗する事を防ぐためにクリーンな仮想環境を利用してビルドする事が推奨される。
また、ビルドは `sonic-slave (slave.mk)` というDockerコンテナを作成しその中で実行されるため、仮想環境としてはコンテナではなく仮想マシンを利用する方法を記載。

### Multipass のセットアップ

仮想マシン作成に Canonical の Multipass を利用した場合の手順を記載。

Multipass は Linux, Windows, macOS で動作する。
それぞれのセットアップ方法は https://multipass.run/  `select OS to get started` をクリックして確認可能。

Linux では `sudo snap install multipass` を実行（Ubuntu 20.04で確認）

### sonic-cloud-init.yaml の編集

Multipass では `cloud-init.yaml` を利用して仮想マシンのビルド後に環境構築を自動化する事が可能なため、[examples/sonic-cloud-init.yaml](examples/sonic-cloud-init.yaml) を利用してビルド環境となる仮想マシンを構築する。

> YAMLファイルは [sonic-builder by Anton Smith](https://github.com/antongisli/sonic-builder) を参考に作成しました。

- sonic-cloud-init.yaml によってビルドスクリプトが用意されるので、有効にしたい機能などに応じて、sonic-cloud-init.yaml を変更
- 特定ブランチや commit のコードをビルドしたい場合は、sonic-cloud-init.yaml の以下箇所を変更

```
> sonic-cloud-init.yaml
     # uncomment below if you want to change the branch
     # cd sonic-buildimage; git checkout 202106; cd ..
```

- 以下箇所でプラットフォーム毎のビルドスクリプトを生成している
- 目的のプラットフォームが無い場合、sonic-cloud-init.yaml の以下部分をコピペ編集して追加可能
- （以下は Docker で動作する SONiC Virtual Switch のビルドスクリプト生成箇所）

> - `PLATFORM=p4` （BMv2） はサポートされてません（2022/06/17 現在）
> - Tofino ASIC 向けビルドは Intel/Barefoot からパッチ＆ドキュメントを入手してください。

```
> sonic-cloud-init.yaml
 - path: /home/ubuntu/build-vs-docker.sh
   permissions: 0744
   owner: root
   content: |
     #!/usr/bin/env bash
     cd sonic-buildimage
     make init
     make configure PLATFORM=vs
     make SONIC_BUILD_JOBS=4 target/docker-sonic-vs.gz

 - chown ubuntu:ubuntu /home/ubuntu/build-vs-docker.sh
```

### 仮想マシンの作成＆ビルドスクリプトの実行

- 仮想マシンの作成＆起動は 'multipass launch' コマンドで実行
- 仮想マシンへ割り当てるリソースは以下オプションで指定（手元の環境に応じて変更）
  - `-c` : CPU Core数
  - `-m` : RAM（メモリ）
  - `-d` : DISKサイズ（300GBが推奨）
- 作成された仮想マシンにログインし、ビルドしたいプラットフォームに応じたスクリプトを実行
- なお、ビルドには２～３時間程度かかるので、リモート環境で実施している場合は `screen` コマンドなどでホストへの接続が切れても再接続可能なようにしておくことを推奨

```
~/nsdevnotes/sonic$ multipass list
No instances found.
~/nsdevnotes/sonic$ screen

~/nsdevnotes/sonic$ time multipass launch 20.04 -n sonic-builder \
-c8 -m10G -d250G --cloud-init examples/sonic-cloud-init.yaml
または
~/nsdevnotes/sonic$ time multipass launch 20.04 -n sonic-builder \
-c32 -m32G -d250G --cloud-init examples/sonic-cloud-init.yaml

~/nsdevnotes/sonic$ time multipass shell sonic-builder

ubuntu@sonic-builder:~$ cat build-broadcom.sh
#!/usr/bin/env bash
cd sonic-buildimage
make init
make configure PLATFORM=broadcom
make SONIC_BUILD_JOBS=4 target/sonic-broadcom.bin

ubuntu@sonic-builder:~$ time ./build-broadcom.sh
...snip...
```
