[![Travis CI](https://travis-ci.org/nnao45/full-routaas.svg?branch=master)](https://travis-ci.org/nnao45/full-routaas)
[![Go Report Card](https://goreportcard.com/badge/github.com/nnao45/jgob)](https://goreportcard.com/report/github.com/nnao45/full-routaas)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/nnao45/full-routaas/master/LICENSE)
# full-routaas
🐣Ultra minimal🐣 BGP IPv4 Unicast Full Route "as a Service" for Docker container.  
Using [GoBGP](https://github.com/osrg/gobgp) (as Golang Library) & [Toml](https://github.com/BurntSushi/toml) config files.

## Ultra minimal images.

full route images's size is "102MB"
```bash
# docker images
REPOSITORY            TAG                 IMAGE ID            CREATED             SIZE
nnao45/full-routaas   latest              f62fc829754b        22 minutes ago      102MB
<none>                <none>              db963eb06590        22 minutes ago      949MB
```

RIB file size's size is "85MB"...so, full-routaas container size is "17MB" !!
```bash
# ls -lh rib.20180418.0000
-rw-r--r-- 1 root root 85M  4月 18 09:00 rib.20180418.0000
```
[You can send this container with Gmail](https://support.google.com/mail/answer/6584?co=GENIE.Platform%3DDesktop&hl=en) 😂

## Getting started(only two or three step😘).
### 0.if you use ubuntu 16.04LTS, you may install this.
```bash
# apt-get install -y docker.io bzip2
```

### 1.You download this repo.
```bash
# git clone https://github.com/nnao45/full-routaas.git
```

### 2.You launch
```bash
make launch
```

if this don't move, use node has memory is not enough or be isn't allowd in && out -bound traffic tcp:179.

## Demo

Nework Diagram.
```bash
+---------------------------------+                            +---------------------------------+
|  Bob                            |                            | Alice                           |
|                                 |                            |                                 |
|   IP Addres   192.168.0.1       |                            |    IP Addres   192.168.0.2      |
|   Subnet mask 255.255.255.0     +----------------------------+    Subnet mask 255.255.255.0    |
|   AS          65555             |                            |    AS          65000            |
|   Device      full-routaaas     |                            |    Device      IOS-XRv          |
|                                 |                            |                                 |
+---------------------------------+                            +---------------------------------+
```

full-routaas Config is config.tml.
if you run your env, change config.tml's paramater
```bash
[bgpdconfig]
  as = 65555
  router-id = "192.168.0.1"

[bgpdconfig.mrt-config]
  best-path = false
  skip-v4 = false
  skip-v6 = true
  next-hop = "nil"

[[bgpdconfig.neighbor-config]]
  peer-as = 65000
  neighbor-address = "192.168.0.2"
  peer-type = "external"
```

Bob Router(full-routaas mem4GB CPU1)  
```bash
Bob# docker run -it --rm --privileged -p 179:179 nnao45/full-routaas:latest
INFO[0000] Add a peer configuration for:192.168.0.2      Topic=Peer
INFO[0000] MRT injection file is                        
INFO[0000] Running full-routaas version 1.0.0 !!        
INFO[0018] Peer Up                                       Key=192.168.0.2 State=BGP_FSM_OPENCONFIRM Topic=Peer
INFO[0073] MRT injection complete!! 
```

Alice Router(IOS-XRv ver5.3.0-1 mem4GB CPU1)
```bash
Alice# sho bgp ipv4 unicast summary | begin Neighbor
Neighbor        Spk    AS MsgRcvd MsgSent   TblVer  InQ OutQ  Up/Down  St/PfxRcd
192.168.0.1       0 65555 2573601      58      112    0    0 00:01:26     696234
```

My env, Total route advertisement time is 1 minutes 26sec. 😉 so fast!!

## Writer & License
full-routaas was writed by nnao45 (WORK:Back-end Engineer, Twitter:@nnao45, MAIL:n4sekai5y@gmail.com).  
This software is released under the MIT License, see LICENSE.
