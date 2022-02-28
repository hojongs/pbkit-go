package cmds

import (
	// "os"
	"testing"

	"github.com/golang/mock/gomock"
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

	NewCmdContextInstall(
		false,
		".pollapo",
		"",
		zd,
		PollapoConfigFileLoader{},
	).installDepsRecursive(
		pollapo.PollapoConfig{Deps: []string{"google/apis@dbfbfdb"}},
	)
}
