package manta

import (
	"flag"
	"testing"
)

var SHARED_FLAGS int = 3
var count int

func tallyFlag(f *flag.Flag) {
	count += 1
}

func TestDefaultFlags(t *testing.T) {
	flags := Flags()
	flags.VisitAll(tallyFlag)
	if count != SHARED_FLAGS {
		t.Errorf("Expected default flag count %i, found %i", SHARED_FLAGS, count)
	}
}
