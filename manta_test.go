package manta

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestParsePrivateKey(t *testing.T) {
	data, err := ioutil.ReadFile("_testdata/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	_, err = parsePrivateKey(data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadPrivateKey(t *testing.T) {
	if _, err := loadPrivateKey("_testdata/id_rsa"); err != nil {
		t.Error(err)
	}
}

func TestPrivateKeySign(t *testing.T) {
	priv, err := loadPrivateKey("_testdata/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	sig, err := priv.Sign([]byte("date: Thu, 05 Jan 2012 21:31:40 GMT"))
	if err != nil {
		t.Fatal(err)
	}
	const want = "ATp0r26dbMIxOopqw0OfABDT7CKMIoENumuruOtarj8n/97Q3htHFYpH8yOSQk3Z5zh8UxUym6FYTb5+A0Nz3NRsXJibnYi7brE/4tx5But9kkFGzG+xpUmimN4c3TMN7OFH//+r8hBf7BT9/GmHDUVZT2JzWGLZES2xDOUuMtA="
	if got := base64.StdEncoding.EncodeToString(sig); got != want {
		t.Fatalf("want: %q, got %q", want, got)
	}
}

func TestSignRequest(t *testing.T) {
	priv, err := loadPrivateKey("_testdata/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	signRequest(req, "Test", priv)
	date := req.Header.Get("date")
	sig, err := priv.Sign([]byte("date: " + date))
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("Signature keyId=%q,algorithm=%q,signature=%q", "Test", "rsa-sha256", base64.StdEncoding.EncodeToString(sig))
	if got := req.Header.Get("Authorization"); got != want {
		t.Fatalf("want: %q, got: %q", want, got)
	}
}

func TestClientNewRequest(t *testing.T) {
	client := Client{
		User:  "test",
		KeyId: "q",
		Key:   "_testdata/id_rsa",
		Url:   "http://example.com",
	}
	req, err := client.NewRequest("GET", "/test/public", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := req.URL.Host, "example.com"; got != want {
		t.Errorf("want: %q, got: %q", want, got)
	}
	if got, want := req.URL.Path, "/test/public"; got != want {
		t.Errorf("want: %q, got: %q", want, got)
	}
}

func TestClientSignRequest(t *testing.T) {
	client := Client{
		User:  "test",
		KeyId: "q",
		Key:   "_testdata/id_rsa",
		Url:   "http://example.com",
	}
	req, err := client.NewRequest("GET", "/test/public", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SignRequest(req); err != nil {
		t.Fatal(err)
	}
}
