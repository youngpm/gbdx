package gbdx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"os"
	"strings"
)

// Catalog holds the ID and all the associated Acquisitions of requested imagery.
type Record struct {
	Owner       string          `json:"owner"`
	ID          string          `json:"identifier"`
	Type        string          `json:"type"`
	Props       Record_Struct   `json:"properties"`
}

// Acquisition holds the ID, state, and location (an s3 path) of the catalog id.
type Record_Struct struct {
	Ordered               string  `json:"ordered"`
	Type                  string  `json:"type"`
	Sun_Elevation         string  `json:"sunElevation"`
	Targert_Azimuth       string  `json:"targetAzimuth"`
	Sensor_Platform_Name  string  `json:"sensorPlatformName"`
	Sun_Azimuth           string  `json:"sunAzimuth"`
	Browse_URL            string  `json:"browseURL"`
	Available             string  `json:"available"`
	Cloud_Cover           string  `json:"cloudCover"`
	Multi_Resolution      string  `json:"multiResolution"`
	Vendor_Name           string  `json:"vendorName"`
	Off_Nadir_Angle       string  `json:"offNadirAngle"`
	Pan_Resolution        string  `json:"panResolution"`
	Catalog_ID            string  `json:"catalogID"`
	Image_Bands           string  `json:"imageBands"`
	FootprintWkt          string  `json:"footprintWkt"`
	Timestamp             string  `json:"timestamp"`
}

// Search criteria holds stardDate, endDate, types and AOI.
type Search_Json struct {
	SearchAreaWkt   string          `json:"searchAreaWkt,omitempty"`
	StartDate	string		`json:"startDate"`
	EndDate		string		`json:"endDate"`
	Types		[]string	`json:"types,omitempty"`
	Filters		[]string	`json:"filters,omitempty"`
	Limit		int		`json:"limit,omitempty"`
}

// Results from search holds results, searchTag and stats.
type Result_Parent struct {
	Records		[]Record	`json:"results"`
	SearchTag	string		`json:"searchTag,omitempty"`
	// Stats		[]string	`json:"stats,omitempty"`
}

// RecordStatus returns the status of the imagery given the string identifying it.
func (a *Api) RecordStatus(catID string) (*Record, error) {

	url := fmt.Sprintf("%s%s", endpoints.record, catID)
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Catalog get record returned a bad status code: %s", resp.Status)
	}

	var record *Record
	err = json.NewDecoder(resp.Body).Decode(&record)
	if err != nil {
		return nil, fmt.Errorf("Catalog get record failed to decode properly: %v", err)
	}
	return record, err
}

func (a *Api) GetRecords(inVal []string) (*Result_Parent, error) {

	var search Search_Json
	if len(inVal) == 5 {
		many_filters := strings.Split(inVal[4], ",")
		search = Search_Json{SearchAreaWkt:inVal[3],StartDate:inVal[0],EndDate:inVal[1],Types:[]string{inVal[2]},Filters:many_filters}
	} else if len(inVal) == 4 {
		search = Search_Json{SearchAreaWkt:inVal[3],StartDate:inVal[0],EndDate:inVal[1],Types:[]string{inVal[2]}}
	} else if len(inVal) == 3 {
		search = Search_Json{StartDate:inVal[0],EndDate:inVal[1],Types:[]string{inVal[2]}}
	} else {
		return nil, fmt.Errorf("Catalog search failed due to missing arguments")
	}

        b := new(bytes.Buffer)
        json.NewEncoder(b).Encode(search)

	resp, err := a.client.Post(endpoints.catalogSearch, "application/json", b)
        if err != nil {
                return nil, fmt.Errorf("GetRecords failed to post: %v", err)
        }

        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
                return nil, fmt.Errorf("Catalog search post returned a bad status code: %s %v", resp.Status, err)
        }

	var parent *Result_Parent
        err = json.NewDecoder(resp.Body).Decode(&parent)
        if err != nil {
                return nil, fmt.Errorf("GetRecords response failed to decode properly: %v", err)
        }
        return parent, err
}

// CatalogHeartbeat checks if the catalog endpoint is alive and well.
func CatalogHeartbeat() error {

	url := endpoints.catalogHeartbeat

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("catalog heartbeat failed %s: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("catalog heartbeat returned a bad status code: %s", resp.Status)
	}
	defer resp.Body.Close()

	return err
}
