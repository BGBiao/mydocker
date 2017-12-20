##
###
### container.go
整个JFDocker容器引擎的核心操作，用来定义容器的多个主操作流程，比如run,delete,update,resize,rebuilt等相关操作，同时进行一些容器信息的格式化以及存储操作

todo:
1.联调多个主操作流程的接口
2.需要尽快将容器相关信息进行落地存储
3.需要将相关输出进行日志存储

### appc.go
- [x] done
主要负责应用容器的创建，销毁，以及更新
`注意:其中Appc在创建应用容器的时候会去调用cpugpu.go中的Cpuget(cpu)去换算相对cpu权重`
//由于docker --hostname 和--net参数的冲突不能同时进行设置，暂时无法设置主机名
func Appc(pauseid,name,appimage,cpus,mems string)  (appcid string,err error)
func DelAppc(name string) (stat,appcid string)
func UpdateAppc(name,image string) (stat,appc string)

todo:
需要对每个逻辑函数进行周全判断


### cpugpu.go
主要是负责分配cpu和gpu逻辑信息，目前只进行cpu的分配逻辑计算
`注意:当前分配策略是按照cpu的时间片进行分配给每个容器cpu计算资源的，通过用户要求的cpu核心数，来来计算cpushares cpuquota两个值`
func Cpuget(cpu string) (cpushares,cpuquota string)

todo:
需要根据用户对gpu的需求来指定gpu卡
func GpuGet()

### pausec.go
主要用来封装pause容器的相关事情，一方面创建pause容器的基本信息，另外一方面调用`net.go`中的AttachNet进行pause容器的网络封装操作，返回给用户pause容器的id
func Pausec(name,ip,mask,gw,vnet string) (pauseid string,err error)

todo:
需要在写一个pause容器删除的函数,该函数首先去操作`net.go`中的DettachNet进行pause容器网络的释放操作，同时进行pause容器的销毁操作
func DelPausec()

### net.go
主要是对容器网络进行相关的操作。目前是使用bridge的方式进行网络构建，通过`ip netns`相关的命令进行容器装配操作
func AttachNet(netns,containerid,ip,mask,gw,vnet string) (*apis.Netns, error)

todo:
需要编写DettachNet相关函数
func DettachNet()

### Logger.go

该文件主要用来定义日志主函数，需要输出和其他相关输出的可以直接使用
该函数默认是.Info()级别的日志
func Logger(msg,data string)

todo:
定义.Error()级别的日志，在日志文件中区分错误以及警告日志



### Conf.go(done)
该文件主要用于定义两个个函数进行配置文件的读取以及写入，配置文件目前主要分两个:container_containername和runConfname，前者用来存储容器的规格，主要是ywid,sn,cpu,mem信息，后者是用来存储容器run参数。因此两个文件的读写操作可以共享以下函数
func WriteToConf(name,content string) (err error)
func ReadFromConf(name string) (runargs string,err error)

`注意:读写容器运行时文件比较简单，只是对字符串的读取. 读写容器配置文件需要注意几点:`
- 写入容器配置文件时需要按照定义的结构体进行内容写入，并使用json.Marshal格式化成字符串写入文件
- 读取容器配置文件时需要将读取的string内容转换成[]byte，并使用json.Unmarshal(b []byte,s interface{}) 将需要读取的内容写入到实际定义的结构体中

### bug list 
场景:clients/container.go 
142 行，数据解析格式正常，但是使用WriteToConf进行字符串文件写入后发现不能对源文件进行重写，可能造成新旧数据混合。需要查查为什么当前的WriteToConf方法为什么会存在该问题。

问题:clients/Conf.go
28 行,func WriteToConf(name,content string) 使用的是文件对象的WriteString(string)方法

