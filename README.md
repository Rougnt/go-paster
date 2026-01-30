# go-paster

## 编译

``` bash
go build -o go-paster main.go
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o go-paster.exe 
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags -H=windowsgui -o go-paster.exe
```
