package gbdx

import (
	"testing"

	"github.com/youngpm/gbdx"
)

func TestListBUcket(t *testing.T) {
	client, err := gbdx.GetHttpClient(gbdx.GBDXConfiguration{
		Username:     "",
		Password:     "",
		ClientID:     "",
		ClientSecret: "",
	})

	if err != nil {
		t.Errorf("gbdx.GetHttpClient(GBDXConfiguration): expected no error, got %q", err)
	}

	contents, err := gbdx.ListBucket(client, "/")
	if err != nil {
		t.Errorf("gbdx.ListBucket(client): expected no error, got %+v", err)
	}
	//t.Errorf("%q", contents)
}
