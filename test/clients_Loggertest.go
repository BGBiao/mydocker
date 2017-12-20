package main
import (
    "os"
    "fmt"
    "encoding/json"
    "github.com/xxbandy/mydocker/clients"
)

type Test struct {
    Name    string `json:name`
    Sex     string `json:sex`
}


func main() {
    data := Test{}
    data.Name = "xxbandy"
    data.Sex = "male"


    datajson,err := json.Marshal(data)
    if err != nil {
        fmt.Println("error"+err.Error())
        os.Exit(1)
    }
    clients.Logger("/export/Logs/","test msg",string(datajson))
}

