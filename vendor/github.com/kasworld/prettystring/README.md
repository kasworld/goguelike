# prettystring
make pretty string(like python) golang struct

reflect 를 사용해서 golang object들을 string으로 만들어 줍니다. 

object가 다른 object를 포함 하고 있는 경우 재귀적으로 처리 하며 
이 회수를 제한 하기 위한 인자를 받습니다. 

## struct field tag 로 printstring 사용법 

    hide : 필드 처리를 생략 
    hidevalue : 필드 내용을 숨김
    simple : 필드 처리를 간단히 - 내부로 재귀 하지 않음. 
    example/main.go 참고 