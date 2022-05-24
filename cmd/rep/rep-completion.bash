#!/usr/bin/env bash

tools="one two three"
options="--help --version"
one_opt="--number --float"
two_opt="--field --slice"

hasword() {
	word=$1
	shift
	args="$@"
	for w in $args; do
		if [ "$word" = "$w" ]; then
			echo 1
			return 0
		fi
	done

	echo 0
}

_rep() {
	COMPREPLY=()
	local cur prev
	_get_comp_words_by_ref cur prev

	if [ ${COMP_CWORD} -eq 1 ]; then
		if [[ ${cur} == -* ]]; then
			COMPREPLY=( $(compgen -W "${options}" -- $cur) )
			return 0
		fi

		if [[ ${tools} == "" ]]; then
			_compopt_o_filenames
			COMPREPLY=( $( compgen -f -- "$cur" ) $( compgen -d -- "$cur" ) )
			return 0
		fi

		COMPREPLY=( $(compgen -W "${tools}" -- $cur) )
		return 0
	fi

	if [[ ${cur} == -* ]]; then
		case ${prev} in
			"one")
				COMPREPLY=( $(compgen -W "${one_opt}" -- $cur) )
				;;
			"two")
				COMPREPLY=( $(compgen -W "${two_opt}" -- $cur) )
				;;
			*)
				if [ $(hasword ${prev} ${one_opt}) == "1" ]; then
					COMPREPLY=( $(compgen -W "${one_opt}" -- $cur) )
					return 0
				fi

				if [ $(hasword ${prev} ${two_opt}) == "1" ]; then
					COMPREPLY=( $(compgen -W "${two_opt}" -- $cur) )
					return 0
				fi

				;;
		esac

		return 0
	fi

	if [[ $(hasword ${prev} ${tools}) == "1" ]]; then
		_compopt_o_filenames
		COMPREPLY+=( $( compgen -f -- "$cur" ) $( compgen -d -- "$cur" ) )
		return 0
	fi

	_compopt_o_filenames
	COMPREPLY=( $( compgen -f -- "$cur" ) $( compgen -d -- "$cur" ) )
	return 0
}

complete -F _rep rep
