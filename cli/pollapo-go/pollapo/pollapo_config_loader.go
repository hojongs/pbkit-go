package pollapo

import (
	"os"
)

type ConfigLoader interface {
	GetPollapoConfig(pollapoYmlPath string) (PollapoConfig, error)
}

type FileConfigLoader struct{}

func (fcl FileConfigLoader) GetPollapoConfig(pollapoYmlPath string) (PollapoConfig, error) {
	// log.Infow("GetPollapoConfig", "pollapoYmlPath", pollapoYmlPath)
	pollapoBytes, err := os.ReadFile(pollapoYmlPath)
	if err != nil {
		return PollapoConfig{}, err
	} else {
		return ParsePollapo(pollapoBytes), nil
	}
}
