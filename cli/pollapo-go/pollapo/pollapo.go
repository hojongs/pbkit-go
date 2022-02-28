package pollapo

import (
	"fmt"
	"regexp"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
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
	Owner string
	Repo  string
	Ref   string
}

func (dep PollapoDep) String() string {
	return fmt.Sprintf("%s/%s@%s", dep.Owner, dep.Repo, dep.Ref)
}

func ParseDep(depTxt string) (PollapoDep, bool) {
	r := regexp.MustCompile(`(?P<owner>.+?)\/(?P<repo>.+?)@(?P<rev>.+)`)
	matches := r.FindStringSubmatch(depTxt)
	if matches == nil {
		return PollapoDep{}, false
	} else {
		groups := map[string]string{}
		for i, name := range r.SubexpNames()[1:] {
			groups[name] = matches[i+1]
		}
		return PollapoDep{
			Owner: groups["owner"],
			Repo:  groups["repo"],
			Ref:   groups["rev"],
		}, true
	}
}
