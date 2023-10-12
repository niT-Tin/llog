#!/bin/bash
#

if [[ $# -eq 0 ]]; then
  go test -v
else
  for t in "$@"; do
    echo -e "\e[1;32mRunning test: $t...\e[0m"
    go test -v -run "$t" -cover
  done
fi
