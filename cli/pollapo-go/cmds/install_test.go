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

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	zd := mocks.NewMockZipDownloader(ctrl)

	// Asserts that the first and only call to Bar() is passed 99.
	// Anything else will fail.
	zd.
		EXPECT().
		GetZipBin(gomock.Eq("google"), gomock.Eq("apis"), gomock.Eq("dbfbfdb")).
		Return([]byte("ASD"))
	// TODO: add mock zip.Unzip(gomock.Eq([]byte("ASD")))

	InstallConfig(
		false,
		".pollapo",
		"",
		pollapo.PollapoConfig{
			Deps: []string{
				"google/apis@dbfbfdb",
			},
		},
		zd,
	)
}
