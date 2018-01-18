package main
import (
    "github.com/xxbandy/mydocker/clients"
    _ "github.com/xxbandy/CephFSOps/fsclient"
    "fmt"
)


func main() {
    // 指定名称进行创建cephfs目录并挂载到相关目录
    name := "xuxuebiao1"
    if err := clients.CreateMountfs(name); err == nil {
        fmt.Println("mount ok")
    }
    

    /*
    uuid := "75a04ae0-f41d-11e7-81da-f000ac192cec"
    if err := fsclient.DeleteCephfs(uuid); err == nil {
        fmt.Println("delete ok")
    } else { fmt.Println(err.Error()) }
    */
}
