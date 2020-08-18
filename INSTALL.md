# 사전 준비 사항

준비물 : linux(debian,ubuntu,mint) , chrome web brower , golang 

버전 string 생성시 사용 : sha256sum, awk

goimports : 소스 코드 정리, import 해결

    go get golang.org/x/tools/cmd/goimports

Packet serializer (MessagePack) : https://github.com/tinylib/msgp

    go get github.com/tinylib/msgp

프로토콜 생성기 : https://github.com/kasworld/genprotocol

    go get github.com/kasworld/genprotocol

Enum 생성기 : https://github.com/kasworld/genenum

    go get github.com/kasworld/genenum

Log 패키지 및 커스텀 로그레벨 로거 생성기 : https://github.com/kasworld/log

    go get github.com/kasworld/log
    install.sh 실행해서 genlog 생성 

tilemaker로 새 타일(merged.json, merged.png)을 만드려면 (있는것을 그냥 사용할때는 필요없음)

    NanumBarunGothic.ttf : /usr/share/fonts/truetype/nanum/NanumBarunGothic.ttf
    Symbola_hint.ttf : /usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf
    의 설치가 필요하다. 


# 소스 코드 설치 

goguelike : https://github.com/kasworld/goguelike

    go get github.com/kasworld/goguelike


# 실행 

gogulike 폴더에서 

코드의 생성 

    ./gencode.sh 

실행바이너리 생성 : goguelike/bin 폴더로 

    ./build.sh

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


towerserver : goguelike game server 

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
    -DisplayName string
            DisplayName (default "Default")
    -GroundRPC string
            GroundRPC (default "localhost:14002")
    -LogLevel uint
            LogLevel (default 7)
    -ServiceHostBase string
            ServiceHostBase (default "http://localhost")
    -ServicePort int
            ServicePort (default 14101)
    -SplitLogLevel uint
            SplitLogLevel
    -StandAlone
            StandAlone (default true)
    -TowerFilename string
            TowerFilename (default "start")
    -TowerNumber int
            TowerNumber (default 1)
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
    -service string
            start,stop,restart,forcestart,logreopen (default "start")


wasmclient : webbrowser 용 GUI client 

    browser 를 통해서 실행된다. 
    
    url에 인자를 주어 option이 선택된 상태로 실행가능하다. 
    LeftInfo=LeftInfoOff LeftInfoOn
    CenterInfo=HelpOff Highscore ClientInfo Help FactionInfo CarryObjectInfo PotionInfo ScrollInfo MoneyColor TileInfo ConditionInfo FieldObjInfo
    RightInfo=RightInfoOff Message DebugInfo InvenList FieldObjList FloorList
    Viewport=PlayVP FloorVP
    Zoom=Zoom0 Zoom1 Zoom2
    Sound=SoundOn SoundOff
    
    authkey 인자를 통해 특수 권한 클라이언트로 실행 가능하다. 


./tool 폴더 아래에 

makechatdata 

    fortune 데이터로 부터 chatdata.txt 를 생성 

tilemaker

    webclient 에서 사용할 merged.json, merged.png 를 생성 

towermaker 

    roguelike 풍 tower script 생성기 
