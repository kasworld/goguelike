
$BUILD_VER=$args[0]


$env:GOOS="js" 
$env:GOARCH="wasm" 
Write-Output "go build -o clientdata/wasmclientgl.wasm -ldflags `"-X main.Ver=${BUILD_VER}`" wasmclientgl.go"
go build -o clientdata/wasmclientgl.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclientgl.go
$env:GOOS=""
$env:GOARCH=""

