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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
       meta index cwr <ens-name> [--count=<count>]
       meta index ern <ens-name> [--count=<count>]
       meta index eidr <ens-name> [--count=<count>]
       meta index musicbrainz <type> <ens-name> [--count=<count>]
       meta dump [--format=<format>] <path>
       meta server [--port=<port>] [--cors-domain=<domain>...] [--index=<index>...]

options:
  --source=<source>           The value to set as @source on all created META objects
                              (defaults to value of the META_SOURCE environment variable).

  --count=<count>             The number of CIDs to index from a stream.

  --format=<format>           The format to dump objects when running 'meta dump'.

  --port=<port>               The port to start the HTTP server on.

  --cors-domain=<domain>...   The allowed CORS domains.

  --index=<index>...          One or more SQLite3 indexes for the HTTP server where <index>
                              has the format <type>:<name>, with <type> being one of
                              'musicbrainz', 'ern', 'cwr' ,'eidr','identity' or 'claim' and <name> being the
                              ENS name of the index. For example:
                              '--index ern:ern.index.meta --index cwr:cwr.index.meta'
`[1:]

type CLI struct {
	store  *meta.Store
	stdin  io.Reader
	stdout io.Writer
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
	var indexFn func(context.Context, *meta.Index, Args) error
	switch {
	case args.Bool("cwr"):
		indexFn = cli.RunCwrIndex
	case args.Bool("ern"):
		indexFn = cli.RunERNIndex
	case args.Bool("eidr"):
		indexFn = cli.RunEIDRIndex
	case args.Bool("musicbrainz"):
		indexFn = cli.RunMusicBrainzIndex
	default:
		return errors.New("unknown index")
	}

	index, err := cli.store.OpenIndex(args.String("<ens-name>"))
	if err != nil {
		return err
	}
	defer index.Close()
	return indexFn(ctx, index, args)
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
	// parse the --index args which have the format <type>:<name>
	// where <type> is one of musicbrainz, ern, eidr or cwr and
	// <name> is the ENS name of the relevant index.
	indexes := make(map[string]*meta.Index)
	for _, index := range args.List("--index") {
		typeName := strings.SplitN(index, ":", 2)
		if len(typeName) != 2 {
			return fmt.Errorf("invalid --index: %q", index)
		}
		typ := typeName[0]
		name := typeName[1]

		switch typ {
		case "musicbrainz", "ern", "eidr", "cwr", "identity", "claim":
		default:
			return fmt.Errorf("invalid --index type %q", typ)
		}
		db, err := cli.store.OpenIndex(name)
		if err != nil {
			return err
		}
		defer db.Close()
		indexes[typ] = db
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
	var convertFn func(context.Context, *meta.StreamWriter, string) error
	switch args.String("<type>") {
	case "artists":
		convertFn = converter.ConvertArtists
	case "recording-work-links":
		convertFn = converter.ConvertRecordingWorkLinks
	default:
		return errors.New("unknown musicbrainz convert command")
	}

	streamName := fmt.Sprintf("%s.musicbrainz.meta", args.String("<type>"))
	stream, err := cli.store.StreamWriter(streamName)
	if err != nil {
		return err
	}
	defer stream.Close()

	return convertFn(ctx, stream, args.String("--source"))
}

func (cli *CLI) RunMusicBrainzIndex(ctx context.Context, index *meta.Index, args Args) error {
	indexer, err := musicbrainz.NewIndexer(index, cli.store)
	if err != nil {
		return err
	}

	var indexFn func(context.Context, *meta.StreamReader) error
	switch args.String("<type>") {
	case "artists":
		indexFn = indexer.IndexArtists
	case "recording-work-links":
		indexFn = indexer.IndexRecordingWorkLinks
	default:
		return errors.New("unknown musicbrainz index command")
	}

	streamName := fmt.Sprintf("%s.musicbrainz.meta", args.String("<type>"))
	var streamOpts []meta.StreamOpts
	if _, ok := args["--count"]; ok {
		streamOpts = append(streamOpts, meta.StreamLimit(args.Int("--count")))
	}
	reader, err := cli.store.StreamReader(streamName, streamOpts...)
	if err != nil {
		return err
	}
	defer reader.Close()

	return indexFn(ctx, reader)
}

func (cli *CLI) RunCwrConvert(ctx context.Context, args Args) error {
	stream, err := cli.store.StreamWriter("cwr.meta")
	if err != nil {
		return err
	}
	defer stream.Close()

	converter := cwr.NewConverter(cli.store)
	for _, file := range args.List("<files>") {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		cid, err := converter.ConvertCWR(f, args.String("--source"))
		if err != nil {
			return err
		}

		if err := stream.Write(cid); err != nil {
			return err
		}
	}
	return nil
}

func (cli *CLI) RunCwrIndex(ctx context.Context, index *meta.Index, args Args) error {
	indexer, err := cwr.NewIndexer(index, cli.store)
	if err != nil {
		return err
	}

	var streamOpts []meta.StreamOpts
	if _, ok := args["--count"]; ok {
		streamOpts = append(streamOpts, meta.StreamLimit(args.Int("--count")))
	}
	reader, err := cli.store.StreamReader("cwr.meta", streamOpts...)
	if err != nil {
		return err
	}
	defer reader.Close()

	return indexer.Index(ctx, reader)
}

func (cli *CLI) RunERNConvert(ctx context.Context, args Args) error {
	stream, err := cli.store.StreamWriter("ern.meta")
	if err != nil {
		return err
	}
	defer stream.Close()

	converter := ern.NewConverter(cli.store)
	for _, file := range args.List("<files>") {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		cid, err := converter.ConvertERN(f, args.String("--source"))
		if err != nil {
			return err
		}

		if err := stream.Write(cid); err != nil {
			return err
		}
	}
	return nil
}

func (cli *CLI) RunERNIndex(ctx context.Context, index *meta.Index, args Args) error {
	indexer, err := ern.NewIndexer(index, cli.store)
	if err != nil {
		return err
	}

	var streamOpts []meta.StreamOpts
	if _, ok := args["--count"]; ok {
		streamOpts = append(streamOpts, meta.StreamLimit(args.Int("--count")))
	}
	reader, err := cli.store.StreamReader("ern.meta", streamOpts...)
	if err != nil {
		return err
	}
	defer reader.Close()

	return indexer.Index(ctx, reader)
}

func (cli *CLI) RunEIDRConvert(ctx context.Context, args Args) error {
	stream, err := cli.store.StreamWriter("eidr.meta")
	if err != nil {
		return err
	}
	defer stream.Close()

	converter := eidr.NewConverter(cli.store)
	for _, file := range args.List("<files>") {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		cid, err := converter.ConvertEIDRXML(f, args.String("--source"))
		if err != nil {
			return err
		}

		if err := stream.Write(cid); err != nil {
			return err
		}
	}

	return nil
}

func (cli *CLI) RunEIDRIndex(ctx context.Context, index *meta.Index, args Args) error {
	indexer, err := eidr.NewIndexer(index, cli.store)
	if err != nil {
		return err
	}

	var streamOpts []meta.StreamOpts
	if _, ok := args["--count"]; ok {
		streamOpts = append(streamOpts, meta.StreamLimit(args.Int("--count")))
	}
	reader, err := cli.store.StreamReader("eidr.meta", streamOpts...)
	if err != nil {
		return err
	}
	defer reader.Close()

	return indexer.Index(ctx, reader)
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
