package clients

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xxbandy/mydocker/apis"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//appcontainer是采用mem和cpu的独占方式还是共享方式?
//一期暂时先采用独占方式和之前的接口完全一致

/*
request: mem:2048m cpu:2
mem: docker run -m 2048m
cpu:
  cpurequest := cpu*1000 (2000)
    --cpu-shares cpurequest
  cpulimit = (cpu-quota/cpu-period)
    --cpu-quota := cpu*100000
    --cpu-period := 100000

  docker run --cpu-shares 2000 --cpu-quota 200000 --cpu-period 100000

*/

//指定pause容器id,appname,image,cpus,mems进行应用容器创建，返回应用容器id以及错误信息
func Appc(pauseid, name, appimage, cpus, mems string) (appcid string, err error) {
	cpushares, cpuquota := Cpuget(cpus)
	var appconid string
	conspecconf := apis.LogPath + apis.ConConfDir + "/ConSpec" + name
	//注意:由于直接使用了exec.Command().Output() 因此命令执行过程中不能输出详细的错误信息
	//在docker中--hostname和当前--net=container:cid 不能共用.因此当前不能手工设置容器内部的容器名
  //-v ConSpecConf:/export/config_info 
  //-v /data/biao:/data(a cephfs)
	if appcid, err := exec.Command("/bin/bash", "-c", `docker run -itd --name `+name+` -v `+conspecconf+`:/export/config_info  -v `+apis.DataPath+name+`:/data --net=container:`+pauseid+` --ipc=container:`+pauseid+` -m `+mems+` --cpu-shares `+cpushares+` --cpu-quota `+cpuquota+` --cpu-period 100000 `+appimage).Output(); err == nil {
		appconid = strings.Replace(string(appcid), "\n", "", -1)
	} else {
		fmt.Println("appcontainer 创建失败，请检查容器参数" + err.Error())
		os.Exit(1)
	}
	return appconid, err

}

//指定appname进行应用容器销毁，返回状态和应用容器的id
func DelAppc(name string) (stat, appcid string) {
	var status, appconid string
	if id, err := exec.Command("/bin/bash", "-c", `docker stop `+name+` && docker rm -f -v `+name).Output(); err == nil {
		status = "ok"
		appcid = strings.Replace(strings.Split(string(id), "\n")[0], "\n", "", -1)
	} else {
		status = "sorry"
		fmt.Println("容器停止删除失败:" + err.Error())
		os.Exit(2)
	}
	return status, appconid
}

//UpdateAppc为更新容器，而容器镜像的更新，内存和cpu都可以认为是一个更新操作.考虑到当前环境中大多数为tomcat环境，如果更改了内存的话,其实也是需要重启应用生效的，所以这里参数将UpdateAppc进行分开构造，如果仅仅更新镜像的话使用如下函数进行应用容器的销毁以及重建。如果是更新容器的内存和mem相关信息，可以直接使用`docker update`相关命令进行重新封装

func UpdateAppc(name, image, cpus, mems string) (stat, appc string) {
	//获取pause容器的id
	var pauseid, status, appcon string
	if out, err := exec.Command("/bin/bash", "-c", `docker inspect -f "{{.HostConfig.NetworkMode}}" `+name).Output(); err == nil {
		//构造获取pause容器的id
		pauseid = strings.Split(strings.Split(string(out), "\n")[0], ":")[1]
		//pauseid := strings.Split(strings.Replace(string(out),"\n","",-1),":")[1]
	} else {
		log.Fatalf("Failed to get the pause container id", err.Error())
	}

	//删除容器成功后进行创建容器
	if stat, _ := DelAppc(name); stat == "ok" {
		Logger("The container of lastimage has been delete", name)
		name, err := Appc(pauseid, name, image, cpus, mems)
		if err != nil {
			log.Fatalf("The container of newimage run failed,please check it ", err.Error())
		}
		status = "ok"
		appcon = name
		Logger("The container of newimage has been created", appcon)

	} else {
		log.Fatalf("The container of lastimage delete failed", "check it")
	}

	return status, appcon
}

/*ResizeAppc 主要用来动态更新容器的基础配置，目前主要对接容器的内存和mem的配置(其实该功能没有必要，因为使用容器方式进行上线基本都是原子性的操作，完全可以使用上述UpdateAppc来重新更新容器的相关基本配置信息) 暂时使用`docker update `用来更新容器的配置信息
 */
//resizetype 变量和type的坑啊
func ResizeAppc(name, resizetype, resizevalue string) {
	if resizetype == "mem" {
		//update时候的mem最大不能超过之前内存的两倍,这个是由于MemorySwap做了限制
		if _, memError := exec.Command("/bin/bash", "-c", `docker update -m `+resizevalue+` `+name).Output(); memError == nil {
			Logger("The container's mem has been resized to:"+resizevalue, name)
		} else {
			log.Fatal("Failed to resize the mem for container:"+name, memError.Error())
		}
	} else if resizetype == "cpu" {
		cpushares, cpuquota := Cpuget(resizevalue)
		if _, cpuError := exec.Command("/bin/bash", "-c", `docker update -c `+cpushares+` --cpu-quota `+cpuquota+` `+name).Output(); cpuError == nil {
			Logger("The container's cpus has been resized to:"+resizevalue, name)
		} else {
			log.Fatal("Failed to resize the cpus for container:"+name, cpuError.Error())
		}

	} else {
		log.Fatal("ResizeAppc's args has some error", "JFDocker resize app,cpu,4/app,mem/2048m")
	}

}

func RunGpuAppc(pauseid, name, appimage, cpus, mems, alloc_cnt string) (appcid string, err error) {
	var (
		alloc_gpus_id string //分配的gpu卡id
		alloc_gpus    int    //分配的gpu卡个数
	)
	cpushares, cpuquota := Cpuget(cpus) //获取CPU的share 和quota

	//转换gpu数量为int型 alloc_gpus
	if num, err := strconv.Atoi(alloc_cnt); err == nil {
		alloc_gpus = num
	}
	//获取当前已经在使用的GPU卡设备id
	used_ids := GetUsedGpu()
	//获取分配的GPU卡id信息[0 2] [4 5 6 7]
	allocGpu := GetFreeGpus(alloc_gpus, used_ids)
	//构造一个可以直接传入nvidia-docker的GPU卡id:0,2  or 4,5,6,7
	for i := 0; i < len(allocGpu); i++ {
		alloc_gpus_id = alloc_gpus_id + "," + allocGpu[i]

	}
	fmt.Println(alloc_gpus_id, cpushares, cpuquota)

	var appconid string
	//conspecconf := apis.LogPath+apis.ConConfDir+"/ConSpec"+name
	//后期是否需要挂载配置文件
  // 挂载CephFS相关目录，需要在-v apis.DataPath+name:/data
	//fmt.Println(NV_GPU="+alloc_gpus_id+" nvidia-docker run -itd  --name  "+name+" --net=container:"+pauseid+" --ipc=container:"+pauseid+" -m "+mems+" --cpu-shares "+cpushares+" --cpu-quota "+cpuquota+" --cpu-period 100000 -p 80:5000 "+appimage)
	if appcid, err := exec.Command("/bin/bash", "-c", `NV_GPU=`+alloc_gpus_id+` nvidia-docker run -itd  --name  `+name+` --net=container:`+pauseid+` --ipc=container:`+pauseid+` -m `+mems+` --cpu-shares `+cpushares+` --cpu-quota `+cpuquota+` --cpu-period 100000  `+appimage).Output(); err == nil {
		appconid = strings.Replace(string(appcid), "\n", "", -1)
	} else {
		fmt.Println("appcontainer 创建失败，请检查容器参数" + err.Error())
		fmt.Println(string(appcid))
		os.Exit(1)
	}

	return appconid, nil

}

func UpdateGpuAppc(name, image, cpus, mems, gpus string) (stat, appc string) {
	var pauseid, status, appcon string
	if out, err := exec.Command("/bin/bash", "-c", `docker inspect -f "{{.HostConfig.NetworkMode}}" `+name).Output(); err == nil {
		pauseid = strings.Split(strings.Split(string(out), "\n")[0], ":")[1]
	} else {
		log.Fatalf("Failed to get the pause container id", err.Error())
	}

	if stat, _ := DelAppc(name); stat == "ok" {
		Logger("The gpu container of lastimage has been delete", name)
		name, err := RunGpuAppc(pauseid, name, image, cpus, mems, gpus)
		if err != nil {
			log.Fatalf("The container of newimage run failed,please check it ", err.Error())
		}
		status = "ok"
		appcon = name
		Logger("The gpu container of newimage has been created", appcon)
	} else {
		log.Fatalf("The container of lastimage delete failed", "check it")
	}
	return status, appcon
}
