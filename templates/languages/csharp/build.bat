:: Java parameters
echo off
SETLOCAL
SET SOURCE_FILE=src/src.csproj
SET BIN_FOLDER=bin
SET BIN_FOLDER_LINUX=linux
SET BIN_FOLDER_WINDOWS=windows
SET BIN_FOLDER_DARWIN=darwin
SET BIN_UNIX=src
SET BIN_WINDOWs=src.exe
SET SH_FILE=%BIN_FOLDER%\run.sh
SET BAT_FILE=%BIN_FOLDER%\run.bat


:build
    mkdir %BIN_FOLDER%/%BIN_FOLDER_LINUX%
    dotnet build %SOURCE_FILE% -o %BIN_FOLDER%/%BIN_FOLDER_LINUX% --configuration Release --runtime linux-x64
    mkdir %BIN_FOLDER%/%BIN_FOLDER_WINDOWS%
	dotnet build %SOURCE_FILE% -o %BIN_FOLDER%/%BIN_FOLDER_WINDOWS% --configuration Release --runtime win10-x64
    mkdir %BIN_FOLDER/%BIN_FOLDER_DARWIN%
	dotnet build %SOURCE_FILE% -o %BIN_FOLDER%/%BIN_FOLDER_DARWIN% --configuration Release --runtime osx-x64
    CALL :BAT_WINDOWS
    CALL :SH_LINUX
    CALL :CP_DOCKER
    GOTO DONE

:BAT_WINDOWS
    echo '@ECHO OFF' > %BAT_FILE%
	echo '%BIN_FOLDER_WINDOWS%/%BIN_WINDOWS%' >> %BAT_FILE%
    GOTO DONE

:SH_LINUX
	echo '#!/bin/sh' > %SH_FILE%
    echo 'if [ $$(uname) = "Darwin" ]; then' >> %SH_FILE%
	echo '  "$$(dirname "$$0")"/%BIN_FOLDER_DARWIN%/%BIN_UNIX%' >> %SH_FILE%
	echo 'else' >> %SH_FILE%
	echo '  "$$(dirname "$$0")"/%BIN_FOLDER_LINUX%/%BIN_UNIX%' >> %SH_FILE%
	echo 'fi' >> %SH_FILE%
	chmod +x %SH_FILE%
    GOTO DONE

:CP_DOCKER
    copy Dockerfile %BIN_FOLDER%
    copy set_umask.sh %BIN_FOLDER%
    GOTO DONE

:DONE
