# SONiC が動作する機材

確認方法は２通り。

- Official Wikiで確認
- Source Code で確認

但し、両方に記載が無くてもデータシートには SONiC サポート有りという機材もあるので、どのようにサポートを提供しているかは要確認。
例えば、ベンダ顧客のみアクセス可能なレポでドライバやドキュメント等を提供している場合、それらに関するノウハウをオープンなコミュニティ内で共有して良いか確認が必要になる。

Table of Contents
- [Official Wikiで確認](#official-wikiで確認)
- [Source Code で確認](#source-code-で確認)
- [Tofino / Tofino 2 をサポートするプラットフォーム](#tofino--tofino-2-をサポートするプラットフォーム)


## Official Wikiで確認

[SONiC Wiki: Supported Platforms](https://github.com/sonic-net/SONiC/blob/sonic_image_md_update/supported_devices_platforms.md)

## Source Code で確認

上記ページに記載無い機材でも、ドライバ等が用意されており動作する場合があり。
例えば Tofino 2 は Wiki に記載が無いが、以下の通りドライバは用意されている。

[GitHub: sonic-buildimage/device](https://github.com/Azure/sonic-buildimage/tree/master/device/) に各社が提供するモデル名一覧がある。

例えば barefoot (Tofino ASIC) の場合は、以下ファイルがある。（2022/006/22現在）

https://github.com/Azure/sonic-buildimage/tree/master/device/barefoot

```
x86_64-accton_as9516_32d-r0
x86_64-accton_wedge100bf_32x-r0
x86_64-accton_wedge100bf_65x-r0
x86_64-accton_as9516bf_32d-r0
```

また、[GitHub: sonic-buildimage/platform](https://github.com/Azure/sonic-buildimage/tree/master/platform) にある makefile (platform-modules-*.mk) を確認すると、モデル番号が記載されている場合がある。

例えば barefoot (Tofino ASIC) の場合は、以下ファイルがある。（2022/006/22現在）

https://github.com/Azure/sonic-buildimage/tree/master/platform/barefoot

```
platform-modules-accton.mk
platform-modules-arista.mk
platform-modules-bfn-montara.mk
platform-modules-bfn-newport.mk
platform-modules-bfn.mk
platform-modules-ingrasys.mk
platform-modules-wnc-osw1800.mk
```

上記ファイル中で、モデル名が明記されているモジュール

```
> sonic-buildimage/platform/barefoot$ grep "_PLATFORM =" platform-modules-*

platform-modules-accton.mk:$(BFN_MONTARA_QS_PLATFORM_MODULE)_PLATFORM = x86_64-accton_wedge100bf_32qs-r0
platform-modules-bfn.mk:$(BFN_PLATFORM_MODULE)_PLATFORM = x86_64-accton_wedge100bf_65x-r0
platform-modules-bfn-montara.mk:$(BFN_MONTARA_PLATFORM_MODULE)_PLATFORM = x86_64-accton_wedge100bf_32x-r0
platform-modules-bfn-newport.mk:$(BFN_NEWPORT_PLATFORM_MODULE)_PLATFORM = x86_64-accton_as9516_32d-r0
platform-modules-bfn-newport.mk:$(BFN_NEWPORT_BF_PLATFORM_MODULE)_PLATFORM = x86_64-accton_as9516bf_32d-r0
platform-modules-ingrasys.mk:$(INGRASYS_S9180_32X_PLATFORM_MODULE)_PLATFORM = x86_64-ingrasys_s9180_32x-r0
platform-modules-ingrasys.mk:$(INGRASYS_S9280_64X_PLATFORM_MODULE)_PLATFORM = x86_64-ingrasys_s9280_64x-r0
platform-modules-wnc-osw1800.mk:$(WNC_OSW1800_PLATFORM_MODULE)_PLATFORM = x86_64-wnc_osw1800-r0
```

## Tofino / Tofino 2 をサポートするプラットフォーム

- Tofino
  - WEDGE100BF-32QS : x86_64-accton_wedge100bf_32qs-r0
  - Wedge100BF-65X  : x86_64-accton_wedge100bf_65x-r0
  - WEDGE100BF-32X  : x86_64-accton_wedge100bf_32x-r0
  - [ufiSpace S9180-32X : x86_64-ingrasys_s9180_32x-r0](https://www.ufispace.com/products/datacenter/100g/s9180-32x)
  - [ufiSpace S9280-64X : x86_64-ingrasys_s9280_64x-r0](https://www.ufispace.com/products/datacenter/100g/s9280-64x)
  - [WSW OSW1800 : x86_64-wnc_osw1800-r0](https://www.wnc.com.tw/index.php?action=pro_detail&top_id=106&scid=142&tid=99&id=353)
- Tofino2
  - [DCS810 (AS9516-32D) : x86_64-accton_as9516_32d-r0](https://www.edge-core.com/productsInfo.php?cls=1&cls2=349&cls3=577&id=916)
  - APS Networks ... ドライバは無いがデータシートではサポートしていると表記有り
    - PDF: [APS6232D (Tofino 2) Data Sheet](https://www.aps-networks.com/wp-content/uploads/2022/04/220421_APSN_APS6232D_V01.pdf)
    - Supported Software (OS) Debian / Ubuntu
    - Supported Applications SONiC / STRATUM
