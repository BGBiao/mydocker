package clients

import (
    "github.com/xxbandy/CephFSOps/fsclient"
    "github.com/xxbandy/mydocker/apis"
    log "github.com/sirupsen/logrus"
    "strings"
    "os/exec"
    "fmt"
    "os"

) 
/*
GetPathName("baoge").Del()
GetPathName("baoge").Mkdirall()
GetPathName("baoge").String()

*/

type MountPoint string
func getPathName(name string) MountPoint {
    dataPath := apis.DataPath+name
    return MountPoint(dataPath)
}

func (dir MountPoint) String() string {
    return string(dir)
}

func (dir MountPoint) Mkdirall() (stat string) {
    var stats,path string
    path = dir.String()
    _,err := os.Stat(path)
    if err != nil {
        if createErr := os.MkdirAll(path, 0755); createErr == nil {
            Logger("CephFS挂载点已经创建",path)
            stats = "yes"
        } else {
            stats = "no"
        }
    } else {
        Logger("数据目录已经存在",path)
        stats = "yes"
    }

    return stats
}

func (dir MountPoint) Umount() error {
    _,umountErr := exec.Command("/bin/bash","-c",`umount `+string(dir)).Output()
    if umountErr != nil {
        log.Fatalf("Cephfs 卸载失败:"+string(dir),umountErr.Error())
    }
    
    return umountErr
}






// cephfs name
func CreateMountfs(name string) error {
    var mountArgs string
    // 使用MountPoint类型的方法
    Logger("创建cephfs挂载点:",string(getPathName(name).Mkdirall()))
    

    // 向CephFS申请默认大小的cephfs,如果创建成功使用该cephfs的uuid进行挂载cephfs
    uuid,_ := fsclient.CreateFS(name)
    
    if mount,parseArgsErr := fsclient.GetFsMsg(uuid);parseArgsErr != nil { 
        log.Fatalf("failed to parse the cephfs's mount args",parseArgsErr.Error()) 
    }else {
        mountArgs = strings.Replace(strings.Replace(strings.Replace(mount,"{your_mount_point}",apis.DataPath+name,1),"{your key}",fsclient.GetSecretkey(),1),"\n","",-1)
    
    }
    fmt.Println(mountArgs)
    // 使用上述格式化的mountArgs进行挂载文件系统
    execOut,err := exec.Command("/bin/bash","-c",` `+mountArgs).Output()
    /*
    if err != nil { 
        fmt.Println("failed to mount the cephfs ,deleting the cephfs because of: ",err.Error())
        if delFsErr := fsclient.DeleteCephfs(uuid);delFsErr == nil {
            fmt.Printf("Delete the cephfs %s with uuid:%s\n",name,uuid)
            os.Exit(0)
        }else {  
            log.Fatalf("failed to delete the cephfs with: ",delFsErr.Error()) 
        }

    }
    fmt.Printf("Successful to mounts the fs for %s %s",name,string(execOut))
    */
    // 挂载失败需要将cephfs直接删除掉以保证整体操作的原子性
    if err != nil {
        fmt.Println("delete the cephfs:"+uuid,fsclient.DeleteCephfs(uuid).Error())
        fmt.Println(string(execOut))
        log.Fatalf(err.Error())
    }
    return nil
}
