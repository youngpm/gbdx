package gbdx

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Config holds the various configuation items we need to interact with GBDX.
type Config struct {
	Username string `mapstructure:"gbdx_username" toml:"gbdx_username"`
	Password string `mapstructure:"gbdx_password" toml:"gbdx_password"`
	//ClientID     string        `mapstructure:"gbdx_client_id" toml:"gbdx_client_id"`
	//ClientSecret string        `mapstructure:"gbdx_client_secret" toml:"gbdx_client_secret"`
	Token *oauth2.Token `mapstructure:"gbdx_token" toml:"gbdx_token"`
}

// Api holds GBDX authorized http clients and tokens.
type Api struct {
	tokenSource oauth2.TokenSource // This guy is kept around so we can cache the token for reuse.
	client      *http.Client
}

// NewApi returns an Api struct for interacting with GBDX.
func NewApi(c Config) (*Api, error) {

	oauth2Conf := &oauth2.Config{
		//ClientID:     c.ClientID,
		//ClientSecret: c.ClientSecret,
		Endpoint: oauth2.Endpoint{TokenURL: endpoints.tokens},
	}

	// Use a pre existing token if we were passed one.
	token := c.Token

	// Get a new token if we need it.
	var err error
	if token == nil {
		ctx, cancel := context.WithTimeout(context.TODO(), GBDX_HTTP_TIMEOUT)
		token, err = oauth2Conf.PasswordCredentialsToken(ctx, c.Username, c.Password)
		defer cancel()
		if err != nil {
			fmt.Println(c)
			return nil, err
		}
	}

	tokenSource := oauth2Conf.TokenSource(context.TODO(), token)

	return &Api{
		tokenSource: tokenSource,
		client:      oauth2.NewClient(context.TODO(), tokenSource),
	}, err
}

// Token returns a GBDX auth token.
func (a *Api) Token() (*oauth2.Token, error) {
	return a.tokenSource.Token()
}
