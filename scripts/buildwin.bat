@echo on

godep go vet github.com/kjk/dbworkbench
@IF ERRORLEVEL 1 goto Error

call ./node_modules/.bin/gulp default
@IF ERRORLEVEL 1 goto Error

go run scripts\build_release.go
@IF ERRORLEVEL 1 goto Error

godep go build -o dbworkbench.exe
@IF ERRORLEVEL 1 goto Error

go run tools\build\main.go tools\build\util.go tools\build\cmd.go tools\build\s3.go tools\build\win.go -no-clean-check
@IF ERRORLEVEL 1 goto Error

goto EndOk

:Error
echo there was an error!
goto End

:EndOk
echo finished ok!

:End
