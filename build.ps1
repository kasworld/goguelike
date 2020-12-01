
################################################################################
cd lib
echo "genlog -leveldatafile ./g2log/g2log.data -packagename g2log "
genlog -leveldatafile ./g2log/g2log.data -packagename g2log 
cd ..

################################################################################
$PROTOCOL_T2G_VERSION=makesha256sum protocol_t2g/*.enum protocol_t2g/t2g_obj/protocol_*.go
echo "Protocol T2G Version: ${PROTOCOL_T2G_VERSION}"
echo "genprotocol -ver=${PROTOCOL_T2G_VERSION} -basedir=protocol_t2g -prefix=t2g -statstype=int"
genprotocol -ver="${PROTOCOL_T2G_VERSION}" -basedir=protocol_t2g -prefix=t2g -statstype=int
goimports -w protocol_t2g

################################################################################
$PROTOCOL_C2T_VERSION=makesha256sum protocol_c2t/*.enum protocol_c2t/c2t_obj/protocol_*.go
echo "Protocol C2T Version: ${PROTOCOL_C2T_VERSION}"
echo "genprotocol -ver=${PROTOCOL_C2T_VERSION} -basedir=protocol_c2t -prefix=c2t -statstype=int"
genprotocol -ver="${PROTOCOL_C2T_VERSION}" -basedir=protocol_c2t -prefix=c2t -statstype=int
goimports -w protocol_c2t

################################################################################
# generate enum
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

$Data_VERSION=makesha256sum config/gameconst/*.go config/gamedata/*.go enum/*.enum
echo "Data Version: ${Data_VERSION}"
mkdir -ErrorAction SilentlyContinue config/dataversion
echo "package dataversion
const DataVersion = `"${Data_VERSION}`" 
" > config/dataversion/dataversion_gen.go 


################################################################################
$DATESTR=Get-Date -UFormat '+%Y-%m-%dT%H:%M:%S%Z:00'
$GITSTR=git rev-parse HEAD
################################################################################
# build bin

$BIN_DIR="bin"
$SRC_DIR="rundriver"

mkdir -ErrorAction SilentlyContinue "${BIN_DIR}"

# build bin here
$BUILD_VER="${DATESTR}_${GITSTR}_release_windows"
echo "Build Version: ${BUILD_VER}"
echo ${BUILD_VER} > ${BIN_DIR}/BUILD_windows
go build -o "${BIN_DIR}\towerserver.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\towerserverwin.go"
go build -o "${BIN_DIR}\multiclient.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\multiclient.go"
go build -o "${BIN_DIR}\textclient.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\textclient.go"

$BUILD_VER="${DATESTR}_${GITSTR}_release_linux"
echo "Build Version: ${BUILD_VER}"
echo ${BUILD_VER} > ${BIN_DIR}/BUILD_linux
$env:GOOS="linux" 
go build -o "${BIN_DIR}\towerserver" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\towerserver.go"
go build -o "${BIN_DIR}\multiclient" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\multiclient.go"
go build -o "${BIN_DIR}\textclient" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\textclient.go"
$env:GOOS=""

$BUILD_VER="${DATESTR}_${GITSTR}_release_wasm"
cd rundriver
./genwasmclient.ps1 ${BUILD_VER}
cd ..

echo "cp -r rundriver/serverdata ${BIN_DIR}"
Copy-Item -Force -r rundriver/serverdata ${BIN_DIR}
echo "cp -r rundriver/clientdata ${BIN_DIR}"
Copy-Item -Force -r rundriver/clientdata ${BIN_DIR}

