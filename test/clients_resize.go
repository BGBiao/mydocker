package main
import (
    _ "github.com/xxbandy/mydocker/apis"
    "github.com/xxbandy/mydocker/clients"

)

func main() {
    cname := "test-conf"
    resizetype := "mem"
    mems := "3000m"
    
    resizetype1 := "cpu"
    values := "4"    

    clients.ResizeAppc(cname,resizetype,mems)
    clients.ResizeAppc(cname,resizetype1,values)


}
