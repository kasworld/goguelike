#!/usr/bin/env bash

$BUILD_VER=$args[0]

# build gl client

Remove-Item wasmclientgl.wasm

Write-Output "GOOS=js GOARCH=wasm go build -o wasmclientgl.wasm -ldflags -X main.Ver=${BUILD_VER}"
$env:GOOS="js" 
$env:GOARCH="wasm" 
go build -o wasmclientgl.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclientgl.go
$env:GOOS=""
$env:GOARCH=""

Write-Output "move files"
Move-Item wasmclientgl.wasm clientdata/
