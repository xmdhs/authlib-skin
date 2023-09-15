package utils

import (
	"testing"
)

func TestUUIDGen(t *testing.T) {
	if UUIDGen("xmdhs") != "6560e064bcfc32baa5fa2aa8831f1298" {
		t.Fail()
	}
}
