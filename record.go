package gbdx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"os"
	"strings"
	"strconv"
)

// Stats holds the number of records returned and the typeCounts
type Stat_Struct struct {
	RecordsReturned		int		`json:"recordsReturned"`
	TotalRecords		int		`json:"totalRecords"`
	StatImpl		Type_Struct	`json:"typeCounts"`
}

// Type_Struct holds the typeCounts
type Type_Struct struct {
	GbdxCatalogRecord	int	`json:"GBDXCatalogRecord"`
	DigitalGlobeAcquisition	int	`json:"DigitalGlobeAcquisition"`
	Acquisition		int	`json:"Acquisition"`
}

// Catalog holds the ID and all the associated Acquisitions of requested imagery.
type Record struct {
	ID          string          `json:"identifier"`
	Type        []string        `json:"type"`
	Props       Record_Struct   `json:"properties"`
}

// Acquisition holds the ID, state, and location (an s3 path) of the catalog id.
type Record_Struct struct {
	Bearing               json.Number  `json:"bearing"`
	BrowseURL             string       `json:"browseURL"`
	CatalogID             string       `json:"catalogID"`
	CloudCover            int          `json:"cloudCover"`
	FootprintWkt          string       `json:"footprintWkt"`
	ImageBands            string       `json:"imageBands"`
	MultiResolution       json.Number  `json:"multiResolution"`
	MultiResolutionMax    json.Number  `json:"multiResolution_max"`
	MultiResolutionStart  json.Number  `json:"multiResolution_start"`
	MultiResolutionEnd    json.Number  `json:"multiResolution_end"`
	MultiResolutionMin    json.Number  `json:"multiResolution_min"`
	OffNadirAngle         json.Number  `json:"offNadirAngle"`
	OffNadirAngleEnd      json.Number  `json:"offNadirAngle_end"`
	OffNadirAngleMax      json.Number  `json:"offNadirAngle_max"`
	OffNadirAngleMin      json.Number  `json:"offNadirAngle_min"`
	OffNadirAngleStart    json.Number  `json:"offNadirAngle_start"`
	PanResolution         json.Number  `json:"panResolution"`
	PanResolutionEnd      json.Number  `json:"panResolution_end"`
	PanResolutionMax      json.Number  `json:"panResolution_max"`
	PanResolutionMin      json.Number  `json:"panResolution_min"`
	PanResolutionStart    json.Number  `json:"panResolution_start"`
	PlatformName          string       `json:"platformName"`
	ScanDirection         string       `json:"scanDirection"`
	SensorPlatformName    string       `json:"sensorPlatformName"`
	StereoPair            json.Number  `json:"stereoPair"`
	SunAzimuth            json.Number  `json:"sunAzimuth"`
	SunAzimuthMax         json.Number  `json:"sunAzimuth_max"`
	SunAzimuthMin         json.Number  `json:"sunAzimuth_min"`
	SunElevation          json.Number  `json:"sunElevation"`
	SunElevationMax       json.Number  `json:"sunElevation_max"`
	SunElevationMin       json.Number  `json:"sunElevation_min"`
	TargetAzimuth         json.Number  `json:"targetAzimuth"`
	TargetAzimuthEnd      json.Number  `json:"targetAzimuth_end"`
	TargetAzimuthMax      json.Number  `json:"targetAzimuth_max"`
	TargetAzimuthMin      json.Number  `json:"targetAzimuth_min"`
	TargetAzimuthStart    json.Number  `json:"targetAzimuth_start"`
	Timestamp             string       `json:"timestamp"`
	Vendor		      string       `json:"vendor"`
}

type Record_1B struct {
        ID          string           `json:"identifier"`
        Type        []string         `json:"type"`
        Props       Record_1B_Struct `json:"properties"`
}

type Record_1B_Struct struct {
        AttFile       string       `json:"attFile"`
        Bands         string       `json:"bands"`
        BandList      string       `json:"bandsList"`
        BrowseJpgFile string       `json:"browseJpgFile"`
        BucketName    string       `json:"bucketName"`
        BucketPrefix  string       `json:"bucketPrefix"`
        CatalogID     string       `json:"catalogID"`
        CloudCover    int          `json:"cloudCover"`
        EphFile       string       `json:"ephFile"`
        FootprintWkt  string       `json:"footprintWkt"`
        GeoFile       string       `json:"geoFile"`
        ImageFile     string       `json:"imageFile"`
        ImdFile       string       `json:"imdFile"`
        OffNadirAngle json.Number  `json:"offNadirAngle"`
        Part          int          `json:"part"`
	PlatformName  string       `json:"platformName"`
	ProductLevel  string       `json:"productLevel"`
	ReadmeTxtFile string       `json:"readmeTxtFile"`
	Resolution    json.Number  `json:"resolution"`
	RpbFile       string       `json:"rpbFile"`
	SensorPlatformName string  `json:"sensorPlatformName"`
        Soli          string       `json:"soli"`
        SunAzimuth    json.Number  `json:"sunAzimuth"`
        SunElevation  json.Number  `json:"sunElevation"`
        TilFile       string       `json:"tilFile"`
        Timestamp     string       `json:"timestamp"`
        Vendor	      string       `json:"vendor"`
        XmlFile       string       `json:"xmlFile"`
}

type Record_Idaho struct {
        ID          string		`json:"identifier"`
        Type        []string		`json:"type"`
        Props       Record_Idaho_Struct	`json:"properties"`
}

type Record_Idaho_Struct struct {
	DGCatalogId		string		`json:"DGCatalogId"`
	BucketName		string		`json:"bucketName"`
	ColorInterpretation	string		`json:"colorInterpretation"`
	DataType		string		`json:"dataType"`
	EpsgCode		string		`json:"epsgCode"`
	FootprintWkt		string		`json:"footprintWkt"`
	GroundSampleDistanceMeters json.Number	`json:"groundSampleDistanceMeters"`
	ImageHeight		json.Number	`json:"imageHeight"`
	ImageId			string		`json:"imageId"`
	ImageWidth		json.Number	`json:"imageWidth"`
	NativeTileFileFormat	string          `json:"nativeTileFileFormat"`
	NumBands		json.Number	`json:"numBands"`
	NumXTiles		json.Number	`json:"numXTiles"`
	NumYTiles		json.Number	`json:"numYTiles"`
	PlatformName		string		`json:"platformName"`
	Pniirs			json.Number	`json:"pniirs"`
	ProfileName		string          `json:"profileName"`
	SatElevation		json.Number     `json:"satElevation"`
	TileBucketName		string		`json:"tileBucketName"`
	TilePartition		string		`json:"tilePartition"`
	TileXOffset		json.Number	`json:"tileXOffset"`
	TileXSize		json.Number	`json:"tileXSize"`
	TileYOffset 		json.Number	`json:"tileYOffset"`
	TileYSize		json.Number	`json:"tileYSize"`
	Timestamp		string		`json:"timestamp"`
	Vendor			string		`json:"vendor"`
	VendorDatasetIdentifier	string		`json:"vendorDatasetIdentifier"`
	VendorDatasetIdentifier1 string		`json:"vendorDatasetIdentifier1"`
	VendorDatasetIdentifier2 string		`json:"vendorDatasetIdentifier2"`
	VendorDatasetIdentifier3 string		`json:"vendorDatasetIdentifier3"`
	VendorDatasetIdentifier4 string		`json:"vendorDatasetIdentifier4"`
	VendorName		string		`json:"vendorName"`
	Version			string		`json:"version"`
}

type Record_Landsat struct {
	ID          string                `json:"identifier"`
	Type        []string              `json:"type"`
	Props       Record_Landsat_Struct `json:"properties"`
}

type Record_Landsat_Struct struct {
	BrowseURL     		string        `json:"browseURL"`
	BucketName		string        `json:"bucketName"`
	BucketPrefix		string        `json:"bucketPrefix"`
	CatalogID		string        `json:"catalogID"`
        CloudCover		int           `json:"cloudCover"`
	FootprintWkt		string        `json:"footprintWkt"`
	MultiResolution		json.Number   `json:"multiResolution"`
	PanResolution		json.Number   `json:"panResolution"`
	Path			int           `json:"path"`
	PlatformName		string        `json:"platformName"`
	Row			int           `json:"row"`
	SensorPlatformName	string        `json:"sensorPlatformName"`
	Vendor			string        `json:"vendor"`
	Timestamp		string        `json:"timestamp"`
}

// Search criteria holds stardDate, endDate, types, filters and AOI.
type Search_Json struct {
	SearchAreaWkt   string          `json:"searchAreaWkt,omitempty"`
	StartDate	string		`json:"startDate"`
	EndDate		string		`json:"endDate"`
	Types		[]string	`json:"types,omitempty"`
	Filters		[]string	`json:"filters,omitempty"`
	Limit		int		`json:"limit,omitempty"`
}

// Results from search holds results and stats.
type Result_Parent struct {
	Stats		Stat_Struct	`json:"stats"`
	Records		[]Record	`json:"results"`
	// SearchTag	string		`json:"searchTag,omitempty"`
}

type Result_1B_Parent struct {
        Stats           Stat_Struct     `json:"stats"`
        Records         []Record_1B     `json:"results"`
}

type Result_Idaho_Parent struct {
        Stats           Stat_Struct	`json:"stats"`
        Records         []Record_Idaho	`json:"results"`
}

type Result_Landsat_Parent struct {
        Stats           Stat_Struct      `json:"stats"`
        Records         []Record_Landsat `json:"results"`
}

// RecordStatus returns the status of the imagery given the string identifying it.
func (a *Api) RecordStatus(catID string) (*Record, *Record_1B, *Record_Landsat, *Record_Idaho, error) {

	url := fmt.Sprintf("%s%s", endpoints.record, catID)
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, nil, fmt.Errorf("Catalog get record returned a bad status code: %s", resp.Status)
	}

	lcSubStr := "LC"
	if len(catID) == 16 {
		var record *Record
		err = json.NewDecoder(resp.Body).Decode(&record)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("Catalog get record failed to decode properly: %v", err)
		}
		return record, nil, nil, nil, err
	} else if strings.HasPrefix(catID,lcSubStr) {
		var record *Record_Landsat
                err = json.NewDecoder(resp.Body).Decode(&record)
                if err != nil {
                        return nil, nil, nil, nil, fmt.Errorf("Catalog get record failed to decode properly: %v", err)
                }
                return nil, nil, record, nil, err
	} else if len(catID) == 36 {
		var record *Record_Idaho
		err = json.NewDecoder(resp.Body).Decode(&record)
                if err != nil {
                        return nil, nil, nil, nil, fmt.Errorf("Catalog get record failed to decode properly: %v", err)
                }
                return nil, nil, nil, record, err

	} else {
		var record *Record_1B
                err = json.NewDecoder(resp.Body).Decode(&record)
                if err != nil {
                        return nil, nil, nil, nil, fmt.Errorf("Catalog get record failed to decode properly: %v", err)
                }
                return nil, record, nil, nil, err
	}
}

func getLimitFromFilterArg(inFilter []string) (int, int) {
	var limit int = 0
	for i := 0; i < len(inFilter); i += 1 {
		if strings.HasPrefix(inFilter[i], "limit=") {
			str_limit := strings.Trim(inFilter[i], "limit=")
			limit, _ = strconv.Atoi(str_limit)
			return limit, i
		}
	}
	return 0, 0
}

func (a *Api) GetRecords(inVal []string) (*Result_Parent, *Result_1B_Parent, *Result_Landsat_Parent, *Result_Idaho_Parent, error) {

	var search Search_Json
        var many_types []string
	if len(inVal) == 5 {
		many_types = strings.Split(inVal[2], ",")
		many_filters := strings.Split(inVal[4], ",")
		limit, pos := getLimitFromFilterArg(many_filters)
		if pos > 0 {
			many_filters = many_filters[:pos]
		}
		if len(many_filters) == 0 {
			search = Search_Json{SearchAreaWkt:inVal[3],StartDate:inVal[0],EndDate:inVal[1],Types:many_types,Limit:limit}
		} else if limit == 0 {
			search = Search_Json{SearchAreaWkt:inVal[3],StartDate:inVal[0],EndDate:inVal[1],Types:many_types,Filters:many_filters}
		} else {
			search = Search_Json{SearchAreaWkt:inVal[3],StartDate:inVal[0],EndDate:inVal[1],Types:many_types,Filters:many_filters,Limit:limit}
		}
	} else if len(inVal) == 4 {
		many_types = strings.Split(inVal[2], ",")
		if strings.HasPrefix(inVal[3], "POLYGON((") {
			search = Search_Json{SearchAreaWkt:inVal[3],StartDate:inVal[0],EndDate:inVal[1],Types:many_types}
		} else {
			many_filters := strings.Split(inVal[3], ",")
			limit, pos := getLimitFromFilterArg(many_filters)
			if pos > 0 {
				many_filters = many_filters[:pos]
			}
			if limit == 0 {
				search = Search_Json{StartDate:inVal[0],EndDate:inVal[1],Types:many_types,Filters:many_filters}
			} else {
				many_filters = many_filters[:pos]
				search = Search_Json{StartDate:inVal[0],EndDate:inVal[1],Types:many_types,Filters:many_filters,Limit:limit}
			}
		}
	} else if len(inVal) == 3 {
		many_types = strings.Split(inVal[2], ",")
		search = Search_Json{StartDate:inVal[0],EndDate:inVal[1],Types:many_types}
	} else {
		return nil, nil, nil, nil, fmt.Errorf("Catalog search failed due to missing arguments")
	}

        b := new(bytes.Buffer)
        json.NewEncoder(b).Encode(search)

	resp, err := a.client.Post(endpoints.catalogSearch, "application/json", b)
        if err != nil {
                return nil, nil, nil, nil, fmt.Errorf("GetRecords failed to post: %v", err)
        }

        defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		return nil, nil, nil, nil, fmt.Errorf("Catalog search post returned a bad status code: %s %v", resp.Status, err)
	}

	var parent *Result_Parent
	dg_acq_exist := false
	for i := 0; i < len(many_types); i += 1 {
		if many_types[i] == "DigitalGlobeAcquisition" {
			dg_acq_exist = true
		}
	}
	if dg_acq_exist == true {
		err = json.NewDecoder(resp.Body).Decode(&parent)
	} else if many_types[0] == "1BProduct" {
		var parent_1b *Result_1B_Parent
		err = json.NewDecoder(resp.Body).Decode(&parent_1b)
        	if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("GetRecords response failed to decode 1B properly: %v", err)
		}
		return nil, parent_1b, nil, nil, err
	} else if many_types[0] == "LandsatAcquisition" {
		var parent_ls *Result_Landsat_Parent
                err = json.NewDecoder(resp.Body).Decode(&parent_ls)
                if err != nil {
                        return nil, nil, nil, nil, fmt.Errorf("GetRecords response failed to decode Landsat properly: %v", err)
                }
                return nil, nil, parent_ls, nil, err
	} else if many_types[0] == "IDAHOImage" {
                var parent_idaho *Result_Idaho_Parent
                err = json.NewDecoder(resp.Body).Decode(&parent_idaho)
                if err != nil {
                        return nil, nil, nil, nil, fmt.Errorf("GetRecords response failed to decode Idaho properly: %v", err)
                }
                return nil, nil, nil, parent_idaho, err
	} else {
		// GE01, WV03 types can return multiple types of catalog records.
		err = json.NewDecoder(resp.Body).Decode(&parent)
	}

	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("GetRecords response failed to decode properly: %v", err)
	}
	return parent, nil, nil, nil, err
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
