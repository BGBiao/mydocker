package main
import (
    _ "github.com/xxbandy/mydocker/apis"
    "fmt"
)

func main() {
    names := []string{"gpu1","gpu2","test-xuxuebiao"}
    for _,name := range names {
      fmt.Println(name)
    }


}
