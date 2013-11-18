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
	"log"
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

func mustHomedir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal("manta: could not determine home directory: %v", err)
	}
	return user.HomeDir
}

// DefaultClient returns a Client instance configured from the
// default Manta environment variables.
func DefaultClient() *Client {
	return &Client{
		User:  MANTA_USER,
		KeyId: MANTA_KEY_ID,
		Key:   filepath.Join(mustHomedir(), ".ssh", "id_rsa"),
		Url:   MANTA_URL,
	}
}

// NewRequest is similar to http.NewRequest except it appends path to
// the API endpoint this client is configured for.
func (c *Client) NewRequest(method, path string, r io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.Url, path)
	return http.NewRequest(method, url, r)
}

func (c *Client) SignRequest(req *http.Request) error {
	if c.signer == nil {
		var err error
		c.signer, err = loadPrivateKey(c.Key)
		if err != nil {
			return fmt.Errorf("could not load private key %q: %v", c.Key, err)
		}
	}
	return signRequest(req, c.User, c.KeyId, c.signer)
}

func signRequest(req *http.Request, MANTA_USER, MANTA_KEY_ID string, priv Signer) error {
	now := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", now)
	signed, err := priv.Sign([]byte(fmt.Sprintf("date: %s", now)))
	if err != nil {
		return fmt.Errorf("could not sign request: %v", err)
	}
	sig := base64.StdEncoding.EncodeToString(signed)
	authz := fmt.Sprintf("Signature keyId=%q,algorithm=%q,signature=%q", fmt.Sprintf("/%s/keys/%s", MANTA_USER, MANTA_KEY_ID), "rsa-sha256", sig)
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
	// PublicKey returns an associated PublicKey instance.
	PublicKey() PublicKey

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
