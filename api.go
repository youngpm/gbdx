package gbdx

import (
	"fmt"
	"net/http"

	"io"

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

// Token returns a GBDX auth token.
func (a *Api) Token() (*oauth2.Token, error) {
	return a.tokenSource.Token()
}

// Browse writes a browse image with catalog id cid and requested dimension dim to w.
func Browse(cid string, dim string, w io.Writer) error {

	url := fmt.Sprintf("%s%s.%s.png", endpoints.browse, cid, dim)

	// Note we use the default http client as no authentication is
	// required (in fact it breaks the endpoint) to get a browse.
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("browse fetch get failure %s: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("browse fetch returned a bad status code: %s", resp.Status)
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("failed copying browse to output: %v", err)
	}
	return err
}
