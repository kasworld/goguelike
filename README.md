# 실행하려면 

[google translated this doc](https://translate.google.co.kr/translate?sl=ko&tl=en&u=https%3A%2F%2Fraw.githubusercontent.com%2Fkasworld%2Fgoguelike%2Fmaster%2FREADME.md)

설치 문서 INSTALL.md

[google translated install.md](https://translate.google.co.kr/translate?sl=ko&tl=en&u=https://raw.githubusercontent.com/kasworld/goguelike/master/INSTALL.md)

지형 스크립트 설명 towerscript.md

# 개요 및 특징  

## (가능한한) 혼자서 만드는 MMO 

서버 관리 서버 ( ground server )

다중 서버 지원 

서버는 linux에서 golang으로 개발/실행 

windows 지원 추가 (groundserver 제외)

클라이언트는 golang으로 webassembly 를 생성 


## 100% 서버 기반 
	
클라이언트는 viewer , 사용자 입력을 서버로 전달 하는 역할 

websocket을 사용 연결 유지형 클라이언트 

클라이언트를 통한 핵킹/치트 가능성을 원천 봉쇄 

클라이언트는 캐릭터의 현재 위치기준으로 시야내의 지형정보만을 받는다. ( 맵핵의 원천 봉쇄)

클라이언트 설치 불필요 

webassembly/html5 canvas를 지원하는 web browser 라면 플랫폼 불문하고 플레이 가능 


## roguelike - roglight?
	
실시간 턴제 

    행동별로 필요한 turn이 다르다. 
    이동, 공격이 대각선인 경우 sqrt2 만큼 더 필요. 

실시간 랭킹 

서버가 시작 될때 마다 random한 지형을 생성 ( 동일한 지형 스크립트를 사용하더라도 )


## 지형 스크립트를 기반으로 지형을 생성 

지형에 따른 공/방 변화 

지형에 따른 시야의 제한 
	
전투 가능/불가능 지역 지원 

다른 지형 스크립트를 사용할 경우 완전히 다른 지형이 가능 

지형 스크립트 생성기가 있다.  - towermaker


## 시간에 따라 변화하는 세계 

시간에 따라 지형이 변화 - 변화되는 지형을 visual 하게 표현 

시간에 따른 환경의 변화가 게임 플레이에 영향을 끼침 ( 유리하게 또는 불리하게 )

서버/지역 별로 시간의 흐름을 다르게 할수 있다. 


## 기본 게임 기능 

시간 기반의 버프/디버프 구현 

다양한 상태 이상(condition) 이 존재 

전투와 탐험에 따른 경험치 획득 및 성장 

함정이 존재 - 데미지, 텔레포트 성향 변화, 기억상실 등 

범위 공격 존재 

지뢰와 범위 공격을 하는 방해물 배치 가능 

사망시 경험치 손실, 아이템, 소지금 드랍 

장비 아이템 - level up 과 병행되는 강함 

스크롤 과 포션, 게임 머니가 존재 

재활용 상점(recycle)이 존재

장비/아이템의 무게 존재, 과중 시 페널티 존재. 

채팅 가능 

성장에 따라 HP(health point) / SP(stamina point) / 시야(sight) 증가
 
권한 키를 사용해 권한을 여러가지로 부여 할수 있다. ( admin 권한)


