# default 를 지원하는 config struct 를 commandline argument 부터 만들어 냅니다. 

config struct from command args with default

golang reflect 로 struct tag 정보를 사용해서 

commandline에서 읽을 인자의 이름과 default 값을 관리 해줍니다. 

example/main.go 참고 

## 사용 tag 

default : 인자가 없는 경우 사용할 기본 값 

argname : 이 태그가 있는 경우 이 필드가 commandline argument로 부터 설정됩니다. 

    빈 값인 경우 : 필드 이름이 argument 이름이 됩니다. 
    있는 경우 : 이 값이 argument 이름이 됩니다. 

## 지원 하는 필드 타입 

    int, int8, int16, int32, int64
    uint, uint8, uint16, uint32, uint64
    float64, float32
    bool
    string

