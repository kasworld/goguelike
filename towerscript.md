# 지형 스크립트 설명 

기본적인 예는 goguelike/rundriver/serverdata/start.tower 를 참고 할것 

스크립트 명령어 

	NewTerrain :             name:string w:int h:int ao:int po:int actturnboost:float
	ResourceMazeWall :       resource:TileRsc_Type amount:int xn:int yn:int connerfill:bool
	ResourceMazeWalk :       resource:TileRsc_Type amount:int xn:int yn:int connerfill:bool
	ResourceRand :           resource:TileRsc_Type mean:int stddev:int repeat:int
	Resource :               resource:TileRsc_Type amount:int x:int y:int
	ResourceHLine :          resource:TileRsc_Type amount:int x:int w:int y:int
	ResourceVLine :          resource:TileRsc_Type amount:int x:int y:int h:int
	ResourceLine :           resource:TileRsc_Type amount:int x1:int y1:int x2:int y2:int
	ResourceRect :           resource:TileRsc_Type amount:int x:int w:int y:int h:int
	ResourceFillRect :       resource:TileRsc_Type amount:int x:int w:int y:int h:int
	ResourceFillEllipses :   resource:TileRsc_Type amount:int x:int w:int y:int h:int
	ResourceFromPNG :        name:string
	ResourceAgeing :         initrun:int msper:int resetaftern:int
	AddRoom :                bgtile:Tile_Type walltile:Tile_Type terrace:bool x:int y:int w:int h:int
	AddRoomMaze :            bgtile:Tile_Type walltile:Tile_Type terrace:bool x:int y:int w:int h:int xn:int yn:int connerfill:bool
	AddRoomsRand :           bgtile:Tile_Type walltile:Tile_Type terrace:bool align:int count:int mean:int stddev:int min:int
	ConnectRooms :           tile:Tile_Type connect:int allconnect:bool diagonal:bool
	FinalizeTerrain :        
	AddPortal :              x:int y:int display:FieldObjDisplay_Type acttype:FieldObjAct_Type PortalID:string DstPortalID:string message:string
	AddPortalRand :         display:FieldObjDisplay_Type acttype:FieldObjAct_Type PortalID:string DstPortalID:string message:string
	AddPortalInRoom :       display:FieldObjDisplay_Type acttype:FieldObjAct_Type PortalID:string DstPortalID:string message:string
	AddRecycler :            x:int y:int display:FieldObjDisplay_Type message:string
	AddRecyclerRand :        display:FieldObjDisplay_Type count:int message:string
	AddRecyclerInRoom :      display:FieldObjDisplay_Type count:int message:string
	AddTrapTeleport :        x:int y:int DstFloor:string message:string 
	AddTrapTeleportsRand :   DstFloor:string count:int message:string
	AddTrapTeleportsInRoom : DstFloor:string count:int message:string
	AddTrap :                x:int y:int display:FieldObjDisplay_Type acttype:FieldObjAct_Type message:string
	AddTrapsRand :           display:FieldObjDisplay_Type acttype:FieldObjAct_Type count:int message:string
	AddTrapsInRoom :         display:FieldObjDisplay_Type acttype:FieldObjAct_Type count:int message:string
