export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

MODULE = fileWarehouse

PORT 			:= 3358
VERSION			:= $(shell git describe --tags --always --match="v*" 2> /dev/null || echo v0.0.0)
VERSION_HASH	:= $(shell git rev-parse --short HEAD)

GOCGO 	:= env CGO_ENABLED=1
LDFLAGS	:= -s -w -X "$(MODULE)/config.Version=$(VERSION)" -X "$(MODULE)/config.CommitSHA=$(VERSION_HASH)"

run: build
	./main --mode debug
	
build:
	cd web && npm i && npm run build && cd ../ 
	$(GOCGO) go build -trimpath -ldflags "$(LDFLAGS)" -o hios

release:
	$(GOCGO) GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./$(MODULE)-$(VERSION)-linux-amd64/$(MODULE)
	$(GOCGO) GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc-10 go build -trimpath -ldflags "$(LDFLAGS)" -o ./$(MODULE)-$(VERSION)-linux-arm64/$(MODULE)
	@for arch in amd64 arm64; \
	do \
		cp install/* $(MODULE)-$(VERSION)-linux-$$arch; \
		tar zcf $(MODULE)-$(VERSION)-linux-$$arch.tar.gz $(MODULE)-$(VERSION)-linux-$$arch; \
	done

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
