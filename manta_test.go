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
