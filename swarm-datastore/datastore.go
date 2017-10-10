package swarmdatastore

import (
	"bytes"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/swarm/api/client"
)

// Datastore struct
type Datastore struct {
	serverURL string
	client    *client.Client
}

// NewDatastore returns a new swarm Datastore
func NewDatastore(serverURL string) (datastore Datastore, err error) {
	datastore = Datastore{
		serverURL: serverURL,
		client:    client.NewClient(serverURL),
	}
	return datastore, err
}

// Put stores the given value and return its hash
func (ds *Datastore) Put(value []byte) (hash string, err error) {
	return ds.client.UploadRaw(bytes.NewReader(value), int64(len(value)))
}

// Get returns the value for given key
func (ds *Datastore) Get(key string) (value interface{}, err error) {
	res, err := ds.client.DownloadRaw(key)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	return ioutil.ReadAll(res)
}
