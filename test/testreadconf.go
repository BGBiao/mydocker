package main

import (
	"encoding/json"
	"fmt"
	"github.com/xxbandy/mydocker/apis"
	"github.com/xxbandy/mydocker/clients"
	"os"
	"strings"
)

func main() {
	//    fmt.Println(Readimage("test-xuxuebiao"))
	ReadConConf("test-xuxuebiao", "testimage")
}

//read runconf
func Readimage(name string) (image, cpu, mem string) {
	runConf := "runConf" + name
	args, err := clients.ReadFromConf(runConf)
	if err != nil {
		os.Exit(2)
	}
	images := strings.Split(args, ",")[2]
	cpus := strings.Split(args, ",")[7]
	mems := strings.Split(args, ",")[8]
	return images, cpus, mems
}

//read con conf

func ReadConConf(name, image string) {
	conContent, err := clients.ReadFromConf(name)
	if err != nil {
		os.Exit(1)
	}
	coninfo := apis.JFDocker{}
	condata := []byte(conContent)
	ConInfoerr := json.Unmarshal(condata, &coninfo)
	if ConInfoerr != nil {
		os.Exit(3)
	}

	coninfo.Image = image
	//构造新函数
	newargs := string(coninfo.Conspec.Ywid) + "," + string(coninfo.Conspec.SN) + "," + string(coninfo.Image) + "," + string(coninfo.Appname) + "," + string(coninfo.Netspec.Ipv4) + "," + string(coninfo.Netspec.Mask) + "," + string(coninfo.Netspec.Gateway) + "," + string(coninfo.Conspec.Cpus) + "," + string(coninfo.Conspec.Mem) + ",1232123," + string(coninfo.Netspec.Vnet)
	fmt.Println(string(coninfo.Image), newargs)

	//测试整个数据
	if conInfo, err := json.Marshal(coninfo); err == nil {
		fmt.Println(string(conInfo))
	}
}
