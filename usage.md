## 
### 1. Creating the config 
```
mkdir -p /etc/JFDocker/
cp config.yml /etc/JFDocker/
```

### 2. Editing the `/etc/JFDocker/config.yml` 

```
#log and config root dir
logpath: /export/Logs/
#container config dir
configdir: containers_config
#pausecontainer image
pauseimage: f66f4bd9b894
#docker container volume dir
datapath: /data/
```

### 3. exec the jfdocker 

```
cp jfdocker /usr/bin/

```
