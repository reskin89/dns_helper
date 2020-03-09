build: clean check-arch
	CGO_ENABLED=0 GOOS=linux GOARCH=$(arch) go build -a -o d2 .

#Build default builds the cli binary for alpine linux to run it as a cron from a container (whether it be crontab or something like the k8s scheduler)
build-default:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o d2 .

test:
	go test

clean:
	go clean
	rm -f d2

check-arch:
ifndef arch
	$(error arch is undefined: expected one of amd64, arm, i386 etc.)
endif
