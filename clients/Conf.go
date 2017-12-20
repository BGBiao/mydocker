package clients
import (
    "os"
    "io/ioutil"
    log "github.com/sirupsen/logrus"
    "github.com/xxbandy/mydocker/apis"
)

//编写容器运行时配置文件
/*
文件读写操作
func (f *File) Write(b []byte) (n int, err error)
func (f *File) WriteAt(b []byte, off int64) (n int, err error)
func (f *File) WriteString(s string) (n int, err error)

将文件读入到[]byte中，当n为0时表示文件读取完成
func (f *File) Read(b []byte) (n int, err error)

使用ioutil包进行文件内容读取
ioutil.FeadFile
ioutil.ReadAll

bufio.ReadString
*/

//write and read runConfname and container_name
//name: runConfname
func WriteToConf(name,content string) (err error){
    fileObj,err := os.OpenFile(apis.LogPath+apis.ConConfDir+"/"+name,os.O_CREATE|os.O_WRONLY|os.O_TRUNC,0644)
    if err != nil {
        log.Fatalf("Failed to write ContainerRunConf:",err.Error())
    }

    _,WriteErr := fileObj.WriteString(content)
    if WriteErr != nil {
        log.Fatalf("Failed to write the ConRunConf "+apis.LogPath+apis.ConConfDir+"/"+name+" :",err.Error())
    }
    defer fileObj.Close()
    return WriteErr
}


func ReadFromConf(name string) (runargs string,err error) {
    var args string
    content, err := ioutil.ReadFile(apis.LogPath+apis.ConConfDir+"/"+name)
    if err != nil {
        log.Fatalf("Failed to read ContainerRunConf: "+apis.LogPath+apis.ConConfDir+"/"+name,err.Error())
    }
    args = string(content)
    return args,err
}



