// This file is part of the go-meta library.
//
// Copyright (C) 2017 JAAK MUSIC LTD
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
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/docopt/docopt-go"
	"github.com/ethereum/go-ethereum/log"
	meta "github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/media"
	"github.com/meta-network/go-meta/media/ern"
	"github.com/rs/cors"
)

var usage = `
usage: meta media import [--url=<url>] [--source=<source>] cwr <file>...
       meta media import [--url=<url>] [--source=<source>] ern <file>...
       meta media import [--url=<url>] [--source=<source>] eidr <file>...
       meta media import [--url=<url>] [--source=<source>] musicbrainz <type> <postgres-uri>
       meta server [--datadir=<dir>] [--port=<port>] [--cors-domain=<domain>...]

options:
  --url=<url>                 The server URL [default: http://localhost:5000]
  --source=<source>           The source name (defaults to value of the META_SOURCE environment variable)
  --datadir=<dir>             The data directory [default: .meta]
  --port=<port>               The port to start the HTTP server on [default: 5000]
  --cors-domain=<domain>...   The allowed CORS domains
`[1:]

func Run(ctx context.Context, cmdArgs ...string) error {
	v, _ := docopt.Parse(usage, cmdArgs, true, "0.0.1", false)
	args := Args(v)

	switch {
	case args.Bool("media"):
		return RunMedia(ctx, args)
	case args.Bool("server"):
		return RunServer(ctx, args)
	default:
		return errors.New("unknown command")
	}
}

func RunMedia(ctx context.Context, args Args) error {
	switch {
	case args.Bool("import"):
		return RunMediaImport(ctx, args)
	default:
		return errors.New("unknown media command")
	}
}

func RunMediaImport(ctx context.Context, args Args) error {
	source := args.String("--source")
	if source == "" {
		source = os.Getenv("META_SOURCE")
	}
	if source == "" {
		return errors.New("missing --source or META_SOURCE")
	}

	client := media.NewClient(
		args.String("--url")+"/media/graphql",
		&media.Source{Name: source},
	)

	switch {
	case args.Bool("ern"):
		return RunMediaImportERN(ctx, client, args)
	case args.Bool("cwr"):
		return RunMediaImportCWR(ctx, client, args)
	case args.Bool("eidr"):
		return RunMediaImportEIDR(ctx, client, args)
	case args.Bool("musicbrainz"):
		return RunMediaImportMusicBrainz(ctx, client, args)
	default:
		return errors.New("unknown media import command")
	}
}

func RunMediaImportERN(ctx context.Context, client *media.Client, args Args) error {
	importer := ern.NewImporter(client)
	for _, file := range args.List("<file>") {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil {
			return err
		}
		log.Info("importing ERN", "path", file, "size", info.Size())
		if err := importer.ImportERN(f); err != nil {
			return err
		}
	}
	return nil
}

func RunMediaImportCWR(ctx context.Context, client *media.Client, args Args) error {
	return nil
}

func RunMediaImportEIDR(ctx context.Context, client *media.Client, args Args) error {
	return nil
}

func RunMediaImportMusicBrainz(ctx context.Context, client *media.Client, args Args) error {
	return nil
}

func RunServer(ctx context.Context, args Args) error {
	dataDir := args.String("--datadir")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	ensDir := filepath.Join(dataDir, "ens")
	if err := os.MkdirAll(ensDir, 0755); err != nil {
		return err
	}
	store, err := meta.NewStore(dataDir, meta.LocalENS(ensDir))
	if err != nil {
		return err
	}

	identityIndex, err := identity.NewIndex(store)
	if err != nil {
		return err
	}
	defer identityIndex.Close()

	mediaIndex, err := media.NewIndex(store)
	if err != nil {
		return err
	}
	defer mediaIndex.Close()

	srv, err := NewServer(store, identityIndex, mediaIndex)
	if err != nil {
		return err
	}
	port := args.String("--port")
	if port == "" {
		port = "5000"
	}
	addr := "0.0.0.0:" + port
	log.Info("starting META HTTP server", "addr", addr)
	httpSrv := http.Server{
		Addr:    addr,
		Handler: srv,
	}
	if _, ok := args["--cors-domain"]; ok {
		httpSrv.Handler = cors.New(cors.Options{
			AllowedOrigins: args.List("--cors-domain"),
			AllowedMethods: []string{"POST", "GET", "DELETE", "PATCH", "PUT"},
			MaxAge:         600,
			AllowedHeaders: []string{"*"},
		}).Handler(httpSrv.Handler)
	}
	go func() {
		<-ctx.Done()
		log.Info("stopping HTTP server")
		httpSrv.Close()
	}()
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
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

func (a Args) Int(name string) int {
	v, ok := a[name]
	if !ok {
		panic(fmt.Sprintf("missing arg: %s", name))
	}
	if v == nil {
		return 0
	}
	s, ok := v.(string)
	if !ok {
		panic(fmt.Sprintf("invalid int arg: %s", name))
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("invalid int arg: %s", name))
	}
	return i
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
