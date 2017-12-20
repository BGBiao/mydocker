package clients 
import (
    "strconv"
    "fmt"
    "os"

)

func Cpuget(cpu string) (cpushares,cpuquota string) {
    //需要注意下int类型的是否会超配
    var Share,Quota int
    if cpuint,err := strconv.Atoi(cpu);err == nil {
          Share = cpuint * 1000
          Quota = cpuint * 100000
    } else {
          fmt.Println("cpu 表示有误，请确定是cpu核心数量"+err.Error())
          os.Exit(1)
    }
    return strconv.Itoa(Share),strconv.Itoa(Quota)

}

