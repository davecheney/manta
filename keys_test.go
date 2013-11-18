package manta

import (
	"encoding/base64"
	"io/ioutil"
	"testing"
)

func TestParsePrivateKey(t *testing.T) {
	data, err := ioutil.ReadFile("_testdata/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParsePrivateKey(data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadPrivateKey(t *testing.T) {
	if _, err := LoadPrivateKey("_testdata/id_rsa"); err != nil {
		t.Error(err)
	}
}

func TestPrivateKeySign(t *testing.T) {
	priv, err := LoadPrivateKey("_testdata/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	sig, err := priv.Sign([]byte("date: Thu, 05 Jan 2012 21:31:40 GMT"))
	if err != nil {
		t.Fatal(err)
	}
	const want = "JldXnt8W9t643M2Sce10gqCh/+E7QIYLiI+bSjnFBGCti7s+mPPvOjVb72sbd1FjeOUwPTDpKbrQQORrm+xBYfAwCxF3LBSSzORvyJ5nRFCFxfJ3nlQD6Kdxhw8wrVZX5nSem4A/W3C8qH5uhFTRwF4ruRjh+ENHWuovPgO/HGQ="
	if got := base64.StdEncoding.EncodeToString(sig); got != want {
		t.Fatalf("want: %q, got %q", want, got)
	}
}

func b64decode(t *testing.T, b64 string) []byte {
	d, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
