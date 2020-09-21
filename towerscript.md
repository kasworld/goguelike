# 지형 스크립트 설명 

기본적인 예는 goguelike/rundriver/serverdata/start.tower 를 참고 할것 

command 목록은 /enum/terraincmd.enum 을 사용해서 자동 생성되며 

실 처리 하는 부분은 /game/terrain/mapmake.go 에서 이루어 진다. 

script를 parse하는 것은 /lib/scriptparse 를 볼것. 

스크립트 명령어 /enum/terraincmd.enum 의 내용.

	# cmd argformat
	# args must same order with GetArgs function args
	# no need to same order with map script

	# create new terrain, must be 1st line except comment
	NewTerrain          name:string w:int h:int actturnboost:float

	# initial ao count 
	ActiveObjectsRand   count:int

	# minimum co count on floor
	CarryObjectsRand    count:int

	# add resource  
	ResourceAt              resource:ResourceType amount:int x:int y:int
	ResourceHLine           resource:ResourceType amount:int x:int w:int y:int
	ResourceVLine           resource:ResourceType amount:int x:int y:int h:int
	ResourceLine            resource:ResourceType amount:int x1:int y1:int x2:int y2:int
	ResourceRect            resource:ResourceType amount:int x:int w:int y:int h:int
	ResourceFillRect        resource:ResourceType amount:int x:int w:int y:int h:int
	ResourceFillEllipses    resource:ResourceType amount:int x:int w:int y:int h:int
	ResourceMazeWall        resource:ResourceType amount:int x:int y:int w:int h:int xn:int yn:int connerfill:bool
	ResourceMazeWalk        resource:ResourceType amount:int x:int y:int w:int h:int xn:int yn:int connerfill:bool
	ResourceRand            resource:ResourceType mean:int stddev:int repeat:int
	ResourceFromPNG         name:string

	# define terrain ageing
	# before start, Millisecond per ageing (0==no ageing), reset terrain to init state after n ageing
	ResourceAgeing          initrun:int msper:int resetaftern:int

	# TileFlag is comma seperrated tile list

	# add room
	AddRoom                 bgtile:TileFlag walltile:TileFlag terrace:bool x:int y:int w:int h:int
	AddRoomMaze             bgtile:TileFlag walltile:TileFlag terrace:bool x:int y:int w:int h:int xn:int yn:int connerfill:bool
	AddRoomsRand            bgtile:TileFlag walltile:TileFlag terrace:bool align:int count:int mean:int stddev:int
	ConnectRooms            tile:TileFlag connect:int allconnect:bool diagonal:bool

	# add tile
	TileAt              tile:TileFlag x:int y:int
	TileHLine           tile:TileFlag x:int w:int y:int
	TileVLine           tile:TileFlag x:int y:int h:int
	TileLine            tile:TileFlag x1:int y1:int x2:int y2:int
	TileRect            tile:TileFlag x:int w:int y:int h:int
	TileFillRect        tile:TileFlag x:int w:int y:int h:int
	TileFillEllipses    tile:TileFlag x:int w:int y:int h:int
	TileMazeWall        tile:TileFlag x:int y:int w:int h:int xn:int yn:int connerfill:bool
	TileMazeWalk        tile:TileFlag x:int y:int w:int h:int xn:int yn:int connerfill:bool

	# end define terrain tiles (above commands)
	FinalizeTerrain         

	# add fieldobj
	AddPortal               x:int y:int display:FieldObjDisplayType acttype:FieldObjActType PortalID:string DstPortalID:string message:string
	AddPortalRand                       display:FieldObjDisplayType acttype:FieldObjActType PortalID:string DstPortalID:string message:string
	AddPortalInRoom                     display:FieldObjDisplayType acttype:FieldObjActType PortalID:string DstPortalID:string message:string

	AddRecycler             x:int y:int display:FieldObjDisplayType message:string
	AddRecyclerRand         count:int   display:FieldObjDisplayType message:string
	AddRecyclerInRoom       count:int   display:FieldObjDisplayType message:string

	AddTrapTeleport         x:int y:int DstFloor:string message:string 
	AddTrapTeleportsRand    count:int   DstFloor:string message:string
	AddTrapTeleportsInRoom  count:int   DstFloor:string message:string

	AddTrap                 x:int y:int display:FieldObjDisplayType acttype:FieldObjActType message:string
	AddTrapsRand            count:int   display:FieldObjDisplayType acttype:FieldObjActType message:string
	AddTrapsInRoom          count:int   display:FieldObjDisplayType acttype:FieldObjActType message:string

	# add fieldobj rotate line attack
	AddRotateLineAttack           x:int y:int display:FieldObjDisplayType winglen:int wingcount:int degree:int perturn:int decay:DecayType message:string
	AddRotateLineAttackRand       count:int   display:FieldObjDisplayType winglen:int wingcount:int degree:int perturn:int decay:DecayType message:string
	AddRotateLineAttackInRoom     count:int   display:FieldObjDisplayType winglen:int wingcount:int degree:int perturn:int decay:DecayType message:string

	# add Mine
	AddMine           x:int y:int display:FieldObjDisplayType decay:DecayType message:string
	AddMineRand       count:int   display:FieldObjDisplayType decay:DecayType message:string
	AddMineInRoom     count:int   display:FieldObjDisplayType decay:DecayType message:string
