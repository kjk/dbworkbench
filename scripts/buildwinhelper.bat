@echo on

godep go vet github.com/kjk/dbworkbench
@IF ERRORLEVEL 1 goto Error

call ./node_modules/.bin/gulp default
@IF ERRORLEVEL 1 goto Error

go run tools\build\cmd.go tools\build\gen_resources.go tools\build\main.go tools\build\s3.go tools\build\util.go tools\build\win.go -gen-resources
@IF ERRORLEVEL 1 goto Error

godep go build -tags embeded_resources -o dbherohelper.exe
@IF ERRORLEVEL 1 goto Error

:Error
