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
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	swarm "github.com/ethereum/go-ethereum/swarm/api/client"
	"github.com/meta-network/go-meta/testutil"
)

// TestCLISwarmServer tests uploading and downloading files via the Swarm API.
func TestCLISwarmServer(t *testing.T) {
	// start the API
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	srv, err := NewServer(store, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	httpSrv := httptest.NewServer(srv)
	defer httpSrv.Close()

	// upload a file
	client := swarm.NewClient(httpSrv.URL)
	data := []byte("some-data")
	hash, err := client.UploadRaw(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}

	// download the file
	res, err := client.DownloadRaw(hash)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Close()
	gotData, err := ioutil.ReadAll(res)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, gotData) {
		t.Fatalf("unexpected data, expected %q but got %q", data, gotData)
	}
}
