package arg

import (
	"text/template"

	"github.com/grimdork/climate/str"
)

// CompList is a list of commands with their options.
type CompList struct {
	// AppName fot references.
	AppName string
	// Options for the main program.
	Options []CompOpt
	// Commands and their options.
	Commands []CompCmd
}

type opt interface {
	CompList | CompCmd
}

// CompCmd is a completion for a tool command.
type CompCmd struct {
	// Name op the command.
	Name string
	// Options list.
	Options []CompOpt
}

// CompOpt is a completion option.
type CompOpt struct {
	// Name of the option.
	Name string
}

// Completions returns a string containing the completion script.
func (opt *Options) Completions(appname string) (string, error) {
	buf := str.NewStringer()
	tpl := template.New("completions")
	tpl = tpl.Delims("##", "@@")
	tpl, err := tpl.Parse(script)
	if err != nil {
		return "", err
	}

	cl := &CompList{AppName: appname}
	for _, o := range opt.short {
		cl.Options = append(cl.Options, CompOpt{Name: "-" + o.ShortName})
	}

	for _, o := range opt.long {
		cl.Options = append(cl.Options, CompOpt{Name: "--" + o.LongName})
	}

	for _, c := range opt.commands {
		cmd := CompCmd{Name: c.Name}
		cl.Commands = append(cl.Commands, cmd)
		if c.Options == nil {
			continue
		}

		for _, o := range c.Options.short {
			cmd.Options = append(cmd.Options, CompOpt{Name: "-" + o.ShortName})
		}

		for _, o := range c.Options.long {
			cmd.Options = append(cmd.Options, CompOpt{Name: "--" + o.LongName})
		}
	}

	err = tpl.Execute(buf, cl)
	if err != nil {
		panic(err)
	}

	return buf.String(), nil
}

const script = `#!/usr/bin/env bash
options="## range .Options @@## .Name @@ ##end@@"
tools="## range .Commands@@ ## .Name @@ ##end@@"
##range .Commands@@## .Name @@_opt="##range .Options@@## .Name @@ ##end@@"
##end@@

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

complete_files() {
	_compopt_o_filenames
	COMPREPLY+=( $( compgen -f -- "$cur" ) $( compgen -d -- "$cur" ) )
}

_## .AppName @@() {
	COMPREPLY=()
	local cur prev
	_get_comp_words_by_ref cur prev

	if [ ${COMP_CWORD} -eq 1 ]; then
		if [[ ${cur} == -* ]]; then
			COMPREPLY=( $(compgen -W "${options}" -- $cur) )
			return 0
		fi

		if [[ ${tools} == "" ]]; then
			complete_files
			return 0
		fi

		COMPREPLY=( $(compgen -W "${tools}" -- $cur) )
		return 0
	fi

	if [[ ${cur} == -* ]]; then
		case ${prev} in
##range .Commands@@			"## .Name @@")
				COMPREPLY=( $(compgen -W "${## .Name @@_opt}" -- $cur) )
				return 0
				;;
##end@@			*)
##range .Commands@@				if [ $(hasword ${prev} ${## .Name @@_opt}) == "1" ]; then
					COMPREPLY=( $(compgen -W "${## .Name @@_opt}" -- $cur) )
					return 0
				fi
##end@@
				if [[ $(hasword ${prev} ${options}) == "1" ]]; then
					COMPREPLY=( $(compgen -W "${options}" -- $cur) )
					return 0
				fi
				;;
		esac
	fi

	if [[ $(hasword ${prev} ${options}) == "1" ]]; then
		if [[ ${tools} == "" ]]; then
			complete_files
			return 0
		fi

		COMPREPLY=( $(compgen -W "${tools}" -- $cur) )
		return 0
	fi

	complete_files
	return 0
}

complete -F _## .AppName @@ ## .AppName @@
`
