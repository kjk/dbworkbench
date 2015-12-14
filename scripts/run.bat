@echo on

godep go vet github.com/kjk/dbworkbench

set GOARCH=386

godep go build -o dbhero.exe
@IF ERRORLEVEL 1 goto Error

dbhero.exe -dev

:Error
