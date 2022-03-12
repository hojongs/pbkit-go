package github

import (
	"os"
	"path/filepath"
	"runtime"

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
	GithubCom GhHost `yaml:"github.com"`
}

type GhHost struct {
	OauthToken  string `yaml:"oauth_token"`
	GitProtocol string `yaml:"git_protocol"`
}

func getDefaultGhHostsPath() string {
	var dir string
	// https://github.com/cli/cli/blob/26d33d6e387857f3d2e34f2529e7b05c7c51535f/internal/config/config_file.go#L29
	if c := os.Getenv("AppData"); runtime.GOOS == "windows" && c != "" {
		dir = filepath.Join(c, "GitHub CLI")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalw("Failed to get home dir", err)
		}
		dir = filepath.Join(homeDir, ".config", "gh")
	}
	return filepath.Join(dir, "hosts.yml")
}
