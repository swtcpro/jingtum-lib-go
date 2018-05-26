@echo off
setlocal
if exist build.bat goto ok
echo build.bat must be run from its folder
goto end

: ok
set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0
gofmt -w src
go install testLib

:end
echo finished

