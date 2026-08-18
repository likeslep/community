# Community Platform — 常用构建 / 测试 / 质量命令

GO      ?= go
BIN_DIR := bin

.PHONY: build test vet lint fmt clean

## build: 编译所有包
build:
	$(GO) build ./...

## test: 运行单元测试
test:
	$(GO) test ./...

## vet: 运行 go vet
vet:
	$(GO) vet ./...

## lint: 运行 golangci-lint（需先安装 golangci-lint）
lint:
	golangci-lint run ./...

## fmt: 格式化代码
fmt:
	gofmt -w .
	goimports -w .

## clean: 清理构建产物
clean:
	rm -rf $(BIN_DIR)
