package cmds

import (
	"os"

	"github.com/hojongs/pbkit-go/cli/pollapo/log"
	"github.com/hojongs/pbkit-go/cli/pollapo/pollapo"
)

func Install(
	clean bool,
	outDir string,
	token string,
	pollapoYmlPath string,
) {
	log.Infow("Params", "clean", clean, "outDir", outDir, "token", token, "config", pollapoYmlPath)

	pollapoBytes, err := os.ReadFile(pollapoYmlPath)
	if err != nil {
		log.Fatalw("Failed to read file", "filename", pollapoYmlPath, "cause", err.Error())
	}

	cfg := pollapo.ParsePollapo(pollapoBytes)
	log.Infow("LoadPollapoYml", "pollapoYml", cfg)
	// install deps in cfg
	q := []string{}
	q = append(q, cfg.Deps...)
	for len(q) > 0 {
		depTxt := q[0]
		q = q[1:]

		dep, b := pollapo.ParseDep(depTxt)
		if !b {
			log.Fatalw("Invalid dep", nil, "dep", depTxt)
		}

	}

	// getToken
	// backoff (validateToken)
	// cacheDir
	// cacheDeps
	// lockTable
	// analyzeDeps(cacheDir, pollapoYml)
	// *emptyDir
	// *recursive installDep
	// stringify sanitizeDeps
	// writeFile
	//
}
