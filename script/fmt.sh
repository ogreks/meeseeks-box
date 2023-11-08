#!/bin/bash

# shellcheck disable=SC2044
for item in $(find . -type f -name '*.go' -not -path './.idea/*' -not -path './.vscode/*'); do
  goimports -l -w "$item";
done