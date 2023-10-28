export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

MODULE = hios

PORT 			:= 3376
VERSION			:= $(shell git tag | tail -1 2> /dev/null || echo v0.0.1)
VERSION_HASH	:= $(shell git rev-parse --short HEAD)
NEXT_VERSION    := $(shell git tag | tail -1 | awk -F. -v OFS=. '{$$NF++; print}')

GOCGO 	:= env CGO_ENABLED=1
LDFLAGS	:= -s -w -X "$(MODULE)/config.Version=$(VERSION)" -X "$(MODULE)/config.CommitSHA=$(VERSION_HASH)"

run: build
	./release/hios

build:
	go build -o ./release/hios

releases: 
    $(GOCGO) CC=x86_64-linux-musl-gcc GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./release/$(MODULE)-linux-amd64/$(MODULE)
	tar zcf ./release/$(MODULE)-linux-amd64.tar.gz ./release/$(MODULE)-linux-amd64 ; \
	rm -r ./release/$(MODULE)-linux-amd64

tag: 
	git tag $(NEXT_VERSION)
	git push origin $(NEXT_VERSION)

docker-releases:
	docker run --rm -v "${PWD}":/myapp -w /myapp golang:1.20 bash -c "make releases"

docker-build:
	docker run --rm -v "${PWD}":/myapp -w /myapp golang:1.20 bash -c "GOOS=linux GOARCH=amd64 go build"

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


# 下载
# wget https://gitee.com/weifashi/hios/raw/v0.0.5/release/hios-linux-amd64.tar.gz
# 	tar -zxf hios-linux-amd64.tar.gz
# 	rm -f hios-linux-amd64.tar.gz
# 	mkdir /usr/lib/weifashi
# 	mv ./release/hios-linux-amd64/hios /usr/lib/weifashi/hios
# 	rm -r ./release
# 	chmod +x /usr/lib/weifashi/hios