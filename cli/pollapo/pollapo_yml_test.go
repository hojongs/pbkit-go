package pollapo

import (
	"reflect"
	"testing"
)

func TestParseYaml(t *testing.T) {
	pollapoYmlText := `
deps:
  - googleapis/googleapis@dbfbfdb
`
	want := PollapoYml{
		Deps: []string{"googleapis/googleapis@dbfbfdb"},
	}
	t.Log(pollapoYmlText)
	pollapoYml := parseYaml(pollapoYmlText)
	t.Log(pollapoYml)
	t.Log(want)
	if !reflect.DeepEqual(pollapoYml, want) {
		t.Fatalf("parseYaml()")
	}
}
