# logistics
国际物流价格查询

build
 ```bigquery
// go build window on mac
brew install mingw-w64
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="x86_64-w64-mingw32-gcc" go build -o ./bin/logistics.exe 

go build ./bin/logistics
```
