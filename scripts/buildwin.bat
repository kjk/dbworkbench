@echo on

@call scripts\buildwinhelper.bat
@IF ERRORLEVEL 1 goto Error

go run tools\build\cmd.go tools\build\gen_resources.go tools\build\main.go tools\build\s3.go tools\build\util.go tools\build\win.go
@IF ERRORLEVEL 1 goto Error

goto EndOk

:Error
echo there was an error!
goto End

:EndOk
echo finished ok!

:End
