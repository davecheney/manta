package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

var (
	MANTA_URL    string
	MANTA_USER   string
	MANTA_KEY_ID string
	MANTA_KEY    string
)

func mustUser() *user.User {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func init() {
	flag.StringVar(&MANTA_URL, "url", "https://us-east.manta.joyent.com", "Manta API endpoint")
	flag.StringVar(&MANTA_USER, "user", os.Getenv("MANTA_USER"), "Your Joyent Cloud account login name")
	flag.StringVar(&MANTA_KEY_ID, "keyid", os.Getenv("MANTA_KEY_ID"), "The fingerprint of your SSH key")
	flag.StringVar(&MANTA_KEY, "key", filepath.Join(mustUser().HomeDir, ".ssh", "id_rsa"), "The fingerprint of your SSH key")

	flag.Parse()
}
