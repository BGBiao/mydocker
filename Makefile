all: build taging push clean
.PHONY: build taging push clean
jfdocker = jfdocker
version = 0.1.1
build:
	go build JFDocker.go && go build GetUsedGpuinfo.go && mkdir $(jfdocker)-$(version) && cp config.yml JFDocker GetUsedGpuinfo usage.md $(jfdocker)-$(version) && tar -zcf $(jfdocker)-$(version).tar.gz $(jfdocker)-$(version)

clean:
	mv $(jfdocker)-$(version).tar.gz /export/ && /bin/rm -rf $(jfdocker)-$(version) JFDocker GetUsedGpuinfo
   


