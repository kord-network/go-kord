// This file is part of the go-meta library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// If you have any questions please contact yo@jaak.io

package cli

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"
)

var usage = `
usage: meta [options] <command> [<args>...]

Options:
        -h, --help        show this usage message
        --version         print the version
        --verbosity <n>   logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail [default: 3]

Commands:
        help     show usage for a specific command
        node     run a META node
        load     load quads into META

See 'meta help <command>' for more information on a specific command.
`[1:]

func Run(ctx *Context, argv ...string) error {
	if ctx.Stdin == nil {
		ctx.Stdin = os.Stdin
	}
	if ctx.Stdout == nil {
		ctx.Stdout = os.Stdout
	}
	if ctx.Stderr == nil {
		ctx.Stderr = os.Stderr
	}

	v, err := docopt.Parse(usage, argv, true, "0.0.1", true)
	if err != nil {
		return err
	}
	args := Args(v)

	cmd := args.String("<command>")
	cmdArgs := args.List("<args>")

	if cmd == "help" {
		if len(cmdArgs) == 0 {
			// 'meta help' so just print usage
			fmt.Println(usage)
			return nil
		}

		// 'meta help <command>' so translate to
		// 'meta <command> --help' and let docopt
		// print the command's usage
		cmd = cmdArgs[0]
		cmdArgs = []string{"--help"}
	}

	if v := args.String("--verbosity"); v != "" {
		if _, err := setLogVerbosity(v); err != nil {
			return err
		}
	}

	return runCommand(ctx, cmd, cmdArgs...)
}

func runCommand(ctx *Context, name string, argv ...string) error {
	argv = append([]string{name}, argv...)

	cmd, ok := commands[name]
	if !ok {
		return fmt.Errorf("%s is not a valid meta command. See 'meta help'", name)
	}

	v, err := docopt.Parse(cmd.usage, argv, true, "", false)
	if err != nil {
		return err
	}
	ctx.Args = Args(v)

	return cmd.run(ctx)
}

type runFn func(*Context) error

type command struct {
	run   runFn
	usage string
}

var commands = make(map[string]*command)

func registerCommand(name string, run runFn, usage string) {
	if _, ok := commands[name]; ok {
		panic(fmt.Sprintf("command already registered: %s", name))
	}
	commands[name] = &command{
		run:   run,
		usage: usage,
	}
}

type Args map[string]interface{}

func (a Args) String(name string) string {
	v, ok := a[name]
	if !ok {
		panic(fmt.Sprintf("missing arg: %s", name))
	}
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		panic(fmt.Sprintf("invalid string arg: %s", name))
	}
	return s
}

func (a Args) List(name string) []string {
	v, ok := a[name]
	if !ok {
		panic(fmt.Sprintf("missing arg: %s", name))
	}
	if v == nil {
		return nil
	}
	l, ok := v.([]string)
	if !ok {
		panic(fmt.Sprintf("invalid list arg: %s", name))
	}
	return l
}

func (a Args) Bool(name string) bool {
	v, ok := a[name]
	if !ok {
		panic(fmt.Sprintf("missing arg: %s", name))
	}
	s, ok := v.(bool)
	if !ok {
		panic(fmt.Sprintf("invalid bool arg: %s", name))
	}
	return s
}
