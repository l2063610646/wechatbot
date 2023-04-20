@REM 64bit
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o bin/wechatbot-amd64.exe main.go

@REM 32-bit
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build -o bin/wechatbot-386.exe main.go

@REM 64-bit
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o bin/wechatbot-amd64-linux main.go

@REM 32-bit
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build -o bin/wechatbot-386-linux main.go

@REM 64-bit
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o bin/wechatbot-amd64-darwin main.go

@REM 32-bit
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=386
go build -o bin/wechatbot-386-darwin main.go
