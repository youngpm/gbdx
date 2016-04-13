package gbdx

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const tokenURL = "https://geobigdata.io/auth/v1/oauth/token/"

// GetAccessToken returns a GBDX Access Token for the credentials in GBDXConfiguration.
func GetAccessToken(conf GBDXConfiguration) (string, error) {
	oauth2Conf := oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientPassword,
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
		},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), GBDX_HTTP_TIMEOUT)
	defer cancel() // releases resources if PasswordCredentialsToken completes before timeout elapses
	token, err := oauth2Conf.PasswordCredentialsToken(ctx, conf.Username, conf.Password)
	if err != nil {
		return "", fmt.Errorf("gbdx.GetAccessToken failed: %s", err)
	}

	return token.AccessToken, nil
}
