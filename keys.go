// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manta

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

// LoadPrivateKey loads an parses a PEM encoded private key file.
func LoadPrivateKey(path string) (Signer, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParsePrivateKey(data)
}

// ParsePublicKey parses a PEM encoded private key.
func ParsePrivateKey(pemBytes []byte) (Signer, error) {
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
	return NewSignerFromKey(rawkey)
}

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// PublicKey returns an associated PublicKey instance.
	PublicKey() PublicKey

	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

// NewPrivateKey takes a pointer to rsa, dsa or ecdsa PrivateKey
// returns a corresponding Signer instance. EC keys should use P256,
// P384 or P521.
func NewSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

// PublicKey is an abstraction of different types of public keys.
type PublicKey interface {
	// PrivateKeyAlgo returns the name of the encryption system.
	PrivateKeyAlgo() string

	// PublicKeyAlgo returns the algorithm for the public key,
	// which may be different from PrivateKeyAlgo for certificates.
	PublicKeyAlgo() string
}

type rsaPublicKey rsa.PublicKey

func (r *rsaPublicKey) PrivateKeyAlgo() string {
	return "ssh-rsa"
}

func (r *rsaPublicKey) PublicKeyAlgo() string {
	return r.PrivateKeyAlgo()
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

func (r *rsaPrivateKey) PublicKey() PublicKey {
	return (*rsaPublicKey)(&r.PrivateKey.PublicKey)
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}
