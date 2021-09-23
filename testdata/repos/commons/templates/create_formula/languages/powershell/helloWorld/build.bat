@ECHO OFF
SETLOCAL

SET BIN_FOLDER=bin
SET BAT_FILE=%BIN_FOLDER%\run.bat
SET SH_FILE=%BIN_FOLDER%\run.sh
SET ENTRY_POINT=main.ps1

:build
  call :checkCommand Powershell
  mkdir %BIN_FOLDER%
  xcopy /E /I src %BIN_FOLDER%
  call :BAT_WINDOWS
  call :SH_LINUX
  call :CP_DOCKER
  GOTO DONE

:BAT_WINDOWS
  echo @ECHO OFF > %BAT_FILE%
  echo SET mypath=%%~dp0 >> %BAT_FILE%
  echo Powershell.exe -executionpolicy remotesigned -File "%%mypath%%"%ENTRY_POINT% >> %BAT_FILE%

:SH_LINUX
    echo #!/bin/sh > %SH_FILE%
    echo pwsh "$(dirname "$0")"/%ENTRY_POINT% >> %SH_FILE%
    GOTO DONE

:CP_DOCKER
  copy Dockerfile %BIN_FOLDER%
  copy set_umask.sh %BIN_FOLDER%
  GOTO DONE

:checkCommand
    WHERE %1 >nul 2>nul
    IF %ERRORLEVEL% NEQ 0 ECHO %1 required 1>&2 && exit 1

:DONE
