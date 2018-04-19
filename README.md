# full-routaas
🐣Ultra minimal🐣 BGP IPv4 Unicast Full Route "as a Service" for Docker container.

## Ultra minimal images.

full route images's size is "168MB"
```bash
root@net-flu:full-routaas # docker images
REPOSITORY                               TAG                 IMAGE ID            CREATED             SIZE
nnao45/full-routaas                      latest              00028f3dd5c9        10 minutes ago      168MB
<none>                                   <none>              d055571d5c0f        10 minutes ago      793MB
```

MIB file size's size is "85MB"...so, full-routaas container size is "83MB" !!
```bash
# ls -lh rib.20180418.0000
-rw-r--r-- 1 root root 85M  4月 18 09:00 rib.20180418.0000
```

## Installed and Started

Nework Diagram.
```bash
+---------------------------------+                            +---------------------------------+
|  Bob                            |                            | Alice                           |
|                                 |                            |                                 |
|   IP Addres   192.168.0.0       |                            |    IP Addres   192.168.1.0      |
|   Subnet mask 255.255.255.0     +----------------------------+    Subnet mask 255.255.255.0    |
|                                 |                            |                                 |
|                                 |                            |                                 |
|                                 |                            |                                 |
+---------------------------------+                            +---------------------------------+
```

Bob Router(full-routaas)  
launch.sh Must use "docker", "bzip2"
```bash
Bob# ./launch.sh
<install process>
</install process>
Successfully built 00028f3dd5c9
Successfully tagged nnao45/full-routaas:latest
INFO[0000] Add a peer configuration for:172.30.1.176     Topic=Peer
INFO[0000] MRT injection file is ./rib.20180418.0000  
INFO[0150] Peer Up                                       Key=172.30.1.176 State=BGP_FSM_OPENCONFIRM Topic=Peer
```

Alice Router(gobgpd)
```bash
Alice# gobgpd -f gobgpd.conf
{"level":"info","msg":"gobgpd started","time":"2018-04-18T17:41:08+09:00"}
{"Topic":"Config","level":"info","msg":"Finished reading the config file","time":"2018-04-18T17:41:08+09:00"}
{"level":"info","msg":"Peer 172.30.1.171 is added","time":"2018-04-18T17:41:08+09:00"}
{"Topic":"Peer","level":"info","msg":"Add a peer configuration for:172.30.1.171","time":"2018-04-18T17:41:08+09:00"}
{"Key":"172.30.1.171","State":"BGP_FSM_OPENCONFIRM","Topic":"Peer","level":"info","msg":"Peer Up","time":"2018-04-18T17:41:26+09:00"}
```
