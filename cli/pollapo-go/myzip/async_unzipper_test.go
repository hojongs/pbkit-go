package myzip

import (
	"os"
	"testing"
)

func TestUnzip(t *testing.T) {
	barr, err := os.ReadFile("temp.zip")
	zipReader := NewZipReader(barr)

	if err != nil {
		t.Fatal(err)
	}
	SyncUnzipper{}.Unzip(zipReader, "pollapo-test")
}
