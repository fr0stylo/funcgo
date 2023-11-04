unet:
	go build -o ./bin/ ./cmd/unet
	sudo chown root:root ./bin/unet
	sudo chmod u+s ./bin/unet
	sudo rm /usr/bin/unet
	sudo ln -s ${PWD}/bin/unet /usr/bin/unet



build:
	go build -o ./bin/ ./cmd/
