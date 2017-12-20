package apis

//注意LogPath必须以/结尾
const LogPath string = "/tmp/Logs/"
const ConConfDir string = "containers_config"

const Pauseimage string = "xxbandy123/k8s-pause"
const Netnamespace string = "/var/run/netns/"

const Usages string = "JFDocker run ywid,sn,image,container,ip,mask,gateway,cpu,mem,label,vnet\n JFDocker update container,images\n JFDocker resize container,cpu,4/container,mem,2048\n JFDocker delete container\n JFDocker  rebuilt container"


type ConSpec struct {
  Ywid  string `json:"ywid"`
  SN    string `json:"sn"`
  Mem   string `json:"mems"`
  Cpus  string `json:"cpus"`

}

type NetSpec struct {
  Ipv4  string `json:"ipaddress"`
  Mask  string `json:"mask"`
  Gateway string `json:gateway`
  Vnet  string `json:"vnet"`

}

type Netns struct {
  Cid     string `json:pausecontainerid,omitempty` 
  Pid     string `json:pausecontainerpid,omitempty`
  Netns   string `json:netns,omitempty`
  Net     NetSpec
}



type JFDocker struct {
  Appname string `json:"containername"`
  Image   string `json:"Image"`
  Conspec ConSpec
  Netspec NetSpec
}

type RspJFDocker struct {
  Result  uint  `json:"result,omitempty"`
  Appname string `json:"containername,omitempty"`
  ConID   string `json:"containerid,omitempty"`
  Ipv4    string `json:"ipaddress,omitempty"`
  ErrMsg  string `json:"errmsg,omitempty"`
  ErrInfo string `json:"errinfo,omitempty"`

}


type JFDockererr struct {
    JFDockerVersion string `json:"JFDockerversion,omitempty"`
    Code            uint   `json:"result,omitempty`
    Msg             string `json:"msg,omitempty"`
    Details         string `json:"details,omitempty"`

}


