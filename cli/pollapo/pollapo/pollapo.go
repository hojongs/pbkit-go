package pollapo

import (
	"regexp"

	"github.com/hojongs/pbkit-go/cli/pollapo/log"
	"gopkg.in/yaml.v3"
)

type PollapoConfig struct {
	Deps []string
	root PollapoRoot
}

type PollapoRoot struct {
	// lock
	// replace file option
}

func ParsePollapo(bytes []byte) PollapoConfig {
	cfg := PollapoConfig{}
	err := yaml.Unmarshal([]byte(bytes), &cfg)
	if err != nil {
		log.Fatalw("Failed to unmarshal yaml", err.Error(), "yaml", bytes)
	}
	return cfg
}

type PollapoDep struct {
	owner string
	repo  string
	rev   string
}

// func deps(pollapoYml PollapoYml) []PollapoDep {
// 	rtv := []PollapoDep{}
// 	for _, dep := range pollapoYml.Deps {
// 		rtv = append(rtv, ParseDep(dep))
// 	}
// 	return rtv
// }
func ParseDep(dep string) (PollapoDep, bool) {
	r := regexp.MustCompile(`(?P<owner>.+?)\/(?P<repo>.+?)@(?P<rev>.+)`)
	matches := r.FindStringSubmatch(dep)
	if matches == nil {
		return PollapoDep{}, false
	} else {
		groups := map[string]string{}
		for i, name := range r.SubexpNames()[1:] {
			groups[name] = matches[i+1]
		}
		return PollapoDep{
			owner: groups["owner"],
			repo:  groups["repo"],
			rev:   groups["rev"],
		}, true
	}
}

// func getPollapoYml(dep PollapoDep, cacheDir string) PollapoYml {
// 	panic("A")
// 	// return LoadPollapoYml(getYmlPath(cacheDir, dep))
// }
// func LoadPollapoYml(ymlPath string) PollapoYml {
// 	return parseYaml(string(pollapoYmlText))
// }

// func parseYaml(pollapoYmlText string) PollapoYml {
// 	pollapoYml := PollapoYml{}
// 	err := yaml.Unmarshal([]byte(pollapoYmlText), &pollapoYml)
// 	if err != nil {
// 		// TODO: PollapoYmlMalformedError(ymlPath)
// 		panic(err)
// 	}
// 	return pollapoYml
// }

// type AnalyzeDepsResultRev struct {
// 	froms []string
// }

// func AnalyzeDeps(
// 	cacheDir string,
// 	pollapoYml PollapoYml,
// ) map[string]map[string]AnalyzeDepsResultRev {
// 	type Dep struct {
// 		PollapoDep
// 		from string
// 	}
// 	result := map[string]map[string]AnalyzeDepsResultRev{}
// 	//   const lockTable = pollapoYml?.root?.lock ?? {};

// 	// parse dependencies from pollapoYml and insert it into queue
// 	temp := []PollapoDep{}
// 	temp = append(temp, deps(pollapoYml)...)
// 	queue := []Dep{} // queue of deps
// 	for _, dep := range temp {
// 		queue = append(queue, Dep{PollapoDep: dep, from: "<root>"})
// 	}

// 	for len(queue) > 0 {
// 		// pop dep from the queue
// 		dep := queue[0]
// 		queue = queue[1:]

// 		// Add dep.from into the revision of result
// 		repoPath := dep.user + "/" + dep.repo
// 		result[repoPath][dep.rev] = AnalyzeDepsResultRev{
// 			froms: append(result[repoPath][dep.rev].froms, dep.from),
// 		}
// 		// Caution! getPollapoYml requires cacheDeps() before run
// 	}
// 	return nil
// }
