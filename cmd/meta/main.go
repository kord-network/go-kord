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

package main

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
	"os/signal"
	"strings"
	"syscall"

	"github.com/docopt/docopt-go"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore/fs"
	_ "github.com/lib/pq"
	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/cwr"
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
`[1:]

type Main struct {
	store *meta.Store
}

func NewMain(store *meta.Store) *Main {
	return &Main{store}
}

func main() {
	log.Root().SetHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(true)))

	store, err := openStore()
	if err != nil {
		log.Crit("error opening meta store", "err", err)
	}

	args, _ := docopt.Parse(usage, os.Args[1:], true, "0.0.1", false)

	m := NewMain(store)
	if err := m.Run(args); err != nil {
		log.Crit("error running meta command", "err", err)
	}
}

func (m *Main) Run(args Args) error {
	switch {
	case args.Bool("import"):
		return m.RunImport(args)
	case args.Bool("dump"):
		return m.RunDump(args)
	case args.Bool("server"):
		return m.RunServer(args)
	case args.Bool("musicbrainz"):
		return m.RunMusicBrainz(args)
	case args.Bool("cwr"):
		return m.RunCwr(args)
	default:
		return errors.New("unknown command")
	}
}

func (m *Main) RunImport(args Args) error {
	switch {
	case args.Bool("xml"):
		return m.RunImportXML(args)
	case args.Bool("xsd"):
		return m.RunImportXMLSchema(args)
	default:
		return errors.New("unknown import format")
	}
}

func (m *Main) RunImportXML(args Args) error {
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

	obj, err := metaxml.EncodeXML(f, context, m.store.Put)
	if err != nil {
		return err
	}

	log.Info("object created", "cid", obj.Cid())

	return nil
}

func (m *Main) RunImportXMLSchema(args Args) error {
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

	if err := m.store.Put(obj); err != nil {
		return err
	}

	log.Info("object created", "cid", obj.Cid())

	return nil
}

func (m *Main) RunDump(args Args) error {
	path := strings.Split(args.String("<path>"), "/")
	cid, err := cid.Decode(path[0])
	if err != nil {
		return err
	}
	obj, err := m.store.Get(cid)
	if err != nil {
		return err
	}
	if len(path) == 1 {
		return json.NewEncoder(os.Stdout).Encode(obj)
	}
	graph := meta.NewGraph(m.store, obj)
	v, err := graph.Get(path[1:]...)
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(v)
}

func (m *Main) RunServer(args Args) error {
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
	srv, err := NewServer(m.store, musicbrainzDB, cwrDB)
	if err != nil {
		return err
	}
	port := args.String("--port")
	if port == "" {
		port = "5000"
	}
	addr := "0.0.0.0:" + port
	log.Info("starting META HTTP server", "addr", addr)
	return http.ListenAndServe(addr, srv)
}

func (m *Main) RunMusicBrainz(args Args) error {
	switch {
	case args.Bool("convert"):
		return m.RunMusicBrainzConvert(args)
	case args.Bool("index"):
		return m.RunMusicBrainzIndex(args)
	default:
		return errors.New("unknown musicbrainz command")
	}
}

func (m *Main) RunMusicBrainzConvert(args Args) error {
	db, err := sql.Open("postgres", args.String("<postgres-uri>"))
	if err != nil {
		return err
	}
	defer db.Close()

	// stream the CIDs to stdout
	stream := make(chan *cid.Cid)
	// shutdown gracefully on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Info("received signal, exiting...")
	}()
	go func() {
		err = musicbrainz.NewConverter(db, m.store).ConvertArtists(ctx, stream)
		close(stream)
	}()
	for cid := range stream {
		fmt.Fprintln(os.Stdout, cid.String())
	}
	return err
}

func (m *Main) RunMusicBrainzIndex(args Args) error {
	db, err := sql.Open("sqlite3", args.String("<sqlite3-uri>"))
	if err != nil {
		return err
	}
	defer db.Close()

	indexer, err := musicbrainz.NewIndexer(db, m.store)
	if err != nil {
		return err
	}

	// stream the CIDs from stdin
	stream := make(chan *cid.Cid)
	go func() {
		defer close(stream)
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			cid, err := cid.Parse(s.Text())
			if err != nil {
				log.Error("error parsing cid", "value", s.Text(), "err", err)
				return
			}
			stream <- cid
		}
	}()

	// shutdown gracefully on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Info("received signal, exiting...")
	}()

	return indexer.IndexArtists(ctx, stream)
}

func (m *Main) RunCwr(args Args) error {
	switch {
	case args.Bool("convert"):
		return m.RunCwrConvert(args)
	case args.Bool("index"):
		return m.RunCwrIndex(args)
	default:
		return errors.New("unknown cwr command")
	}
}

func (m *Main) RunCwrConvert(args Args) error {
	// stream the CIDs to stdout
	stream := make(chan *cid.Cid)
	// shutdown gracefully on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Info("received signal, exiting...")
	}()
	cwrFileReader, err := os.Open(args.String("<file>"))
	if err != nil {
		return err
	}
	defer cwrFileReader.Close()
	go func() {
		err = cwr.NewConverter(m.store).ConvertRegisteredWork(ctx, stream, cwrFileReader, args.String("<cwr-python-dir>"))
		close(stream)
	}()
	for cid := range stream {
		fmt.Fprintln(os.Stdout, cid.String())
	}
	return err
}

func (m *Main) RunCwrIndex(args Args) error {

	db, err := sql.Open("sqlite3", args.String("<sqlite3-uri>"))
	if err != nil {
		return err
	}
	defer db.Close()

	indexer, err := cwr.NewIndexer(db, m.store)
	if err != nil {
		return err
	}

	// stream the CIDs from stdin
	stream := make(chan *cid.Cid)
	go func() {
		defer close(stream)
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			cid, err := cid.Parse(s.Text())
			if err != nil {
				log.Error("error parsing cid", "value", s.Text(), "err", err)
				return
			}
			stream <- cid
		}
	}()

	// shutdown gracefully on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Info("received signal, exiting...")
	}()

	return indexer.IndexRegisteredWorks(ctx, stream)
}

func openStore() (*meta.Store, error) {
	metaDir := ".meta"
	if err := os.MkdirAll(metaDir, 0755); err != nil {
		return nil, err
	}
	store, err := fs.NewDatastore(metaDir)
	if err != nil {
		return nil, err
	}
	return meta.NewStore(store), nil
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
