package github

import (
	"context"

	"github.com/google/go-github/v42/github"
)

func GetZipLink(owner string, repo string, ref string) string {
	client := github.NewClient(nil)
	opts := github.RepositoryContentGetOptions{Ref: ref}
	url, _, _ := client.Repositories.GetArchiveLink(context.Background(), owner, repo, github.Zipball, &opts, true)
	return url.String()
}
