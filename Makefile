DOCKER_RELEASE_TAG = "kjk/dbworkbench:$(shell git describe --abbrev=0 --tags | sed 's/v//')"
BINDATA_IGNORE = $(shell git ls-files -io --exclude-standard $< | sed 's/^/-ignore=/;s/[.]/[.]/g')

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make setup           : Install all necessary dependencies"
	@echo "make dev             : Generate development build"
	@echo "make test            : Run tests"
	@echo "make build           : Generate production build for current OS"
	@echo "make bootstrap       : Install cross-compilation toolchain"
	@echo "make release         : Generate binaries for all supported OSes"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make docker          : Build docker image"
	@echo "make docker-release  : Build and tag docker image"
	@echo ""

test:
	godep go test -cover

build:
	godep go build
	@echo "You can now execute ./pgweb"

release:
	gox -osarch="darwin/amd64 darwin/386 linux/amd64 linux/386 windows/amd64 windows/386" -output="./bin/pgweb_{{.OS}}_{{.Arch}}"

bootstrap:
	gox -build-toolchain

setup:
	go get github.com/tools/godep
	go get golang.org/x/tools/cmd/cover
	godep get github.com/mitchellh/gox
	godep get github.com/jteeuwen/go-bindata/...
	godep restore

clean:
	rm -f ./dbworkbench
	rm -f ./bin/*

docker:
	docker build -t dbworkbench .

docker-release:
	docker build -t $(DOCKER_RELEASE_TAG) .
