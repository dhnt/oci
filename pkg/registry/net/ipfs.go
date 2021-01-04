package net

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	api "github.com/ipfs/go-ipfs-api"
	files "github.com/ipfs/go-ipfs-files"
)

// IpfsClient is the IPFS client
type IpfsClient struct {
	client *api.Shell
}

type object struct {
	Hash string
}

// AddImage adds components of an image recursively
func (client *IpfsClient) AddImage(manifest map[string][]byte, layers map[string][]byte) (string, error) {
	mf := make(map[string]files.Node)
	for k, v := range manifest {
		mf[k] = files.NewBytesFile(v)
	}

	bf := make(map[string]files.Node)
	for k, v := range layers {
		bf[k] = files.NewBytesFile(v)
	}

	sf := files.NewMapDirectory(map[string]files.Node{
		"blobs":     files.NewMapDirectory(bf),
		"manifests": files.NewMapDirectory(mf),
	})
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("image", sf)})

	reader := files.NewMultiFileReader(slf, true)
	resp, err := client.client.Request("add").
		Option("recursive", true).
		Option("cid-version", 1).
		Body(reader).
		Send(context.Background())
	if err != nil {
		return "", err
	}

	defer resp.Close()

	if resp.Error != nil {
		return "", resp.Error
	}

	dec := json.NewDecoder(resp.Output)
	var final string
	for {
		var out object
		err = dec.Decode(&out)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		final = out.Hash
	}

	if final == "" {
		return "", errors.New("no results received")
	}

	return final, nil
}

// Cat the content at the given path. Callers need to drain and close the returned reader after usage.
func (client *IpfsClient) Cat(path string) (io.ReadCloser, error) {
	return client.client.Cat(path)
}

// List entries at the given path
func (client *IpfsClient) List(path string) ([]*api.LsLink, error) {
	return client.client.List(path)
}

func NewIpfsClient(host string) *IpfsClient {
	client := api.NewShell(host)
	return &IpfsClient{
		client: client,
	}
}
