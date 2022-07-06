# SONiC FAQ

> 分類できないモノをメモする場所

## Install and Upgrade

Reference: [Edgecore SONiC: Installation & Upgrade image](https://support.edge-core.com/hc/en-us/articles/900000208626--Edgecore-SONiC-Installation-Upgrade-image)


### Upgrade via HTTP

Get install image via HTTP and reboot.

```
> (optional) Start HTTP Server on server
$ ufw allow 8000
~/sonic-img$ python3 -m http.server
Serving HTTP on 0.0.0.0 port 8000 (http://0.0.0.0:8000/) ...

> Download and Install SONiC image

admin@sonic:~$ sudo sonic-installer install http://172.20.105.171:8000/sonic-barefoot-20220630-vanilla-rulesmk.bin -y

admin@sonic:~$ sudo sonic-installer list
Current: SONiC-OS-HEAD.0-dirty-20220614.110556
Next: SONiC-OS-HEAD.0-dirty-20220630.003623
Available:
SONiC-OS-HEAD.0-dirty-20220630.003623
SONiC-OS-HEAD.0-dirty-20220614.110556

admin@sonic:~$ sudo reboot
requested COLD shutdown
/var/log: 0 B (0 bytes) trimmed on /dev/loop1
/host: 774.7 MiB (812380160 bytes) trimmed on /dev/sda4
Mon 04 Jul 2022 07:41:53 AM UTC Issuing OS-level reboot ...
> Available images will be listed in GRUB menu

> (optional) Change default image to boot with
$ sudo sonic_installer set_default <image>
```

### Difference between `sonic_installer` and `sonic-installer`

`sonic_installer` is deprecated. One should use `sonic-installer`.

However, both commands are exactly the same and only command name is different. (as of 2022/06/14)

```
admin@sonic:~$ sudo sonic_installer list
Warning: 'sonic_installer' command is deprecated and will be removed in the future
Please use 'sonic-installer' instead
Current: SONiC-OS-HEAD.0-dirty-20220614.110556
Next: SONiC-OS-HEAD.0-dirty-20220614.110556
Available:
SONiC-OS-HEAD.0-dirty-20220614.110556

admin@sonic:~$ sudo sonic-installer list
Current: SONiC-OS-HEAD.0-dirty-20220614.110556
Next: SONiC-OS-HEAD.0-dirty-20220614.110556
Available:
SONiC-OS-HEAD.0-dirty-20220614.110556

admin@sonic:~$ diff /usr/local/bin/sonic_installer /usr/local/bin/sonic-installer
admin@sonic:~$

admin@sonic:~$ cat /usr/local/bin/sonic-installer
#!/usr/bin/python3
# -*- coding: utf-8 -*-
import re
import sys
from sonic_installer.main import sonic_installer
if __name__ == '__main__':
    sys.argv[0] = re.sub(r'(-script\.pyw|\.exe)?$', '', sys.argv[0])
    sys.exit(sonic_installer())
```
