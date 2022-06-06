# SONiC Commands

- SONiC では bash 内で実行する `show` `config` 等のコマンドが用意されており、以下の特長（機能）を持つ
  - `config --help` `config interface --help` `show ip int --help` でヘルプを表示可能
  - TABキーで補完される： `show int<tab>` => `show interfaces `
  - 完全一致する場合、オプションを完全に入力しなくても受け付ける： `show ip int` で `show ip interfaces` が実行される
- 各コマンドの詳細はヘルプ（`--help`）及びオフィシャル Wiki を参照
  - [sonic-utilities/doc/Command-Reference.md](https://github.com/Azure/sonic-utilities/blob/master/doc/Command-Reference.md)
- 本ページには、特に理解が難しかった or 利用頻度の高い or アウトプットを後日参照したい or 覚書が必要なコマンドをメモしていく。

Table of Contents
- [SONiC コマンドの実装について](#sonic-コマンドの実装について)
  - [Source Code の場所](#source-code-の場所)
  - [SONiC 上の実行ファイルの場所](#sonic-上の実行ファイルの場所)
- [show vlan](#show-vlan)
- [xx](#xx)

## SONiC コマンドの実装について

### Source Code の場所

Python3 で実装されており、Source Code は以下レポジトリに存在。

- [https://github.com/Azure/sonic-utilities](https://github.com/Azure/sonic-utilities)
  - [config](https://github.com/Azure/sonic-utilities/tree/master/config)
  - [show](https://github.com/Azure/sonic-utilities/tree/master/show)

### SONiC 上の実行ファイルの場所

- config コマンド

```
$ cat /usr/local/bin/config
#!/usr/bin/python3
# -*- coding: utf-8 -*-
import re
import sys
from config.main import config
if __name__ == '__main__':
    sys.argv[0] = re.sub(r'(-script\.pyw|\.exe)?$', '', sys.argv[0])
    sys.exit(config())

admin@sonic:/usr/local/lib/python3.9/dist-packages/config$ ls
aaa.py              config_mgmt.py  feature.py        __init__.py  kube.py  mclag.py     nat.py   __pycache__  vlan.py
chassis_modules.py  console.py      flow_counters.py  kdump.py     main.py  muxcable.py  plugins  utils.py     vxlan.py
```

- show コマンド

```
$ cat /usr/local/bin/show
!/usr/bin/python3
# -*- coding: utf-8 -*-
import re
import sys
from show.main import cli
if __name__ == '__main__':
    sys.argv[0] = re.sub(r'(-script\.pyw|\.exe)?$', '', sys.argv[0])
    sys.exit(cli())

>>> 実装
admin@sonic:/usr/local/lib/python3.9/dist-packages/show$ ls
acl.py         bgp_frr_v6.py       dropcounters.py   gearbox.py   kube.py      platform.py   reboot_cause.py   vnet.py
aliases.ini    bgp_quagga_v4.py    feature.py        __init__.py  main.py      plugins       sflow.py          vxlan.py
bgp_common.py  bgp_quagga_v6.py    fgnhg.py          interfaces   muxcable.py  processes.py  system_health.py  warm_restart.py
bgp_frr_v4.py  chassis_modules.py  flow_counters.py  kdump.py     nat.py       __pycache__   vlan.py
```

## show vlan

```
$ show vlan config

Name        VID  Member      Mode
--------  -----  ----------  --------
Vlan1001   1001  Ethernet0   untagged
Vlan1001   1001  Ethernet4   untagged
Vlan1002   1002  Ethernet8   untagged
Vlan1002   1002  Ethernet12  tagged

$ show vlan brief
+-----------+----------------+------------+----------------+-------------+-----------------------+
|   VLAN ID | IP Address     | Ports      | Port Tagging   | Proxy ARP   | DHCP Helper Address   |
+===========+================+============+================+=============+=======================+
|      1001 | 192.168.1.1/24 | Ethernet0  | untagged       | disabled    |                       |
|           |                | Ethernet4  | untagged       |             |                       |
+-----------+----------------+------------+----------------+-------------+-----------------------+
|      1002 | 192.168.2.1/24 | Ethernet8  | untagged       | disabled    |                       |
|           |                | Ethernet12 | tagged         |             |                       |
+-----------+----------------+------------+----------------+-------------+-----------------------+
```

## xx

