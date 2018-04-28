package main

import (
  "fmt"
  "bytes"
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

  imagename := "xxbandy.github.io/grafana"
  dockerfile := "/export/Dockerfile-images/grafana/Dockerfile"


/*
inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
tr := tar.NewWriter(inputbuf)
tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
tr.Write([]byte("FROM base\n"))
tr.Close()
*/

  outputbuf := bytes.NewBuffer(nil)
  buildopts := docker.BuildImageOptions{
        Name:           imagename,
        Dockerfile:     dockerfile,
        NoCache:        true,
        Pull:           true,
        RmTmpContainer: true,
        NetworkMode:    "host",
        //InputStream:    inputbuf,
        OutputStream:   outputbuf,
    }
  //docker build -t $imagename -f $dockerfile 
  tagErr := jfdocker.BuildImages(client,buildopts)
  if tagErr != nil {
    fmt.Println(tagErr)
    panic(tagErr)
  }
  fmt.Println("tag ok")
}
