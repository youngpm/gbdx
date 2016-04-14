package gbdx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchFilters struct {
	Owner         interface{} `json:"owner"`
	State         interface{} `json:"state"`
	LookbackHours interface{} `json:"lookback_h"`
}

type SearchResults struct {
	Workflows []string
}

// Search workflows by filter.  Go interface to
// https://gbdxdocs.digitalglobe.com/docs/search-for-workflows-by-filter
//
//FIXME: Probably should _not_ be a method on SearchFilter!  Should
//this take an argument which implements something like io.Reader to
//define the filters?
func (params SearchFilters) Search(client *http.Client) (SearchResults, error) {
	var results = SearchResults{}

	// FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME
	//       Create an Interface Type that SearchFilters struct can have
	//       which will properly marshal JSON expected by GBDX.
	// FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME FIXME
	paramsMap := make(map[string]interface{})
	paramsMap["owner"] = params.Owner

	msgBody, err := json.Marshal(paramsMap)
	if err != nil {
		return results, fmt.Errorf("marshalling %q: %v", params, err)
	}

	url := "https://geobigdata.io/workflows/v1/workflows/search"
	response, err := client.Post(url, "application/json", bytes.NewReader(msgBody))
	if err != nil {
		return results, fmt.Errorf("HTTP POST %s: %v", url, err)

	}
	defer response.Body.Close()

	if response.Status != "200 OK" {
		var byteSlice []byte
		response.Request.Body.Read(byteSlice)
		return results, fmt.Errorf("HTTP POST %s;  returned status %s; request body %q", url, response.Status, byteSlice)
	}

	if err = json.NewDecoder(response.Body).Decode(&results); err != nil {
		return results, fmt.Errorf("Decoding search response %q: %v", response.Body, err)
	}

	return results, nil
}
