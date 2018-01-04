package main

import (
	"fmt"
	_ "github.com/xxbandy/mydocker/apis"
)

func main() {
	names := []string{"gpu1", "gpu2", "test-xuxuebiao"}
	for _, name := range names {
		fmt.Println(name)
	}

}
