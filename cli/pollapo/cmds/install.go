package cmds

import (
	"fmt"

	"github.com/hojongs/pbkit-go/cli/pollapo/yml"
)

func Install(
	clean bool,
	outDir string,
	token string,
	// ymlPath
	config string,
) {
	fmt.Println(clean, outDir, token, config)
	// etToken
	// ackoff (validateToken)
	// cacheDir
	pollapoYml := yml.LoadPollapoYml(config)
	fmt.Println(pollapoYml)
	// cacheDeps
	// lockTable
	// analyzeDeps
	// *emptyDir
	// *recursive installDep
	// stringify sanitizeDeps
	// writeFile
	//
}
