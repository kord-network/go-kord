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
	"path/filepath"
	"runtime"
	"strings"
)

// ParseCWRFile parse a give cwr file and returns an array of registeredWorks
func ParseCWRFile(cwrFileReader io.Reader) (registeredWorks []*RegisteredWork, err error) {
	//phase 1 : Transform cwr formatted file to cwr-json file using CWR-DataApi python script.
	//get the absolute path the cwrfile.
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("error getting base path")
	}
	basePath := filepath.Dir(b)
	cwr2jsonpy := basePath + "/CWR-DataApi/cwr2json.py"

	cmd := exec.Command("python3", cwr2jsonpy)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		buf := new(bytes.Buffer)
		// file, _ := os.Open(cwrFilePath)
		buf.ReadFrom(cwrFileReader)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			line = line + "\n"
			io.WriteString(stdin, line)
		}
	}()
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(outbuf.Bytes(), &data); err != nil {
		return nil, err
	}
	for _, group := range data["transmission"].(map[string]interface{})["groups"].([]interface{}) {
		for _, tx := range group.(map[string]interface{})["transactions"].([]interface{}) {
			var registeredWork *RegisteredWork

			for _, record := range tx.([]interface{}) {

				if record.(map[string]interface{})["record_type"] == "REV" ||
					record.(map[string]interface{})["record_type"] == "NWR" {

					registeredWorkBytes, err := json.Marshal(record)
					if err != nil {
						return nil, err
					}
					dec := json.NewDecoder(strings.NewReader(string(registeredWorkBytes)))
					err = dec.Decode(&registeredWork)
					if err != nil {
						return nil, err
					}
					registeredWorks = append(registeredWorks, registeredWork)
				}
			}
		}
	}
	return
}
