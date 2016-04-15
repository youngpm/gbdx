package gbdx

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type Api struct {
	tokenSource oauth2.TokenSource
	client      *http.Client
}

// NewApi returns an Api struct for interacting with GBDX.
func NewApi(c Config) (*Api, error) {

	oauth2Conf := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     oauth2.Endpoint{TokenURL: endpoints.tokens},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), GBDX_HTTP_TIMEOUT)
	token, err := oauth2Conf.PasswordCredentialsToken(ctx, c.Username, c.Password)
	defer cancel()
	if err != nil {
		return nil, err
	}

	tokenSource := oauth2Conf.TokenSource(context.TODO(), token)
	return &Api{
			tokenSource: tokenSource,
			client:      oauth2.NewClient(context.TODO(), tokenSource)},
		nil
}

func (a *Api) Token() (*oauth2.Token, error) {
	return a.tokenSource.Token()
}
