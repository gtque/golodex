package main

import (
	"internal/golodexdata"
	"testing"
)

func TestTokenPassword(t *testing.T) {
	token := "sample"
	password := golodexdata.GetSessionPassword(token)
	if password != token {
		t.Fatalf("%q != %q", token, password)
	}
}
