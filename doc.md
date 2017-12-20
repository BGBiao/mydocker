##

### requirements

- ip netns 
- docker > 1.10



### Usage 

#### run docker
```
$ ./JFDocker run YW-D-TPBQD2@af41f7,4Y8K742@af41f7,806ddafce40a,test-xuxuebiao,172.25.63.22,22,172.25.63.254,4,2048m,121121,br0
ip:172.25.63.22 mask:22 gw:172.25.63.254 vnet:br0
hveth: jfdc@9587a,vnet: br0,pid : 32190,ip: 172.25.63.22,mask :22,gw :172.25.63.254
{
  "containername": "test-xuxuebiao",
  "containerid": "cf87e67033a470356c2a22a12e30ce5e9ef16b43f860d561c3d8a442ae2f8543",
  "ipaddress": "172.25.63.22",
  "errmsg": "no",
  "errinfo": "no"
}
$ docker ps
CONTAINER ID        IMAGE                                      COMMAND                  CREATED             STATUS              PORTS                                NAMES
cf87e67033a4        806ddafce40a                               "supervisord -c /etc/"   7 seconds ago       Up 6 seconds                                             test-xuxuebiao
9587a3978adc        xxbandy123/k8s-pause                       "/pod"                   8 seconds ago       Up 7 seconds                                             test-xuxuebiao2017-12-19-12-29-54
$ curl 172.25.63.22
Active connections: 1
server accepts handled requests
 1 1 1
Reading: 0 Writing: 1 Waiting: 0


```
#### conf info

```

$ cat /tmp/Logs/containers_config/runConftest-xuxuebiao
YW-D-TPBQD2@af41f7,4Y8K742@af41f7,806ddafce40a,test-xuxuebiao,172.25.63.22,22,172.25.63.254,4,2048m,121121,br0

$ cat /tmp/Logs/containers_config/ConSpectest-xuxuebiao | python -m json.tool
{
    "cpus": "4",
    "mems": "2048m",
    "sn": "4Y8K742@af41f7",
    "ywid": "YW-D-TPBQD2@af41f7"
}

$ cat /tmp/Logs/containers_config/test-xuxuebiao | python -m json.tool
{
    "Conspec": {
        "cpus": "4",
        "mems": "2048m",
        "sn": "4Y8K742@af41f7",
        "ywid": "YW-D-TPBQD2@af41f7"
    },
    "Image": "806ddafce40a",
    "Netspec": {
        "Gateway": "172.25.63.254",
        "ipaddress": "172.25.63.22",
        "mask": "22",
        "vnet": "br0"
    },
    "containername": "test-xuxuebiao"
}
```

#### log info 

```
$ ls  /tmp/Logs/JFDocker.log

$ tail -n 3 /tmp/Logs/JFDocker.log
time="2017-12-19T12:29:58+08:00" level=info msg="Successful to create pause container for test-xuxuebiao" data=9587a3978adcea077f7ba4ff3a61bd5ef62ff40c79bcc709039d557fa390f56d
time="2017-12-19T12:29:58+08:00" level=info msg="Creating the app container test-xuxuebiao" data=....
time="2017-12-19T12:29:59+08:00" level=info msg="Successful to create app container for test-xuxuebiao" data=cf87e67033a470356c2a22a12e30ce5e9ef16b43f860d561c3d8a442ae2f8543
```

#### update docker with spec image

```
$ ./JFDocker update test-xuxuebiao,74c0b9f60e14
{"containername":"test-xuxuebiao","Image":"74c0b9f60e14","Conspec":{"ywid":"YW-D-TPBQD2@af41f7","sn":"4Y8K742@af41f7","mems":"2048m","cpus":"4"},"Netspec":{"ipaddress":"172.25.63.22","mask":"22","Gateway":"172.25.63.254","vnet":"br0"}}

$ curl 172.25.63.22:9090
Hello ,I'm biaoge.
My Container name is 9587a3978adc

Currently Date:2017-12-18 23:40:53

$ cat /tmp/Logs/containers_config/test-xuxuebiao | python -m json.tool
{
    "Conspec": {
        "cpus": "4",
        "mems": "2048m",
        "sn": "4Y8K742@af41f7",
        "ywid": "YW-D-TPBQD2@af41f7"
    },
    "Image": "74c0b9f60e14",
    "Netspec": {
        "Gateway": "172.25.63.254",
        "ipaddress": "172.25.63.22",
        "mask": "22",
        "vnet": "br0"
    },
    "containername": "test-xuxuebiao"
}

```

#### resize the cpu or mem for app container
`注意:不会去修改原始的配置`
`内存的resize不能超过初始的2倍`

```
# docker inspect test-xuxuebiao | grep -i cpu
            "CpuShares": 4000,
            "CpuPeriod": 100000,
            "CpuQuota": 400000,
            "CpusetCpus": "",
            "CpusetMems": "",
            "CpuCount": 0,
            "CpuPercent": 0,
# ./JFDocker resize test-xuxuebiao,cpu,2
Resize  the docker container's cpu or mem
# docker inspect test-xuxuebiao | grep -i cpu
            "CpuShares": 2000,
            "CpuPeriod": 100000,
            "CpuQuota": 200000,
            "CpusetCpus": "",
            "CpusetMems": "",
            "CpuCount": 0,
            "CpuPercent": 0,


# docker inspect test-xuxuebiao | grep -i mem
            "Memory": 2147483648,
            "CpusetMems": "",
            "KernelMemory": 0,
            "MemoryReservation": 0,
            "MemorySwap": 4294967296,
            "MemorySwappiness": -1,
# ./JFDocker resize test-xuxuebiao,mem,3000m
Resize  the docker container's cpu or mem
# docker inspect test-xuxuebiao | grep -i mem
            "Memory": 3145728000,
            "CpusetMems": "",
            "KernelMemory": 0,
            "MemoryReservation": 0,
            "MemorySwap": 4294967296,
            "MemorySwappiness": -1,


```

#### delete container with app-container

```
# ./JFDocker delete test-xuxuebiao
{
  "containername": "test-xuxuebiao"
}

# ping 172.25.63.22
PING 172.25.63.22 (172.25.63.22) 56(84) bytes of data.
^C
--- 172.25.63.22 ping statistics ---
2 packets transmitted, 0 received, 100% packet loss, time 999ms

# docker ps -a | grep test-xuxuebiao
```
