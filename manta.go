package manta

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	user := os.Getenv("MANTA_USER")
	if user == "" {
		log.Fatal("manta: MANTA_USER not defined or empty")
	}
	keyid := os.Getenv("MANTA_KEY_ID")
	if keyid == "" {
		log.Fatal("manta: MANTA_KEY_ID not defined or empty")
	}
	url := os.Getenv("MANTA_URL")
	if url == "" {
		log.Fatal("manta: MANTA_URL not defined or empty")
	}
	return &Client{
		User:  user,
		KeyId: keyid,
		Key:   filepath.Join(mustHomedir(), ".ssh", "id_rsa"),
		Url:   url,
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
		c.signer, err = LoadPrivateKey(c.Key)
		if err != nil {
			return fmt.Errorf("could not load private key %q: %v", c.Key, err)
		}
	}
	return SignRequest(req, c.User, c.KeyId, c.signer)
}

func SignRequest(req *http.Request, MANTA_USER, MANTA_KEY_ID string, priv Signer) error {
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
