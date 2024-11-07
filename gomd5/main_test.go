package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintMD5(t *testing.T) {
	in := strings.NewReader("go")
	out := &bytes.Buffer{}
	printMD5(in, out)
	want := "34d1f91fb2e514b8576fab1a75a89a6b"
	got := out.String()
	if got != want {
		t.Errorf("got wrong md5 hash")
	}
}
