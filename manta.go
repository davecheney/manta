// Copyright 2013-2014 David Cheney and Contributors.
// All rights reserved. Use of this source code is governed by a
// BSD-style license that can be found in the LICENSE file.

////////////////////////////////////////////////////////////
// Manta implements a client for the Joyent Manta API.
// http://apidocs.joyent.com/manta/index.html.
//
// Included in the package is an incomplete implementation of the
// CLI Utilities.
// http://apidocs.joyent.com/manta/commands-reference.html
package manta

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/user"
	"path/filepath"
	"time"
)

// Client is a Manta client. Client is not safe for concurrent use.
type Client struct {
	User   string
	KeyId  string
	Key    string
	Url    string
	signer Signer
}

func homeDir() (dir string, err error) {
	u, err := user.Current()
	if err != nil {
		return dir, fmt.Errorf("manta: could not determine home directory: %v", err)
	}
	dir = u.HomeDir
	return
}

// DefaultClient returns a Client instance configured from the
// default Manta environment variables.
func DefaultClient() (c *Client, err error) {
	dir, err := homeDir()
	if err != nil {
		return
	}
	c = &Client{
		User:  MANTA_USER,
		KeyId: MANTA_KEY_ID,
		Key:   filepath.Join(dir, ".ssh", "id_rsa"),
		Url:   MANTA_URL,
	}
	return
}

// NewRequest is similar to http.NewRequest except it appends path to
// the API endpoint this client is configured for.
func (c *Client) NewRequest(method, path string, r io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.Url, path)
	return http.NewRequest(method, url, r)
}

// SignRequest signs the 'date' field of req.
func (c *Client) SignRequest(req *http.Request) error {
	if c.signer == nil {
		var err error
		c.signer, err = loadPrivateKey(c.Key)
		if err != nil {
			return fmt.Errorf("could not load private key %q: %v", c.Key, err)
		}
	}
	return signRequest(req, fmt.Sprintf("/%s/keys/%s", MANTA_USER, MANTA_KEY_ID), c.signer)
}

// Get executes a GET request and returns the response.
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.Do("DELETE", path, nil)
}

// Get executes a GET request and returns the response.
func (c *Client) Get(path string) (*http.Response, error) {
	return c.Do("GET", path, nil)
}

// Put executes a PUT request and returns the response.
func (c *Client) Put(path string, r io.Reader) (*http.Response, error) {
	return c.Do("PUT", path, r)
}

// Do executes a method request and returns the response.
func (c *Client) Do(method, path string, r io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(method, path, r)
	if err != nil {
		return nil, err
	}
	if err := c.SignRequest(req); err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func signRequest(req *http.Request, keyid string, priv Signer) error {
	now := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", now)
	signed, err := priv.Sign([]byte(fmt.Sprintf("date: %s", now)))
	if err != nil {
		return fmt.Errorf("could not sign request: %v", err)
	}
	sig := base64.StdEncoding.EncodeToString(signed)
	authz := fmt.Sprintf("Signature keyId=%q,algorithm=%q,signature=%q", keyid, "rsa-sha256", sig)
	req.Header.Set("Authorization", authz)
	return nil
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPrivateKey(path string) (Signer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(data)
}

// parsePublicKey parses a PEM encoded private key.
func parsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type rsaPublicKey rsa.PublicKey

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}
