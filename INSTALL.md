# 사전 준비 사항

준비물 : linux(debian,ubuntu,mint) , chrome web brower , golang 

버전 string 생성시 사용 : sha256sum, awk

release.sh 에서 사용 : 7z, scp, ssh

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

goguelike : https://github.com/kasworld/goguelike2

    go get github.com/kasworld/goguelike2


# 실행 

gogulike 폴더에서 

코드의 생성 

    ./gencode.sh 

실행바이너리 생성 : goguelike/bin 폴더로 

    ./build.sh

# 구성 

./rundriver 속에 

groundserver 

    towerserver 관리 서버 

towerserver 

    goguelike game server 

multiclinet 

    load test용 다중 client 
  
textclient 

    debug/test 용 ui없는 client 

wasmclient 

    webbrowser 용 GUI client 

./tool 폴더 아래에 

makechatdata 

    fortune 데이터로 부터 chatdata.txt 를 생성 

tilemaker

    webclient 에서 사용할 merged.json, merged.png 를 생성 

towermaker 

    roguelike 풍 tower script 생성기 
