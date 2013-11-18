package manta

// shared flagset support

import (
	"flag"
	"os"
)

var (
	Flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	MANTA_USER string
)

func init() {
	Flags.StringVar(&MANTA_USER, "a", os.Getenv("MANTA_USER"), "Authenticate as account. Defaults to $MANTA_USER")
}
