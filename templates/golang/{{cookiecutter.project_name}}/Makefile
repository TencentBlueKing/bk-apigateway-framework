.PHONY: tidy build test

ifdef VERSION
    VERSION=${VERSION}
else
    VERSION=$(shell git describe --always 2>/dev/null || echo "--")
endif

GITCOMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "--")
BUILDTIME=${shell date +%Y-%m-%dT%H:%M:%S%z}

LDFLAGS="-X bk.tencent.com/{{cookiecutter.project_name}}/pkg/version.AppVersion=${VERSION} \
	-X bk.tencent.com/{{cookiecutter.project_name}}/pkg/version.GitCommit=${GITCOMMIT} \
	-X bk.tencent.com/{{cookiecutter.project_name}}/pkg/version.BuildTime=${BUILDTIME}"

# go mod tidy
tidy:
	go mod tidy

# build executable binary
build: tidy
	CGO_ENABLED=0 go build -ldflags ${LDFLAGS} -o {{cookiecutter.project_name}} ./main.go

# generate swagger api doc with swag
doc: swag
	$(SWAG) fmt
	$(SWAG) init --parseDependency --parseDepth 3

# run unittest
test: tidy
	go test ./...

# build docker image
docker-build:
	docker build -f ./Dockerfile -t {{cookiecutter.project_name}}:${VERSION} .

.PHONY: fmt
fmt: golines gofumpt ## 执行 golines，gofumpt ...
	$(GOLINES) ./ -m 119 -w --base-formatter gofmt --no-reformat-tags
	$(GOFUMPT) -l -w .

.PHONY: vet
vet: ## 执行 go vet ./...
	go vet ./...

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLINES ?= $(LOCALBIN)/golines
GOFUMPT ?= $(LOCALBIN)/gofumpt
SWAG ?= $(LOCALBIN)/swag

.PHONY: golines
golines: $(GOLINES) ## 安装 golines 二进制
$(GOLINES): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/segmentio/golines@v0.12.2

.PHONY: gofumpt
gofumpt: $(GOFUMPT) ## 安装 gofumpt 二进制
$(GOFUMPT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install mvdan.cc/gofumpt@v0.6.0

.PHONY: swag
swag: $(SWAG) ## 安装 swag 二进制
$(SWAG): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/swaggo/swag/cmd/swag@v1.16.3
