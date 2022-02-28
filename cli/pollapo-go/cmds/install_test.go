package cmds

import (
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/myzip"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
)

// TODO: installDepsRecursive is hard to test. too big.
func TestInstallConfig(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	zd := myzip.NewMockZipDownloader(ctrl)
	zd.EXPECT().
		GetZip(gomock.Eq("google"), gomock.Eq("apis"), gomock.Eq("dbfbfdb")).
		Return(nil, []byte("ASD"))
	uz := myzip.NewMockUnzipper(ctrl)
	uz.EXPECT().
		Unzip(gomock.Any(), gomock.Any()).
		Return()
	loader := pollapo.NewMockConfigLoader(ctrl)
	loader.EXPECT().
		GetPollapoConfig(gomock.Eq(filepath.Join(cache.GetDefaultCacheRoot(), "pollapo.yml"))).
		Return(pollapo.PollapoConfig{Deps: []string{"my/apis@abcd"}}, nil)

	// TODO: mock PollapoConfigFileLoader

	NewCmdInstall(
		false,
		".pollapo",
		"",
		zd,
		uz,
		pollapo.FileConfigLoader{},
		cache.EmptyCache{}, // Don't use cache
	).installDepsRecursive(
		pollapo.PollapoConfig{Deps: []string{"google/apis@dbfbfdb"}},
	)
}
