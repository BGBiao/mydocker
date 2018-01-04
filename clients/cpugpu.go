package clients

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func Cpuget(cpu string) (cpushares, cpuquota string) {
	//需要注意下int类型的是否会超配
	var Share, Quota int
	if cpuint, err := strconv.Atoi(cpu); err == nil {
		Share = cpuint * 1000
		Quota = cpuint * 100000
	} else {
		fmt.Println("cpu 表示有误，请确定是cpu核心数量" + err.Error())
		os.Exit(1)
	}
	return strconv.Itoa(Share), strconv.Itoa(Quota)

}

func GetFreeGpus(alloc int, usedGpus string) (freeGpus []string) {
	var (
		midn            int      //定义GPU列表的中值以及中值在可用列表的下标
		midv            string   //定义freegpus 的slice中的中值
		gpus_pool       []string //获取宿主机的gpu卡个数并构造gpus卡池
		freeGpus_pool   []string //定义可分配的gpus卡池
		alloc_gpus_pool []string //按照分配的需求分配的gpus列表
	)

	//获取宿主机的GPU个数并构造GPU卡列表
	//gpu_pool = [0 1 2 3 4 5 6 7]
	allGpus := getGpus()
	for i := 0; i < allGpus; i++ {
		gpus_pool = append(gpus_pool, strconv.Itoa(i))
	}

	//对已经使用的gpusid进行排序并且构造gpu卡已使用资源池
	used_pool := strings.Split(usedGpus, ",")
	sort.Strings(used_pool)

	//根据gpus_pool 和used_pool 对比构造出freeGpus_pool,可分配的gpu卡池
	//同时构造出两个字符串变量midn(gpus_pool中值在freeGpus中的位置) 和midv(gpus_pool中值)
	if isslice, diff := checkSlice(used_pool, gpus_pool); isslice == false {
		//构造可用的GPU卡列表
		freeGpus_pool = diff
		//计算gpu卡列表中的中值在可用gpu卡列表中的位置
		for n, diffmid := range freeGpus_pool {
			//字符串转int
			if num, err := strconv.Atoi(diffmid); err == nil {
				//计算gpu卡中值以及中值在可用gpu列表的位置
				if num < allGpus/2 {
					midn = n
					midv = diffmid
				}
			}
		}
	}

	//获取当前可用gpu卡池中的数量
	freeGpuNums := len(freeGpus_pool)

	//全部GPU[0 1 2 3 4 5 6 7],midv是3，该值在可用GPU列表的位置为midn
	//如果在可用列表的前midn个元素中可以满足分配的卡，就从前半段分。如果前半段不够就判断后半段是否有足够的卡来分给用户，如果也不够，则退出，不能自动分配到近亲缘性的卡，如果够直接分配
	//alloc 传入需要的卡数量
	if alloc > freeGpuNums {
		fmt.Println("无可分配的GPU资源")
		os.Exit(1)
	}

	//当有8卡的时候,midn相当于gpu3在可用列表里的位置
	//当midn不等于0时，说明第一组GPU卡没有分配完
	//fmt.Println(midn)
	if midn > 0 {
		if alloc <= midn+1 {
			alloc_gpus_pool = freeGpus_pool[0:alloc]
		} else if len(freeGpus_pool[midn+1:]) >= alloc {
			alloc_gpus_pool = freeGpus_pool[midn+1 : midn+1+alloc]
		} else {
			alloc_gpus_pool = freeGpus_pool[midn-1 : midn-1+alloc]
			fmt.Println("没有亲缘性的GPU节点,请检查!如有需要请手动指定多跨CPU的卡设备", alloc_gpus_pool)
			os.Exit(2)
		}
		//midn = 0 midv 为中值GPU 即为3
	} else if midn = 0; freeGpus_pool[0] == midv {
		if alloc != 1 {
			//如果分配数量不等于1的话需要判断亲和性
			if md, _ := strconv.Atoi(midv); md < getGpus()/2 && len(freeGpus_pool[midn+1:]) >= alloc {
				alloc_gpus_pool = freeGpus_pool[midn+1 : midn+1+alloc]
			} else {
				alloc_gpus_pool = freeGpus_pool[midn : midn+alloc]
				fmt.Println("没有亲缘性的GPU节点,请检查!如有需要请手动指定多跨CPU的卡设备", alloc_gpus_pool)
				os.Exit(1)
			}
		} else {
			alloc_gpus_pool = freeGpus_pool[midn : midn+alloc]
		}
	} else {
		alloc_gpus_pool = freeGpus_pool[midn : midn+alloc]
	}

	//fmt.Println("分配到的GPU信息如下:",alloc_gpus_pool)
	return alloc_gpus_pool
}

//获取物理机的GPU个数
func getGpus() (gpu int) {

	gpus, err := exec.Command("/bin/bash", "-c", `nvidia-smi -L | wc -l`).Output()
	if err != nil {
		fmt.Println("Failed to get the gpu info with nvidia-smi command", err.Error())
		os.Exit(1)
	}
	gpunums, err := strconv.Atoi(strings.Replace(string(gpus), "\n", "", -1))
	if err != nil {
		os.Exit(1)
	}
	return gpunums

}

//对比已使用GPU卡以及全部GPU卡列表，计算出可分配的GPU列表
func checkSlice(a, b []string) (isIn bool, c []string) {
	for _, valueOfb := range b {
		temp := valueOfb
		for j := 0; j < len(a); j++ {
			if temp == a[j] {
				break
			} else {
				if len(a) == (j + 1) {
					c = append(c, temp)
				}
			}

		}
	}
	if len(c) == 0 {
		isIn = true
	} else {
		isIn = false
	}
	return isIn, c
}

func GetUsedGpu() (userid string) {
	//定义设备列表,通过docker inspect获取所有容器的gpu设备列表并拼接设备名称/dev/nvidia0,/dev/nvidia1,/dev/nvidia3
	var allUsedDevice string
	used_ids := ","
	//获取本地的docker容器列表(gpu卡都是使用docker container创建)
	GetConName, err := exec.Command("/bin/bash", "-c", `docker ps --format="{{.Names}}"`).Output()
	if err != nil {
		fmt.Println("no container")
		return used_ids
	}
	GetList := strings.Split(string(GetConName), "\n")
	ConNameList := GetList[0 : len(GetList)-1]

	//获取已经在使用的gpu列表信息
	for _, name := range ConNameList {
		//获取容器的挂载设备返回带换行符的[]byte
		GetDeivceOfName, err := exec.Command("/bin/bash", "-c", `docker inspect -f "{{.HostConfig.Devices}}" `+name).Output()
		if err != nil {
			fmt.Println("no get container device")
			return used_ids
		}
		//GetDeivceOfName [{/dev/nvidiactl /dev/nvidiactl rwm} {/dev/nvidia-uvm /dev/nvidia-uvm rwm} {/dev/nvidia-uvm-tools /dev/nvidia-uvm-tools rwm} {/dev/nvidia5 /dev/nvidia5 rwm} {/dev/nvidia6 /dev/nvidia6 rwm}]
		devicestrings := strings.Replace(string(GetDeivceOfName), "\n", "", -1)
		//过滤普通的docker容器(设备列表为kong)
		if devicestrings != "[]" {
			var usedDevice string
			if re := regexp.MustCompile("k8s*"); re.MatchString(name) == true {
				usedDevice = getK8sGpuDevice(name, devicestrings)
			} else {
				usedDevice = getDockerGpuDevice(name, devicestrings)
			}
			allUsedDevice = allUsedDevice + usedDevice
		}

	}

	//从gpu设备列表中获取设备编号/dev/nvidia0,/dev/nvidia1,/dev/nvidia3 -> 0,1,3
	for _, id := range strings.Split(allUsedDevice, ",") {
		used_ids = used_ids + strings.Replace(id, "/dev/nvidia", ",", -1)
	}
	return used_ids

}

//return ,/dev/nvidia0,/dev/nvidia1
func getDockerGpuDevice(name, devices string) (useddevice string) {
	var gpus []string
	devicelist := strings.Split(devices, " ") //[]string
	gpulist := devicelist[9:len(devicelist)]
	for i := 1; i < len(gpulist); i = i + 3 {
		gpus = append(gpus, gpulist[i])
	}

	var alldevice string
	for i := 0; i < len(gpus); i = i + 1 {
		alldevice = alldevice + "," + gpus[i]

	}

	return alldevice
}

func getK8sGpuDevice(name, devices string) (useddevice string) {
	var gpus []string
	devicelist := strings.Split(devices, " ")
	gpulist := devicelist[0 : len(devicelist)-9]
	for i := 1; i < len(gpulist); i = i + 3 {
		gpus = append(gpus, gpulist[i])
	}

	var alldevice string
	for i := 0; i < len(gpus); i = i + 1 {
		alldevice = alldevice + "," + gpus[i]

	}

	return alldevice
}
