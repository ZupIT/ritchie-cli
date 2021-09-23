#!/bin/sh

runFormula() {
  echo "Hello World!"
  echoColor "green" "My name is $RIT_INPUT_TEXT."

  if [ "$RIT_INPUT_BOOLEAN" = "true" ]; then
    echoColor "blue" "I've already created formulas using Ritchie."
  else
    echoColor "red" "I'm excited in creating new formulas using Ritchie."
  fi

  echoColor "yellow" "Today, I want to automate $RIT_INPUT_LIST."
  echoColor "cyan"  "My secret is $RIT_INPUT_PASSWORD."
}

echoColor() {
  case $1 in
    red)
      echo "$(printf '\033[31m')$2$(printf '\033[0m')"
      ;;
    green)
      echo "$(printf '\033[32m')$2$(printf '\033[0m')"
      ;;
    yellow)
      echo "$(printf '\033[33m')$2$(printf '\033[0m')"
      ;;
    blue)
      echo "$(printf '\033[34m')$2$(printf '\033[0m')"
      ;;
    cyan)
      echo "$(printf '\033[36m')$2$(printf '\033[0m')"
      ;;
    esac
}
