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
	"github.com/meta-network/go-meta/eidr"
	"github.com/meta-network/go-meta/ern"
	"github.com/meta-network/go-meta/musicbrainz"
	"github.com/meta-network/go-meta/xml"
	"github.com/rs/cors"
)

var usage = `
usage: meta convert [--source=<source>] xml <file> [<context>...]
       meta convert [--source=<source>] xsd <name> <uri> [<file>]
       meta convert [--source=<source>] cwr <files>...
       meta convert [--source=<source>] ern <files>...
       meta convert [--source=<source>] eidr <files>...
       meta convert [--source=<source>] musicbrainz <type> <postgres-uri>
       meta index cwr <sqlite3-filename> [--bzzapi=<bzzuri>] [--bzzdir=<bzzdirhash>]
       meta index ern <sqlite3-filename> [--bzzapi=<bzzuri>] [--bzzdir=<bzzdirhash>]
       meta index eidr <sqlite3-filename> [--bzzapi=<bzzuri>] [--bzzdir=<bzzdirhash>]
       meta index musicbrainz <type> <sqlite3-filename> [--bzzapi=<bzzuri>] [--bzzdir=<bzzdirhash>]
       meta dump [--format=<format>] <path>
       meta server [--port=<port>] [--cors-domain=<domain>...] [--index=<index>...] [--bzzapi=<bzzuri>] [--bzzdir=<bzzdirhash>]

options:
  --source=<source>           The value to set as @source on all created META objects
                              (defaults to value of the META_SOURCE environment variable).

  --format=<format>           The format to dump objects when running 'meta dump'.

  --port=<port>               The port to start the HTTP server on.

  --cors-domain=<domain>...   The allowed CORS domains.

  --index=<index>...          One or more SQLite3 indexes for the HTTP server where <index>
                              has the format <name>:<path>, with <name> being one of
                              'musicbrainz', 'ern', 'cwr' or 'eidr' and <path> being the
                              path to the index. For example:
                              '--index ern:path/to/ern.db --index cwr:path/to/cwr.db'
`[1:]

type CLI struct {
	store  *meta.Store
	stdin  io.Reader
	stdout io.Writer
	bzz    *SwarmBackend
}

func New(store *meta.Store, stdin io.Reader, stdout io.Writer) *CLI {
	return &CLI{
		store:  store,
		stdin:  stdin,
		stdout: stdout,
	}
}

func (cli *CLI) Run(ctx context.Context, cmdArgs ...string) error {
	v, _ := docopt.Parse(usage, cmdArgs, true, "0.0.1", false)
	args := Args(v)

	if v, ok := args["--bzzapi"]; ok && v != nil {
		cli.bzz = &SwarmBackend{}
		if err := cli.bzz.OpenIndex(args.String("--bzzapi"), args.String("--bzzdir")); err != nil {
			return err
		}
		defer cli.bzz.CloseIndex()
	}

	switch {
	case args.Bool("convert"):
		if args.String("--source") == "" {
			source := os.Getenv("META_SOURCE")
			if source == "" {
				return errors.New("missing --source or META_SOURCE")
			}
			args["--source"] = source
		}
		return cli.RunConvert(ctx, args)
	case args.Bool("index"):
		return cli.RunIndex(ctx, args)
	case args.Bool("dump"):
		return cli.RunDump(ctx, args)
	case args.Bool("server"):
		return cli.RunServer(ctx, args)
	default:
		return errors.New("unknown command")
	}
}

func (cli *CLI) RunConvert(ctx context.Context, args Args) error {
	switch {
	case args.Bool("xml"):
		return cli.RunConvertXML(ctx, args)
	case args.Bool("xsd"):
		return cli.RunConvertXMLSchema(ctx, args)
	case args.Bool("cwr"):
		return cli.RunCwrConvert(ctx, args)
	case args.Bool("ern"):
		return cli.RunERNConvert(ctx, args)
	case args.Bool("eidr"):
		return cli.RunEIDRConvert(ctx, args)
	case args.Bool("musicbrainz"):
		return cli.RunMusicBrainzConvert(ctx, args)
	default:
		return errors.New("unknown convert format")
	}
}

func (cli *CLI) RunIndex(ctx context.Context, args Args) (err error) {
	filename := args.String("<sqlite3-filename>")
	filepath := filename
	if cli.bzz != nil {
		filepath, err = cli.bzz.GetIndexFile(filename, false)
		if err != nil {
			return err
		}
	}
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	defer db.Close()

	switch {
	case args.Bool("cwr"):
		err = cli.RunCwrIndex(ctx, db, args)
	case args.Bool("ern"):
		err = cli.RunERNIndex(ctx, db, args)
	case args.Bool("eidr"):
		err = cli.RunEIDRIndex(ctx, db, args)
	case args.Bool("musicbrainz"):
		err = cli.RunMusicBrainzIndex(ctx, db, args)
	default:
		err = errors.New("unknown index")
	}
	if err != nil {
		return err
	}

	if cli.bzz != nil {
		hash, err := cli.bzz.PutIndexFile(filename)
		if err != nil {
			return err
		}
		cli.stdout.Write([]byte(hash))
	}
	return
}

func (cli *CLI) RunConvertXML(ctx context.Context, args Args) error {
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

	converter := metaxml.NewConverter(cli.store)
	obj, err := converter.ConvertXML(f, context, args.String("--source"))
	if err != nil {
		return err
	}

	log.Info("object created", "cid", obj.Cid())

	return nil
}

func (cli *CLI) RunConvertXMLSchema(ctx context.Context, args Args) error {
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

	converter := metaxml.NewConverter(cli.store)
	obj, err := converter.ConvertXMLSchema(src, args.String("<name>"), args.String("<uri>"), args.String("--source"))
	if err != nil {
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
	if cli.bzz != nil {
		err := cli.bzz.OpenIndex(args.String("--bzzapi"), args.String("--bzzdir"))
		if err != nil {
			return err
		}
		defer cli.bzz.CloseIndex()
	}

	// parse the --index args which have the format <name>:<path>
	// where <name> is one of musicbrainz, ern, eidr or cwr and
	// <path> is the path to the relevant index.
	indexes := make(map[string]*sql.DB)
	for _, index := range args.List("--index") {
		namePath := strings.SplitN(index, ":", 2)
		if len(namePath) != 2 {
			return fmt.Errorf("invalid --index: %q", index)
		}
		name := namePath[0]
		path := namePath[1]

		switch name {
		case "musicbrainz", "ern", "eidr", "cwr":
		default:
			return fmt.Errorf("invalid --index name %q", name)
		}
		if cli.bzz != nil {
			var err error
			path, err = cli.bzz.GetIndexFile(path, true)
			if err != nil {
				return err
			}
		}
		db, err := sql.Open("sqlite3", path)
		if err != nil {
			return err
		}
		defer db.Close()
		indexes[name] = db
	}

	srv, err := NewServer(cli.store, indexes)
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

func (cli *CLI) RunMusicBrainzConvert(ctx context.Context, args Args) error {
	db, err := sql.Open("postgres", args.String("<postgres-uri>"))
	if err != nil {
		return err
	}
	defer db.Close()

	converter := musicbrainz.NewConverter(db, cli.store)
	var convertFn func(context.Context, chan *cid.Cid, string) error
	switch args.String("<type>") {
	case "artists":
		convertFn = converter.ConvertArtists
	case "recording-work-links":
		convertFn = converter.ConvertRecordingWorkLinks
	default:
		return errors.New("unknown musicbrainz convert command")
	}

	// run the converter in a goroutine so that we only exit once all CIDs
	// have been read from the stream
	stream := make(chan *cid.Cid)
	errC := make(chan error, 1)
	go func() {
		defer close(stream)
		errC <- convertFn(ctx, stream, args.String("--source"))
	}()

	// output the resulting CIDs to stdout
	for cid := range stream {
		fmt.Fprintln(cli.stdout, cid.String())
	}

	return <-errC
}

func (cli *CLI) RunMusicBrainzIndex(ctx context.Context, db *sql.DB, args Args) error {
	indexer, err := musicbrainz.NewIndexer(db, cli.store)
	if err != nil {
		return err
	}

	var indexFn func(context.Context, chan *cid.Cid) error
	switch args.String("<type>") {
	case "artists":
		indexFn = indexer.IndexArtists
	case "recording-work-links":
		indexFn = indexer.IndexRecordingWorkLinks
	default:
		return errors.New("unknown musicbrainz index command")
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
	return indexFn(ctx, stream)
}

func (cli *CLI) RunCwrConvert(ctx context.Context, args Args) error {
	converter := cwr.NewConverter(cli.store)
	files := args.List("<files>")
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		cid, err := converter.ConvertCWR(f, args.String("--source"))
		if err != nil {
			return err
		}
		fmt.Fprintln(cli.stdout, cid.String())
	}
	return nil
}

func (cli *CLI) RunCwrIndex(ctx context.Context, db *sql.DB, args Args) error {
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
	return indexer.Index(ctx, stream)
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
		cid, err := converter.ConvertERN(f, args.String("--source"))
		if err != nil {
			return err
		}
		fmt.Fprintln(cli.stdout, cid.String())
	}
	return nil
}

func (cli *CLI) RunERNIndex(ctx context.Context, db *sql.DB, args Args) error {
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

func (cli *CLI) RunEIDRConvert(ctx context.Context, args Args) error {
	converter := eidr.NewConverter(cli.store)
	files := args.List("<files>")
	stream := make(chan *cid.Cid)
	go func() {
		defer close(stream)
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				continue
			}
			defer f.Close()
			cid, err := converter.ConvertEIDRXML(f, args.String("--source"))
			if err != nil {
				continue
			}
			stream <- cid
		}
	}()

	// output the resulting CIDs to stdout
	for cid := range stream {
		fmt.Fprintln(cli.stdout, cid.String())
	}
	return nil
}

func (cli *CLI) RunEIDRIndex(ctx context.Context, db *sql.DB, args Args) error {
	indexer, err := eidr.NewIndexer(db, cli.store)
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
