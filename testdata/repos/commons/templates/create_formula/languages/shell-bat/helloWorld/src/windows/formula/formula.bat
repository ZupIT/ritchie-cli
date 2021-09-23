@echo off
SETLOCAL

call:%~1
goto exit

:runFormula
  echo Hello World!
  echo My name is %RIT_INPUT_TEXT%.

  if "%RIT_INPUT_BOOLEAN%" == "true" (
    echo I've already created formulas using Ritchie.
  ) else (
    echo I'm excited in creating new formulas using Ritchie.
  )

  echo Today, I want to automate %RIT_INPUT_LIST%.
  echo My secret is %RIT_INPUT_PASSWORD%.

  goto exit

:exit
  ENDLOCAL
  exit /b 0
