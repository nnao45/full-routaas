# full-routaas
🐣Ultra minimal🐣 BGP IPv4 Unicast Full Route "as a Service" for Docker container.

## Ultra minimal images.

full route images's size is "168MB"
```bash
# docker images
REPOSITORY                               TAG                 IMAGE ID            CREATED             SIZE
nnao45/full-routaas                      latest              00028f3dd5c9        10 minutes ago      168MB
<none>                                   <none>              d055571d5c0f        10 minutes ago      793MB
```

MIB file size's size is "85MB"...so, full-routaas container size is "83MB" !!
```bash
# ls -lh rib.20180418.0000
-rw-r--r-- 1 root root 85M  4月 18 09:00 rib.20180418.0000
```

## Demo

Nework Diagram.
```bash
+---------------------------------+                            +---------------------------------+
|  Bob                            |                            | Alice                           |
|                                 |                            |                                 |
|   IP Addres   192.168.0.1       |                            |    IP Addres   192.168.0.2      |
|   Subnet mask 255.255.255.0     +----------------------------+    Subnet mask 255.255.255.0    |
|                                 |                            |                                 |
|                                 |                            |                                 |
|                                 |                            |                                 |
+---------------------------------+                            +---------------------------------+
```

Bob Router(full-routaas mem4GB CPU1)  
```bash
Bob# docker run -it --rm --privileged -p 179:179 nnao45/full-routaas:latest
INFO[0000] Add a peer configuration for:192.168.0.2      Topic=Peer
INFO[0000] MRT injection file is ./rib.20180419.0000    
INFO[0015] Peer Up                                       Key=192.168.0.2 State=BGP_FSM_OPENCONFIRM Topic=Peer
INFO[0073] MRT injection complete!! 
```

Alice Router(IOS-XRv ver5.3.0-1 mem4GB CPU1)
```bash
Alice# sho bgp ipv4 unicast summary | begin Neighbor
Neighbor        Spk    AS MsgRcvd MsgSent   TblVer  InQ OutQ  Up/Down  St/PfxRcd
192.168.0.1       0 65555 1238838      31      112    0    0 00:01:37     696234
```

My env, Total route advertisement time is 1 minutes 37sec. 😉 so fast!!
