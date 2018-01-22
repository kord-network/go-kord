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

package meta

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/tent/canonical-json-go"

	metasql "github.com/meta-network/go-meta/sql"
)

type ENS interface {
	Content(name string) (common.Hash, error)
}

type Storage struct {
	dir string
	dpa *storage.DPA
	ens ENS

	notifyMtx sync.RWMutex
	notifiers map[string]map[*notifier]struct{}
}

func NewStorage(dir string, dpa *storage.DPA, ens ENS) *Storage {
	return &Storage{
		dir:       dir,
		dpa:       dpa,
		ens:       ens,
		notifiers: make(map[string]map[*notifier]struct{}),
	}
}

func (s *Storage) GetDB(name string, notify chan struct{}) (string, metasql.Notifier, error) {
	path := filepath.Join(s.dir, name)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = s.fetchDB(name, path)
	}
	if err != nil {
		return "", nil, err
	}
	notifier := newNotifier(name, s, notify)
	s.addNotifier(notifier)
	return path, notifier, nil
}

func (s *Storage) UpdateDB(name string, hash common.Hash) error {
	tmp, err := ioutil.TempFile("", "meta-db")
	if err != nil {
		return err
	}
	defer tmp.Close()
	if err := s.fetchHash(hash, tmp); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	path := filepath.Join(s.dir, name)
	if err := os.Rename(tmp.Name(), path); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	s.notify(name)
	return nil
}

func (s *Storage) SaveDB(name string) (common.Hash, error) {
	path := filepath.Join(s.dir, name)
	f, err := os.Open(path)
	if err != nil {
		return common.Hash{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return common.Hash{}, err
	}
	key, err := s.dpa.Store(f, info.Size(), &sync.WaitGroup{}, &sync.WaitGroup{})
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(key[:]), nil
}

func (s *Storage) fetchDB(name, path string) error {
	hash, err := s.ens.Content(name)
	if err != nil {
		return err
	}
	if common.EmptyHash(hash) {
		return nil
	}
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	return s.fetchHash(hash, dst)
}

func (s *Storage) fetchHash(hash common.Hash, dst io.Writer) error {
	reader := s.dpa.Retrieve(storage.Key(hash[:]))
	size, err := reader.Size(nil)
	if err != nil {
		return err
	}
	n, err := io.Copy(dst, io.LimitReader(reader, size))
	if err != nil {
		return err
	} else if n != size {
		return fmt.Errorf("failed to fetch database, expected %d bytes, copied %d", size, n)
	}
	return nil
}

func (s *Storage) addNotifier(n *notifier) {
	s.notifyMtx.Lock()
	defer s.notifyMtx.Unlock()
	notifiers, ok := s.notifiers[n.name]
	if !ok {
		notifiers = make(map[*notifier]struct{})
		s.notifiers[n.name] = notifiers
	}
	notifiers[n] = struct{}{}
}

func (s *Storage) removeNotifier(n *notifier) {
	s.notifyMtx.Lock()
	defer s.notifyMtx.Unlock()
	if notifiers, ok := s.notifiers[n.name]; ok {
		delete(notifiers, n)
	}
}

func (s *Storage) notify(name string) {
	s.notifyMtx.RLock()
	defer s.notifyMtx.RUnlock()
	for n := range s.notifiers[name] {
		n.Notify()
	}
}

type notifier struct {
	name     string
	storage  *Storage
	ch       chan struct{}
	done     chan struct{}
	doneOnce sync.Once
	err      error
}

func newNotifier(name string, storage *Storage, ch chan struct{}) *notifier {
	return &notifier{
		name:    name,
		storage: storage,
		ch:      ch,
		done:    make(chan struct{}),
	}
}

func (n *notifier) Notify() {
	select {
	case n.ch <- struct{}{}:
	case <-n.done:
	}
}

func (n *notifier) Done() chan struct{} {
	return n.done
}

func (n *notifier) Close() {
	n.CloseWithErr(nil)
}

func (n *notifier) CloseWithErr(err error) {
	n.err = err
	n.storage.removeNotifier(n)
	n.doneOnce.Do(func() { close(n.done) })
}

func (n *notifier) Err() error {
	return n.err
}

type Backend interface {
	Apply(tx *Tx) (common.Hash, error)
}

type Signer interface {
	SignHash(address common.Address, hash []byte) ([]byte, error)
}

type QuadStore struct {
	graph.QuadStore

	address common.Address
	name    string
	backend Backend
	signer  Signer
	storage *Storage
}

func NewQuadStore(address common.Address, name string, backend Backend, signer Signer, storage *Storage) (*QuadStore, error) {
	if err := graph.InitQuadStore("meta", name, graph.Options{}); err != nil {
		return nil, err
	}
	qs, err := graph.NewQuadStore("meta", name, graph.Options{})
	if err != nil {
		return nil, err
	}
	return &QuadStore{
		QuadStore: qs,
		address:   address,
		name:      name,
		backend:   backend,
		signer:    signer,
		storage:   storage,
	}, nil
}

func (qs *QuadStore) ApplyDeltas(in []graph.Delta, opts graph.IgnoreOpts) error {
	deltas := make([]graph.Delta, len(in))
	for i, delta := range in {
		delta.Quad.Label = quad.String(qs.address.Hex())
		deltas[i] = delta
	}
	req := &Request{
		Name:   qs.name,
		Deltas: deltas,
	}
	sig, err := qs.signer.SignHash(qs.address, req.Hash())
	if err != nil {
		return err
	}
	tx := &Tx{
		Address: qs.address,
		Data:    req.Bytes(),
		Sig:     sig,
	}
	hash, err := qs.backend.Apply(tx)
	if err != nil {
		return err
	}
	return qs.storage.UpdateDB(qs.name, hash)
}

type State struct {
	mtx     sync.Mutex
	stores  map[string]graph.QuadStore
	storage *Storage
}

func NewState(storage *Storage) *State {
	return &State{
		stores:  make(map[string]graph.QuadStore),
		storage: storage,
	}
}

// Apply applies a transaction.
func (s *State) Apply(tx *Tx) (common.Hash, error) {
	// check the signature
	hash := crypto.Keccak256(tx.Data)
	pubKey, err := crypto.SigToPub(hash, tx.Sig)
	if err != nil {
		return common.Hash{}, err
	}
	address := crypto.PubkeyToAddress(*pubKey)
	if address != tx.Address {
		return common.Hash{}, errors.New("invalid signature")
	}

	// decode the request
	var req Request
	if err := json.Unmarshal(tx.Data, &req); err != nil {
		return common.Hash{}, err
	}

	// check all quads labelled with address
	for _, delta := range req.Deltas {
		label := delta.Quad.Label
		if label == nil {
			return common.Hash{}, errors.New("quad label must be set")
		}
		labelString, ok := label.Native().(string)
		if !ok {
			return common.Hash{}, fmt.Errorf("invalid label type: %T", label.Native())
		}
		if !common.IsHexAddress(labelString) {
			return common.Hash{}, fmt.Errorf("invalid address in quad label: %s", label)
		}
		addr := common.HexToAddress(labelString)
		if addr != address {
			return common.Hash{}, fmt.Errorf("invalid address in quad label: %s", addr)
		}
	}

	// ApplyDeltas to the store
	store, err := s.store(req.Name)
	if err != nil {
		return common.Hash{}, err
	}
	if err := store.ApplyDeltas(req.Deltas, graph.IgnoreOpts{}); err != nil {
		return common.Hash{}, err
	}

	// TODO: rollback ApplyDeltas if saving fails
	return s.storage.SaveDB(req.Name)
}

func (s *State) Close() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, store := range s.stores {
		store.Close()
	}
}

func (s *State) store(name string) (graph.QuadStore, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if store, ok := s.stores[name]; ok {
		return store, nil
	}
	store, err := graph.NewQuadStore("meta", name, graph.Options{})
	if err != nil {
		return nil, err
	}
	s.stores[name] = store
	return store, nil
}

type Request struct {
	Name   string        `json:"name"`
	Deltas []graph.Delta `json:"deltas"`
}

func (r *Request) Bytes() []byte {
	data, _ := cjson.Marshal(r)
	return data
}

func (r *Request) Hash() []byte {
	return crypto.Keccak256(r.Bytes())
}

// Tx is a transaction which can be applied to the META state.
type Tx struct {
	Address common.Address
	Data    []byte
	Sig     []byte
}

type txJSON struct {
	Address string `json:"address"`
	Data    string `json:"data"`
	Sig     string `json:"sig"`
}

func (tx *Tx) MarshalJSON() ([]byte, error) {
	return json.Marshal(txJSON{
		Address: tx.Address.Hex(),
		Data:    hexutil.Encode(tx.Data),
		Sig:     hexutil.Encode(tx.Sig),
	})
}

func (tx *Tx) UnmarshalJSON(data []byte) error {
	var v txJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	txData, err := hexutil.Decode(v.Data)
	if err != nil {
		return err
	}
	txSig, err := hexutil.Decode(v.Sig)
	if err != nil {
		return err
	}
	*tx = Tx{
		Address: common.HexToAddress(v.Address),
		Data:    txData,
		Sig:     txSig,
	}
	return nil
}
