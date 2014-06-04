::===========================================================
:: Curl  : Simple http download and readline lib by Golang
:: HOST  : https://github.com/kenshin/curl
:: Author: Kenshin <kenshin@ksria.com>
::===========================================================

@ECHO off

IF "%1" == "doc" GOTO doc
IF "%1" == "install" GOTO install
IF "%1" == "test" GOTO test

:doc
@ECHO godoc http://127.0.0.1:6060
godoc -http=:6060 -server=:6060
IF "%1" == "doc" GOTO exit

:install
@ECHO go install
go install
IF "%1" == "install" GOTO exit

:test
@ECHO go test
go test
GOTO exit

:exit
@ECHO complete.