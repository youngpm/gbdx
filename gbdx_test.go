package gbdx

import (
	"testing"

	"github.com/youngpm/gbdx"
)

func TestIncorrectCredentials(t *testing.T) {
	client, err := gbdx.GetHttpClient(gbdx.GBDXConfiguration{
		Username:     "foo",
		Password:     "bar",
		ClientID:     "baz",
		ClientSecret: "bam",
	})
	if err.Error() != "gbdx.GetHttpClient failed: oauth2: cannot fetch token: 401 UNAUTHORIZED\nResponse: {\"error\": \"invalid_client\"}" {
		t.Errorf("GetHttpClient(GBDXConfiguration): expected 401 UNAUTHORIZED, actual %q", err)
	}
	if client != nil {
		t.Errorf("GetHttpClient(GBDXConfiguration): expected nil client, actual %q", client)
	}
}
