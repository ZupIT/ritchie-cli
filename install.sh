#!/bin/sh

usage () {
	echo "./install.sh"
}


idempotent_config() {
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

rit_install () {
	echo "Downloading rit..."
	STABLE_VERSION=$(curl -s https://commons-repo.ritchiecli.io/stable.txt)
	curl -SLO "https://commons-repo.ritchiecli.io/${STABLE_VERSION}/${OPERATIONAL_SYSTEM}/rit"


	chmod +x ./rit

	INSTALL_PATH="/usr/local/bin"

	if [ ! -d "$INSTALL_PATH" ]; then
		sudo mkdir -p $INSTALL_PATH
	fi

	sudo mv ./rit $INSTALL_PATH/rit
	$INSTALL_PATH/rit --version

	idempotent_config
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

rit_identify_os () {

  if [ $(uname) = "Linux" ]; then
      echo "Installing Ritchie for Linux"
      OPERATIONAL_SYSTEM="linux"
  elif [ $(uname) = "Darwin" ]; then
      echo "Installing Ritchie for Mac"
      OPERATIONAL_SYSTEM="darwin"
  else
    echo "Unable to identify which OS you're using"
    exit 1
  fi
}

rit_identify_os

rit_identify_shell

rit_install
