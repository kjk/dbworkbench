@echo on

call ./node_modules/.bin/gulp default
@IF ERRORLEVEL 1 goto Error

godep go vet github.com/kjk/dbworkbench
@IF ERRORLEVEL 1 goto Error

go run scripts\build_release.go
@IF ERRORLEVEL 1 goto Error

godep go build -o dbworkbench.exe
@IF ERRORLEVEL 1 goto Error

go run tools\buildwin\main.go tools\buildwin\util.go tools\buildwin\cmd.go 
@IF ERRORLEVEL 1 goto Error
	
goto EndOk

:Error
echo there was an error!
goto End

:EndOk
echo finished ok!

:End
