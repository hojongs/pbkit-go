package github

import (
	"os"
	"path/filepath"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/yaml"
)

func WriteTokenGhHosts(token string) []byte {
	path := getDefaultGhHostsPath()
	hostsYml := GhHosts{
		GithubCom: GhHost{
			OauthToken:  token,
			GitProtocol: "ssh",
		},
	}
	err := os.MkdirAll(filepath.Base(path), 0755)
	if err != nil {
		log.Fatalw("Failed to mkdir", err)
	}
	barr, err := yaml.Marshal(hostsYml)
	if err != nil {
		log.Fatalw("Failed to marshal yml", err)
	}
	err = os.WriteFile(path, barr, 0644)
	if err != nil {
		log.Fatalw("Failed to write file", err, "path", path)
	}
	return barr
}

func GetTokenFromGhHosts() string {
	path := getDefaultGhHostsPath()
	barr, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	cfg := GhHosts{}
	err = yaml.Unmarshal([]byte(barr), &cfg)
	if err != nil {
		log.Fatalw("Failed to parse gh hosts.yaml", err)
	}
	return cfg.GithubCom.OauthToken
}

type GhHosts struct {
	GithubCom GhHost `yaml:"githubcom"`
}

type GhHost struct {
	OauthToken  string `yaml:"oauth_token"`
	GitProtocol string `yaml:"git_protocol"`
}

func getDefaultGhHostsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalw("Failed to get home dir", err)
	}
	dir := filepath.Join(homeDir, ".config/gh")
	return filepath.Join(dir, "hosts.yml")
}
