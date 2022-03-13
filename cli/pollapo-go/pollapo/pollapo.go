package pollapo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
	"gopkg.in/yaml.v3"
)

type PollapoConfig interface {
	GetDeps(verbose bool) []PollapoDep
}

type PollapoConfigYml struct {
	Deps []string
	Root PollapoRoot
}

type PollapoRoot struct {
	Lock map[string]string
	// replace file option
}

var logName = "Pollapo"

func ParsePollapo(barr []byte) PollapoConfig {
	cfg := PollapoConfigYml{}
	err := yaml.Unmarshal([]byte(barr), &cfg)
	if err != nil {
		util.Sugar.Fatalw("Failed to unmarshal yaml", err.Error(), "yaml", barr)
	}
	return cfg
}

// Get parsed deps
// resolve deps as commit hashes if there are corresponding lock in root.lock
// TODO: impl store root.lock
func (cfg PollapoConfigYml) GetDeps(verbose bool) []PollapoDep {
	util.PrintfVerbose(logName, verbose, "deps: %s\n", util.Yellow(cfg.Deps))
	lockedDeps := []PollapoDep{}
	for _, depTxt := range cfg.Deps {
		dep, isOk := ParseDep(depTxt)
		if !isOk {
			util.Sugar.Fatalw("Invalid dep", nil, "dep", depTxt)
		}
		hasLock := false
		for k, v := range cfg.Root.Lock {
			if depTxt[:strings.Index(depTxt, "@")] == k {
				util.PrintfVerbose(logName, verbose, "%s locked by %s\n", util.Yellow(dep), util.Yellow(v))
				dep.Ref = v
				lockedDeps = append(lockedDeps, dep)
				hasLock = true
				break
			}
		}
		if !hasLock {
			util.PrintfVerbose(logName, verbose, "No lock for %s\n", util.Yellow(dep))
			lockedDeps = append(lockedDeps, dep)
		}
	}
	util.PrintfVerbose(logName, verbose, "Locked deps: %v\n", util.Yellow(lockedDeps))
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
