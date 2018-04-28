package client
import (
    docker "github.com/fsouza/go-dockerclient"
)

func NewDockerClient(addr,version string) (*docker.Client,error) {
    return docker.NewVersionedClient(addr,version)
}


/*
ldImageOptions struct {
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
*/


