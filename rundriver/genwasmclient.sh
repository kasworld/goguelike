#!/usr/bin/env bash

BUILD_VER=${1}

# build 2d client

# rm wasmclient.wasm

# echo "GOOS=js GOARCH=wasm go build -o wasmclient.wasm -ldflags -X main.Ver=${BUILD_VER}"
# GOOS=js GOARCH=wasm go build -o wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go

# echo "move files"
# mv wasmclient.wasm clientdata/

# build gl client

rm wasmclientgl.wasm

echo "GOOS=js GOARCH=wasm go build -o wasmclientgl.wasm -ldflags -X main.Ver=${BUILD_VER}"
GOOS=js GOARCH=wasm go build -o wasmclientgl.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclientgl.go

echo "move files"
mv wasmclientgl.wasm clientdata/
