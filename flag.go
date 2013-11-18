package manta

// shared flagset support

import (
	"flag"
	"os"
)

var (
	Flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	MANTA_USER   string
	MANTA_URL    string
	MANTA_KEY_ID string
)

func init() {
	Flags.StringVar(&MANTA_USER, "a", os.Getenv("MANTA_USER"), "Authenticate as account. Defaults to $MANTA_USER")
	Flags.StringVar(&MANTA_URL, "u", "https://us-east.manta.joyent.com", "Manta base URL.")
	Flags.StringVar(&MANTA_KEY_ID, "k", os.Getenv("MANTA_KEY_ID"), "Authenticate using the SSH key described by FINGERPRINT. Defaults to $MANTA_KEY_ID")
}
