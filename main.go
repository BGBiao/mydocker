package main

import (
  "fmt"
  jfdocker "mydocker/client"
  docker "github.com/fsouza/go-dockerclient"
)

func main() {
  //endpoint := "unix:///var/run/docker.sock"
  endpoint := "tcp://172.25.60.149:5256"
  api := "1.21"
  //client, err := mydocker.NewClient(endpoint)
  client, err := jfdocker.NewDockerClient(endpoint,api)
  if err != nil {
    panic(err)
  }
  imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
  if err != nil {
    panic(err)
  }
  for _, img := range imgs {
    fmt.Println("ID: ", img.ID)
    fmt.Println("RepoTags: ", img.RepoTags)
  }
}
