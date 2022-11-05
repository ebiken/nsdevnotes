# gonlmon

Golang Netlink Monitoring Example Code.

Most code was taken from: [gwind/go-netlink-socket-monitor.go](https://gist.github.com/gwind/05f5f649d93e6015cf47ffa2b2fd9713)

## How to use

```
sudo go get github.com/vishvananda/netlink/nl
sudo go get github.com/sirupsen/logrus

nsdevnotes/examples/gonlmon$ go run .
Starting gonlmon.go
Family,       State, InetDiagSockId
--------------------------------------------------------
   tcp,      listen, 127.0.0.1:37987 -> 0.0.0.0:0
   tcp,      listen, 127.0.0.1:33939 -> 0.0.0.0:0
   tcp,      listen, 10.28.65.1:53 -> 0.0.0.0:0
   tcp,      listen, 192.168.122.1:53 -> 0.0.0.0:0
   tcp,      listen, 127.0.0.53:53 -> 0.0.0.0:0
   tcp,      listen, 0.0.0.0:22 -> 0.0.0.0:0
   tcp, established, 127.0.0.1:33939 -> 127.0.0.1:38144
   tcp, established, 172.20.105.171:22 -> 192.168.120.35:53048
   tcp, established, 127.0.0.1:38148 -> 127.0.0.1:33939
   tcp, established, 127.0.0.1:33939 -> 127.0.0.1:38148
   tcp, established, 127.0.0.1:38144 -> 127.0.0.1:33939
   tcp, established, 172.20.105.171:22 -> 192.168.120.35:53067
   tcp, established, 172.20.105.171:22 -> 192.168.120.35:53545
   tcp, established, 172.20.105.171:22 -> 192.168.120.35:53033
```


## Memo

### How to init go module

```
~/sandbox/nsdevnotes/examples/gonlmon$

go mod init example/gonlmon
```
