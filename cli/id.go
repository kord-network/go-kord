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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/log"
	"github.com/moby/moby/pkg/term"
)

func init() {
	registerCommand("id", RunID, `
usage: meta id new [options]

Create a new META ID.

options:
	-k, --keystore <dir>   Keystore directory [default: dev/keystore]
`[1:])
}

func RunID(ctx *Context, args Args) error {
	switch {
	case args.Bool("new"):
		return RunIDNew(ctx, args)
	default:
		return errors.New("unknown id command")
	}
}

func RunIDNew(ctx *Context, args Args) error {
	log.Info("creating new META ID")
	passphrase, err := getPassphrase(ctx, true)
	if err != nil {
		return fmt.Errorf("error reading passphrase: %s", err)
	}
	ks := keystore.NewKeyStore(
		args.String("--keystore"),
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)
	account, err := ks.NewAccount(string(passphrase))
	if err != nil {
		return err
	}
	fmt.Fprintln(ctx.Stdout, account.Address.Hex())
	return nil
}

func getPassphrase(ctx *Context, confirm bool) ([]byte, error) {
	if stdin, ok := ctx.Stdin.(*os.File); ok && term.IsTerminal(stdin.Fd()) {
		state, err := term.SaveState(stdin.Fd())
		if err != nil {
			return nil, err
		}
		term.DisableEcho(stdin.Fd(), state)
		defer term.RestoreTerminal(stdin.Fd(), state)
	}
	stdin := bufio.NewReader(ctx.Stdin)
	fmt.Fprint(ctx.Stderr, "Passphrase: ")
	passphrase, err := stdin.ReadBytes('\n')
	fmt.Fprintln(ctx.Stderr)
	if err != nil {
		return nil, err
	}
	passphrase = passphrase[0 : len(passphrase)-1]
	if !confirm {
		return passphrase, nil
	}

	fmt.Fprintf(ctx.Stderr, "Repeat passphrase: ")
	confirmation, err := stdin.ReadBytes('\n')
	fmt.Fprintln(ctx.Stderr)
	if err != nil {
		return nil, err
	}
	confirmation = confirmation[0 : len(confirmation)-1]

	if !bytes.Equal(passphrase, confirmation) {
		return nil, errors.New("The entered passphrases do not match")
	}
	return passphrase, nil
}
