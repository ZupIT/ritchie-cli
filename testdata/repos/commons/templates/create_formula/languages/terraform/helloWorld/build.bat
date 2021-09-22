@echo off

SETLOCAL
SET BIN_FOLDER=bin
SET BINARY_NAME=run.bat
SET BINARY_NAME_UNIX=run.sh
SET ENTRY_POINT=main.bat
SET ENTRY_POINT_UNIX=main.sh

:BUILD
  call :checkCommand terraform

  mkdir %BIN_FOLDER%
  xcopy /e/h/i/c src %BIN_FOLDER%
  cd %BIN_FOLDER%
  call :SH_UNIX
  call :BAT_WINDOWS
  GOTO EXIT

:SH_UNIX
  rename %ENTRY_POINT_UNIX% %BINARY_NAME_UNIX%

:BAT_WINDOWS
  rename %ENTRY_POINT% %BINARY_NAME%

:CP_DOCKER
  cd ..
  copy Dockerfile %BIN_FOLDER%
  copy set_umask.sh %BIN_FOLDER%

:checkCommand
    WHERE %1 >nul 2>nul
    IF %ERRORLEVEL% NEQ 0 ECHO %1 required 1>&2 && exit 1

:EXIT
  ENDLOCAL
  exit /b
