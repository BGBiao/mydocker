// 容器网络相关的操作,主要包含AttactNet和DetachNet
//

package clients

import (
	_ "encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xxbandy/mydocker/apis"
	"os"
	"os/exec"
	"strings"
)

func AttachNet(netns, containerid, ip, mask, gw, vnet string) (*apis.Netns, error) {

	//cpid 用来获取pause容器的pid，但是是uint8格式的，默认有个换行符
	cpid, err := exec.Command("/bin/bash", "-c", `docker inspect -f '{{ .State.Pid }}' `+containerid).Output()
	if err != nil {
		log.Fatalf("Failed to get pause container pid", err.Error())
		//fmt.Println("pause containerid get error")
		//os.Exit(1)
	}
	pid := strings.Replace(string(cpid), "\n", "", -1)

	//构造进程ns相关路径以及容器的虚拟网卡 /proc/$pidnum/ns/net /var/run/netns/$pidnum jfdc@contaild[0-5]
	netnspath := "/proc/" + pid + "/ns/net"
	netpath := netns + pid
	hveth := "jfdc@" + containerid[0:5]

	//判断/var/run/netns目录是否存在，不存在创建
	netfile, err := os.Stat(netns)
	if err != nil {
		fmt.Printf("%s does not exist,creating ...\n", netns)
		fileerr := os.MkdirAll(netns, 0755)
		if fileerr != nil {
			fmt.Printf("%s creation failure: %s", netns, fileerr.Error())
			os.Exit(1)
		}
	} else {
		Logger("Linux net namespace path exist", netfile.Name())
	}
	//fmt.Println(netfile.Name()) //netfile是一个Filepath结构体

	//创建netlink软连接
	testnetlink, err := os.Lstat(netpath)
	if err == nil {
		Logger("netlink aleady exist", testnetlink.Name())
		//fmt.Printf("netlink %s is exits\n",testnetlink.Name())
	} else {
		Logger("Creating the netlink", netpath)
		//fmt.Println("netlink creation: "+netpath)
		netnserr := os.Symlink(netnspath, netpath)
		if netnserr != nil {
			log.Fatalf("netlink creation failure", err.Error())
			//fmt.Println("netlink creation failure:"+err.Error())
			//os.Exit(1)
		}
		Logger("Successful to create the netlink", netpath)
	}

	//创建网络命名空间,并创建相关的容器网络
	fmt.Printf("hveth: %s,vnet: %s,pid : %s,ip: %s,mask :%s,gw :%s\n", hveth, vnet, pid, ip, mask, gw)
	Logger("Creating the vnet and net configuration", "...")
	out, err := exec.Command("/bin/bash", "-c", `ip link add `+hveth+` type veth peer name eth0`+hveth).Output()
	if err != nil {
		Logger("ip link add "+hveth+" type veth peer eth0"+hveth, "fail")
		log.Fatalf("Failed to create net pair", err.Error())
		//fmt.Println(err.Error())
		///fmt.Printf("ip link add %s type veth peer eth0%s\n failed!",hveth,hveth)
		//os.Exit(1)
	}
	//fmt.Println(string(out))
	Logger("Successful to create net pairs", string(out))

	netbr, err := exec.Command("/bin/bash", "-c", `brctl addif `+vnet+` `+hveth).Output()
	if err != nil {
		Logger("brctl addif "+vnet+" "+hveth, "fail")
		log.Fatalf("Failed to create bridge interface", err.Error())
		//fmt.Println(err.Error())
		//fmt.Printf("brctl addif %s %s failed!\n",vnet,hveth)
		//os.Exit(1)
	}
	//fmt.Println(string(netbr))
	Logger("Successful to create bridge interface ", string(netbr))

	netup, err := exec.Command("/bin/bash", "-c", `ip link set `+hveth+` up`).Output()
	if err != nil {
		Logger("ip link set "+hveth+" up", "fail")
		log.Fatalf("Failed to up the vnet", err.Error())
		// fmt.Println(err.Error())
		// fmt.Printf("ip link set %s up failed!\n",hveth)
		// os.Exit(1)
	}
	Logger("Successful to up the vnet", string(netup))
	//fmt.Println(string(netup))

	netlink, err := exec.Command("/bin/bash", "-c", `ip  link set eth0`+hveth+` netns `+pid).Output()
	if err != nil {
		Logger("ip link set eth0"+hveth+" netns "+pid, "fail")
		log.Fatalf("Failed to create relation with vnet and containerpid", err.Error())
		// fmt.Println(err.Error())
		// fmt.Printf("ip link set eth0%s netns %s failed!\n",hveth,pid)
		// os.Exit(1)
	}
	Logger("Successful to create the relation", string(netlink))
	//fmt.Println(string(netlink))

	netdevup, err := exec.Command("/bin/bash", "-c", `ip netns exec `+pid+` ip link set eth0`+hveth+` up`).Output()
	if err != nil {
		Logger("ip netns exec "+pid+" ip link set eth0"+hveth+" up", "fail")
		log.Fatalf("Failed to set up the net dev", err.Error())
		// fmt.Println(err.Error())
		// fmt.Printf("ip netns exec %s ip link set eth0%s up failed!\n",pid,hveth)
		// os.Exit(1)
	}
	Logger("Successful to set up the container net dev", string(netdevup))
	//fmt.Println(string(netdevup))

	netip, err := exec.Command("/bin/bash", "-c", `ip netns exec `+pid+` ip addr add `+ip+`/`+mask+` dev eth0`+hveth).Output()
	if err != nil {
		Logger("ip netns exec "+pid+" ip addr add "+ip+"/"+mask+" dev eth0"+hveth, "fail")
		log.Fatalf("Failed to setup the con ip", err.Error())
		// fmt.Println(err.Error())
		// fmt.Printf("ip netns exec %s ip addr add %s/%s dev eth0 failed!\n",pid,ip,mask)
		// os.Exit(1)
	}
	Logger("Successful to setup the container ip", string(netip))
	//fmt.Println(string(netip))

	netgw, err := exec.Command("/bin/bash", "-c", `ip netns exec `+pid+` ip route add default via `+gw).Output()
	if err != nil {
		Logger("ip netns exec "+pid+" ip route add default via "+gw, "fail")
		log.Fatalf("Failed to setup the container gateway", err.Error())
		//fmt.Println(err.Error())
		//fmt.Printf("ip netns exec %s ip route add default via %s fialed!\n",pid,gw)
		//os.Exit(1)
	}
	Logger("Successful to setup the container gateway", string(netgw))
	//fmt.Println(string(netgw))

	//构造网络相关返回信息
	netdata := apis.Netns{}
	netdata.Cid = containerid
	netdata.Pid = pid
	netdata.Netns = netpath
	netdata.Net = apis.NetSpec{ip, mask, gw, vnet}
	//直接在函数内部返回结构体指针，在调用方进行json格式化
	//b,err := json.Marshal(netdata)
	//if err != nil {
	//    fmt.Printf("pause container net info to json failed:%s \n",err.Error())
	//}

	return &netdata, nil
}

func DetachNet(pauseid, ip, mask string) error {
	cpid, err := exec.Command("/bin/bash", "-c", `docker inspect -f '{{ .State.Pid }}' `+pauseid).Output()
	if err != nil {
		log.Fatalf("Failed to get pause container pid", err.Error())
	}
	//获取pause容器的pid和容器的虚拟网卡
	pid := strings.Replace(string(cpid), "\n", "", -1)
	hveth := "jfdc@" + pauseid[0:5]

	//开始进行netns 卸载操作
	Logger("Detaching the container netns:"+pid, hveth)

	delip, err := exec.Command("/bin/bash", "-c", `ip netns exec `+pid+` ip addr del `+ip+`/`+mask+` dev eth0`+hveth).Output()
	if err != nil {
		Logger("ip netns exec "+pid+" ip addr del "+ip+"/"+mask+" dev eth0"+hveth, "fail")
		log.Fatalf("Failed to delete the container ip", err.Error())
	}
	Logger("Successful to detach a ip for pause container:"+pauseid, string(delip))

	delvdev, err := exec.Command("/bin/bash", "-c", `ip netns exec `+pid+` ip link set eth0`+hveth+` down && ip netns exec `+pid+` ip link del dev eth0`+hveth).Output()
	if err != nil {
		Logger("ip netns exec "+pid+" ip link set eth0"+hveth+" down && ip netns exec "+pid+" ip link del dev eth0"+hveth, "fail")
		log.Fatalf("Failed to delete virtual dev.", err.Error())
	}

	Logger("Successful to detach  a virtual net device  for pause container:"+pauseid, string(delvdev))

	//删除pause容器网络的netns "/var/run/netns/$pid"
	testNetlink := apis.Netnamespace + pid
	if removeErr := os.Remove(testNetlink); removeErr == nil {
		Logger("Successful to unlink the netns. ", testNetlink)
	} else {
		log.Fatalf("Failed to unlink the netns ", testNetlink)
	}

	return nil

}
