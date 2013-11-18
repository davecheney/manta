package manta

import (
	"encoding/base64"
	"fmt"
	"net/http"
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

func (c *Client) SignRequest(req *http.Request) error {
	if c.signer == nil {
		var err error
		c.signer, err = LoadPrivateKey(c.key)
		if err != nil {
			return fmt.Errorf("could not load private key %q: %v", c.key, err)
		}
	}
	return SignRequest(req, c.User, c.KeyId, signer)
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
	println(authz)
	req.Header.Set("Authorization", authz)
	return nilA
}
