:: Php parameters
echo off
SETLOCAL
SET BIN_FOLDER=bin
SET BAT_FILE=%BIN_FOLDER%\run.bat
SET SH_FILE=%BIN_FOLDER%\run.sh

:build
    call checkCommand composer
    call checkCommand php
    mkdir %BIN_FOLDER%
    xcopy /E /I src %BIN_FOLDER%
    cd %BIN_FOLDER%
    call composer install -q
    cd ..
    CALL :BAT_WINDOWS
    CALL :SH_LINUX
    CALL :CP_DOCKER
    GOTO DONE

:BAT_WINDOWS
    echo @ECHO OFF > %BAT_FILE%
    echo SET mypath=%%~dp0 >> %BAT_FILE%
    echo start /B /D "%%mypath%%" /WAIT php -f index.php >> %BAT_FILE%

:SH_LINUX
    echo php -f "$(dirname "$0")"/index.php > %SH_FILE%
    GOTO DONE


:CP_DOCKER
    copy Dockerfile %BIN_FOLDER%
    copy set_umask.sh %BIN_FOLDER%
    GOTO DONE

:checkCommand
    WHERE %1 >nul 2>nul
    IF %ERRORLEVEL% NEQ 0 ECHO %1 required 1>&2 && exit 1

:DONE
