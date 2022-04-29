#!/bin/bash

exe() { echo "\$ ${@/eval/}" ; "$@" ; }

for path in `find -type f -name "*.proto"`
do
    exe eval "protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $path"
done
