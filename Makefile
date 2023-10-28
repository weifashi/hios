export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

MODULE = hios

VERSION			:= $(shell git tag | tail -1 2> /dev/null || echo v0.0.1)
VERSION_HASH	:= $(shell git rev-parse --short HEAD)
NEXT_VERSION    := $(shell git tag | tail -1 | awk -F. -v OFS=. '{$$NF++; print}')

GOCGO 			:= env CGO_ENABLED=1
LDFLAGS			:= -s -w -X "$(MODULE)/config.Version=$(VERSION)" -X "$(MODULE)/config.CommitSHA=$(VERSION_HASH)"
OS_ARCHS		:=linux:amd64
# OS_ARCHS		:=darwin:amd64 darwin:arm64 linux:amd64 linux:arm64

## run
.PHONY: run
run:
	go run main.go

## dev
.PHONY: dev
dev:
	lsof -i :3377 | grep node | awk '{print $$2}' | xargs kill -9
	cd web && nohup npm run dev > ../output.log >&1 & cd ../ 
	${HOME}/go/bin/fresh -c ./fresh.conf

## build
.PHONY: build
build:
	go build -o ./$(MODULE)

## build-run
.PHONY: build-run
build-run: build
	./$(MODULE)

## build-all
.PHONY: build-all
build-all: | ; $(info $(M) building all…)
	$(shell mkdir -p release)
	@$(foreach n, $(OS_ARCHS),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		$(GOCGO) GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build -trimpath -ldflags "$(LDFLAGS)" -o ./release/$(MODULE)_$${target_suffix};\
	)

## release
.PHONY: release
release: | ; $(info $(M) release all…)
	@$(foreach n, $(OS_ARCHS),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		target_suffix=$${os}_$${arch};\
		$(GOCGO) GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build  -o ./release/$(MODULE);\
		tar zcf ./release/$(MODULE)_$${target_suffix}.tar.gz ./release/$(MODULE);\
		rm -r ./release/$(MODULE);\
	)

## docker-release
.PHONY: docker-release
docker-release: | ; $(info $(M) release all…)
	docker run --rm -v "${PWD}":/myapp -w /myapp --platform linux/amd64 golang:1.20 bash -c "make release"

## translate
.PHONY: translate
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
# wget https://gitee.com/weifashi/hios/raw/v0.0.5/release/hios_linux_amd64.tar.gz
# 	tar -zxf hios_linux_amd64.tar.gz
# 	rm -f hios_linux_amd64.tar.gz
# 	mkdir /usr/lib/weifashi
# 	mv ./release/hios /usr/lib/weifashi/hios
# 	rm -r ./release
# 	chmod +x /usr/lib/weifashi/hios