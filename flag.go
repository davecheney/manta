// Copyright 2013-2014 David Cheney and Contributors.
// All rights reserved. Use of this source code is governed by a
// BSD-style license that can be found in the LICENSE file.

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
