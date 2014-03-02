package manta

// shared flagset support

import (
	"flag"
	"os"
)

var (
	MANTA_USER   string
	MANTA_URL    string
	MANTA_KEY_ID string
)

// Flags returns a flag.FlagSet containing the shared flags required for DefaultClient.
func Flags() *flag.FlagSet {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&MANTA_USER, "a", os.Getenv("MANTA_USER"), "Authenticate as account. Defaults to $MANTA_USER")
	flags.StringVar(&MANTA_URL, "u", "https://us-east.manta.joyent.com", "Manta base URL.")
	flags.StringVar(&MANTA_KEY_ID, "k", os.Getenv("MANTA_KEY_ID"), "Authenticate using the SSH key described by FINGERPRINT. Defaults to $MANTA_KEY_ID")
	return flags
}
