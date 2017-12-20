package main

import (
    _ "os/exec"
    "fmt"
    _ "encoding/json"
    _ "github.com/urfave/cli"
    "os"
    "github.com/xxbandy/mydocker/clients"
    "github.com/xxbandy/mydocker/apis"
    log "github.com/sirupsen/logrus"
)

func init() {
    _,err := os.Stat(apis.LogPath)
    //判断LogPath不存在则进行创建
    if err != nil {
        logpatherr := os.MkdirAll(apis.LogPath,0755)
        //如果创建成功继续创建容器目录
        if logpatherr == nil {
            if conpatherr := os.MkdirAll(apis.LogPath+"containers_config",0755);conpatherr == nil {
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
    if _,err := os.Stat(apis.LogPath+"containers_config"); err != nil {
        conpatherr := os.MkdirAll(apis.LogPath+"containers_config",0755)
        if conpatherr != nil { log.Fatalf("Failed to create container_conf dir:",conpatherr.Error()) }
    }
}

func main() {
    fmt.Println("测试目录是否创建成功")
}
