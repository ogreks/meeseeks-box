#!/bin/bash

echo "安装 golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

echo "安装 goimports..."
go install golang.org/x/tools/cmd/goimports@latest