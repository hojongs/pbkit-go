package github

import (
	"github.com/cli/oauth"
)

const (
	// https://github.com/cli/cli/blob/91c4a5d828265823f6ed29129713778e2c839470/internal/authflow/flow.go#L19-L22
	// The "GitHub CLI" OAuth app
	oauthClientID = "178c6fc778ccc68e1d6a"
	// This value is safe to be embedded in version control
	oauthClientSecret = "34ddeff2b558a23d38fba8a6de74f086ede1cc0b"
)

// Try initiating OAuth Device flow on the server and fall back to OAuth Web application flow if
// Device flow seems unsupported. This approach isn't strictly needed for github.com, as its Device
// flow support is globally available, but enables logging in to hosted GitHub instances as well.
func TryOauthFlow() string {
	flow := &oauth.Flow{
		Host:         oauth.GitHubHost("https://github.com"),
		ClientID:     oauthClientID,
		ClientSecret: oauthClientSecret,           // only applicable to web app flow
		CallbackURI:  "http://127.0.0.1/callback", // only applicable to web app flow
		Scopes:       []string{"repo", "read:org", "gist"},
	}

	accessToken, err := flow.DetectFlow()
	if err != nil {
		panic(err)
	}
	return accessToken.Token
}
