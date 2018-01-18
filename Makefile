all: build taging push clean
.PHONY: build taging push clean
jfdocker = jfdocker
version = 0.1.1
DOCKER_REGISTRY = dockerhub.jd.com/deeplearning/
build:
	go build jfdocker.go && tar -zcf $(jfdocker)-$(version).tar.gz config.yml jfdocker usage.md 

clean:
	mv $(jfdocker)-$(version).tar.gz /export/ && /bin/rm -rf jfdocker 
   


