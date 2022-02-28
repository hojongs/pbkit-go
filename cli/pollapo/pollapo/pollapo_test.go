package pollapo

import (
	// "reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestParseYaml(t *testing.T) {
// 	pollapoYmlText := `
// deps:
//   - googleapis/googleapis@dbfbfdb
// `
// 	want := PollapoYml{
// 		Deps: []string{"googleapis/googleapis@dbfbfdb"},
// 	}
// 	t.Log(pollapoYmlText)
// 	pollapoYml := parseYaml(pollapoYmlText)
// 	t.Log(pollapoYml)
// 	t.Log(want)
// 	if !reflect.DeepEqual(pollapoYml, want) {
// 		t.Fatalf("parseYaml()")
// 	}
// }

func TestParseDep(t *testing.T) {
	rtv, _ := ParseDep("google/apis@dbfbfdb")
	assert.Equal(t, "google", rtv.owner)
	assert.Equal(t, "apis", rtv.repo)
	assert.Equal(t, "dbfbfdb", rtv.rev)
}
