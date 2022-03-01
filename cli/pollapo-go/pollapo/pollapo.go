package pollapo

import (
	"fmt"
	"regexp"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/yaml"
)

type PollapoConfig struct {
	Deps []string
	// root PollapoRoot
}

type PollapoRoot struct {
	// lock
	// replace file option
}

func ParsePollapo(barr []byte) PollapoConfig {
	cfg := PollapoConfig{}
	err := yaml.Unmarshal([]byte(barr), &cfg)
	if err != nil {
		log.Fatalw("Failed to unmarshal yaml", err.Error(), "yaml", barr)
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
		dep := PollapoDep{
			Owner: groups["owner"],
			Repo:  groups["repo"],
			Ref:   groups["rev"],
		}
		if dep.Owner == "" || dep.Repo == "" || dep.Ref == "" {
			return PollapoDep{}, false
		}
		return dep, true
	}
}
