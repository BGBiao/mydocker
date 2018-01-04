package clients

import (
	"encoding/json"
	_ "fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xxbandy/mydocker/apis"
	"os"
	"os/exec"
	"strings"
	"time"
)

//pausec 其实可以写成一个方法，即拥有了一个pauseid 之后可以对该pause container直接做DetachNet和AttachNet

func Pausec(name, ip, mask, gw, vnet string) (pauseid string, err error) {
	timenow := time.Now().Format("2006-01-02-15-04-05")
	name = name + string(timenow)
	//docker run -d --net=none --name name+time pauseimage
	//containerid返回的是一个[]byte类型的，并且后面带有空格
	containerid, err := exec.Command("/bin/bash", "-c", `docker run -d --net=none --name `+name+` `+apis.Pauseimage).Output()
	if err != nil {
		log.Fatalf("Failed to create pause container for "+name, err.Error())
		//fmt.Println(err.Error())
	}
	//pause容器创建成功，需要添加网络相关信息
	netns := apis.Netnamespace
	//截取pause容器的id
	pausecid := strings.Replace(string(containerid), "\n", "", -1)
	Logger("Susseeded to create pause container for  "+name, netns+" id:"+pausecid)

	//fmt.Println("pause容器命名空间以及容器id:"+netns,pausecid)
	//fmt.Println("开始装配pause容器网络:")
	Logger("Creating the pause container netconfig:", pausecid)

	//开始装配pause容器的网络
	if netresult, err := AttachNet(netns, pausecid, ip, mask, gw, vnet); err == nil {
		if result, err := json.Marshal(netresult); err == nil {
			//容器网络装配成功，打印pause容器相关基本信息
			//fmt.Println("pause容器网络基本信息:"+string(result))
			Logger("Successful to create network for pause container", string(result))
		} else {
			log.Fatalf("Failed to json pause container network conf", err.Error())
			//fmt.Println("pause容器网络信息解析失败:"+err.Error())
		}
		return pausecid, err
	} else {
		log.Fatalf("Failed to create network for pause container", err.Error())
		//fmt.Println("pause容器网络装配失败:"+err.Error())
		//os.Exit(1)
	}

	return pausecid, err
}

func DelPausec(pauseid, appname string) {
	//path在删除文件的时候需要指定绝对路径
	path := apis.LogPath + apis.ConConfDir + "/"
	runConf := "runConf" + appname
	conSpecConf := "ConSpec" + appname
	conConf := appname

	runArgs, err := ReadFromConf(runConf)
	if err != nil {
		log.Fatalf("Failed to get the runConf for container"+appname, err.Error())
	}

	argslist := strings.Split(runArgs, ",")
	ip := argslist[4]
	mask := argslist[5]

	//摘除pause容器网络
	delNetErr := DetachNet(pauseid, ip, mask)
	if delNetErr != nil {
		log.Fatalf("Failed to detach net", err.Error())
	}
	//销毁pause容器(其实网络摘取之后pause容器销毁可以不用)
	pausecontainer, err := exec.Command("/bin/bash", "-c", `docker stop `+pauseid+` && docker rm -f -v `+pauseid).Output()
	if err != nil {
		log.Fatalf("Failed to remove pause container", err.Error())
	}
	Logger("Successful to detach and delete pause container:"+appname, string(pausecontainer))

	Logger("Waiting to delete some conf", appname)

	if conSpecConfErr := os.Remove(path + conSpecConf); conSpecConfErr == nil {
		if conConfErr := os.Remove(path + conConf); conConfErr == nil {
			if runConfErr := os.Remove(path + runConf); runConfErr == nil {
				Logger("Successful to delete all the conf for "+appname, "")
			} else {
				log.Fatalf("Failed to delete container conf: "+path+runConf, runConfErr.Error())
			}
		} else {
			log.Fatalf("Failed to delete container conf: "+path+conConf, conConfErr.Error())
		}

	} else {
		log.Fatalf("Failed to delete container conf: "+path+conSpecConf, conSpecConfErr.Error())
	}

}
