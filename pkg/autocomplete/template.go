/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package autocomplete

const (
	command string = `_{{FunctionName}}()
{
    last_command="{{LastCommand}}"

    commands=()
{{Commands}}
}
`

	autoCompletionBash string = `
{{DynamicCode}}

__{{BinaryName}}_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__{{BinaryName}}_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__{{BinaryName}}_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__{{BinaryName}}_handle_command()
{

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_{{BinaryName}}_root"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    declare -F "$next_command" >/dev/null && $next_command
}


__{{BinaryName}}_handle_reply()
{
    local completions
    completions=("${commands[@]}")
    COMPREPLY=( $(compgen -W "${completions[*]}" -- "$cur") )

}

__{{BinaryName}}_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __{{BinaryName}}_handle_reply
        return
    fi
    if __{{BinaryName}}_contains_word "${words[c]}" "${commands[@]}"; then
        __{{BinaryName}}_handle_command
        fi

    __{{BinaryName}}_handle_word
}

__start_{{BinaryName}}()
{
    local cur prev words cword
    __{{BinaryName}}_init_completion -n "=" || return

    local c=0
    local commands=("{{BinaryName}}")
    local last_command

    __{{BinaryName}}_handle_word
}

complete -F __start_{{BinaryName}} {{BinaryName}}`

	autoCompletionZsh string = `#compdef {{BinaryName}}

__{{BinaryName}}_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand

	source "$@"
}

__{{BinaryName}}_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?

	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}

__{{BinaryName}}_compopt() {
	true
}

__{{BinaryName}}_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

autoload -U +X bashcompinit && bashcompinit

LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__{{BinaryName}}_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__{{BinaryName}}_get_comp_words_by_ref/g" \
	-e "s/${LWORD}compgen${RWORD}/__{{BinaryName}}_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__{{BinaryName}}_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	<<'BASH_COMPLETION_EOF'

{{AutoCompleteBash}}

# ex: ts=4 sw=4 et filetype=sh
BASH_COMPLETION_EOF
}

__{{BinaryName}}_bash_source <(__{{BinaryName}}_convert_bash_to_zsh)

_complete {{BinaryName}}  2>/dev/null`
)
