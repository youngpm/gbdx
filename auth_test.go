package gbdx

import "testing"

func TestIncorrectCredentials(t *testing.T) {
	token, err := GetAccessToken(GBDXConfiguration{
		Username:       "foo",
		Password:       "bar",
		ClientID:       "baz",
		ClientPassword: "bam",
	})
	if err.Error() != "gbdx.GetAccessToken failed: oauth2: cannot fetch token: 401 UNAUTHORIZED\nResponse: {\"error\": \"invalid_client\"}" {
		t.Errorf("GetAccessToken(GBDXConfiguration): expected 401 UNAUTHORIZED, actual %q", err)
	}
	if token != "" {
		t.Errorf("GetAccessToken(GBDXConfiguration): expected token \"\", actual %q", token)
	}
}
