# golang signal handler for server/service for linux   

리눅스용 서버/서비스를 만드는데 사용하는 signal handler 입니다. 

# 기능 

pid lockfile을 사용해서 서비스의 제어를 가능하게 해줍니다. 

서비스 시작 : 중복 실행을 막으며 서비스 시작이 가능합니다. 

    -service=start 

서비스 종료 : ctrl-c 등으로 종료할때 서비스 클린업을 하고 종료 할 수 있습니다. 

    -service=stop


로그 로테이트 : 서비스 실행중에 로그 파일의 교체가 가능합니다. 

    -service=logreopen


서비스 재시작 : 실행중인 서비스를 종료하고 다시 시작 합니다. 

    -service=restart 

서비스 강제 재시작 : 전원 중단 등으로 비정상 종료된 경우 pid lockfile을 지우고 서비스를 강재로 시작 합니다. 

    -service=forcestart


# 사용시 추가 되는 인자 

    -service string
    	start,stop,restart,forcestart,logreopen (default "start")

# 사용자 프로그램에서 구현 해야 하는 interface

서비스로 만드려는 프로그램 (struct) 에서 다음 인터페이스를 구현 합니다. 

    type ServiceI interface {
        GetServiceLockFilename() string
        ServiceInit() error
        ServiceMain(ctx context.Context)
        ServiceCleanup()
        GetLogger() interface{} // LoggerI 를 구현 해야 합니다. 
    }

    type LoggerI interface {
        Reload() error
        Fatal(format string, v ...interface{})
        Error(format string, v ...interface{})
        Debug(format string, v ...interface{})
    }

# 사용 예제 

example/example.go 참고 

# 사용 프로젝트 

https://github.com/kasworld/goonlinescaffolding

https://github.com/kasworld/gowasm3dgame

https://github.com/kasworld/gowasm2dgame

https://github.com/kasworld/goguelike

