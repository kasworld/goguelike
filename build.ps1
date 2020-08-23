
################################################################################
cd lib
genlog -leveldatafile ./g2log/g2log.data -packagename g2log 
cd ..

################################################################################
$PROTOCOL_T2G_VERSION=makesha256sum protocol_t2g/*.enum protocol_t2g/t2g_obj/protocol_noti.go protocol_t2g/t2g_obj/protocol_cmd.go
echo  "genprotocol -ver=${PROTOCOL_T2G_VERSION} -basedir=protocol_t2g -prefix=t2g -statstype=int"
genprotocol -ver="${PROTOCOL_T2G_VERSION}" -basedir=protocol_t2g -prefix=t2g -statstype=int
cd protocol_t2g
goimports -w .
cd ..
echo "Protocol T2G Version:" ${PROTOCOL_T2G_VERSION}

################################################################################
$PROTOCOL_C2T_VERSION=makesha256sum protocol_c2t/*.enum protocol_c2t/c2t_obj/protocol_objects.go protocol_c2t/c2t_obj/protocol_noti.go protocol_c2t/c2t_obj/protocol_admin.go protocol_c2t/c2t_obj/protocol_aoact.go protocol_c2t/c2t_obj/protocol_cmd.go
echo "genprotocol -ver=${PROTOCOL_C2T_VERSION} -basedir=protocol_c2t -prefix=c2t -statstype=int"
genprotocol -ver="${PROTOCOL_C2T_VERSION}" -basedir=protocol_c2t -prefix=c2t -statstype=int
cd protocol_c2t
goimports -w .
cd ..
echo "Protocol C2T Version:" ${PROTOCOL_C2T_VERSION}

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
genenum -typename=TerrainCmd -packagename=terraincmd -basedir=enum -vectortype=int

cd enum
goimports -w .
cd ..


$Data_VERSION=makesha256sum config/gameconst/gameconst.go config/gameconst/serviceconst.go config/gamedata/*.go enum/*.enum

echo "Data Version:" ${Data_VERSION}

################################################################################
