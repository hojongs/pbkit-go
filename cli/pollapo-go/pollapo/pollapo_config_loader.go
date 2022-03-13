package pollapo

import "os"

type ConfigLoader interface {
	GetPollapoConfig(pollapoYmlPath string) (PollapoConfig, error)
}

type FileConfigLoader struct{}

func (cl FileConfigLoader) GetPollapoConfig(pollapoYmlPath string) (PollapoConfig, error) {
	pollapoBytes, err := os.ReadFile(pollapoYmlPath)
	if err != nil {
		return PollapoConfigYml{}, err
	} else {
		return ParsePollapo(pollapoBytes), nil
	}
}
