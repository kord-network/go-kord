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

package cwr

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
)

// ParseCWRFile parse a give cwr file and returns an array of registeredWorks
func ParseCWRFile(cwrFileReader io.Reader, CWRDataApiPath string) (registeredWorks []*RegisteredWork, err error) {
	//phase 1 : Transform cwr formatted file to cwr-json file using CWR-DataApi python script.
	cmd := exec.Command("python3", CWRDataApiPath+"/cwr2json.py")
	cmd.Stdin = cwrFileReader

	var pythonStdout, pythonErrbuf bytes.Buffer
	cmd.Stdout = &pythonStdout
	cmd.Stderr = &pythonErrbuf

	err = cmd.Start()
	if err != nil {
		return nil, errors.New(pythonErrbuf.String())
	}
	err = cmd.Wait()
	if err != nil {
		return nil, errors.New(pythonErrbuf.String())
	}
	var cwr *Cwr
	if err := json.NewDecoder(bytes.NewReader(pythonStdout.Bytes())).Decode(&cwr); err != nil {
		return nil, err
	}
	for _, group := range cwr.Transmission.Groups {
		for _, tx := range group.Transactions {
			var registeredWork *RegisteredWork
			for _, record := range tx {
				if record.RecordType == "REV" ||
					record.RecordType == "NWR" {
					registeredWorkBytes, err := json.Marshal(record)
					if err != nil {
						return nil, err
					}
					if err := json.NewDecoder(bytes.NewReader(registeredWorkBytes)).Decode(&registeredWork); err != nil {
						return nil, err
					}
					registeredWorks = append(registeredWorks, registeredWork)
				}
			}
		}
	}
	return
}
