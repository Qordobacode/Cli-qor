
build:
		build-windows
		build-linux
		build-arm

build-windows:
	GOOS=windows go build -o cli.exe

build-linux:
    GOOS=linux go build -o cli
    
build-arm:
    GOARCH=armv7 GOOS=linux go build -o cli-rpi
