package main
import (
    "github.com/xxbandy/mydocker/clients"
    "github.com/xxbandy/mydocker/apis"
    "encoding/json"
    log "github.com/sirupsen/logrus"

)



func main() {

    //读写容器运行时配置文件
    /*
    runname := "runConfxxb"
    content := "1236,1235,172.25.46.9:5001/centos6.5-test-vm-app-biao,woman2,172.25.66.222,24,172.25.66.254,4,4096m,20008662,br0"
    //写入容器运行时配置文件
    RunConferr := clients.WriteToConf(name,content)
    if RunConferr == nil {
        clients.Logger("容器运行时文件写入成功",name)
    }
    
    //读取容器运行的配置文件
    args,err := clients.ReadFromConf(name)
    if err == nil {
        clients.Logger("容器:"+name+" 运行参数读取成功",args)
    }
    */


    //读写容器配置文件
    conname := "container_xxb"
    containerInfo := apis.ConSpec{}
    /*
    //写入容器配置文件(ywid,sn,mem,cpu)
    containerInfo.Ywid = "1236"
    containerInfo.SN = "sn12345"
    containerInfo.Mem = "2048m"
    containerInfo.Cpus = "4"
    if conInfo,err := json.Marshal(containerInfo);err == nil {
        if ConConferr := clients.WriteToConf(conname,string(conInfo));ConConferr == nil {
            clients.Logger("容器配置文件写入成功",string(conInfo)) 
        }
    } else {
        log.Fatalf("Failed to json the container conf:",err.Error())
    }
    */

    //读取容器配置文件(读取ywid,sn,mem,cpu信息到apis.ConSpec{}中)
    //读取的conContent是字符串信息
    conContent,err := clients.ReadFromConf(conname)
    if err == nil {
        condata := []byte(conContent) 
        if ConInfoerr := json.Unmarshal(condata,&containerInfo);ConInfoerr == nil {
            clients.Logger("容器配置文件读取成功运维id:",string(containerInfo.Ywid))
        } else {
            log.Fatalf("Failed to unjson the ConInfo:",ConInfoerr.Error())   
        }
    }
    if conInfo,err := json.Marshal(containerInfo);err == nil {
        clients.Logger("读取的容器配置内容",string(conInfo))
    }
}


