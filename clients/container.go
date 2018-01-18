package clients

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xxbandy/mydocker/apis"
	"os"
	"os/exec"
	"strings"
)

func JFDockerRun(args string) (*apis.RspJFDocker, error) {
	//构造函数返回信息
	data := apis.RspJFDocker{}

	//需要将传进来的args(以,分割的参数)进行分割提取，并且判断参数长度是否符合需求
	//[YW-D-TPBQD2@af41f7 4Y8K742@af41f7 172.25.46.9:5001/centos6.8-jdjr-test-app 172.25.47.21.h.chinabank.com.cn 172.25.47.21 24 172.25.47.254 2 4096m 20007772 br0]
	argslist := strings.Split(args, ",")

	if len(argslist) != 11 {
		fmt.Println("The JFDocker run args have some error,please check it")
		Logger("JFDocker run 参数有误请检查", string(args))
		os.Exit(1)
	}

	//JFDocker run参数落地存储
	runconf := "runConf" + argslist[3]
	Logger("JFDocker run args are storing to the "+runconf, args)
	if RunConfErr := WriteToConf(runconf, args); RunConfErr == nil {
		Logger("Successful store the JFDocker run args", args)
	}

	//格式化容器相关信息
	container := apis.JFDocker{}
	container.Appname = argslist[3]
	container.Image = argslist[2]
	//container.Conspec = apis.ConSpec{argslist[0],argslist[1],argslist[8],argslist[7]}
	container.Conspec = apis.ConSpec{Ywid: argslist[0], SN: argslist[1], Mem: argslist[8], Cpus: argslist[7]}
	container.Netspec = apis.NetSpec{argslist[4], argslist[5], argslist[6], argslist[10]}

  //创建Cephfs目录并根据相关目录结构进行挂载
  /*
  if err := CreateMountfs(container.Appname); err == nil {
      Logger("Successful to create and mount the cephfs ",apis.DataPath+container.Appname)
  }
  */




	//b as a []byte
	b, err := json.Marshal(container)
	if err != nil {
		log.Fatalf("Failed to json the container datainfo", err.Error())
	}

	//JFDocker container详细信息参数落地存储
	//fmt.Printf("标准的json串:%s\n",b)
	if ConDataInfoErr := WriteToConf(container.Appname, string(b)); ConDataInfoErr == nil {
		Logger("Successful store the Container data info:"+container.Appname, string(b))
	}

	//JFDocker 容器规格信息参数落地存储
	//conspec := apis.ConSpec{argslist[0],argslist[1],argslist[8],argslist[7]}
	conspec := &container.Conspec
	if conspecdata, err := json.Marshal(conspec); err == nil {
		conconf := "ConSpec" + container.Appname
		if ConSpecErr := WriteToConf(conconf, string(conspecdata)); ConSpecErr == nil {
			Logger("Successful store the Container Spec info:"+container.Appname, string(conspecdata))
		}
	} else {
		log.Fatalf("Failed to write the conspec file", err.Error())
	}

	fmt.Printf("ip:%s mask:%s gw:%s vnet:%s\n", container.Netspec.Ipv4, container.Netspec.Mask, container.Netspec.Gateway, container.Netspec.Vnet)

	//创建pause容器，并且同时需要增加网络操作

	pauseid, err := Pausec(container.Appname, container.Netspec.Ipv4, container.Netspec.Mask, container.Netspec.Gateway, container.Netspec.Vnet)
	if err != nil {
		log.Fatalf("Failed to run the pausecontainer", err.Error())
		/*
		   fmt.Println("Pause container created filed!")
		   os.Exit(1)
		*/
	}
	//创建pause容器成功
	//fmt.Println("pause 容器创建成功"+pauseid)
	Logger("Successful to create pause container for "+container.Appname, pauseid)
	Logger("Creating the app container "+container.Appname, "....")

	fmt.Println(pauseid, container.Appname, container.Image, container.Conspec.Cpus, container.Conspec.Mem)
	//根据pause容器id进行创建业务容器appcon
	//func Appc(pauseid,name,appimage,cpus,mems string)  (appcid string,err error)

	appcid, err := Appc(pauseid, container.Appname, container.Image, container.Conspec.Cpus, container.Conspec.Mem)
	if err != nil {
		log.Fatalf("Failed to run the appcontainer", err.Error())
		/*
		   fmt.Println("应用容器创建失败:"+err.Error())
		   os.Exit(1)
		*/
	}
	//fmt.Printf("应用容器:%s 创建成功,容器id为:%s",container.Appname,appcid)
	Logger("Successful to create app container for "+container.Appname, appcid)

	//构造响应
	data.Result = 0
	data.Appname = argslist[3]
	data.ConID = appcid
	data.Ipv4 = argslist[4]
	data.ErrMsg = "no"
	data.ErrInfo = "no"

	return &data, nil

}

func JFDockerUpdate(args string) (*apis.JFDocker, error) {
	data := apis.JFDocker{}
	argslist := strings.Split(args, ",")
	//更新镜像的操作参数应该为: name,image
	if len(argslist) != 2 {
		log.Fatalf("The JFDocker update args have some error", apis.Usages)
	}

	name := argslist[0]
	image := argslist[1]

	//读取当前容器配置信息到data接口中
	ConConfContent, err := ReadFromConf(name)
	if err != nil {
		os.Exit(1)
	}
	ConData := []byte(ConConfContent)
	ConInfoErr := json.Unmarshal(ConData, &data)
	if ConInfoErr != nil {
		os.Exit(2)
	}
	//读取容器配置文件成功并进行配置更新，方便后期更新配置文件
	data.Image = image
	cpus := data.Conspec.Cpus
	mems := data.Conspec.Mem

	//构造新的JFDocker run参数
	newargs := string(data.Conspec.Ywid) + "," + string(data.Conspec.SN) + "," + string(data.Image) + "," + string(data.Appname) + "," + string(data.Netspec.Ipv4) + "," + string(data.Netspec.Mask) + "," + string(data.Netspec.Gateway) + "," + string(data.Conspec.Cpus) + "," + string(data.Conspec.Mem) + ",1232123," + string(data.Netspec.Vnet)

	//容器镜像更新成功
	s, c := UpdateAppc(name, image, cpus, mems)
	if s == "ok" {
		Logger("Successful to update the appcontainer:"+name+" with image:"+image, "容器id:"+c)
	}

	Logger("Updating  the Runconf", "...")
	runConf := "runConf" + name
	if RunConferr := WriteToConf(runConf, newargs); RunConferr == nil {
		Logger("Successful to update the runConf for container "+name, "")
	}
	Logger("Updating the container conf", "...")
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to update the container datainfo", err.Error())
	}
	fmt.Println(string(b))
	//这里写这个string(b)数据的时候可能造成数据写入不正确，不能直接清空之前的内容进行重写
	if ConDataInfoErr := WriteToConf(name, string(b)); ConDataInfoErr == nil {
		Logger("Successful to update the container data info:"+name, string(b))
	}

	//func UpdateAppc(name,image string) (stat,appc string)
	return &data, nil
}

//理论上为保证容器以及配置状态的一致性，应该在resize之后进行状态同步,但是由于该操作是临时操作，希望在下一次恢复之后使用原来的配置规格进行重启，此处配置没有进行同步修改，另外如果需要更改规格大小，后期可以通过UpdateAppc传入相关的mem和cpu参数进行更新容器配置
func JFDockerResize(args string) {
	argslist := strings.Split(args, ",")
	if len(argslist) != 3 {
		log.Fatalf("The JFDocker resize args have some error,please check it", string(args))
	}

	name := argslist[0]
	resizetype := argslist[1] //cpu or mem
	values := argslist[2]

	ResizeAppc(name, resizetype, values)
	//读取三个配置文件并将更新的相关配置内容同步到文件
	Logger("Successful to resize the container's "+resizetype+" to "+values, name)

}

func JFDockerDelete(args string) (*apis.RspJFDocker, error) {
	//JFDocker delete appname
	delData := apis.RspJFDocker{}
	argslist := strings.Split(args, ",")
	if len(argslist) != 1 {
		log.Fatalf("The JFDocker delete args have some error,pleace check it", apis.Usages)
	}

	delData.Appname = argslist[0]
	//获取pauseid
	netmode, err := exec.Command("/bin/bash", "-c", `docker inspect -f "{{.HostConfig.NetworkMode}}" `+delData.Appname).Output()
	if err != nil {
		log.Fatalf("Failed to get the pauseid,please check it", err.Error())
	}

	pauseid := strings.Replace(strings.Split(string(netmode), ":")[1], "\n", "", -1)

	//先销毁业务容器
	status, _ := DelAppc(delData.Appname)
	if status == "ok" {
		Logger("Successful to delete the app container", delData.Appname)
	}

	//销毁pause容器
	//需要销毁pause容器的同时删除相关配置文件
	DelPausec(pauseid, delData.Appname)

	return &delData, nil
}

func JFDockerRebuilt(args string) error {
	argslist := strings.Split(args, ",")
	if len(argslist) != 1 {
		log.Fatalf("The JFDocker rebuilt args have some error,pleace check it", apis.Usages)
	}
	rebuiltname := argslist[0]
	runConf := "runConf" + rebuiltname
	runArgs, err := ReadFromConf(runConf)
	if err != nil {
		log.Fatalf("Failed to get the runConf for container"+rebuiltname, err.Error())
	}

	//JFDocker delete
	if result, err := JFDockerDelete(rebuiltname); err == nil {
		Logger("Successful to delete the container", string(result.Appname))
		//run the delete
		if runResult, runError := JFDockerRun(runArgs); runError == nil {
			if data, err := json.Marshal(runResult); err == nil {
				Logger("Successful to rebuilt the container for "+rebuiltname, string(data))
			}
		} else {
			log.Fatalf("Failed to run container", rebuiltname)
		}

	} else {
		log.Fatalf("Failed to delete container", rebuiltname)
	}
	return nil
}

func JFDockerRunGpu(args string) (*apis.RspJFDocker, error) {
	data := apis.RspJFDocker{}
	argslist := strings.Split(args, ",")
	//[YW-D-TPBQD2@af41f7 4Y8K742@af41f7 172.25.46.9:5001/centos6.8-jdjr-test-app 172.25.47.21.h.chinabank.com.cn 172.25.47.21 24 172.25.47.254 2 4096m 20007772 br0]YW-D-TPBQD2@af41f7 4Y8K742@af41f7 172.25.46.9:5001/centos6.8-jdjr-test-app 172.25.47.21.h.chinabank.com.cn 172.25.47.21 24 172.25.47.254 2 4096m 20007772 br0 4}
	if len(argslist) != 12 {
		log.Fatalf("The JFDocker rungpu args have some error,please check it", string(args))
	}
	//runConf+name
	runconf := "runConf" + argslist[3]
	Logger("JFDocker run args are storing to the "+runconf, args)
	if RunConfErr := WriteToConf(runconf, args); RunConfErr == nil {
		Logger("Successful store the JFDocker run args", args)
	}

	container := apis.JFDocker{}
	container.Appname = argslist[3]
	container.Image = argslist[2]
	container.Conspec = apis.ConSpec{argslist[0], argslist[1], argslist[8], argslist[7], argslist[11]}
	container.Netspec = apis.NetSpec{argslist[4], argslist[5], argslist[6], argslist[10]}

  // 检查并创建容器共享存储. GPU当前环境没有可用CephFs集群
  /*
  if err := CreateMountfs(container.Appname); err == nil {
      Logger("Successful to create and mount the cephfs ",apis.DataPath+container.Appname)
  }

  */

	b, err := json.Marshal(container)
	if err != nil {
		log.Fatalf("Failed to json the container datainfo", err.Error())
	}
	//container sepc conf
	if ConDataInfoErr := WriteToConf(container.Appname, string(b)); ConDataInfoErr == nil {
		Logger("Successful store the Container data info:"+container.Appname, string(b))
	}

	//ConSpec+name
	conspec := &container.Conspec
	if conspecdata, err := json.Marshal(conspec); err == nil {
		conconf := "ConSpec" + container.Appname
		if ConSpecErr := WriteToConf(conconf, string(conspecdata)); ConSpecErr == nil {
			Logger("Successful store the Container Spec info:"+container.Appname, string(conspecdata))
		}
	} else {
		log.Fatalf("Failed to write the conspec file", err.Error())
	}

	fmt.Printf("ip:%s mask:%s gw:%s vnet:%s\n", container.Netspec.Ipv4, container.Netspec.Mask, container.Netspec.Gateway, container.Netspec.Vnet)

	//创建pause容器
	pauseid, err := Pausec(container.Appname, container.Netspec.Ipv4, container.Netspec.Mask, container.Netspec.Gateway, container.Netspec.Vnet)
	if err != nil {
		log.Fatalf("Failed to run the pausecontainer", err.Error())
	}

	Logger("Successful to create pause container for "+container.Appname, pauseid)
	Logger("Creating the GPU app container "+container.Appname, "....")
	Logger("gpu args"+pauseid+" "+container.Appname+" "+container.Image+"  "+container.Conspec.Cpus+" "+container.Conspec.Mem+" ", container.Conspec.Gpus)

	//创建GPU应用容器
	appcid, err := RunGpuAppc(pauseid, container.Appname, container.Image, container.Conspec.Cpus, container.Conspec.Mem, container.Conspec.Gpus)
	if err != nil {
		fmt.Println("应用容器创建失败:" + err.Error())
		os.Exit(1)
	}
	fmt.Println(container.Appname, appcid)
	return &data, nil
}

func JFDockerUpdateGpu(args string) (*apis.JFDocker, error) {
	data := apis.JFDocker{}
	argslist := strings.Split(args, ",")
	if len(argslist) != 2 {
		log.Fatalf("The JFDocker update args have some error", apis.Usages)
	}

	name := argslist[0]
	image := argslist[1]
	ConConfContent, err := ReadFromConf(name)
	if err != nil {
		os.Exit(1)
	}
	ConData := []byte(ConConfContent)
	ConInfoErr := json.Unmarshal(ConData, &data)
	if ConInfoErr != nil {
		os.Exit(2)
	}
	data.Image = image
	cpus := data.Conspec.Cpus
	mems := data.Conspec.Mem
	gpus := data.Conspec.Gpus

	newargs := string(data.Conspec.Ywid) + "," + string(data.Conspec.SN) + "," + string(data.Image) + "," + string(data.Appname) + "," + string(data.Netspec.Ipv4) + "," + string(data.Netspec.Mask) + "," + string(data.Netspec.Gateway) + "," + string(data.Conspec.Cpus) + "," + string(data.Conspec.Mem) + ",1232123," + string(data.Netspec.Vnet) + "," + string(data.Conspec.Gpus)
	fmt.Println(newargs)

	//update gpu container with gpu image
	s, c := UpdateGpuAppc(name, image, cpus, mems, gpus)
	if s == "ok" {
		Logger("Successful to update the appcontainer:"+name+" with image:"+image, "容器id:"+c)
	}

	Logger("Updating  the Runconf", "...")
	runConf := "runConf" + name
	if RunConferr := WriteToConf(runConf, newargs); RunConferr == nil {
		Logger("Successful to update the runConf for gpu container "+name, "")
	}

	Logger("Updating the container conf", "...")
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to update the container datainfo", err.Error())
	}
	fmt.Println(string(b))
	if ConDataInfoErr := WriteToConf(name, string(b)); ConDataInfoErr == nil {
		Logger("Successful to update the container data info:"+name, string(b))
	}
	return &data, nil
}

func JFDockerDeleteGpu(args string) (*apis.RspJFDocker, error) {
	delData := apis.RspJFDocker{}
	argslist := strings.Split(args, ",")
	if len(argslist) != 1 {
		log.Fatalf("The JFDocker deletegpu args have some error,pleace check it", apis.Usages)
	}
	delData.Appname = argslist[0]
	netmode, err := exec.Command("/bin/bash", "-c", `docker inspect -f "{{.HostConfig.NetworkMode}}" `+delData.Appname).Output()
	if err != nil {
		log.Fatalf("Failed to get the pauseid,please check it", err.Error())
	}
	pauseid := strings.Replace(strings.Split(string(netmode), ":")[1], "\n", "", -1)
	status, _ := DelAppc(delData.Appname)
	if status == "ok" {
		Logger("Successful to delete the gpu app container", delData.Appname)
	}
	DelPausec(pauseid, delData.Appname)
	return &delData, nil
}
