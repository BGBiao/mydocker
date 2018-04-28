package main
import (
    docker "github.com/fsouza/go-dockerclient"
    "time"
    "bytes"
    "log"
    "archive/tar"
    "io/ioutil"
)

func main() {
client, err := docker.NewClient("unix:///var/run/docker.sock")
if err != nil {
    log.Fatal(err)
}


imagename := "xxbandy.github.io/grafana"
//需要注意的是Dockerfile中的首行不能为注释,且默认是去公网读取镜像
//应该需要做一个认证就可以了
//dockerfile := "/export/Dockerfile-images/grafana/Dockerfile"
dockerfile := "/export/Dockerfile-images/GPU-images/jdjrbase/Dockerfile"

t := time.Now()
inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
tr := tar.NewWriter(inputbuf)
tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
filecontent,_ := ioutil.ReadFile(dockerfile)
//tr.Write([]byte("FROM nginx\n"))
tr.Write(filecontent)
tr.Close()

  buildopts := docker.BuildImageOptions{
        Name:           imagename,
        NoCache:        true,
        Pull:           true,
        RmTmpContainer: true,
        NetworkMode:    "host",
        InputStream:    inputbuf,
        OutputStream:   outputbuf,
    }

if err := client.BuildImage(buildopts); err != nil {
    log.Fatal(err)
}
}
