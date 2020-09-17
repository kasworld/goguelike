# 지형 스크립트 설명 

기본적인 예는 goguelike/rundriver/serverdata/start.tower 를 참고 할것 

command 목록은 /enum/terraincmd.enum 을 사용해서 자동 생성되며 

실 처리 하는 부분은 /game/terrain/mapmake.go 에서 이루어 진다. 

script를 parse하는 것은 /lib/scriptparse 를 볼것. 

스크립트 명령어 /enum/terraincmd.enum 의 내용.

	# cmd argformat

	NewTerrain  name:string w:int h:int actturnboost:float

	# initial ao count 
	AddActiveObjectRand count:int

	# co count on floor
	AddCarryObjectRand count:int

	# add resource  
	ResourceMazeWall        resource:ResourceType amount:int x:int y:int w:int h:int xn:int yn:int connerfill:bool
	ResourceMazeWalk        resource:ResourceType amount:int x:int y:int w:int h:int xn:int yn:int connerfill:bool
	ResourceRand            resource:ResourceType mean:int stddev:int repeat:int
	ResourceAt              resource:ResourceType amount:int x:int y:int
	ResourceHLine           resource:ResourceType amount:int x:int w:int y:int
	ResourceVLine           resource:ResourceType amount:int x:int y:int h:int
	ResourceLine            resource:ResourceType amount:int x1:int y1:int x2:int y2:int
	ResourceRect            resource:ResourceType amount:int x:int w:int y:int h:int
	ResourceFillRect        resource:ResourceType amount:int x:int w:int y:int h:int
	ResourceFillEllipses    resource:ResourceType amount:int x:int w:int y:int h:int
	ResourceFromPNG         name:string
	ResourceAgeing          initrun:int msper:int resetaftern:int

	# add room
	AddRoom                 bgtile:TileType walltile:TileType terrace:bool x:int y:int w:int h:int
	AddRoomMaze             bgtile:TileType walltile:TileType terrace:bool x:int y:int w:int h:int xn:int yn:int connerfill:bool
	AddRoomsRand            bgtile:TileType walltile:TileType terrace:bool align:int count:int mean:int stddev:int
	ConnectRooms            tile:TileType connect:int allconnect:bool diagonal:bool

	# add tile
	TileMazeWall        tile:TileType x:int y:int w:int h:int xn:int yn:int connerfill:bool
	TileMazeWalk        tile:TileType x:int y:int w:int h:int xn:int yn:int connerfill:bool
	TileAt              tile:TileType x:int y:int
	TileHLine           tile:TileType x:int w:int y:int
	TileVLine           tile:TileType x:int y:int h:int
	TileLine            tile:TileType x1:int y1:int x2:int y2:int
	TileRect            tile:TileType x:int w:int y:int h:int
	TileFillRect        tile:TileType x:int w:int y:int h:int
	TileFillEllipses    tile:TileType x:int w:int y:int h:int

	FinalizeTerrain         

	# add fieldobj
	AddPortal               x:int y:int display:FieldObjDisplayType acttype:FieldObjActType PortalID:string DstPortalID:string message:string
	AddPortalRand           display:FieldObjDisplayType acttype:FieldObjActType PortalID:string DstPortalID:string message:string
	AddPortalInRoom         display:FieldObjDisplayType acttype:FieldObjActType PortalID:string DstPortalID:string message:string
	AddRecycler             x:int y:int display:FieldObjDisplayType message:string
	AddRecyclerRand         display:FieldObjDisplayType count:int message:string
	AddRecyclerInRoom       display:FieldObjDisplayType count:int message:string
	AddTrapTeleport         x:int y:int DstFloor:string message:string 
	AddTrapTeleportsRand    DstFloor:string count:int message:string
	AddTrapTeleportsInRoom  DstFloor:string count:int message:string
	AddTrap                 x:int y:int display:FieldObjDisplayType acttype:FieldObjActType message:string
	AddTrapsRand            display:FieldObjDisplayType acttype:FieldObjActType count:int message:string
	AddTrapsInRoom          display:FieldObjDisplayType acttype:FieldObjActType count:int message:string

	# add fieldobj LightHouse, GateKeeper
	AddRotateLineAttack           x:int y:int display:FieldObjDisplayType winglen:int wingcount:int degree:int perturn:int decay:DecayType message:string
	AddRotateLineAttackRand       display:FieldObjDisplayType winglen:int wingcount:int degree:int perturn:int decay:DecayType count:int message:string
	AddRotateLineAttackInRoom     display:FieldObjDisplayType winglen:int wingcount:int degree:int perturn:int decay:DecayType count:int message:string

	# add Mine
	AddMine           x:int y:int display:FieldObjDisplayType decay:DecayType message:string
	AddMineRand       display:FieldObjDisplayType decay:DecayType count:int message:string
	AddMineInRoom     display:FieldObjDisplayType decay:DecayType count:int message:string
