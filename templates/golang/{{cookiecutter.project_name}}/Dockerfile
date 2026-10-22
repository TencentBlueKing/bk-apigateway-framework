# -------------- builder container --------------
FROM golang:1.26 as builder

WORKDIR /go/src/

ARG VERSION

COPY go.mod .
COPY go.sum .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.cloud.tencent.com/go/,direct

RUN go mod download

COPY . .

RUN make build

# -------------- runner container --------------
FROM alpine:3.20 AS runner

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories

RUN apk --update --no-cache add bash

WORKDIR /app

COPY --from=builder /go/src/{{cookiecutter.project_name}} /usr/bin/{{cookiecutter.project_name}}

COPY --from=builder /go/src/sync_apigateway.sh /app/sync_apigateway.sh
# docs
RUN mkdir -p /app/docs

COPY --from=builder /go/src/docs /app/docs

# DOC_FILE_BASE_DIR：generate_*_yaml / sync_apigateway 命令读取 swagger.json 的目录
ENV DOC_FILE_BASE_DIR=/app/docs
# API_DOC_FILE_BASE_DIR：webserver 提供 /swagger.json 的目录（blueapps-go 配置项）
ENV API_DOC_FILE_BASE_DIR=/app/docs

# logs
RUN mkdir -p /app/v3logs

ENV LOG_FILE_BASE_DIR=/app/v3logs

