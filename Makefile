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


nsnetgo:
	wget "https://github.com/teddyking/netsetgo/releases/download/0.0.1/netsetgo"
	sudo mv netsetgo /usr/local/bin/
	sudo chown root:root /usr/local/bin/netsetgo
	sudo chmod 4755 /usr/local/bin/netsetgo


bridge:
	sudo iptables -tnat -N netsetgo
	sudo iptables -tnat -A PREROUTING -m addrtype --dst-type LOCAL -j netsetgo
	sudo iptables -tnat -A OUTPUT ! -d 127.0.0.0/8 -m addrtype --dst-type LOCAL -j netsetgo
	sudo iptables -tnat -A POSTROUTING -s 10.10.10.0/24 ! -o brg0 -j MASQUERADE
	sudo iptables -tnat -A netsetgo -i brg0 -j RETURN
	sudo sysctl -w net.ipv4.ip_forward=1

net: nsnetgo bridge
