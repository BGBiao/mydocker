## 


### env
```
$ go version
go version go1.8.3 linux/amd64

$ docker version 
Client:
 Version: 18.03.0-ce
 API version: 1.37
 Go version:  go1.9.4
Server:
 Engine:
  Version:  18.03.0-ce
  API version:  1.37 (minimum version 1.12)
  Go version: go1.9.4

# docker1.12.x版本详情
Client:
 Version:      1.12.1
 API version:  1.24
 Go version:   go1.6.3

# docker1.13.x版本详情
Client:
 Version:      1.13.1
 API version:  1.26
 Go version:   go1.7.5

# docker1.9.x 版本详情(需要特殊的dockerclient版本Docker 1.9 and Go 1.4)
Client:
 Version:      1.9.1
 API version:  1.21

# https://github.com/fsouza/go-dockerclient/ 最新版为1.2.0 
当前使用golang1.8构建


# 本地开发环境
# docker version
Client:
 Version:      1.12.1
 API version:  1.24
 Go version:   go1.6.3
 Git commit:   23cf638

Server:
 Version:      1.12.1
 API version:  1.24
 Go version:   go1.6.3
 Git commit:   23cf638

# go version
go version go1.9.2 linux/amd64
# go run list-image.go
ID:  sha256:c8c29d842c09d6c61f537843808e01c0af4079e9e74079616f57dfcfa91d4e25
RepoTags:  [nginx:1.9]
Created:  1464063409
Size:  182722590
VirtualSize:  182722590
ParentId:
```

## go-dockerclient apis

客户端初始化
```
type Client struct {
    SkipServerVersionCheck bool
    HTTPClient             *http.Client
    TLSConfig              *tls.Config
    Dialer                 Dialer
    // contains filtered or unexported fields
}

//将使用最新的API
func NewClient(endpoint string) (*Client, error)

//使用指定的apiversion构建客户端
func NewVersionedClient(endpoint string, apiVersionString string) (*Client, error)

//从默认的环境中读取并创建客户端实例 DOCKER_HOST, DOCKER_TLS_VERIFY, and DOCKER_CERT_PATH.
//https://github.com/moby/moby/blob/1f963af697e8df3a78217f6fdbf67b8123a7db94/docker/docker.go#L68
func NewClientFromEnv() (*Client, error)
//同上，但是可以指定apiversion
func NewVersionedClientFromEnv(apiVersionString string) (*Client, error)

```


### 客户端相关方法

```
//添加新的listen用来监听docker的事件信息，主要通过channel进行事件通讯
func (c *Client) AddEventListener(listener chan<- *APIEvents) error

//attach到容器内部
func (c *Client) AttachToContainer(opts AttachToContainerOptions) error

//构建镜像
func (c *Client) BuildImage(opts BuildImageOptions) error


//从容器拷贝内容
func (c *Client) CopyFromContainer(opts CopyFromContainerOptions) error

//创建配置
func (c *Client) CreateConfig(opts CreateConfigOptions) (*swarm.Config, error)

//创建容器
func (c *Client) CreateContainer(opts CreateContainerOptions) (*Container, error)

//exec到容器,容器运行参数必须是-id
func (c *Client) CreateExec(opts CreateExecOptions) (*Exec, error)


//创建volume
func (c *Client) CreateVolume(opts CreateVolumeOptions) (*Volume, error)

//创建插件
func (c *Client) CreatePlugin(opts CreatePluginOptions) (string, error)

//禁用插件
func (c *Client) DisablePlugin(opts DisablePluginOptions) error

//启用插件
func (c *Client) EnablePlugin(opts EnablePluginOptions) error

//返回当前的endpoint
func (c *Client) Endpoint() string


//返回系统级别的信息
func (c *Client) Info() (*DockerInfo, error)

//解析容器配置信息
func (c *Client) InspectContainer(id string) (*Container, error)

//解析镜像相关信息
func (c *Client) InspectImage(name string) (*Image, error)

//部署插件
func (c *Client) InstallPlugins(opts InstallPluginOptions) error


//停止容器
func (c *Client) KillContainer(opts KillContainerOptions) error

//列出容器
func (c *Client) ListContainers(opts ListContainersOptions) ([]APIContainers, error)

//列出镜像
func (c *Client) ListImages(opts ListImagesOptions) ([]APIImages, error)

//ping docker服务
func (c *Client) Ping() error


//下载镜像
func (c *Client) PullImage(opts PullImageOptions, auth AuthConfiguration) error

//上传镜像
func (c *Client) PushImage(opts PushImageOptions, auth AuthConfiguration) error

//删除容器
func (c *Client) RemoveContainer(opts RemoveContainerOptions) error

//删除镜像
func (c *Client) RemoveImage(name string) error

//重命名容器
func (c *Client) RenameContainer(opts RenameContainerOptions) error

//调整容器的tty
func (c *Client) ResizeContainerTTY(id string, height, width int) error

//重启容器
func (c *Client) RestartContainer(id string, timeout uint) error

//设置超时时间
func (c *Client) SetTimeout(t time.Duration)


//启动容器
func (c *Client) StartContainer(id string, hostConfig *HostConfig) error


//使用上下文启动容器(docker1.10.x支持hostConfig 1.12.x之后不支持)
func (c *Client) StartContainerWithContext(id string, hostConfig *HostConfig, ctx context.Context) error

//查看容器状态
func (c *Client) Stats(opts StatsOptions) (retErr error)

//停止容器
func (c *Client) StopContainer(id string, timeout uint) error

//镜像标签
func (c *Client) TagImage(name string, opts TagImageOptions) error


//更新容器配置
func (c *Client) UpdateConfig(id string, opts UpdateConfigOptions) error

//更新容器
func (c *Client) UpdateContainer(id string, opts UpdateContainerOptions) error

//查看版本
func (c *Client) Version() (*Env, error)


```

### 核心的结构体

```
// 延迟资源释放的接口
type CloseWaiter interface {
    io.Closer
    Wait() error
}

//容器提交的选项 CommitContainer方法的选项
type CommitContainerOptions struct {
    Container  string
    Repository string `qs:"repo"`
    Tag        string
    Message    string `qs:"comment"`
    Author     string
    Changes    []string `qs:"changes"`
    Run        *Config  `qs:"-"`
    Context    context.Context
}

//创建容器的配置信息
type Config struct {
    Hostname          string              `json:"Hostname,omitempty" yaml:"Hostname,omitempty" toml:"Hostname,omitempty"`
    Domainname        string              `json:"Domainname,omitempty" yaml:"Domainname,omitempty" toml:"Domainname,omitempty"`
    User              string              `json:"User,omitempty" yaml:"User,omitempty" toml:"User,omitempty"`
    Memory            int64               `json:"Memory,omitempty" yaml:"Memory,omitempty" toml:"Memory,omitempty"`
    MemorySwap        int64               `json:"MemorySwap,omitempty" yaml:"MemorySwap,omitempty" toml:"MemorySwap,omitempty"`
    MemoryReservation int64               `json:"MemoryReservation,omitempty" yaml:"MemoryReservation,omitempty" toml:"MemoryReservation,omitempty"`
    KernelMemory      int64               `json:"KernelMemory,omitempty" yaml:"KernelMemory,omitempty" toml:"KernelMemory,omitempty"`
    CPUShares         int64               `json:"CpuShares,omitempty" yaml:"CpuShares,omitempty" toml:"CpuShares,omitempty"`
    CPUSet            string              `json:"Cpuset,omitempty" yaml:"Cpuset,omitempty" toml:"Cpuset,omitempty"`
    PortSpecs         []string            `json:"PortSpecs,omitempty" yaml:"PortSpecs,omitempty" toml:"PortSpecs,omitempty"`
    ExposedPorts      map[Port]struct{}   `json:"ExposedPorts,omitempty" yaml:"ExposedPorts,omitempty" toml:"ExposedPorts,omitempty"`
    PublishService    string              `json:"PublishService,omitempty" yaml:"PublishService,omitempty" toml:"PublishService,omitempty"`
    StopSignal        string              `json:"StopSignal,omitempty" yaml:"StopSignal,omitempty" toml:"StopSignal,omitempty"`
    StopTimeout       int                 `json:"StopTimeout,omitempty" yaml:"StopTimeout,omitempty" toml:"StopTimeout,omitempty"`
    Env               []string            `json:"Env,omitempty" yaml:"Env,omitempty" toml:"Env,omitempty"`
    Cmd               []string            `json:"Cmd" yaml:"Cmd" toml:"Cmd"`
    Shell             []string            `json:"Shell,omitempty" yaml:"Shell,omitempty" toml:"Shell,omitempty"`
    Healthcheck       *HealthConfig       `json:"Healthcheck,omitempty" yaml:"Healthcheck,omitempty" toml:"Healthcheck,omitempty"`
    DNS               []string            `json:"Dns,omitempty" yaml:"Dns,omitempty" toml:"Dns,omitempty"` // For Docker API v1.9 and below only
    Image             string              `json:"Image,omitempty" yaml:"Image,omitempty" toml:"Image,omitempty"`
    Volumes           map[string]struct{} `json:"Volumes,omitempty" yaml:"Volumes,omitempty" toml:"Volumes,omitempty"`
    VolumeDriver      string              `json:"VolumeDriver,omitempty" yaml:"VolumeDriver,omitempty" toml:"VolumeDriver,omitempty"`
    WorkingDir        string              `json:"WorkingDir,omitempty" yaml:"WorkingDir,omitempty" toml:"WorkingDir,omitempty"`
    MacAddress        string              `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty" toml:"MacAddress,omitempty"`
    Entrypoint        []string            `json:"Entrypoint" yaml:"Entrypoint" toml:"Entrypoint"`
    SecurityOpts      []string            `json:"SecurityOpts,omitempty" yaml:"SecurityOpts,omitempty" toml:"SecurityOpts,omitempty"`
    OnBuild           []string            `json:"OnBuild,omitempty" yaml:"OnBuild,omitempty" toml:"OnBuild,omitempty"`
    Mounts            []Mount             `json:"Mounts,omitempty" yaml:"Mounts,omitempty" toml:"Mounts,omitempty"`
    Labels            map[string]string   `json:"Labels,omitempty" yaml:"Labels,omitempty" toml:"Labels,omitempty"`
    AttachStdin       bool                `json:"AttachStdin,omitempty" yaml:"AttachStdin,omitempty" toml:"AttachStdin,omitempty"`
    AttachStdout      bool                `json:"AttachStdout,omitempty" yaml:"AttachStdout,omitempty" toml:"AttachStdout,omitempty"`
    AttachStderr      bool                `json:"AttachStderr,omitempty" yaml:"AttachStderr,omitempty" toml:"AttachStderr,omitempty"`
    ArgsEscaped       bool                `json:"ArgsEscaped,omitempty" yaml:"ArgsEscaped,omitempty" toml:"ArgsEscaped,omitempty"`
    Tty               bool                `json:"Tty,omitempty" yaml:"Tty,omitempty" toml:"Tty,omitempty"`
    OpenStdin         bool                `json:"OpenStdin,omitempty" yaml:"OpenStdin,omitempty" toml:"OpenStdin,omitempty"`
    StdinOnce         bool                `json:"StdinOnce,omitempty" yaml:"StdinOnce,omitempty" toml:"StdinOnce,omitempty"`
    NetworkDisabled   bool                `json:"NetworkDisabled,omitempty" yaml:"NetworkDisabled,omitempty" toml:"NetworkDisabled,omitempty"`

    // This is no longer used and has been kept here for backward
    // compatibility, please use HostConfig.VolumesFrom.
    VolumesFrom string `json:"VolumesFrom,omitempty" yaml:"VolumesFrom,omitempty" toml:"VolumesFrom,omitempty"`
}

//另外有个HostConfig结构体是docker daemon(dockerd)的配置信息


//容器的配置信息
type Container struct {
    ID  string `json:"Id" yaml:"Id" toml:"Id"`

    Created time.Time `json:"Created,omitempty" yaml:"Created,omitempty" toml:"Created,omitempty"`

    Path string   `json:"Path,omitempty" yaml:"Path,omitempty" toml:"Path,omitempty"`
    Args []string `json:"Args,omitempty" yaml:"Args,omitempty" toml:"Args,omitempty"`

    Config *Config `json:"Config,omitempty" yaml:"Config,omitempty" toml:"Config,omitempty"`
    State  State   `json:"State,omitempty" yaml:"State,omitempty" toml:"State,omitempty"`
    Image  string  `json:"Image,omitempty" yaml:"Image,omitempty" toml:"Image,omitempty"`

    Node *SwarmNode `json:"Node,omitempty" yaml:"Node,omitempty" toml:"Node,omitempty"`

    NetworkSettings *NetworkSettings `json:"NetworkSettings,omitempty" yaml:"NetworkSettings,omitempty" toml:"NetworkSettings,omitempty"`

    SysInitPath    string  `json:"SysInitPath,omitempty" yaml:"SysInitPath,omitempty" toml:"SysInitPath,omitempty"`
    ResolvConfPath string  `json:"ResolvConfPath,omitempty" yaml:"ResolvConfPath,omitempty" toml:"ResolvConfPath,omitempty"`
    HostnamePath   string  `json:"HostnamePath,omitempty" yaml:"HostnamePath,omitempty" toml:"HostnamePath,omitempty"`
    HostsPath      string  `json:"HostsPath,omitempty" yaml:"HostsPath,omitempty" toml:"HostsPath,omitempty"`
    LogPath        string  `json:"LogPath,omitempty" yaml:"LogPath,omitempty" toml:"LogPath,omitempty"`
    Name           string  `json:"Name,omitempty" yaml:"Name,omitempty" toml:"Name,omitempty"`
    Driver         string  `json:"Driver,omitempty" yaml:"Driver,omitempty" toml:"Driver,omitempty"`
    Mounts         []Mount `json:"Mounts,omitempty" yaml:"Mounts,omitempty" toml:"Mounts,omitempty"`

    Volumes     map[string]string `json:"Volumes,omitempty" yaml:"Volumes,omitempty" toml:"Volumes,omitempty"`
    VolumesRW   map[string]bool   `json:"VolumesRW,omitempty" yaml:"VolumesRW,omitempty" toml:"VolumesRW,omitempty"`
    HostConfig  *HostConfig       `json:"HostConfig,omitempty" yaml:"HostConfig,omitempty" toml:"HostConfig,omitempty"`
    ExecIDs     []string          `json:"ExecIDs,omitempty" yaml:"ExecIDs,omitempty" toml:"ExecIDs,omitempty"`
    GraphDriver *GraphDriver      `json:"GraphDriver,omitempty" yaml:"GraphDriver,omitempty" toml:"GraphDriver,omitempty"`

    RestartCount int `json:"RestartCount,omitempty" yaml:"RestartCount,omitempty" toml:"RestartCount,omitempty"`

    AppArmorProfile string `json:"AppArmorProfile,omitempty" yaml:"AppArmorProfile,omitempty" toml:"AppArmorProfile,omitempty"`
}


//运行态的容器
type ContainerAlreadyRunning struct {
    ID string
}

//容器网络接口
type ContainerNetwork struct {
    Aliases             []string `json:"Aliases,omitempty" yaml:"Aliases,omitempty" toml:"Aliases,omitempty"`
    MacAddress          string   `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty" toml:"MacAddress,omitempty"`
    GlobalIPv6PrefixLen int      `json:"GlobalIPv6PrefixLen,omitempty" yaml:"GlobalIPv6PrefixLen,omitempty" toml:"GlobalIPv6PrefixLen,omitempty"`
    GlobalIPv6Address   string   `json:"GlobalIPv6Address,omitempty" yaml:"GlobalIPv6Address,omitempty" toml:"GlobalIPv6Address,omitempty"`
    IPv6Gateway         string   `json:"IPv6Gateway,omitempty" yaml:"IPv6Gateway,omitempty" toml:"IPv6Gateway,omitempty"`
    IPPrefixLen         int      `json:"IPPrefixLen,omitempty" yaml:"IPPrefixLen,omitempty" toml:"IPPrefixLen,omitempty"`
    IPAddress           string   `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty" toml:"IPAddress,omitempty"`
    Gateway             string   `json:"Gateway,omitempty" yaml:"Gateway,omitempty" toml:"Gateway,omitempty"`
    EndpointID          string   `json:"EndpointID,omitempty" yaml:"EndpointID,omitempty" toml:"EndpointID,omitempty"`
    NetworkID           string   `json:"NetworkID,omitempty" yaml:"NetworkID,omitempty" toml:"NetworkID,omitempty"`
}

//容器拷贝内容选项
type CopyFromContainerOptions struct {
    OutputStream io.Writer `json:"-"`
    Container    string    `json:"-"`
    Resource     string
    Context      context.Context `json:"-"`
}

//创建容器参数
type CreateContainerOptions struct {
    Name             string
    Config           *Config           `qs:"-"`
    HostConfig       *HostConfig       `qs:"-"`
    NetworkingConfig *NetworkingConfig `qs:"-"`
    Context          context.Context
}

//创建CreateExecContainer参数
type CreateExecOptions struct {
    AttachStdin  bool            `json:"AttachStdin,omitempty" yaml:"AttachStdin,omitempty" toml:"AttachStdin,omitempty"`
    AttachStdout bool            `json:"AttachStdout,omitempty" yaml:"AttachStdout,omitempty" toml:"AttachStdout,omitempty"`
    AttachStderr bool            `json:"AttachStderr,omitempty" yaml:"AttachStderr,omitempty" toml:"AttachStderr,omitempty"`
    Tty          bool            `json:"Tty,omitempty" yaml:"Tty,omitempty" toml:"Tty,omitempty"`
    Env          []string        `json:"Env,omitempty" yaml:"Env,omitempty" toml:"Env,omitempty"`
    Cmd          []string        `json:"Cmd,omitempty" yaml:"Cmd,omitempty" toml:"Cmd,omitempty"`
    Container    string          `json:"Container,omitempty" yaml:"Container,omitempty" toml:"Container,omitempty"`
    User         string          `json:"User,omitempty" yaml:"User,omitempty" toml:"User,omitempty"`
    Context      context.Context `json:"-"`
    Privileged   bool            `json:"Privileged,omitempty" yaml:"Privileged,omitempty" toml:"Privileged,omitempty"`
}


//创建volume参数
type CreateVolumeOptions struct {
    Name       string
    Driver     string
    DriverOpts map[string]string
    Context    context.Context `json:"-"`
    Labels     map[string]string
}

//容器和宿主机之间的映射关系
type Device struct {
    PathOnHost        string `json:"PathOnHost,omitempty" yaml:"PathOnHost,omitempty" toml:"PathOnHost,omitempty"`
    PathInContainer   string `json:"PathInContainer,omitempty" yaml:"PathInContainer,omitempty" toml:"PathInContainer,omitempty"`
    CgroupPermissions string `json:"CgroupPermissions,omitempty" yaml:"CgroupPermissions,omitempty" toml:"CgroupPermissions,omitempty"`
}

//查看docker server上的全部信息
type DockerInfo struct {
    ID                 string
    Containers         int
    ContainersRunning  int
    ContainersPaused   int
    ContainersStopped  int
    Images             int
    Driver             string
    DriverStatus       [][2]string
    SystemStatus       [][2]string
    Plugins            PluginsInfo
    MemoryLimit        bool
    SwapLimit          bool
    KernelMemory       bool
    CPUCfsPeriod       bool `json:"CpuCfsPeriod"`
    CPUCfsQuota        bool `json:"CpuCfsQuota"`
    CPUShares          bool
    CPUSet             bool
    IPv4Forwarding     bool
    BridgeNfIptables   bool
    BridgeNfIP6tables  bool `json:"BridgeNfIp6tables"`
    Debug              bool
    OomKillDisable     bool
    ExperimentalBuild  bool
    NFd                int
    NGoroutines        int
    SystemTime         string
    ExecutionDriver    string
    LoggingDriver      string
    CgroupDriver       string
    NEventsListener    int
    KernelVersion      string
    OperatingSystem    string
    OSType             string
    Architecture       string
    IndexServerAddress string
    RegistryConfig     *ServiceConfig
    SecurityOptions    []string
    NCPU               int
    MemTotal           int64
    DockerRootDir      string
    HTTPProxy          string `json:"HttpProxy"`
    HTTPSProxy         string `json:"HttpsProxy"`
    NoProxy            string
    Name               string
    Labels             []string
    ServerVersion      string
    ClusterStore       string
    ClusterAdvertise   string
    Isolation          string
    InitBinary         string
    DefaultRuntime     string
    LiveRestoreEnabled bool
    Swarm              swarm.Info
}

//endpoint结构体
type Endpoint struct {
    Name        string
    ID          string `json:"EndpointID"`
    MacAddress  string
    IPv4Address string
    IPv6Address string
}


//endpoint配置结构体
type EndpointConfig struct {
    IPAMConfig          *EndpointIPAMConfig `json:"IPAMConfig,omitempty" yaml:"IPAMConfig,omitempty" toml:"IPAMConfig,omitempty"`
    Links               []string            `json:"Links,omitempty" yaml:"Links,omitempty" toml:"Links,omitempty"`
    Aliases             []string            `json:"Aliases,omitempty" yaml:"Aliases,omitempty" toml:"Aliases,omitempty"`
    NetworkID           string              `json:"NetworkID,omitempty" yaml:"NetworkID,omitempty" toml:"NetworkID,omitempty"`
    EndpointID          string              `json:"EndpointID,omitempty" yaml:"EndpointID,omitempty" toml:"EndpointID,omitempty"`
    Gateway             string              `json:"Gateway,omitempty" yaml:"Gateway,omitempty" toml:"Gateway,omitempty"`
    IPAddress           string              `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty" toml:"IPAddress,omitempty"`
    IPPrefixLen         int                 `json:"IPPrefixLen,omitempty" yaml:"IPPrefixLen,omitempty" toml:"IPPrefixLen,omitempty"`
    IPv6Gateway         string              `json:"IPv6Gateway,omitempty" yaml:"IPv6Gateway,omitempty" toml:"IPv6Gateway,omitempty"`
    GlobalIPv6Address   string              `json:"GlobalIPv6Address,omitempty" yaml:"GlobalIPv6Address,omitempty" toml:"GlobalIPv6Address,omitempty"`
    GlobalIPv6PrefixLen int                 `json:"GlobalIPv6PrefixLen,omitempty" yaml:"GlobalIPv6PrefixLen,omitempty" toml:"GlobalIPv6PrefixLen,omitempty"`
    MacAddress          string              `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty" toml:"MacAddress,omitempty"`
}


//容器的环境
type Env []string
环境变量会有几个比较重要的方法
func (env *Env) Map() map[string]string
func (env *Env) Set(key, value string)
func (env *Env) SetJSON(key string, value interface{}) error

func (env *Env) GetJSON(key string, iface interface{}) error


//健康检查
type Health struct {
    Status        string        `json:"Status,omitempty" yaml:"Status,omitempty" toml:"Status,omitempty"`
    FailingStreak int           `json:"FailingStreak,omitempty" yaml:"FailingStreak,omitempty" toml:"FailingStreak,omitempty"`
    Log           []HealthCheck `json:"Log,omitempty" yaml:"Log,omitempty" toml:"Log,omitempty"`
}


type HealthCheck struct {
    Start    time.Time `json:"Start,omitempty" yaml:"Start,omitempty" toml:"Start,omitempty"`
    End      time.Time `json:"End,omitempty" yaml:"End,omitempty" toml:"End,omitempty"`
    ExitCode int       `json:"ExitCode,omitempty" yaml:"ExitCode,omitempty" toml:"ExitCode,omitempty"`
    Output   string    `json:"Output,omitempty" yaml:"Output,omitempty" toml:"Output,omitempty"`
}


//HostConfig
type HostConfig struct {
    Binds                []string               `json:"Binds,omitempty" yaml:"Binds,omitempty" toml:"Binds,omitempty"`
    CapAdd               []string               `json:"CapAdd,omitempty" yaml:"CapAdd,omitempty" toml:"CapAdd,omitempty"`
    CapDrop              []string               `json:"CapDrop,omitempty" yaml:"CapDrop,omitempty" toml:"CapDrop,omitempty"`
    GroupAdd             []string               `json:"GroupAdd,omitempty" yaml:"GroupAdd,omitempty" toml:"GroupAdd,omitempty"`
    ContainerIDFile      string                 `json:"ContainerIDFile,omitempty" yaml:"ContainerIDFile,omitempty" toml:"ContainerIDFile,omitempty"`
    LxcConf              []KeyValuePair         `json:"LxcConf,omitempty" yaml:"LxcConf,omitempty" toml:"LxcConf,omitempty"`
    PortBindings         map[Port][]PortBinding `json:"PortBindings,omitempty" yaml:"PortBindings,omitempty" toml:"PortBindings,omitempty"`
    Links                []string               `json:"Links,omitempty" yaml:"Links,omitempty" toml:"Links,omitempty"`
    DNS                  []string               `json:"Dns,omitempty" yaml:"Dns,omitempty" toml:"Dns,omitempty"` // For Docker API v1.10 and above only
    DNSOptions           []string               `json:"DnsOptions,omitempty" yaml:"DnsOptions,omitempty" toml:"DnsOptions,omitempty"`
    DNSSearch            []string               `json:"DnsSearch,omitempty" yaml:"DnsSearch,omitempty" toml:"DnsSearch,omitempty"`
    ExtraHosts           []string               `json:"ExtraHosts,omitempty" yaml:"ExtraHosts,omitempty" toml:"ExtraHosts,omitempty"`
    VolumesFrom          []string               `json:"VolumesFrom,omitempty" yaml:"VolumesFrom,omitempty" toml:"VolumesFrom,omitempty"`
    UsernsMode           string                 `json:"UsernsMode,omitempty" yaml:"UsernsMode,omitempty" toml:"UsernsMode,omitempty"`
    NetworkMode          string                 `json:"NetworkMode,omitempty" yaml:"NetworkMode,omitempty" toml:"NetworkMode,omitempty"`
    IpcMode              string                 `json:"IpcMode,omitempty" yaml:"IpcMode,omitempty" toml:"IpcMode,omitempty"`
    PidMode              string                 `json:"PidMode,omitempty" yaml:"PidMode,omitempty" toml:"PidMode,omitempty"`
    UTSMode              string                 `json:"UTSMode,omitempty" yaml:"UTSMode,omitempty" toml:"UTSMode,omitempty"`
    RestartPolicy        RestartPolicy          `json:"RestartPolicy,omitempty" yaml:"RestartPolicy,omitempty" toml:"RestartPolicy,omitempty"`
    Devices              []Device               `json:"Devices,omitempty" yaml:"Devices,omitempty" toml:"Devices,omitempty"`
    DeviceCgroupRules    []string               `json:"DeviceCgroupRules,omitempty" yaml:"DeviceCgroupRules,omitempty" toml:"DeviceCgroupRules,omitempty"`
    LogConfig            LogConfig              `json:"LogConfig,omitempty" yaml:"LogConfig,omitempty" toml:"LogConfig,omitempty"`
    SecurityOpt          []string               `json:"SecurityOpt,omitempty" yaml:"SecurityOpt,omitempty" toml:"SecurityOpt,omitempty"`
    Cgroup               string                 `json:"Cgroup,omitempty" yaml:"Cgroup,omitempty" toml:"Cgroup,omitempty"`
    CgroupParent         string                 `json:"CgroupParent,omitempty" yaml:"CgroupParent,omitempty" toml:"CgroupParent,omitempty"`
    Memory               int64                  `json:"Memory,omitempty" yaml:"Memory,omitempty" toml:"Memory,omitempty"`
    MemoryReservation    int64                  `json:"MemoryReservation,omitempty" yaml:"MemoryReservation,omitempty" toml:"MemoryReservation,omitempty"`
    KernelMemory         int64                  `json:"KernelMemory,omitempty" yaml:"KernelMemory,omitempty" toml:"KernelMemory,omitempty"`
    MemorySwap           int64                  `json:"MemorySwap,omitempty" yaml:"MemorySwap,omitempty" toml:"MemorySwap,omitempty"`
    MemorySwappiness     int64                  `json:"MemorySwappiness,omitempty" yaml:"MemorySwappiness,omitempty" toml:"MemorySwappiness,omitempty"`
    CPUShares            int64                  `json:"CpuShares,omitempty" yaml:"CpuShares,omitempty" toml:"CpuShares,omitempty"`
    CPUSet               string                 `json:"Cpuset,omitempty" yaml:"Cpuset,omitempty" toml:"Cpuset,omitempty"`
    CPUSetCPUs           string                 `json:"CpusetCpus,omitempty" yaml:"CpusetCpus,omitempty" toml:"CpusetCpus,omitempty"`
    CPUSetMEMs           string                 `json:"CpusetMems,omitempty" yaml:"CpusetMems,omitempty" toml:"CpusetMems,omitempty"`
    CPUQuota             int64                  `json:"CpuQuota,omitempty" yaml:"CpuQuota,omitempty" toml:"CpuQuota,omitempty"`
    CPUPeriod            int64                  `json:"CpuPeriod,omitempty" yaml:"CpuPeriod,omitempty" toml:"CpuPeriod,omitempty"`
    CPURealtimePeriod    int64                  `json:"CpuRealtimePeriod,omitempty" yaml:"CpuRealtimePeriod,omitempty" toml:"CpuRealtimePeriod,omitempty"`
    CPURealtimeRuntime   int64                  `json:"CpuRealtimeRuntime,omitempty" yaml:"CpuRealtimeRuntime,omitempty" toml:"CpuRealtimeRuntime,omitempty"`
    BlkioWeight          int64                  `json:"BlkioWeight,omitempty" yaml:"BlkioWeight,omitempty" toml:"BlkioWeight,omitempty"`
    BlkioWeightDevice    []BlockWeight          `json:"BlkioWeightDevice,omitempty" yaml:"BlkioWeightDevice,omitempty" toml:"BlkioWeightDevice,omitempty"`
    BlkioDeviceReadBps   []BlockLimit           `json:"BlkioDeviceReadBps,omitempty" yaml:"BlkioDeviceReadBps,omitempty" toml:"BlkioDeviceReadBps,omitempty"`
    BlkioDeviceReadIOps  []BlockLimit           `json:"BlkioDeviceReadIOps,omitempty" yaml:"BlkioDeviceReadIOps,omitempty" toml:"BlkioDeviceReadIOps,omitempty"`
    BlkioDeviceWriteBps  []BlockLimit           `json:"BlkioDeviceWriteBps,omitempty" yaml:"BlkioDeviceWriteBps,omitempty" toml:"BlkioDeviceWriteBps,omitempty"`
    BlkioDeviceWriteIOps []BlockLimit           `json:"BlkioDeviceWriteIOps,omitempty" yaml:"BlkioDeviceWriteIOps,omitempty" toml:"BlkioDeviceWriteIOps,omitempty"`
    Ulimits              []ULimit               `json:"Ulimits,omitempty" yaml:"Ulimits,omitempty" toml:"Ulimits,omitempty"`
    VolumeDriver         string                 `json:"VolumeDriver,omitempty" yaml:"VolumeDriver,omitempty" toml:"VolumeDriver,omitempty"`
    OomScoreAdj          int                    `json:"OomScoreAdj,omitempty" yaml:"OomScoreAdj,omitempty" toml:"OomScoreAdj,omitempty"`
    PidsLimit            int64                  `json:"PidsLimit,omitempty" yaml:"PidsLimit,omitempty" toml:"PidsLimit,omitempty"`
    ShmSize              int64                  `json:"ShmSize,omitempty" yaml:"ShmSize,omitempty" toml:"ShmSize,omitempty"`
    Tmpfs                map[string]string      `json:"Tmpfs,omitempty" yaml:"Tmpfs,omitempty" toml:"Tmpfs,omitempty"`
    Privileged           bool                   `json:"Privileged,omitempty" yaml:"Privileged,omitempty" toml:"Privileged,omitempty"`
    PublishAllPorts      bool                   `json:"PublishAllPorts,omitempty" yaml:"PublishAllPorts,omitempty" toml:"PublishAllPorts,omitempty"`
    ReadonlyRootfs       bool                   `json:"ReadonlyRootfs,omitempty" yaml:"ReadonlyRootfs,omitempty" toml:"ReadonlyRootfs,omitempty"`
    OOMKillDisable       bool                   `json:"OomKillDisable,omitempty" yaml:"OomKillDisable,omitempty" toml:"OomKillDisable,omitempty"`
    AutoRemove           bool                   `json:"AutoRemove,omitempty" yaml:"AutoRemove,omitempty" toml:"AutoRemove,omitempty"`
    StorageOpt           map[string]string      `json:"StorageOpt,omitempty" yaml:"StorageOpt,omitempty" toml:"StorageOpt,omitempty"`
    Sysctls              map[string]string      `json:"Sysctls,omitempty" yaml:"Sysctls,omitempty" toml:"Sysctls,omitempty"`
    CPUCount             int64                  `json:"CpuCount,omitempty" yaml:"CpuCount,omitempty"`
    CPUPercent           int64                  `json:"CpuPercent,omitempty" yaml:"CpuPercent,omitempty"`
    IOMaximumBandwidth   int64                  `json:"IOMaximumBandwidth,omitempty" yaml:"IOMaximumBandwidth,omitempty"`
    IOMaximumIOps        int64                  `json:"IOMaximumIOps,omitempty" yaml:"IOMaximumIOps,omitempty"`
    Mounts               []HostMount            `json:"Mounts,omitempty" yaml:"Mounts,omitempty" toml:"Mounts,omitempty"`
    Init                 bool                   `json:",omitempty" yaml:",omitempty"`
}

//IPAM配置
type IPAMConfig struct {
    Subnet     string            `json:",omitempty"`
    IPRange    string            `json:",omitempty"`
    Gateway    string            `json:",omitempty"`
    AuxAddress map[string]string `json:"AuxiliaryAddresses,omitempty"`
}


//image相关
type Image struct {
    ID              string    `json:"Id" yaml:"Id" toml:"Id"`
    RepoTags        []string  `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty" toml:"RepoTags,omitempty"`
    Parent          string    `json:"Parent,omitempty" yaml:"Parent,omitempty" toml:"Parent,omitempty"`
    Comment         string    `json:"Comment,omitempty" yaml:"Comment,omitempty" toml:"Comment,omitempty"`
    Created         time.Time `json:"Created,omitempty" yaml:"Created,omitempty" toml:"Created,omitempty"`
    Container       string    `json:"Container,omitempty" yaml:"Container,omitempty" toml:"Container,omitempty"`
    ContainerConfig Config    `json:"ContainerConfig,omitempty" yaml:"ContainerConfig,omitempty" toml:"ContainerConfig,omitempty"`
    DockerVersion   string    `json:"DockerVersion,omitempty" yaml:"DockerVersion,omitempty" toml:"DockerVersion,omitempty"`
    Author          string    `json:"Author,omitempty" yaml:"Author,omitempty" toml:"Author,omitempty"`
    Config          *Config   `json:"Config,omitempty" yaml:"Config,omitempty" toml:"Config,omitempty"`
    Architecture    string    `json:"Architecture,omitempty" yaml:"Architecture,omitempty"`
    Size            int64     `json:"Size,omitempty" yaml:"Size,omitempty" toml:"Size,omitempty"`
    VirtualSize     int64     `json:"VirtualSize,omitempty" yaml:"VirtualSize,omitempty" toml:"VirtualSize,omitempty"`
    RepoDigests     []string  `json:"RepoDigests,omitempty" yaml:"RepoDigests,omitempty" toml:"RepoDigests,omitempty"`
    RootFS          *RootFS   `json:"RootFS,omitempty" yaml:"RootFS,omitempty" toml:"RootFS,omitempty"`
    OS              string    `json:"Os,omitempty" yaml:"Os,omitempty" toml:"Os,omitempty"`
}

//镜像列表
type ListImagesOptions struct {
    Filters map[string][]string
    All     bool
    Digests bool
    Filter  string
    Context context.Context
}

//构建镜像参数
type BuildImageOptions struct {
    Name                string             `qs:"t"`
    Dockerfile          string             `qs:"dockerfile"`
    NoCache             bool               `qs:"nocache"`
    CacheFrom           []string           `qs:"-"`
    SuppressOutput      bool               `qs:"q"`
    Pull                bool               `qs:"pull"`
    RmTmpContainer      bool               `qs:"rm"`
    ForceRmTmpContainer bool               `qs:"forcerm"`
    RawJSONStream       bool               `qs:"-"`
    Memory              int64              `qs:"memory"`
    Memswap             int64              `qs:"memswap"`
    CPUShares           int64              `qs:"cpushares"`
    CPUQuota            int64              `qs:"cpuquota"`
    CPUPeriod           int64              `qs:"cpuperiod"`
    CPUSetCPUs          string             `qs:"cpusetcpus"`
    Labels              map[string]string  `qs:"labels"`
    InputStream         io.Reader          `qs:"-"`
    OutputStream        io.Writer          `qs:"-"`
    Remote              string             `qs:"remote"`
    Auth                AuthConfiguration  `qs:"-"` // for older docker X-Registry-Auth header
    AuthConfigs         AuthConfigurations `qs:"-"` // for newer docker X-Registry-Config header
    ContextDir          string             `qs:"-"`
    Ulimits             []ULimit           `qs:"-"`
    BuildArgs           []BuildArg         `qs:"-"`
    NetworkMode         string             `qs:"networkmode"`
    InactivityTimeout   time.Duration      `qs:"-"`
    CgroupParent        string             `qs:"cgroupparent"`
    SecurityOpt         []string           `qs:"securityopt"`
    Target              string             `gs:"target"`
    Context             context.Context
}

//停止容器选项
type KillContainerOptions struct {
    // The ID of the container.
    ID  string `qs:"-"`

    // The signal to send to the container. When omitted, Docker server
    // will assume SIGKILL.
    Signal  Signal
    Context context.Context
}



//挂载参数
type Mount struct {
    Name        string
    Source      string
    Destination string
    Driver      string
    Mode        string
    RW          bool
}
//网络
type Network struct {
    Name       string
    ID         string `json:"Id"`
    Scope      string
    Driver     string
    IPAM       IPAMOptions
    Containers map[string]Endpoint
    Options    map[string]string
    Internal   bool
    EnableIPv6 bool `json:"EnableIPv6"`
    Labels     map[string]string
}


//异常信息
type NoSuchContainer struct {
    ID  string
    Err error
}

//端口映射
type Port string

func (p Port) Port() string
func (p Port) Proto() string

type PortBinding struct {
    HostIP   string `json:"HostIp,omitempty" yaml:"HostIp,omitempty" toml:"HostIp,omitempty"`
    HostPort string `json:"HostPort,omitempty" yaml:"HostPort,omitempty" toml:"HostPort,omitempty"`
}
type PortMapping map[string]string

//pull image 选项
type PullImageOptions struct {
    Repository string `qs:"fromImage"`
    Tag        string

    // Only required for Docker Engine 1.9 or 1.10 w/ Remote API < 1.21
    // and Docker Engine < 1.9
    // This parameter was removed in Docker Engine 1.11
    Registry string

    OutputStream      io.Writer     `qs:"-"`
    RawJSONStream     bool          `qs:"-"`
    InactivityTimeout time.Duration `qs:"-"`
    Context           context.Context
}

//push image 选项
type PushImageOptions struct {
    // Name of the image
    Name string

    // Tag of the image
    Tag string

    // Registry server to push the image
    Registry string

    OutputStream      io.Writer     `qs:"-"`
    RawJSONStream     bool          `qs:"-"`
    InactivityTimeout time.Duration `qs:"-"`

    Context context.Context
}



//删除配置以及删除容器
type RemoveConfigOptions struct {
    ID      string `qs:"-"`
    Context context.Context
}

type RemoveContainerOptions struct {
    // The ID of the container.
    ID  string `qs:"-"`

    // A flag that indicates whether Docker should remove the volumes
    // associated to the container.
    RemoveVolumes bool `qs:"v"`

    // A flag that indicates whether Docker should remove the container
    // even if it is currently running.
    Force   bool
    Context context.Context
}

//删除镜像
type RemoveImageOptions struct {
    Force   bool `qs:"force"`
    NoPrune bool `qs:"noprune"`
    Context context.Context
}


//重命名容器
type RenameContainerOptions struct {
    // ID of container to rename
    ID  string `qs:"-"`

    // New name
    Name    string `json:"name,omitempty" yaml:"name,omitempty"`
    Context context.Context
}

//容器状态
type State struct {
    Status            string    `json:"Status,omitempty" yaml:"Status,omitempty" toml:"Status,omitempty"`
    Running           bool      `json:"Running,omitempty" yaml:"Running,omitempty" toml:"Running,omitempty"`
    Paused            bool      `json:"Paused,omitempty" yaml:"Paused,omitempty" toml:"Paused,omitempty"`
    Restarting        bool      `json:"Restarting,omitempty" yaml:"Restarting,omitempty" toml:"Restarting,omitempty"`
    OOMKilled         bool      `json:"OOMKilled,omitempty" yaml:"OOMKilled,omitempty" toml:"OOMKilled,omitempty"`
    RemovalInProgress bool      `json:"RemovalInProgress,omitempty" yaml:"RemovalInProgress,omitempty" toml:"RemovalInProgress,omitempty"`
    Dead              bool      `json:"Dead,omitempty" yaml:"Dead,omitempty" toml:"Dead,omitempty"`
    Pid               int       `json:"Pid,omitempty" yaml:"Pid,omitempty" toml:"Pid,omitempty"`
    ExitCode          int       `json:"ExitCode,omitempty" yaml:"ExitCode,omitempty" toml:"ExitCode,omitempty"`
    Error             string    `json:"Error,omitempty" yaml:"Error,omitempty" toml:"Error,omitempty"`
    StartedAt         time.Time `json:"StartedAt,omitempty" yaml:"StartedAt,omitempty" toml:"StartedAt,omitempty"`
    FinishedAt        time.Time `json:"FinishedAt,omitempty" yaml:"FinishedAt,omitempty" toml:"FinishedAt,omitempty"`
    Health            Health    `json:"Health,omitempty" yaml:"Health,omitempty" toml:"Health,omitempty"`
}

func (s *State) String() string
func (s *State) StateString() string

//镜像tag
type TagImageOptions struct {
    Repo    string
    Tag     string
    Force   bool
    Context context.Context
}

//更新容器
type UpdateConfigOptions struct {
    Auth AuthConfiguration `qs:"-"`
    swarm.ConfigSpec
    Context context.Context
    Version uint64
}


type UpdateContainerOptions struct {
    BlkioWeight        int           `json:"BlkioWeight"`
    CPUShares          int           `json:"CpuShares"`
    CPUPeriod          int           `json:"CpuPeriod"`
    CPURealtimePeriod  int64         `json:"CpuRealtimePeriod"`
    CPURealtimeRuntime int64         `json:"CpuRealtimeRuntime"`
    CPUQuota           int           `json:"CpuQuota"`
    CpusetCpus         string        `json:"CpusetCpus"`
    CpusetMems         string        `json:"CpusetMems"`
    Memory             int           `json:"Memory"`
    MemorySwap         int           `json:"MemorySwap"`
    MemoryReservation  int           `json:"MemoryReservation"`
    KernelMemory       int           `json:"KernelMemory"`
    RestartPolicy      RestartPolicy `json:"RestartPolicy,omitempty"`
    Context            context.Context
}


```

