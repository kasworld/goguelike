#!/usr/bin/env bash

# find -name \*.go ! -name \*_gen.go ! -name \*_string.go ! -name \*_test.go | xargs wc
# find -name \*.go ! -name \*_gen.go ! -name \*_string.go ! -name \*_test.go ! -path ./vendor/\* | xargs wc

BuildBin() {
    local srcfile=${1}
    local dstdir=${2}
    local dstfile=${3}
    local args="-X main.Ver=${BUILD_VER}"

    echo "go build -i -o ${dstdir}/${dstfile} -ldflags "${args}" ${srcfile}"

    mkdir -p ${dstdir}
    go build -i -o ${dstdir}/${dstfile} -ldflags "${args}" ${srcfile}

    if [ ! -f "${dstdir}/${dstfile}" ]; then
        echo "${dstdir}/${dstfile} build fail, build file: ${srcfile}"
        exit 1
    fi
    strip "${dstdir}/${dstfile}"
}

################################################################################
DATESTR=`date -Iseconds`
GITSTR=`git rev-parse HEAD`
BUILD_VER=${DATESTR}_${GITSTR}_release
echo "Build Version:" ${BUILD_VER}


BIN_DIR="bin"
SRC_DIR="rundriver"

echo ${BUILD_VER} > ${BIN_DIR}/BUILD

BuildBin ${SRC_DIR}/towerserver.go ${BIN_DIR} towerserver
BuildBin ${SRC_DIR}/groundserver.go ${BIN_DIR} groundserver
BuildBin ${SRC_DIR}/multiclient.go ${BIN_DIR} multiclient
BuildBin ${SRC_DIR}/textclient.go ${BIN_DIR} textclient

cd rundriver
./genwasmclient.sh ${BUILD_VER}
cd ..

echo cp -r rundriver/serverdata ${BIN_DIR}
cp -r rundriver/serverdata ${BIN_DIR}
echo cp -r rundriver/clientdata ${BIN_DIR}
cp -r rundriver/clientdata ${BIN_DIR}

