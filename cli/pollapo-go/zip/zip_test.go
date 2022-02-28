package zip

import (
	"os"
	"testing"
)

func TestUnzip(t *testing.T) {
	barr, err := os.ReadFile("temp.zip")
	if err != nil {
		t.Fatal(err)
	}
	Unzip(barr, "pollapo-test")
}
