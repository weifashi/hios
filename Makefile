export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

MODULE = hios

PORT 			:= 3376
VERSION			:= $(shell git describe --tags --always --match="v*" 2> /dev/null || echo v0.0.1)
VERSION_HASH	:= $(shell git rev-parse --short HEAD)

GOCGO 	:= env CGO_ENABLED=1
LDFLAGS	:= -s -w -X "$(MODULE)/config.Version=$(VERSION)" -X "$(MODULE)/config.CommitSHA=$(VERSION_HASH)"

run: build
	./release/hios

build:
	go build -o ./release/hios

releases: 
	$(GOCGO) CC=x86_64-linux-musl-gcc GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./release/$(MODULE)-linux-amd64/$(MODULE)
	@for arch in amd64; \
	do \
		tar zcf ./release/$(MODULE)-linux-$$arch.tar.gz ./release/$(MODULE)-linux-$$arch; \
	done
	rm -r ./release/$(MODULE)-linux-amd64

docker-build:
	docker run --rm -v "${PWD}":/myapp -w /myapp golang:1.20 bash -c "make build"

dev:
	lsof -i :3377 | grep node | awk '{print $$2}' | xargs kill -9
	cd web && nohup npm run dev > ../output.log >&1 & cd ../ 
	${HOME}/go/bin/fresh -c ./fresh.conf

translate:
	cd web && npm run translate $(text) && cd ../


# 提示 fresh: No such file or directory 时解決辦法
# go install github.com/pilu/fresh@latest

# 提示 air: No such file or directory 时解決辦法
# go install github.com/cosmtrek/air@latest

# 提示 swag: No such file or directory 时解決辦法
# go get -u github.com/swaggo/swag/cmd/swag
# go install github.com/swaggo/swag/cmd/swag@latest
