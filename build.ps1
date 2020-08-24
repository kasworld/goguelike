
################################################################################
$DATESTR=Get-Date -UFormat '+%Y-%m-%dT%H:%M:%S%Z:00'
$GITSTR=git rev-parse HEAD
$BUILD_VER="${DATESTR}_${GITSTR}_release_windows"
Write-Output "Build Version: ${BUILD_VER}"

################################################################################

$BIN_DIR="bin"
$SRC_DIR="rundriver"

Write-Output ${BUILD_VER} > ${BIN_DIR}/BUILD_windows

# build bin here
go build -o "${BIN_DIR}\towerserver.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\towerserverwin.go"


Set-Location rundriver
./genwasmclient.ps1 ${BUILD_VER}
Set-Location ..

Write-Output "cp -r rundriver/serverdata ${BIN_DIR}"
Copy-Item -Force -r rundriver/serverdata ${BIN_DIR}
Write-Output "cp -r rundriver/clientdata ${BIN_DIR}"
Copy-Item -Force -r rundriver/clientdata ${BIN_DIR}

