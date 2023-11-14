#!/bin/bash

echo "install golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

echo "install goimports..."
go install golang.org/x/tools/cmd/goimports@latest

echo "install wire..."
go install github.com/google/wire/cmd/wire@latest

## exec wire
make api-wire