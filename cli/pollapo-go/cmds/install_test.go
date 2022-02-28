package cmds

import (
	// "os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/mocks"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
)

func TestInstallConfig(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	zd := mocks.NewMockZipDownloader(ctrl)

	zd.
		EXPECT().
		GetZipBin(gomock.Eq("google"), gomock.Eq("apis"), gomock.Eq("dbfbfdb")).
		Return([]byte("ASD"))
	// TODO: mock PollapoConfigFileLoader

	NewCmdInstall(
		false,
		".pollapo",
		"",
		zd,
		PollapoConfigFileLoader{},
		cache.EmptyCache{}, // Don't use cache
	).installDepsRecursive(
		pollapo.PollapoConfig{Deps: []string{"google/apis@dbfbfdb"}},
	)
}
