# 사전 준비 사항

[google translated install.md](https://translate.google.co.kr/translate?sl=ko&tl=en&u=https://raw.githubusercontent.com/kasworld/goguelike/master/INSTALL.md)


# 준비물 
    
    linux(debian,ubuntu,mint) 
    또는 window (아마도 10) 제한적 지원(groundserver 미지원) + powershell(5.1, 7.x)
        windows에서는 sig_usr1을 사용할수 없어서 logrotate가 불가능해 졌습니다.
    chrome web brower ( 또는 websocket, webassembly, webgl 을 지원하는 브라우저)
    golang
    webgl을 지원할 그래픽 카드 

goimports : 소스 코드 정리, import 해결

    go get golang.org/x/tools/cmd/goimports


Packet serializer (MessagePack) : https://github.com/tinylib/msgp - gob 로 변경되어 불필요

    go get github.com/tinylib/msgp

버전 string 생성시 사용 : windows, linux 간에 같은 string생성

    go get github.com/kasworld/makesha256sum

프로토콜 생성기 : https://github.com/kasworld/genprotocol

    go get github.com/kasworld/genprotocol

Enum 생성기 : https://github.com/kasworld/genenum

    go get github.com/kasworld/genenum

Log 패키지 및 커스텀 로그레벨 로거 생성기 : https://github.com/kasworld/log

    go get github.com/kasworld/log
    install.sh 실행해서 genlog 생성 


# 소스 코드 설치 

goguelike : https://github.com/kasworld/goguelike

    go get github.com/kasworld/goguelike


# 실행 

gogulike 폴더에서 

코드의 생성 

    linux : ./gencode.sh 
    windows : ./gencode.ps1

실행바이너리 생성 : goguelike/bin 폴더로 

    linux : ./build.sh
    windows : ./build.ps1

wasm 바이너리 생성 : ./rundriver 폴더에서 

    linux : ./genwasmclient.sh 
    windows : ./genwasmclient.ps1

# 실행파일 및 인자 

실행파일 -h 로 실행하면 간단한 도움말이 나온다. 

./rundriver 속에 

groundserver : towerserver 관리 서버 

    Usage of ./groundserver:
    -AdminAuthKey string
        AdminAuthKey (default "6e9456cf-ab29-99b2-f223-1459e00cfcd5")
    -BaseLogDir string
        BaseLogDir (default "/tmp/")
    -ClientDataFolder string
        ClientDataFolder (default "./clientdata")
    -DataFolder string
        DataFolder (default "./serverdata")
    -ExeFolder string
        ExeFolder (default "./")
    -GroundAdminWebPort int
        GroundAdminWebPort (default 14001)
    -GroundHost string
        GroundHost (default "localhost")
    -GroundRPCPort int
        GroundRPCPort (default 14002)
    -GroundServiceWebPort int
        GroundServiceWebPort (default 14003)
    -HighScoreFile string
        HighScoreFile (default "highscore.json")
    -LogLevel uint
        LogLevel (default 7)
    -SplitLogLevel uint
        SplitLogLevel
    -TowerAdminHostBase string
        TowerAdminHostBase (default "http://localhost")
    -TowerAdminWebPortBase int
        TowerAdminWebPortBase (default 14200)
    -TowerBin string
        TowerBin (default "towerserver")
    -TowerDataFile string
        TowerDataFile (default "towerdata.json")
    -TowerServiceHostBase string
        TowerServiceHostBase (default "http://localhost")
    -TowerServicePortBase int
        TowerServicePortBase (default 14100)
    -WebAdminID string
        WebAdminID (default "root")
    -WebAdminPass string
        WebAdminPass (default "password")
    -i string
        server config file or url
    -service string
        start,stop,restart,forcestart,logreopen (default "start")



towerserver : goguelike game server 또는 windows용 towerserverwin 

    Usage of ./towerserver:
    -AdminAuthKey string
        AdminAuthKey (default "6e9456cf-ab29-99b2-f223-1459e00cfcd5")
    -AdminPort int
        AdminPort (default 14201)
    -BaseLogDir string
        BaseLogDir (default "/tmp/")
    -ClientDataFolder string
        ClientDataFolder (default "./clientdata")
    -ConcurrentConnections int
        ConcurrentConnections (default 10000)
    -DataFolder string
        DataFolder (default "./serverdata")
    -GroundRPC string
        GroundRPC (default "localhost:14002")
    -LogLevel uint
        LogLevel (default 7)
    -ScriptFilename string
        ScriptFilename (default "start")
    -ServiceHostBase string
        ServiceHostBase (default "http://localhost")
    -ServicePort int
        ServicePort (default 14101)
    -SplitLogLevel uint
        SplitLogLevel
    -StandAlone
        StandAlone (default true)
    -TowerName string
        TowerName (default "Default")
    -TurnPerSec float
        TurnPerSec (default 2)
    -WebAdminID string
        WebAdminID (default "root")
    -WebAdminPass string
        WebAdminPass (default "password")
    -cpuprofilename string
        cpu profile filename
    -i string
        server config file or url
    -memprofilename string
        memory profile filename
    -service string
        start,stop,restart,forcestart,logreopen (default "start")

multiclinet : load test용 다중 client 

    Usage of ./multiclient:
    -AccountOverlap int
        AccountOverlap
    -AccountPool int
        AccountPool
    -BaseLogDir string
        BaseLogDir (default "/tmp/")
    -Concurrent int
        Concurrent (default 1000)
    -ConnectToTower string
        ConnectToTower (default "localhost:14101")
    -DisconnectOnDeath
        DisconnectOnDeath
    -LimitEndCount int
        LimitEndCount
    -LimitStartCount int
        LimitStartCount
    -ListenWebInfoPort string
        ListenWebInfoPort (default ":14011")
    -LogLevel uint
        LogLevel
    -Net string
        Net (default "web")
    -PlayerNameBase string
        PlayerNameBase (default "MC_")
    -RetryDelayTimeOut int
        RetryDelayTimeOut (default -1)
    -SplitLogLevel uint
        SplitLogLevel
    -cpuprofilename string
        cpu profile filename
    -i string
        client config file or url
    -memprofilename string
        memory profile filename
  
textclient : debug/test 용 ui없는 client 

    Usage of ./textclient:
    -ConnectToTower string
        ConnectToTower (default "localhost:14101")
    -LogDir string
        LogDir (default "/tmp/textclient.logfile")
    -LogLevel uint
        LogLevel
    -PidFilename string
        PidFilename (default "/tmp/textclient.pid")
    -PlayerName string
        PlayerName (default "Player")
    -SplitLogLevel uint
        SplitLogLevel
    -i string
        client config file or url

여러가지 tower script 생성기 

    Usage of ./towermaker:
    -floorcount int
        roguelike,big tower floor count (default 100)
    -towername string
        all,start,roguelike,big,sight,objmax tower to make

wasmclientgl : webbrowser 용 webgl client 

    browser 를 통해서 실행된다. 
    
    url에 인자를 주어 option이 선택된 상태로 실행가능하다. 
    LeftInfo=LeftInfoOff LeftInfoOn
    CenterInfo=HelpOff Highscore ClientInfo Help FactionInfo CarryObjectInfo PotionInfo ScrollInfo MoneyColor TileInfo ConditionInfo FieldObjInfo
    RightInfo=RightInfoOff Message DebugInfo InvenList FieldObjList FloorList
    Viewport=PlayVP FloorVP
    Zoom=Zoom0 Zoom1 Zoom2
    Angle=Angle0 Angle1 Angle2
    Sound=SoundOn SoundOff
    
    authkey 인자를 통해 특수 권한 클라이언트로 실행 가능하다. 


./tool 폴더 아래에 

makechatdata 

    fortune 데이터로 부터 chatdata.txt 를 생성 


