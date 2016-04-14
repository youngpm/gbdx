package gbdx

import (
	"testing"

	"github.com/youngpm/gbdx"
)

func TestSearch(t *testing.T) {
	client, err := gbdx.GetHttpClient(gbdx.GBDXConfiguration{
		Username:     "",
		Password:     "",
		ClientID:     "",
		ClientSecret: "",
	})

	if err != nil {
		t.Errorf("gbdx.GetHttpClient(GBDXConfiguration): expected no error, got %q", err)
	}

	results, err := gbdx.SearchFilters{Owner: "pschmitt"}.Search(client)
	if err != nil {
		t.Errorf("SearchFilter.Search(client): expected no error, got %q", err)
	}
	t.Log(results)
}
