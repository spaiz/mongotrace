VERSION=0.0.2

PATH_BUILD=build
FILE_COMMAND=mongotrace
FILE_ARCH=darwin_amd64

clean:
	@rm -rf ./build

build: clean
	GOOS=darwin GOARCH=amd64 go build -mod vendor -ldflags "-X main.VERSION=$(VERSION)" -o $(PATH_BUILD)/$(VERSION)/$(FILE_ARCH)/$(FILE_COMMAND)

version:
	@echo $(VERSION)

install:
	install -d -m 755 '/usr/local/bin'
	install $(PATH_BUILD)/$(VERSION)/$(FILE_ARCH)/$(FILE_COMMAND) '/usr/local/bin/$(FILE_COMMAND)'