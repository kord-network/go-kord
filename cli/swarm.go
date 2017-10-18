package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/swarm/api/client"
)

type SwarmBackend struct {
	indexDirHash string
	api          *client.Client
	indexDir     string
}

func (bzz *SwarmBackend) OpenIndex(bzzuri string, indexDirHash string) (err error) {
	if indexDirHash != "" {
		bzz.indexDirHash = indexDirHash
	} else {
		bzz.indexDirHash = bzz.getIndexManifestHash()
	}
	bzz.indexDir, err = ioutil.TempDir("", "meta-index-proxy")
	if err != nil {
		return err
	}
	if bzzuri != "" {
		bzz.api = client.NewClient(bzzuri)
	} else {
		bzz.api = client.DefaultClient
	}
	_, err = bzz.api.DownloadManifest(bzz.indexDirHash)
	if err != nil {
		return err
	}
	return nil
}

// TODO: will match if file name start is unique, and will not return pathname due to bzz bugs
func (bzz *SwarmBackend) GetIndexFile(path string, mustfind bool) (string, error) {
	olddir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	err = os.Chdir(bzz.indexDir)
	if err != nil {
		return "", err
	}
	fw, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	defer fw.Close()
	fr, err := bzz.api.Download(bzz.indexDirHash, path)
	if err == nil {
		defer fr.Close()
		_, err = io.Copy(fw, fr)
		if err != nil {
			return "", err
		}
	} else if mustfind {
		return "", err
	}
	err = os.Chdir(olddir)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{bzz.indexDir, path}, "/"), nil
}

// WARNING: if the file was initially retrieved from bzz, there is no guarantee it wasn't changed in the meantime
func (bzz *SwarmBackend) PutIndexFile(filename string) (string, error) {
	olddir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Getwd failed before uploading file '%s', panicking to avoid file deletion: %v", filename, err))
	}
	err = os.Chdir(bzz.indexDir)
	if err != nil {
		panic(fmt.Sprintf("could not change to database dir to upload file '%s', panicking to avoid file deletion: %v", filename, err))
	}
	bzzfile, err := client.Open(filename)
	if err != nil {
		return "", err
	}
	bzzfile.ContentType = "application/x-sqlite3"
	bzzfile.Path = filename
	log.Info("bzzfile", "file", bzzfile, "filename", filename)
	mhash, err := bzz.api.Upload(bzzfile, bzz.indexDirHash)
	if err != nil {
		panic(fmt.Sprintf("could not upload file '%s', panicking to avoid file deletion: %v", filename, err))
	}
	err = os.Chdir(olddir)
	if err != nil {
		return "", err
	}
	txr, err := bzz.updateENS(mhash)
	return txr, err
}

// TODO: replace with ens resolve
func (bzz *SwarmBackend) getIndexManifestHash() string {
	return "71e7f9ef19240b0716cf0efdece1c1dc64d478c132f749482fa56ca1a9a67c71"
}

func (bzz *SwarmBackend) updateENS(hash string) (string, error) {
	log.Info("ens updater not implemented. please note hash manually", "hash", hash)
	return hash, nil
}
