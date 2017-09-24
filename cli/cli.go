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
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ipfs/go-cid"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/cwr"
	"github.com/meta-network/go-meta/ern"
	"github.com/meta-network/go-meta/musicbrainz"
	"github.com/meta-network/go-meta/xml"
)

var usage = `
usage: meta import xml <file> [<context>...]
       meta import xsd <name> <uri> [<file>]
       meta dump [--format=<format>] <path>
       meta server [--port=<port>] [--musicbrainz-index=<sqlite3-uri>] [--cwr-index=<sqlite3-uri>]
       meta musicbrainz convert <postgres-uri>
       meta musicbrainz index <sqlite3-uri>
       meta cwr convert <file> <cwr-python-dir>
       meta cwr index <sqlite3-uri>
       meta ern convert <files>...
       meta ern index <sqlite3-uri>
`[1:]

type CLI struct {
	store  *meta.Store
	stdin  io.Reader
	stdout io.Writer
}

func New(store *meta.Store, stdin io.Reader, stdout io.Writer) *CLI {
	return &CLI{store, stdin, stdout}
}

func (cli *CLI) Run(ctx context.Context, cmdArgs ...string) error {
	v, _ := docopt.Parse(usage, cmdArgs, true, "0.0.1", false)
	args := Args(v)

	switch {
	case args.Bool("import"):
		return cli.RunImport(ctx, args)
	case args.Bool("dump"):
		return cli.RunDump(ctx, args)
	case args.Bool("server"):
		return cli.RunServer(ctx, args)
	case args.Bool("musicbrainz"):
		return cli.RunMusicBrainz(ctx, args)
	case args.Bool("cwr"):
		return cli.RunCwr(ctx, args)
	case args.Bool("ern"):
		return cli.RunERN(ctx, args)
	default:
		return errors.New("unknown command")
	}
}

func (cli *CLI) RunImport(ctx context.Context, args Args) error {
	switch {
	case args.Bool("xml"):
		return cli.RunImportXML(ctx, args)
	case args.Bool("xsd"):
		return cli.RunImportXMLSchema(ctx, args)
	default:
		return errors.New("unknown import format")
	}
}

func (cli *CLI) RunImportXML(ctx context.Context, args Args) error {
	f, err := os.Open(args.String("<file>"))
	if err != nil {
		return err
	}
	defer f.Close()

	var context []*cid.Cid
	if contextArg, ok := args["<context>"]; ok {
		for _, v := range contextArg.([]string) {
			cid, err := cid.Decode(v)
			if err != nil {
				return fmt.Errorf("invalid CID in --context value %q: %s", v, err)
			}
			context = append(context, cid)
		}
	}

	obj, err := metaxml.EncodeXML(f, context, cli.store.Put)
	if err != nil {
		return err
	}

	log.Info("object created", "cid", obj.Cid())

	return nil
}

func (cli *CLI) RunImportXMLSchema(ctx context.Context, args Args) error {
	var src io.Reader
	if path := args.String("<file>"); path != "" {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		src = f
	} else {
		uri := args.String("<uri>")
		res, err := http.Get(uri)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status from %s: %s", uri, res.Status)
		}
		src = res.Body
	}

	obj, err := metaxml.EncodeXMLSchema(src, args.String("<name>"), args.String("<uri>"))
	if err != nil {
		return err
	}

	if err := cli.store.Put(obj); err != nil {
		return err
	}

	log.Info("object created", "cid", obj.Cid())

	return nil
}

func (cli *CLI) RunDump(ctx context.Context, args Args) error {
	path := strings.Split(args.String("<path>"), "/")
	cid, err := cid.Decode(path[0])
	if err != nil {
		return err
	}
	obj, err := cli.store.Get(cid)
	if err != nil {
		return err
	}
	if len(path) == 1 {
		return json.NewEncoder(cli.stdout).Encode(obj)
	}
	graph := meta.NewGraph(cli.store, obj)
	v, err := graph.Get(path[1:]...)
	if err != nil {
		return err
	}
	return json.NewEncoder(cli.stdout).Encode(v)
}

func (cli *CLI) RunServer(ctx context.Context, args Args) error {
	var musicbrainzDB *sql.DB = nil
	var cwrDB *sql.DB = nil
	if uri := args.String("--musicbrainz-index"); uri != "" {
		db, err := sql.Open("sqlite3", uri)
		if err != nil {
			return err
		}
		defer db.Close()
		musicbrainzDB = db
	}
	if uri := args.String("--cwr-index"); uri != "" {
		db, err := sql.Open("sqlite3", uri)
		if err != nil {
			return err
		}
		defer db.Close()
		cwrDB = db
	}
	srv, err := NewServer(cli.store, musicbrainzDB, cwrDB)
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

func (cli *CLI) RunMusicBrainz(ctx context.Context, args Args) error {
	switch {
	case args.Bool("convert"):
		return cli.RunMusicBrainzConvert(ctx, args)
	case args.Bool("index"):
		return cli.RunMusicBrainzIndex(ctx, args)
	default:
		return errors.New("unknown musicbrainz command")
	}
}

func (cli *CLI) RunMusicBrainzConvert(ctx context.Context, args Args) error {
	db, err := sql.Open("postgres", args.String("<postgres-uri>"))
	if err != nil {
		return err
	}
	defer db.Close()

	// run the converter in a goroutine so that we only exit once all CIDs
	// have been read from the stream
	stream := make(chan *cid.Cid)
	errC := make(chan error, 1)
	go func() {
		defer close(stream)
		errC <- musicbrainz.NewConverter(db, cli.store).ConvertArtists(ctx, stream)
	}()

	// output the resulting CIDs to stdout
	for cid := range stream {
		fmt.Fprintln(cli.stdout, cid.String())
	}

	return <-errC
}

func (cli *CLI) RunMusicBrainzIndex(ctx context.Context, args Args) error {
	db, err := sql.Open("sqlite3", args.String("<sqlite3-uri>"))
	if err != nil {
		return err
	}
	defer db.Close()

	indexer, err := musicbrainz.NewIndexer(db, cli.store)
	if err != nil {
		return err
	}

	// stream the CIDs from stdin
	stream := make(chan *cid.Cid)
	go func() {
		defer close(stream)
		s := bufio.NewScanner(cli.stdin)
		for s.Scan() {
			cid, err := cid.Parse(s.Text())
			if err != nil {
				log.Error("error parsing cid", "value", s.Text(), "err", err)
				return
			}
			stream <- cid
		}
	}()

	return indexer.IndexArtists(ctx, stream)
}

func (cli *CLI) RunCwr(ctx context.Context, args Args) error {
	switch {
	case args.Bool("convert"):
		return cli.RunCwrConvert(ctx, args)
	case args.Bool("index"):
		return cli.RunCwrIndex(ctx, args)
	default:
		return errors.New("unknown cwr command")
	}
}

func (cli *CLI) RunCwrConvert(ctx context.Context, args Args) error {
	file, err := os.Open(args.String("<file>"))
	if err != nil {
		return err
	}
	defer file.Close()

	// run the converter in a goroutine so that we only exit once all CIDs
	// have been read from the stream
	stream := make(chan *cid.Cid)
	errC := make(chan error, 1)
	go func() {
		defer close(stream)
		errC <- cwr.NewConverter(cli.store).ConvertRegisteredWork(ctx, stream, file, args.String("<cwr-python-dir>"))
	}()

	// output the resulting CIDs to stdout
	for cid := range stream {
		fmt.Fprintln(cli.stdout, cid.String())
	}

	return <-errC
}

func (cli *CLI) RunCwrIndex(ctx context.Context, args Args) error {

	db, err := sql.Open("sqlite3", args.String("<sqlite3-uri>"))
	if err != nil {
		return err
	}
	defer db.Close()

	indexer, err := cwr.NewIndexer(db, cli.store)
	if err != nil {
		return err
	}

	// stream the CIDs from stdin
	stream := make(chan *cid.Cid)
	go func() {
		defer close(stream)
		s := bufio.NewScanner(cli.stdin)
		for s.Scan() {
			cid, err := cid.Parse(s.Text())
			if err != nil {
				log.Error("error parsing cid", "value", s.Text(), "err", err)
				return
			}
			stream <- cid
		}
	}()

	return indexer.IndexRegisteredWorks(ctx, stream)
}

func (cli *CLI) RunERN(ctx context.Context, args Args) error {
	switch {
	case args.Bool("convert"):
		return cli.RunERNConvert(ctx, args)
	case args.Bool("index"):
		return cli.RunERNIndex(ctx, args)
	default:
		return errors.New("unknown ern command")
	}
}

func (cli *CLI) RunERNConvert(ctx context.Context, args Args) error {
	converter := ern.NewConverter(cli.store)
	files := args.List("<files>")
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		cid, err := converter.ConvertERN(f)
		if err != nil {
			return err
		}
		fmt.Fprintln(cli.stdout, cid.String())
	}
	return nil
}

func (cli *CLI) RunERNIndex(ctx context.Context, args Args) error {
	db, err := sql.Open("sqlite3", args.String("<sqlite3-uri>"))
	if err != nil {
		return err
	}
	defer db.Close()

	indexer, err := ern.NewIndexer(db, cli.store)
	if err != nil {
		return err
	}

	// stream the CIDs from stdin
	stream := make(chan *cid.Cid)
	go func() {
		defer close(stream)
		s := bufio.NewScanner(cli.stdin)
		for s.Scan() {
			cid, err := cid.Parse(s.Text())
			if err != nil {
				log.Error("error parsing cid", "value", s.Text(), "err", err)
				return
			}
			stream <- cid
		}
	}()

	return indexer.Index(ctx, stream)
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
