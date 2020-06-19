#!/usr/bin/env bash

# find -name \*.go ! -name \*_gen.go ! -name \*_string.go ! -name \*_test.go | xargs wc
# find -name \*.go ! -name \*_gen.go ! -name \*_string.go ! -name \*_test.go ! -path ./vendor/\* | xargs wc

GenMSGP() {
    local gosrc="${2}"
    local basedir="${1}"
    rm ${basedir}/"${gosrc}"_gen.go
    # rm ${basedir}/"${gosrc}"_gen_test.go
    msgp -file ${basedir}/"${gosrc}".go -o ${basedir}/"${gosrc}"_gen.go -tests=0 
}

################################################################################
cd lib
genlog -leveldatafile ./g2log/g2log.data -packagename g2log 
cd ..

################################################################################
ProtocolT2GFiles="protocol_t2g/*.enum \
protocol_t2g/t2g_obj/protocol_noti.go \
protocol_t2g/t2g_obj/protocol_cmd.go \
"
PROTOCOL_T2G_VERSION=`cat ${ProtocolT2GFiles}| sha256sum | awk '{print $1}'`

cd protocol_t2g
genprotocol -ver=${PROTOCOL_T2G_VERSION} \
    -basedir=. \
    -prefix=t2g -statstype=int
goimports -w .
cd ..

################################################################################
ProtocolC2TFiles="protocol_c2t/*.enum \
protocol_c2t/c2t_obj/protocol_objects.go \
protocol_c2t/c2t_obj/protocol_noti.go \
protocol_c2t/c2t_obj/protocol_admin.go \
protocol_c2t/c2t_obj/protocol_aoact.go \
protocol_c2t/c2t_obj/protocol_cmd.go \
"
PROTOCOL_C2T_VERSION=`cat ${ProtocolC2TFiles}| sha256sum | awk '{print $1}'`

cd protocol_c2t
genprotocol -ver=${PROTOCOL_C2T_VERSION} \
    -basedir=. \
    -prefix=c2t -statstype=int

goimports -w .
cd ..

################################################################################
echo genenum

genenum -typename=Way9Type -packagename=way9type -basedir=enum 
genenum -typename=ActiveObjType -packagename=aotype -basedir=enum -vectortype=int
genenum -typename=CarryingObjectType -packagename=carryingobjecttype -basedir=enum -vectortype=int
genenum -typename=FieldObjActType -packagename=fieldobjacttype -basedir=enum -vectortype=int
genenum -typename=FieldObjDisplayType -packagename=fieldobjdisplaytype -basedir=enum
genenum -typename=Condition -packagename=condition -basedir=enum -flagtype=uint16 -vectortype=int
genenum -typename=PotionType -packagename=potiontype -basedir=enum -vectortype=int
genenum -typename=ScrollType -packagename=scrolltype -basedir=enum -vectortype=int
genenum -typename=AchieveType -packagename=achievetype -basedir=enum -vectortype=float64
genenum -typename=ResourceType -packagename=resourcetype -basedir=enum -vectortype=int
genenum -typename=TileOpType -packagename=tileoptype -basedir=enum 
genenum -typename=EquipSlotType -packagename=equipslottype -basedir=enum -vectortype=int
genenum -typename=StatusOpType -packagename=statusoptype -basedir=enum
genenum -typename=TurnResultType -packagename=turnresulttype -basedir=enum
genenum -typename=Tile -packagename=tile -basedir=enum -flagtype=uint16 -vectortype=int
genenum -typename=TowerAchieve -packagename=towerachieve -basedir=enum -vectortype=float64
genenum -typename=ClientControlType -packagename=clientcontroltype -basedir=enum 
genenum -typename=FactionType -packagename=factiontype -basedir=enum -vectortype=int
genenum -typename=AIPlan -packagename=aiplan -basedir=enum -vectortype=int

cd enum
goimports -w .
cd ..

# change to use gob

# GenMSGP "enum/way9type" way9type_gen
# GenMSGP "enum/carryingobjecttype" carryingobjecttype_gen
# GenMSGP "enum/fieldobjacttype" fieldobjacttype_gen
# GenMSGP "enum/fieldobjdisplaytype" fieldobjdisplaytype_gen
# GenMSGP "enum/potiontype" potiontype_gen
# GenMSGP "enum/scrolltype" scrolltype_gen
# GenMSGP "enum/equipslottype" equipslottype_gen
# GenMSGP "enum/turnresulttype" turnresulttype_gen
# GenMSGP "enum/factiontype" factiontype_gen
# GenMSGP "enum/aiplan" aiplan_gen
# GenMSGP "enum/tile_flag" tile_flag_gen
# GenMSGP "enum/condition_flag" condition_flag_gen

################################################################################
# GenMSGP "vendor/github.com/kasworld/htmlcolors" color24

# GenMSGP "protocol_c2t/c2t_error" error_gen
# GenMSGP "protocol_c2t/c2t_idcmd" command_gen
# GenMSGP "protocol_c2t/c2t_idnoti" noti_gen
# GenMSGP "protocol_c2t/c2t_obj" protocol_objects
# GenMSGP "protocol_c2t/c2t_obj" protocol_noti
# GenMSGP "protocol_c2t/c2t_obj" protocol_admin
# GenMSGP "protocol_c2t/c2t_obj" protocol_aoact
# GenMSGP "protocol_c2t/c2t_obj" protocol_cmd
# GenMSGP "config/viewportdata" viewportdata
# GenMSGP "lib/g2id" g2id
# GenMSGP "game/aoactreqrsp" aoactreqrsp
# GenMSGP "game/bias" bias
# GenMSGP "game/tilearea" tilearea

GameDataFiles="
config/gameconst/gameconst.go \
config/gameconst/serviceconst.go \
config/gamedata/*.go \
enum/*.enum \
"
Data_VERSION=`cat ${GameDataFiles}| sha256sum | awk '{print $1}'`

echo "
package gameconst

const DataVersion = \"${Data_VERSION}\"
" > config/gameconst/dataversion_gen.go 

echo "Protocol T2G Version:" ${PROTOCOL_T2G_VERSION}
echo "Protocol C2T Version:" ${PROTOCOL_C2T_VERSION}
echo "Data Version:" ${Data_VERSION}



DATESTR=`date -Iseconds`
GITSTR=`git rev-parse HEAD`
BUILD_VER=${DATESTR}_${GITSTR}_release
echo "Build Version:" ${BUILD_VER}


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

