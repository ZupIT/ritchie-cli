#!/bin/sh -xe

idempotent_config() {
  if [ -d ${HOME}"/.rit" ]; then
	  echo "~/.rit directory already exists, moving on"
  else
	  mkdir ${HOME}"/.rit"
  fi
  if [ -n "$SHELL_TYPE" ]; then

    if [ -f ${HOME}"/."${SHELL_TYPE}"rc" ]; then

        sed -i'' -e '/source <(rit completion/d' ${HOME}"/."${SHELL_TYPE}"rc"

        if grep "[[ -r "/usr/local/bin/rit" ]] && rit completion $SHELL_TYPE > ~/.rit_completion" ~/."${SHELL_TYPE}"rc > /dev/null
        then
          echo ".${SHELL_TYPE}rc IS ALREADY CONFIGURED TO SUPPORT RITCHIE AUTO COMPLETION"
        else
          echo "[[ -r "/usr/local/bin/rit" ]] && rit completion $SHELL_TYPE > ~/.rit_completion" >> ~/."${SHELL_TYPE}"rc
          echo "source ~/.rit_completion" >> ~/."${SHELL_TYPE}"rc
        fi

    elif [ -f ${HOME}"/."${SHELL_TYPE}"_profile" ]; then

        sed -i'' -e '/source <(rit completion/d' ${HOME}"/."${SHELL_TYPE}"_profile"

        if grep "[[ -r "/usr/local/bin/rit" ]] && rit completion $SHELL_TYPE > ~/.rit_completion" ~/."${SHELL_TYPE}"_profile > /dev/null
        then
          echo ".${SHELL_TYPE}_profile IS ALREADY CONFIGURED TO SUPPORT RITCHIE AUTO COMPLETION"
        else
          echo "[[ -r "/usr/local/bin/rit" ]] && rit completion $SHELL_TYPE > ~/.rit_completion" >> ~/."${SHELL_TYPE}"_profile
          echo "source ~/.rit_completion" >> ~/."${SHELL_TYPE}"_profile
        fi

    else
      echo "Installed Ritchie without autocomplete"
      return
    fi


      echo "Installed Ritchie with autocomplete"
  fi
}

rit_identify_shell () {
  if [ -n "$($SHELL -c 'echo $ZSH_VERSION')" ]; then
    echo "Going to install autocomplete for zsh"
    SHELL_TYPE="zsh"
  elif [ -n "$($SHELL -c 'echo $BASH_VERSION')" ]; then
    echo "Going to install autocomplete for bash"
    SHELL_TYPE="bash"
  else
    echo "Please consider using bash or zsh to improve your ritchie experience"
    SHELL_TYPE=""
  fi
}

rit_identify_shell
idempotent_config