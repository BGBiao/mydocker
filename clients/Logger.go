package clients
import (
    "os"
    log "github.com/sirupsen/logrus"
    "github.com/xxbandy/mydocker/apis"
)

var logs = log.New()

//定义日志输出格式
func Logger(msg,data string) {

    //判断logpath是否存在,不存在进行创建
    _,err := os.Stat(apis.LogPath)
    if err != nil {
      fileerr := os.MkdirAll(apis.LogPath,0755)
      if fileerr != nil {
          log.Fatalf("Failed to create logpath %s. Error: %s", apis.LogPath,fileerr)
      }        
    }

    //创建日志文件并进行日志内容写入
    if file,err := os.OpenFile(apis.LogPath+"JFDocker.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666);err == nil {
        logs.Out = file
        defer file.Close()
        logs.WithFields(log.Fields{"data": string(data)}).Info(msg)
    } else {
        log.Fatalf("Failed to write logfile %s. Error:%s",apis.LogPath+"JFDocker.log",err)
    }

}
