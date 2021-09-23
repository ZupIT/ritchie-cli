:: Java parameters
echo off
SETLOCAL
SET BIN_FOLDER=bin
SET BIN_NAME=Main.jar
SET BAT_FILE=%BIN_FOLDER%\run.bat
SET SH_FILE=%BIN_FOLDER%\run.sh
:build
    call mvn clean install 1>&2
    if %errorlevel% neq 0 exit /b %errorlevel%
    mkdir %BIN_FOLDER%
    copy target\%BIN_NAME% %BIN_FOLDER%\%BIN_NAME%
    rmdir /Q /S target
    CALL :BAT_WINDOWS
    CALL :SH_LINUX
    CALL :CP_DOCKER
    GOTO DONE

:BAT_WINDOWS
    echo @ECHO OFF > %BAT_FILE%
    echo java -jar %BIN_NAME% >> %BAT_FILE%
    GOTO DONE

:SH_LINUX
    echo java -jar "$(dirname "$0")"/%BIN_NAME% > %SH_FILE%
    GOTO DONE

:CP_DOCKER
    copy Dockerfile %BIN_FOLDER%
    copy set_umask.sh %BIN_FOLDER%
    GOTO DONE

:DONE
