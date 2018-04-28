package main

import (
  "fmt"

  mydocker "github.com/fsouza/go-dockerclient"
)

func main() {
  endpoint := "unix:///var/run/docker.sock"
  client, err := mydocker.NewClient(endpoint)
  if err != nil {
    panic(err)
  }
  imgs, err := client.ListImages(mydocker.ListImagesOptions{All: false})
  if err != nil {
    panic(err)
  }
  for _, img := range imgs {
    fmt.Println("ID: ", img.ID)
    fmt.Println("RepoTags: ", img.RepoTags)
    fmt.Println("Created: ", img.Created)
    fmt.Println("Size: ", img.Size)
    fmt.Println("VirtualSize: ", img.VirtualSize)
    fmt.Println("ParentId: ", img.ParentID)
  }
}
