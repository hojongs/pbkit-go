package cmds

import (
	"github.com/hojongs/pbkit-go/cli/pollapo/yml"
)

func Install(
	clean bool,
	outDir string,
	token string,
	// ymlPath
	config string,
) {
	// etToken
	// ackoff (validateToken)
	// cacheDir
	pollapoYml := yml.LoadPollapoYml(config)
	// cacheDeps
	// lockTable
	// analyzeDeps
	// *emptyDir
	// *recursive installDep
	// stringify sanitizeDeps
	// writeFile
	//
}
