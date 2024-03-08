FROM golang:1.22 AS builder

ENV GO111MODULE=on \
    GOPROXY="https://goproxy.cn,direct" \
    APP_NAME="wechatbot"

WORKDIR $GOPATH/src/$APP_NAME

# manage dependencies
COPY go.mod ./
COPY go.sum ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /$APP_NAME ./cmd/wechatbot/wechatbot.go

FROM alpine:3.19 as runtime
# 换源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# Change TimeZone
RUN --mount=type=cache,id=apk,target=/var/cache/apk \
    apk add --no-cache tzdata nmap
ENV TZ=Asia/Shanghai
# Clean APK cache
RUN rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /wechatbot ./
COPY /cmd/wechatbot/wechatbot.yaml ./wechatbot.yaml
COPY /flyway/wechatbot ./flyway
RUN mkdir -p ./logs
CMD ["./wechatbot"]  