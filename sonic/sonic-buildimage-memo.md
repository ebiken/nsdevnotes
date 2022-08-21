# SONiC Build Image Memo

Memo of how SONiC Build Image works.

https://github.com/Azure/sonic-buildimage/

Table of Contents
- [Reference](#reference)
- [これは変更しておけ、というビルド設定](#これは変更しておけというビルド設定)
- [ビルド高速化方法（試行錯誤中）](#ビルド高速化方法試行錯誤中)
  - [並列度アップ](#並列度アップ)
  - [利用しない submodule init の省略](#利用しない-submodule-init-の省略)
  - [apt キャッシュサーバの活用 apt-cacher-ng](#apt-キャッシュサーバの活用-apt-cacher-ng)
  - [Avoid copying files of all targets, which are not asked to build (TBD)](#avoid-copying-files-of-all-targets-which-are-not-asked-to-build-tbd)
- [ビルド設定：rules/config](#ビルド設定rulesconfig)
  - [有効にする機能の設定場所](#有効にする機能の設定場所)
  - [デフォルト USERNAME PASSWORD の変更](#デフォルト-username-password-の変更)
  - [DEBUG 設定](#debug-設定)
- [各モジュールで利用されている Debian version 確認方法](#各モジュールで利用されている-debian-version-確認方法)


## Reference

Official Doc/Blog

- [SONiC Buildimage Guide](https://github.com/Azure/sonic-buildimage/blob/master/README.buildsystem.md) from GitHub: Azure / sonic-buildimage

Other Doc/Blog

- [Tim's Blog
  - [Hack - How SONiC Builds? (2022/02/12)](https://timmy00274672.wordpress.com/2020/02/12/hack-how-sonic-builds/)
  - [SONiC – How are SONIC_DPKG_DEBS built? (2022/02/13)](https://timmy00274672.wordpress.com/2020/02/13/sonic-how-are-sonic_dpkg_debs-built/)
  - [SONiC – How are SONIC_MAKE_DEBS built? (2022/02/13)](https://timmy00274672.wordpress.com/2020/02/13/sonic-how-are-sonic_make_debs-built/)


## これは変更しておけ、というビルド設定

```
DEFAULT_BUILD_LOG_TIMESTAMP = none # simple, none
```

## ビルド高速化方法（試行錯誤中）

### 並列度アップ

[rules/config] の `SONIC_CONFIG_BUILD_JOBS` `SONIC_CONFIG_MAKE_JOBS` を変更する事で並列度アップしてビルド高速化可能

`SONIC_CONFIG_BUILD_JOBS` の最適値はビルド環境（サーバ性能等）依存と考えられる。
SONIC_BUILD_JOBS=4 を 8, 16 と増やした時のビルド時間は [sonic-image.md：参考：仮想マシンのリソース割り当てによるビルド時間比較](sonic-image.md#cheat-sheet) を参照

```
# SONIC_CONFIG_BUILD_JOBS - set number of jobs for parallel build.
# Corresponding -j argument will be passed to make command inside docker
# container.
SONIC_CONFIG_BUILD_JOBS = 1
```

`SONIC_CONFIG_MAKE_JOBS` は core 数に応じて自動設定されるため変更不要。

```
# SONIC_CONFIG_MAKE_JOBS - set number of parallel make jobs per package.
# Corresponding -j argument will be passed to make/dpkg commands that build separate packages
SONIC_CONFIG_MAKE_JOBS = $(shell nproc)

> Example:
ubuntu@sonic-builder:~/sonic-buildimage$ lscpu | grep "CPU(s)"
CPU(s):                          8
On-line CPU(s) list:             0-7
NUMA node0 CPU(s):               0-7
ubuntu@sonic-builder:~/sonic-buildimage$ nproc
8
```

### 利用しない submodule init の省略

`git submodule update --init --recursive` を実施する際、ビルドしたい Target とは無関係にすべての platform が init される。
init は 4min31sec 程度かかるため、利用しない platform を省略することで、数十秒短縮できることが期待される。

```
2022/06/30 Vanilla
time git submodule update --init --recursive
> real    4m38.462s

2022/06/30 reduced submodules
# Delete the relevant lines (platform/xxx) from `.gitmodules` and `.git/config`
#Run git rm --cached path_to_submodule (no trailing slash).
git rm --cached -r platform/marvell
git rm --cached -r platform/nephos
git rm --cached -r platform/centec
git rm --cached -r platform/centec-arm64
git rm --cached -r platform/cavium
git rm --cached -r platform/marvell-arm64
git rm --cached -r platform/marvell-armhf
git rm --cached -r platform/innovium
git rm --cached -r platform/barefoot/sonic-platform-modules-arista
rm -rf platform/barefoot/sonic-platform-modules-arista
git add .gitmodules
time git submodule update --init --recursive
> real    0m58.321s
```

### apt キャッシュサーバの活用 apt-cacher-ng

（APRESIA 桑田さんからの情報）

SONiCビルド時に、何度もapt installが発生するので、apt-cacher-ng（aptのキャッシュサーバ）を立てて、ビルドしていました。ただ、SONiCのビルドスクリプトが改善された覚えがありますので、apt キャッシュサーバを立てても、今ではビルド時間の短縮の効果は小さいかもしれません。

### Avoid copying files of all targets, which are not asked to build (TBD)

`PLATFORM=broadcom` と指定したのに `mellanox` や `p4` 等もビルドされている？？

```
# make configure PLATFORM=broadcom
# make NOJESSIE=1 NOSTRETCH=1 NOBUSTER=1 SONIC_BUILD_JOBS=4 target/sonic-broadcom.bin
> 20220616-01-sonic-builder-broadcom.log
...snip...
loning into '/home/ubuntu/sonic-buildimage/platform/barefoot/sonic-platform-modules-arista'...
Cloning into '/home/ubuntu/sonic-buildimage/platform/broadcom/saibcm-modules-dnx'...
Cloning into '/home/ubuntu/sonic-buildimage/platform/broadcom/sonic-platform-modules-arista'...
Cloning into '/home/ubuntu/sonic-buildimage/platform/broadcom/sonic-platform-modules-nokia'...
Cloning into '/home/ubuntu/sonic-buildimage/platform/mellanox/hw-management/hw-mgmt'...
Cloning into '/home/ubuntu/sonic-buildimage/platform/mellanox/mlnx-sai/SAI-Implementation'...
Cloning into '/home/ubuntu/sonic-buildimage/platform/mellanox/sdk-src/sx-kernel/Switch-SDK-drivers'...
Cloning into '/home/ubuntu/sonic-buildimage/platform/p4/SAI-P4-BM'...
```

## ビルド設定：rules/config

ビルド関係の設定はこちらで定義されている：[rules/config](https://github.com/Azure/sonic-buildimage/blob/master/rules/config) <br />
また、変更したい項目のみを `rules/config.user` に記載することで、オリジナルのファイルを変更することなくビルドオプションを上書きできる：[Makefile.workの当該箇所](https://github.com/sonic-net/sonic-buildimage/blob/946bc3b969d1960261f806b16db6d2f3d6c4ecee/Makefile.work#L124-L130)

### 有効にする機能の設定場所

SFLOW, FRR, NAT, RESTAPI, P4 等々、SONiCのどの機能を有効にするか設定可能。
以下サンプル抜粋：

```
# INCLUDE_SFLOW - build docker-sflow for sFlow support
INCLUDE_SFLOW = y

# INCLUDE_RESTAPI - build docker-sonic-restapi for configuring the switch using REST APIs
INCLUDE_RESTAPI = n

# INCLUDE_P4RT - build docker-p4rt for P4RT support
INCLUDE_P4RT = n
```

### デフォルト USERNAME PASSWORD の変更

デフォルトの `YourPaSsWoRd` って左手小指疲れる... と思ったら変更。

```
> rules/config.user
# DEFAULT_USERNAME - default username for installer build
DEFAULT_USERNAME = admin

# DEFAULT_PASSWORD - default password for installer build
DEFAULT_PASSWORD = YourPaSsWoRd
```

### DEBUG 設定

詳細はオフィシャル文書参照：[Build debug dockers and debug SONiC installer image](https://github.com/Azure/sonic-buildimage/blob/master/README.buildsystem.md#build-debug-dockers-and-debug-sonic-installer-image)


## 各モジュールで利用されている Debian version 確認方法

各モジュール（TODO:用語確認）で利用される Debian version は指定されている事が多くかつそれぞれ異なる。よってビルド効率化のためにビルドする Debian version を限定する場合、各モジュールが利用している Debian version は含める必要がある。

> 以下のように NOXXXX を加えるとそれらのコンテナは生成されない。
> 但し、BUSTERに依存しているモジュールがあるため以下ビルドは途中で失敗する。（2022/06/17 現在）
> `make NOJESSIE=1 NOSTRETCH=1 NOBUSTER=1 SONIC_BUILD_JOBS=4 target/sonic-broadcom.bin`

各モジュールの Docker File は `dockers/` に保存されている。
Docker File 先頭の `FROM` を読むと、ベースとなるパッケージが確認できる。
（以下、サンプルとして抜粋）

```
e.g. FRR
> dockers/docker-fpm-frr/Dockerfile
FROM docker-swss-layer-buster-sonic:latest
```
