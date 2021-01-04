package registry

import (
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	cid "github.com/ipfs/go-cid"
	base58 "github.com/jbenet/go-base58"
	mbase "github.com/multiformats/go-multibase"
)

// Default http client with timeout
// https://golang.org/pkg/net/http/#pkg-examples
// Clients and Transports are safe for concurrent use by multiple goroutines.
var defaultClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: defaultTransport,
}

// https://golang.org/src/net/http/transport.go
var defaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 10 * time.Second,
		DualStack: true,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          10,
	IdleConnTimeout:       10 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// Get issues a GET to the specified URL - a drop-in replacement for http.Get with timeouts.
func Get(url string) (resp *http.Response, err error) {
	return defaultClient.Get(url)
}

func getContent(gw string, cid string, s []string) ([]byte, error) {
	uri := IpfsURL(gw, append([]string{cid}, s...))
	resp, err := Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cid: %s %s", cid, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

// IpfsifyHash does base32 to base58 conversion
func IpfsifyHash(base32Hash string) string {
	decodedB32, err := base32.StdEncoding.DecodeString(strings.ToUpper(base32Hash) + "=")
	if err != nil {
		return ""
	}

	return base58.Encode(decodedB32)
}

// IpfsURL constructs IPFS gateway URL
func IpfsURL(gw string, s []string) string {
	return fmt.Sprintf("%s/ipfs/%s", gw, strings.Join(s, "/"))
}

func toCidV0(c cid.Cid) (cid.Cid, error) {
	if c.Type() != cid.DagProtobuf {
		return cid.Cid{}, fmt.Errorf("can't convert non-protobuf nodes to cidv0")
	}
	return cid.NewCidV0(c.Hash()), nil
}

func toCidV1(c cid.Cid) (cid.Cid, error) {
	return cid.NewCidV1(c.Type(), c.Hash()), nil
}

// ToB58 returns base58 encoded string if s is a valid cid
func ToB58(s string) string {
	c, err := cid.Decode(s)
	if err == nil {
		return c.Hash().B58String()
	}
	return ""
}

// ToB32 returns base32 encoded string if s is a valid cid
func ToB32(s string) string {
	c, err := cid.Decode(s)
	if err != nil {
		return ""
	}
	c1, err := toCidV1(c)
	if err != nil {
		return ""
	}
	b32, err := c1.StringOfBase(mbase.Base32)
	if err != nil {
		return ""
	}
	return b32
}

// Dig interrogates registry server. It performs CID lookups and shows the response.
func Dig(gw string, short bool, name string) (string, error) {
	uri := fmt.Sprintf("http://%s/dig?q=%s&short=%v", gw, name, short)

	resp, err := Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
