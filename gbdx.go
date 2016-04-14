package gbdx

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const GBDX_HTTP_TIMEOUT = 60 * time.Second

type GBDXConfiguration struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
}

// GetHttpClient returns an HTTP client pre-configured with a GBDX oauth2
// token. https://gbdxdocs.digitalglobe.com/docs/authentication-course
func GetHttpClient(conf GBDXConfiguration) (*http.Client, error) {
	oauth2Conf := &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Endpoint:     oauth2.Endpoint{TokenURL: "https://geobigdata.io/auth/v1/oauth/token/"},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), GBDX_HTTP_TIMEOUT)
	defer cancel() // releases resources if PasswordCredentialsToken completes before timeout elapses
	token, err := oauth2Conf.PasswordCredentialsToken(ctx, conf.Username, conf.Password)
	if err != nil {
		return nil, fmt.Errorf("gbdx.GetHttpClient failed: %s", err)
	}
	return oauth2Conf.Client(oauth2.NoContext, token), nil
}
