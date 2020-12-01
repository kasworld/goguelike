#!/usr/bin/env bash

# find -name \*.go ! -name \*_gen.go ! -name \*_string.go ! -name \*_test.go | xargs wc
# find -name \*.go ! -name \*_gen.go ! -name \*_string.go ! -name \*_test.go ! -path ./vendor/\* | xargs wc

################################################################################
echo "genlog -leveldatafile ./g2log/g2log.data -packagename g2log "
cd lib
genlog -leveldatafile ./g2log/g2log.data -packagename g2log 
cd ..

################################################################################
ProtocolT2GFiles="protocol_t2g/*.enum protocol_t2g/t2g_obj/protocol_*.go"
PROTOCOL_T2G_VERSION=`makesha256sum ${ProtocolT2GFiles}`
echo "Protocol T2G Version: ${PROTOCOL_T2G_VERSION}"
echo "genprotocol -ver=${PROTOCOL_T2G_VERSION} -basedir=protocol_t2g -prefix=t2g -statstype=int"
genprotocol -ver=${PROTOCOL_T2G_VERSION} -basedir=protocol_t2g -prefix=t2g -statstype=int
goimports -w protocol_t2g

################################################################################
ProtocolC2TFiles="protocol_c2t/*.enum protocol_c2t/c2t_obj/protocol_*.go"
PROTOCOL_C2T_VERSION=`makesha256sum ${ProtocolC2TFiles}`
echo "Protocol C2T Version: ${PROTOCOL_C2T_VERSION}"
echo "genprotocol -ver=${PROTOCOL_C2T_VERSION} -basedir=protocol_c2t -prefix=c2t -statstype=int"
genprotocol -ver=${PROTOCOL_C2T_VERSION} -basedir=protocol_c2t -prefix=c2t -statstype=int
goimports -w protocol_c2t

################################################################################
echo "generate enums"
genenum -typename=AchieveType -packagename=achievetype -basedir=enum -vectortype=float64
genenum -typename=AIPlan -packagename=aiplan -basedir=enum -vectortype=int
genenum -typename=ActiveObjType -packagename=aotype -basedir=enum -vectortype=int
genenum -typename=CarryingObjectType -packagename=carryingobjecttype -basedir=enum -vectortype=int
genenum -typename=ClientControlType -packagename=clientcontroltype -basedir=enum 
genenum -typename=Condition -packagename=condition -basedir=enum -flagtype=uint16 -vectortype=int
genenum -typename=DangerType -packagename=dangertype -basedir=enum -vectortype=int
genenum -typename=DecayType -packagename=decaytype -basedir=enum
genenum -typename=EquipSlotType -packagename=equipslottype -basedir=enum -vectortype=int
genenum -typename=FactionType -packagename=factiontype -basedir=enum -vectortype=int
genenum -typename=FieldObjActType -packagename=fieldobjacttype -basedir=enum -vectortype=int
genenum -typename=FieldObjDisplayType -packagename=fieldobjdisplaytype -basedir=enum
genenum -typename=PotionType -packagename=potiontype -basedir=enum -vectortype=int
genenum -typename=ResourceType -packagename=resourcetype -basedir=enum -vectortype=int
genenum -typename=RespawnType -packagename=respawntype -basedir=enum 
genenum -typename=ScrollType -packagename=scrolltype -basedir=enum -vectortype=int
genenum -typename=StatusOpType -packagename=statusoptype -basedir=enum
genenum -typename=TerrainCmd -packagename=terraincmd -basedir=enum -vectortype=int
genenum -typename=Tile -packagename=tile -basedir=enum -flagtype=uint16 -vectortype=int
genenum -typename=TowerAchieve -packagename=towerachieve -basedir=enum -vectortype=float64
genenum -typename=TurnResultType -packagename=turnresulttype -basedir=enum
genenum -typename=Way9Type -packagename=way9type -basedir=enum 

goimports -w enum

################################################################################
# change to use gob
# GenMSGP() {
#     local gosrc="${2}"
#     local basedir="${1}"
#     rm ${basedir}/"${gosrc}"_gen.go
#     # rm ${basedir}/"${gosrc}"_gen_test.go
#     msgp -file ${basedir}/"${gosrc}".go -o ${basedir}/"${gosrc}"_gen.go -tests=0 
# }
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

GameDataFiles="config/gameconst/*.go config/gamedata/*.go enum/*.enum"
Data_VERSION=`makesha256sum ${GameDataFiles}`
echo "Data Version: ${Data_VERSION}"
mkdir -p config/dataversion
echo "package dataversion
const DataVersion = \"${Data_VERSION}\"
" > config/dataversion/dataversion_gen.go 



DATESTR=`date -Iseconds`
GITSTR=`git rev-parse HEAD`

BIN_DIR="bin"
SRC_DIR="rundriver"

mkdir -p ${BIN_DIR}

BUILD_VER=${DATESTR}_${GITSTR}_release_linux
echo "Build Version:" ${BUILD_VER}
echo ${BUILD_VER} > ${BIN_DIR}/BUILD_linux
go build -o "${BIN_DIR}/groundserver" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}/groundserver.go"
go build -o "${BIN_DIR}/towerserver" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}/towerserver.go"
go build -o "${BIN_DIR}/multiclient" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}/multiclient.go"
go build -o "${BIN_DIR}/textclient" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}/textclient.go"


BUILD_VER=${DATESTR}_${GITSTR}_release_windows
echo "Build Version:" ${BUILD_VER}
echo ${BUILD_VER} > ${BIN_DIR}/BUILD_windows
GOOS=windows go build -o "${BIN_DIR}/towerserver.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}/towerserverwin.go"
GOOS=windows go build -o "${BIN_DIR}/multiclient.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}/multiclient.go"
GOOS=windows go build -o "${BIN_DIR}/textclient.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}/textclient.go"


BUILD_VER=${DATESTR}_${GITSTR}_release_wasm
echo "Build Version:" ${BUILD_VER}
cd rundriver
./genwasmclient.sh ${BUILD_VER}
cd ..

echo cp -r rundriver/serverdata ${BIN_DIR}
cp -r rundriver/serverdata ${BIN_DIR}
echo cp -r rundriver/clientdata ${BIN_DIR}
cp -r rundriver/clientdata ${BIN_DIR}

