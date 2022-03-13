package pollapo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/yaml"
)

type PollapoConfig struct {
	deps []string
	root PollapoRoot
}

type PollapoRoot struct {
	lock map[string]string
	// replace file option
}

func ParsePollapo(barr []byte) PollapoConfig {
	cfg := PollapoConfig{}
	err := yaml.Unmarshal([]byte(barr), &cfg)
	if err != nil {
		log.Sugar.Fatalw("Failed to unmarshal yaml", err.Error(), "yaml", barr)
	}
	return cfg
}

// Get parsed deps
// resolve deps as commit hashes if there are corresponding lock in root.lock
// TODO: impl store root.lock
func (cfg PollapoConfig) GetDeps() []PollapoDep {
	lockedDeps := []PollapoDep{}
	for _, depTxt := range cfg.deps {
		dep, isOk := ParseDep(depTxt)
		if !isOk {
			log.Sugar.Fatalw("Invalid dep", nil, "dep", depTxt)
		}
		hasLock := false
		for k, v := range cfg.root.lock {
			if depTxt[:strings.Index(depTxt, "@")] == k {
				dep.Ref = v
				lockedDeps = append(lockedDeps, dep)
				hasLock = true
				break
			}
		}
		if !hasLock {
			lockedDeps = append(lockedDeps, dep)
		}
	}
	return lockedDeps
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
