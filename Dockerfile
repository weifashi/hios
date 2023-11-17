# docker build --pull --rm -f "Dockerfile" -t hios:latest "."

# Stage 1: Build webui
FROM --platform=${BUILDPLATFORM} docker.io/node:16.2.0 as web-builder

ENV NODE_ENV=production

WORKDIR /work/web

RUN --mount=type=bind,target=/work/web/package.json,src=./web/package.json \
    --mount=type=bind,target=/work/web/package-lock.json,src=./web/package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci --include=dev

COPY ./web /work/web/

RUN npm run build

# Stage 2: Build go proxy
FROM --platform=linux/amd64 docker.io/golang:1.18 AS go-builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

ARG GOOS=$TARGETOS
ARG GOARCH=$TARGETARCH

WORKDIR /go/src/hios.io

RUN --mount=type=bind,target=/go/src/hios.io/go.mod,src=./go.mod \
    --mount=type=bind,target=/go/src/hios.io/go.sum,src=./go.sum \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./ /go/src/hios.io/

ENV CGO_ENABLED=1

CMD make release OS_ARCHS=linux/amd64


# Stage 3: Build go proxy
FROM --platform=linux/arm64 docker.io/golang:1.18 AS go-builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

ARG GOOS=$TARGETOS
ARG GOARCH=$TARGETARCH

WORKDIR /go/src/hios.io

RUN --mount=type=bind,target=/go/src/hios.io/go.mod,src=./go.mod \
    --mount=type=bind,target=/go/src/hios.io/go.sum,src=./go.sum \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./ /go/src/hios.io/

ENV CGO_ENABLED=1

CMD make release OS_ARCHS=linux/arm64