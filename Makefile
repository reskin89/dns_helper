region?=us-east-1

build: check-appname clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(appname) .

test:
	go test

clean: check-appname
	go clean
	rm -f $(appname)

check-appname:
ifndef appname
	$(error appname is undefined)
endif
