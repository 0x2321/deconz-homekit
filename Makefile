build:
	go build -o app main.go

build_rpi:
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o app

run:
	DECONZ_IP=phoscon.home go run main.go

watch:
	reflex -s -r '\.go$$' make run