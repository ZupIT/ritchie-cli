:: Go parameters
echo off
SETLOCAL
SET BINARY_NAME=main
SET GOCMD=go
SET GOBUILD=%GOCMD% build
SET CMD_PATH=main.go
SET BIN_FOLDER=..\bin
SET DIST_WIN_DIR=%BIN_FOLDER%\windows
SET DIST_LINUX_DIR=%BIN_FOLDER%\linux
SET BIN_WIN=%BINARY_NAME%.exe
SET BAT_FILE=%BIN_FOLDER%\run.bat
SET SH_FILE=%BIN_FOLDER%\run.sh

:build
    cd src
    mkdir %DIST_WIN_DIR%
    SET GO111MODULE=on
    for /f %%i in ('go list -m') do set MODULE=%%i
    CALL :windows
    CALL :linux
    if %errorlevel% neq 0 exit /b %errorlevel%
    GOTO CP_DOCKER
    GOTO DONE
    cd ..

:windows
    SET CGO_ENABLED=
	SET GOOS=windows
    SET GOARCH=amd64
    %GOBUILD% -tags release -o %DIST_WIN_DIR%\%BIN_WIN% %CMD_PATH%
    echo @ECHO OFF > %BAT_FILE%
    echo SET mypath=%%~dp0 >> %BAT_FILE%
    echo start /B /D "%%mypath%%" /WAIT windows\main.exe >> %BAT_FILE%
    GOTO DONE

:linux
    SET CGO_ENABLED=0
	SET GOOS=linux
    SET GOARCH=amd64
    %GOBUILD% -tags release -o %DIST_LINUX_DIR%\%BINARY_NAME% %CMD_PATH%
    echo "$(dirname "$0")"/linux/%BINARY_NAME% > %SH_FILE%
    GOTO DONE

:CP_DOCKER
    copy ..\Dockerfile %BIN_FOLDER%
    copy ..\set_umask.sh %BIN_FOLDER%
    GOTO DONE
:DONE
