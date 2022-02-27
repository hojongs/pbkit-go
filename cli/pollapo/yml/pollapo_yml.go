package yml

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PollapoYml struct {
	Deps []string
	root PollapoRoot
}

type PollapoRoot struct {
	// lock
	// replace file option
}

func LoadPollapoYml(ymlPath string) PollapoYml {
	pollapoYmlText, err := os.ReadFile(ymlPath)
	if err != nil {
		// TODO: PollapoYmlNotFoundError(ymlPath)
		panic(err)
	}
	return parseYaml(string(pollapoYmlText))
}

func parseYaml(pollapoYmlText string) PollapoYml {
	pollapoYml := PollapoYml{}
	err := yaml.Unmarshal([]byte(pollapoYmlText), &pollapoYml)
	if err != nil {
		// TODO: PollapoYmlMalformedError(ymlPath)
		panic(err)
	}
	return pollapoYml
}
