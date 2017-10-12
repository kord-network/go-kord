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

package identity

import (
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestAPI tests creating and updating identities via the API.
func TestAPI(t *testing.T) {
	// start the identity API
	store := NewMemoryStore()
	api := NewAPI(store)
	srv := httptest.NewServer(api)
	defer srv.Close()

	// create an identity
	client := NewClient(srv.URL)
	name := "lmars"
	identity, err := client.CreateIdentity(name)
	if err != nil {
		t.Fatal(err)
	}

	// check the identity can be loaded from the API
	gotIdentity, err := client.GetIdentity(name)
	if err != nil {
		t.Fatal(err)
	}
	if gotIdentity.Name != identity.Name {
		t.Fatalf("expected identity to have name %q, got %q", identity.Name, gotIdentity.Name)
	}

	// update the identity
	updates := map[string]string{"foo": "bar"}
	if _, err := client.UpdateIdentity(name, updates); err != nil {
		t.Fatal(err)
	}

	// check that the updated identity can be loaded from the API
	gotIdentity, err = client.GetIdentity(name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotIdentity.Aux, updates) {
		t.Fatalf("expected identity to have aux %v, got %v", updates, gotIdentity.Aux)
	}
}
