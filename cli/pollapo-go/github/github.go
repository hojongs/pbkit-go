package github

import (
	"context"

	"github.com/google/go-github/v42/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
)

func GetZipLink(dep pollapo.PollapoDep) string {
	client := github.NewClient(nil)
	opts := github.RepositoryContentGetOptions{Ref: dep.Ref}
	url, _, _ := client.Repositories.GetArchiveLink(context.Background(), dep.Owner, dep.Repo, github.Zipball, &opts, true)
	return url.String()
}
