package main

import (
    _ "os/exec"
    "fmt"
    "encoding/json"
    _ "github.com/urfave/cli"
    "os"
    "github.com/xxbandy/mydocker/clients"
    "github.com/xxbandy/mydocker/apis"
    log "github.com/sirupsen/logrus"
)

var Usage string

//init函数当前用来初始化相关的目录结构
func init() {
    _,err := os.Stat(apis.LogPath)
    //判断LogPath不存在则进行创建
    if err != nil {
        logpatherr := os.MkdirAll(apis.LogPath,0755)
        //如果创建成功继续创建容器目录
        if logpatherr == nil {
            if conpatherr := os.MkdirAll(apis.LogPath+apis.ConConfDir,0755);conpatherr == nil {
                //fmt.Println("容器配置目录创建成功")
                clients.Logger("容器日志配置相关目录创建成功",apis.LogPath)
            } else {
                log.Fatalf("Failed to create container_conf dir:",conpatherr.Error())
            }


        } else {
            log.Fatalf("Failed to create logpath dir:",logpatherr.Error())
        }
    }
    //判断容器配置目录不存在并且进行创建
    if _,err := os.Stat(apis.LogPath+apis.ConConfDir); err != nil {
        conpatherr := os.MkdirAll(apis.LogPath+apis.ConConfDir,0755)
        if conpatherr != nil { log.Fatalf("Failed to create container_conf dir:",conpatherr.Error()) }
    }
}

func main() {
    Usage = os.Args[0]+" run ywid,sn,image,container,ip,mask,gateway,cpu,mem,label,vnet\n"+os.Args[0]+" update container,images\n"+os.Args[0]+" resize container,cpu,4/container,mem,2048\n"+os.Args[0]+" delete container\n"+os.Args[0]+" rebuilt container"

    //获取JFDocker所有的参数
    args := os.Args
    
    if len(args) == 1 {
        fmt.Println(Usage)
        os.Exit(1)
    }

    /*每个操作类型需要执行的参数:
    JFDocker run YW-D-TPBQD2@af41f7,4Y8K742@af41f7,172.25.46.9:5001/centos6.8-jdjr-test-app,172.25.47.21.h.chinabank.com.cn,172.25.47.21,24,172.25.47.254,2,4096m,20007772,br0
    JFDocker update appname,imagername
    JFDocker resize appname,cpu,4
    JFDocker delete appname
    JFDocker rebuilt appnem
    */

    //操作类型:run resize update delete rebuilt
    JFDockertype := args[1]

    //获取所有操作参数
    JFDockerArgs := args[len(args)-1]

    //JFDocker run操作
    if JFDockertype == "run" {
        //开始运行JFDocker run
        clients.Logger("<JFDockerRun:解析并进行容器创建操作",JFDockerArgs)
        
        //JFDockerRun(args)将传run操作的所有参数解析并进行容器创建，最终返回容器对象的结构体RspJFDocker
        result,err := clients.JFDockerRun(JFDockerArgs)
        if err != nil || result.Result != 0 {
            cerr := apis.JFDockererr{}
            cerr.JFDockerVersion = "0.0.1"
            if result != nil {
                cerr.Code = result.Result
                cerr.Msg = "JFDocker:" + result.ErrMsg
                cerr.Details = result.ErrInfo
            } else {
                cerr.Code = 1
                cerr.Msg = "JFDocker:" + err.Error()
            }
            eOut, err := json.Marshal(&cerr)
            if err == nil {
                fmt.Printf("%s", eOut)
            } else {
                fmt.Println(err)
            }
            os.Exit(1)
        }

        if data,err := json.MarshalIndent(result,"","  "); err == nil {
            fmt.Println(string(data))
        }

/*
        //本身上层接口返回的是interface接口类型，可以直接返回相关的数据，如果需要重新构造可以再次初始化结构体接口
        out := apis.RspJFDocker{}
        //问题:无法解析Result中的内容
        out.Result = result.Result
        out.Appname = result.Appname
        out.Ipv4 = result.Ipv4
        out.ConID = result.ConID

        data, err := json.MarshalIndent(out, "", "    ")
        if err != nil{
            fmt.Println(err)
            return
        }
        fmt.Printf(fmt.Sprintf("%s", data))
        
*/

    } else if JFDockertype == "resize" {
        fmt.Println("Resize  the docker container's cpu or mem ")
        clients.JFDockerResize(JFDockerArgs)
    } else if JFDockertype == "update" {
        //使用新镜像更新容器内容
        clients.Logger("<JFDockerUpdate:解析参数并进行容器更新",JFDockerArgs) 
        //func JFDockerUpdate(args string) (*apis.JFDocker,error) 
        result,err := clients.JFDockerUpdate(JFDockerArgs)
        if err != nil { os.Exit(2) }
        if data,err := json.Marshal(result); err == nil {
            clients.Logger("Successful to update the app container with JFDockerUpdate",string(data))
        }

    } else if JFDockertype == "delete" {
        //删除容器相关的信息 
        //JFDocker delete appname 
       
        result,err := clients.JFDockerDelete(JFDockerArgs)
        if err != nil {
            log.Fatalf("Failed to delete the container","")
        }
        if data,err := json.MarshalIndent(result,"","  "); err == nil {
            fmt.Println(string(data))
        }
    } else if JFDockertype == "rebuilt" {
        //rebuilt container appcontainer
        if rebuiltErr := clients.JFDockerRebuilt(JFDockerArgs); rebuiltErr == nil {
            clients.Logger("Rebuilt container done!","")
        }

    } else {
        fmt.Println("The operations has no opstype,please check it!")
        fmt.Println(Usage)
    }

}


