:: Rust parameters
echo off
SETLOCAL
SET BIN_NAME=formula
SET BIN_FOLDER=bin
SET SH_FILE=%BIN_FOLDER%\run.sh
SET BAT_FILE=%BIN_FOLDER%\run.bat
:build
    call :checkCommand cargo

    mkdir %BIN_FOLDER%
    xcopy /E /I src %BIN_FOLDER%
    cd %BIN_FOLDER%
    call cargo build --release
    cd ..
    call :BAT_WINDOWS
    call :SH_LINUX
    call :CP_DOCKER
    GOTO DONE

:BAT_WINDOWS
    echo @ECHO OFF > %BAT_FILE%
    echo SET mypath=%%~dp0 >> %BAT_FILE%
    echo start /B /D "%%mypath%%" /WAIT target\release\%BIN_NAME% >> %BAT_FILE%

:CP_DOCKER
    copy Dockerfile %BIN_FOLDER%
    copy set_umask.sh %BIN_FOLDER%
    GOTO DONE

:checkCommand
    WHERE %1 >nul 2>nul
    IF %ERRORLEVEL% NEQ 0 ECHO %1 required 1>&2 && exit 1

:DONE
