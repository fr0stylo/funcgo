unet:
	go build -o ./bin/ ./cmd/unet
	sudo chown root:root ./bin/unet
	sudo chmod u+s ./bin/unet
	sudo rm /usr/bin/unet
	sudo ln -s ${PWD}/bin/unet /usr/bin/unet

all:
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s"  -v  -o ./bin/ ./...

build: fs all unet

fs:
	sudo rm -rf ./fs
	cp -r ${shell docker image inspect alpine:latest -f \ {{.GraphDriver.Data.UpperDir}}} ./fs
	echo "168.0.0.1 host.funcgo.internal" >> ./fs/etc/resolv.conf


docker: all
	ln -s ${PWD}/bin/unet /usr/bin/unet
