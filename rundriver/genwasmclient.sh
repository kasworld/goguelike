#!/usr/bin/env bash

BUILD_VER=${1}

echo "GOOS=js GOARCH=wasm go build -o clientdata/wasmclientgl.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclientgl.go"
GOOS=js GOARCH=wasm go build -o clientdata/wasmclientgl.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclientgl.go
