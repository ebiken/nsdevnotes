# SONiC Switch Image

- [Getting pre-built image](#getting-pre-built-image)
- [Building SONiC image](#building-sonic-image)

## Getting pre-built image

You can find Pre-built packages here: https://sonic-build.azurewebsites.net/ui/sonic/pipelines

- Find `Platform`, `BranchName` and click link on `Builds` (e.g. `vs`, `master`)
- Find the build you want to use and click `Artifacts`
  - make sure the `Result` is `succeeded`
- Click `Name` (e.g. `sonic-buildimage.vs1`)
- Download `*.img.gz` or `*.bin` (e.g. `target/sonic-vs.img.gz`, `target/sonic-broadcom.bin`)

## Building SONiC image

> TODO: KVM, Docker, BCM ASIC, Tofino ASIC 各々をビルドする方法を記載