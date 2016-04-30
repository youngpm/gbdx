package gbdx

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Browse writes a browse image with catalog id cid and requested size to w.
func Browse(cid string, size string, json bool, w io.Writer) error {

	var endpoint string
	if json {
		endpoint = endpoints.browseJSON
	} else {
		endpoint = endpoints.browse
	}
	url := fmt.Sprintf("%s%s.%s.png", endpoint, cid, size)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("browse fetch get failure %s: %v", resp.Status, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("browse fetch returned a bad status code: %s", resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("failed copying browse to output: %v", err)
	}
	return err
}

// BrowseMetadata writes the GeoJSON metadata of the catalog id cid to w.
func BrowseMetadata(cid string, w io.Writer) error {

	url := fmt.Sprintf("%s%s.json", endpoints.browseMetadata, cid)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("browse metadata fetch get failure %s: %v", resp.Status, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("browse metadata fetch returned a bad status code: %s", resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("failed copying browse metadata to output: %v", err)
	}
	return err
}

// Thumbnail writes a thumbnail image with catalog id cid, dimension dim, and orientation to w.
func Thumbnail(cid string, dim int, orientation string, w io.Writer) error {

	u, err := url.Parse(fmt.Sprintf("%s%s/%d", endpoints.thumbnail, cid, dim))
	if err != nil {
		return err
	}
	q := u.Query()
	q.Add("orientation", orientation)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("thumbnail fetch get failure %s: %v", resp.Status, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("thumbnail fetch returned a bad status code: %s", resp.Status)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("failed copying thumbnail to output: %v", err)
	}
	return err
}
