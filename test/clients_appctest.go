package main
import (
  "os"
  _ "os/exec"
  _ "fmt"
  _ "strings"
  _ "strconv"
  "github.com/xxbandy/mydocker/clients"
)

func main() {
    /*指定appname和image名称进行更新应用容器
    func UpdateAppc(name,image string) (stat,app string)
    */
    s,c := clients.UpdateAppc("test-xuxuebiao","e740f4a4a24d","4","1024m")
    if s == "ok"{
      clients.Logger("应用容器更新成功,容器id:",c)
    } else {
      clients.Logger("应用容器更新失败",s)
      os.Exit(1)
    }

    /* 指定appname进行容器销毁操作
    func DelAppc(name string) (status,appcid string)
    s,c := clients.DelAppc("test-biao")
    if s == "ok" {
        clients.Logger("应用容器销毁成功",c)
    } else {
        os.Exit(1)
    }
   */ 

    /*指定相关参数对容器进行创建
    func Appc(pauseid,name,appimage,cpus,mems string)  (appid string,err error)
    s,err := clients.Appc("589963836cba","test-biao","172.25.46.8:5001/sandbox-jdk7-tomcat6","4","1000m")
    if err != nil {
        clients.Logger("应用容器创建失败",err.Error())
        os.Exit(1)
    }
    clients.Logger("应用容器创建成功",s)
    */
}

