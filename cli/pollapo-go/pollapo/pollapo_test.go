package pollapo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYaml(t *testing.T) {
	pollapoYmlText := `
deps:
  - googleapis/googleapis@dbfbfdb
`
	want := PollapoConfig{
		Deps: []string{"googleapis/googleapis@dbfbfdb"},
	}
	t.Log(pollapoYmlText)
	pollapoYml := ParsePollapo([]byte(pollapoYmlText))
	t.Log(pollapoYml)
	t.Log(want)
	if !reflect.DeepEqual(pollapoYml, want) {
		t.Fatalf("parseYaml()")
	}
}

func TestParseDep(t *testing.T) {
	rtv, _ := ParseDep("google/apis@dbfbfdb")
	assert.Equal(t, "google", rtv.Owner)
	assert.Equal(t, "apis", rtv.Repo)
	assert.Equal(t, "dbfbfdb", rtv.Ref)
}
