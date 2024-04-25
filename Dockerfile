##
## Build
##
FROM golang:latest AS build

WORKDIR /src

RUN go env -w GOPROXY=https://proxy.golang.com.cn,direct

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /file-server .

##
## Deploy
##
FROM alpine:latest 
# FROM golang:latest 
WORKDIR /

# 复制编译后的程序
COPY --from=build /file-server /file-server

# 复制配置文件
COPY config.toml config.toml

EXPOSE 11888

ENTRYPOINT ["/file-server"]