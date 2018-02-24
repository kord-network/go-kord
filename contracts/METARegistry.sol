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

pragma solidity ^0.4.0;

// The META registry contract.
contract METARegistry {
    mapping(address=>bytes32) graphs;

    function graph(address metaID) constant returns (bytes32) {
        return graphs[metaID];
    }

    // ref: https://gist.github.com/axic/5b33912c6f61ae6fd96d6c4a47afde6d
    //
    // TODO: prevent replay attacks with sha3(hash || nonce)
    function setGraph(bytes32 hash, bytes sig) {
        uint8 v;
        bytes32 r;
        bytes32 s;

        if (sig.length != 65) throw;

        assembly {
            r := mload(add(sig, 32))
            s := mload(add(sig, 64))
            v := byte(0, mload(add(sig, 96)))
        }

        if (v < 27) v += 27;

        if (v != 27 && v != 28) throw;

        address metaID = ecrecover(hash, v, r, s);

        graphs[metaID] = hash;
    }
}
