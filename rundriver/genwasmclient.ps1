#!/usr/bin/env bash

BUILD_VER=${1}

# build gl client

rm wasmclientgl.wasm

echo "GOOS=js GOARCH=wasm go build -o wasmclientgl.wasm -ldflags -X main.Ver=${BUILD_VER}"
GOOS=js GOARCH=wasm go build -o wasmclientgl.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclientgl.go

echo "move files"
mv wasmclientgl.wasm clientdata/
