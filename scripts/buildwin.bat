@echo on

@rem flags:
@rem   -upload : upload to s3
@rem   -beta   : build beta version

@call scripts\buildwinhelper.bat
@IF ERRORLEVEL 1 goto Error

go run tools\build\cmd.go tools\build\gen_resources.go tools\build\main.go tools\build\s3.go tools\build\util.go tools\build\win.go %1 %2 %3
@IF ERRORLEVEL 1 goto Error

goto EndOk

:Error
echo there was an error!
goto End

:EndOk
echo finished ok!

:End
