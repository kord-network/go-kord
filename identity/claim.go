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
	"github.com/ethereum/go-ethereum/crypto"
)

// Claim structure
type Claim struct {
	Issuer    string `json:"issuer"`
	Subject    string `json:"subject"`
	Claim     string `json:"claim"`
	Signature string `json:"signature"`
	ID        string `json:"id"`
}

// NewClaim create and returns new Claim.
func NewClaim(issuer string, subject string, claim string, signature string) *Claim {
	return &Claim{
		Issuer:    issuer,
		Subject:    subject,
		Claim:     claim,
		Signature: signature,
		ID:        crypto.Keccak256Hash([]byte(issuer), []byte(subject), []byte(claim), []byte(signature)).String(),
	}
}
