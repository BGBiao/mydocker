package main

import (
  "fmt"
  jfdocker "mydocker/client"
  docker "github.com/fsouza/go-dockerclient"
)

func main() {
  endpoint := "tcp://172.25.60.149:5256"
  api := "1.37"

  client, err := jfdocker.NewDockerClient(endpoint,api)
  if err != nil {
    panic(err)
  }

  imagename := "xxbandy123/k8s-pause"
  repo := "dockerhub.jd.com/jdjr/k8s-pause"
  tag := "18-04-20"

  tagopts := docker.TagImageOptions{Repo:repo,Tag:tag}
  //docker tag $imagename $repo:$tag
  tagErr := jfdocker.TagImages(client,imagename,tagopts)
  if tagErr != nil {
    panic(tagErr)
  }
  fmt.Println("tag ok")
}
