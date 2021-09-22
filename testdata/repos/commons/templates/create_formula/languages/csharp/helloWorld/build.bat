:: Csharp parameters
echo off
SETLOCAL
SET SOURCE_FILE=src/src.csproj
SET BIN_FOLDER=bin
SET BIN_FOLDER_WINDOWS=windows
SET BIN_WINDOWs=src.exe
SET BAT_FILE=%BIN_FOLDER%\run.bat

:build
    call :checkCommand dotnet

    mkdir %BIN_FOLDER%
    xcopy /E /I src %BIN_FOLDER%
	dotnet build %SOURCE_FILE% -o %BIN_FOLDER%/%BIN_FOLDER_WINDOWS% --configuration Release
    CALL :BAT_WINDOWS
    CALL :CP_DOCKER
    GOTO DONE

:BAT_WINDOWS
    echo @ECHO OFF > %BAT_FILE%
	echo %%~dp0%BIN_FOLDER_WINDOWS%/src.exe >> %BAT_FILE%
    GOTO DONE

:CP_DOCKER
    copy Dockerfile %BIN_FOLDER%
    copy set_umask.sh %BIN_FOLDER%
    GOTO DONE

:checkCommand
    WHERE %1 >nul 2>nul
    IF %ERRORLEVEL% NEQ 0 ECHO %1 required 1>&2 && exit 1

:DONE
