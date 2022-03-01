package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteTokenGhHosts(t *testing.T) {
	barr := WriteTokenGhHosts("token")
	t.Logf("yaml: %s", string(barr))
	assert.Greater(t, len(barr), 10)
}
func TestGetTokenFromGhHosts(t *testing.T) {
	token := GetTokenFromGhHosts()
	t.Logf("token: %s", token)
	assert.Greater(t, len(token), 10)
}
