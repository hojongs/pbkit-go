package pollapo

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
	"gopkg.in/yaml.v3"
)

type PollapoConfig interface {
	GetDeps(verbose bool) []PollapoDep
	GetLock(dep PollapoDep) (string, bool)
	SetLock(dep PollapoDep, lockedRef string)
	SaveFile(filepath string) error
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
	if cfg.Root.Lock == nil {
		cfg.Root.Lock = make(map[string]string)
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

		if lockedRef, found := cfg.getLock(dep, verbose); found {
			lockedDep := dep
			lockedDep.Ref = lockedRef
			lockedDeps = append(lockedDeps, lockedDep)
		} else {
			util.PrintfVerbose(logName, verbose, "No lock for %s\n", util.Yellow(dep))
			lockedDeps = append(lockedDeps, dep)
		}
	}
	util.PrintfVerbose(logName, verbose, "Locked deps: %v\n", util.Yellow(lockedDeps))
	return lockedDeps
}

func (cfg PollapoConfigYml) GetLock(dep PollapoDep) (string, bool) {
	return cfg.getLock(dep, false)
}

func (cfg PollapoConfigYml) SetLock(dep PollapoDep, lockedRef string) {
	cfg.Root.Lock[lockKey(dep)] = lockedRef
}

func (cfg PollapoConfigYml) getLock(dep PollapoDep, verbose bool) (string, bool) {
	lockedRef, found := cfg.Root.Lock[lockKey(dep)]
	if found {
		util.PrintfVerbose(logName, verbose, "%s locked by %s\n", util.Yellow(dep), util.Yellow(lockedRef))
		return lockedRef, true
	}
	return dep.Ref, found
}

func (cfg PollapoConfigYml) SaveFile(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func lockKey(dep PollapoDep) string {
	return fmt.Sprintf("%s/%s@%s", dep.Owner, dep.Repo, dep.Ref)
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
