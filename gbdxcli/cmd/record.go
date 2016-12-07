// Copyright © 2016 Michael Smith <mike.s.smith@digitalglobe.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"errors"

	"bufio"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/youngpm/gbdx"
)

// Records is an aggregation of the content from catalog searches.
// Currently, the GBDX api will only let you search 100 items at a time,
// but this is inconvienent to deal with.
type Records struct {
	IDs        []string      `json:"search_ids"`
	Collection []gbdx.Record `json:"search_records"`
}

func (r *Records) Append(parent *gbdx.Result_Parent) {
        if len(parent.SearchTag) == 0 {
		r.IDs = append(r.IDs, string(len(r.IDs)))
        } else {
		r.IDs = append(r.IDs, parent.SearchTag)
        }
	r.Collection = append(r.Collection, parent.Records...)
}

func dedup_records(r Records) ([]gbdx.Record) {
	if len(r.IDs) == 0 {
		return []gbdx.Record{}
	}

	mRecords := make(map[string]gbdx.Record)
	for i := 0; i < len(r.Collection); i++ {
		rec := r.Collection[i]
		_, exist := mRecords[rec.ID]
		if exist == false {
			mRecords[rec.ID] = rec
		}
	}

	retArr := make([]gbdx.Record, 0, len(mRecords))
	for _, value := range mRecords {
		retArr = append(retArr, value)
	}
	return retArr
}

func search(cmd *cobra.Command, args []string) (err error) {

	// Read record ids from stdin (line separated)  if given no arguments.
	addAreas := []string{}
	foundAddAreas := false
	
	if len(args) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 {
				break
			} else if len(line) <= 8 {
				foundAddAreas = true
			} else if foundAddAreas == false {
				args = append(args, line)
			} else {
				addAreas = append(addAreas, line)
			}
		}
	}
	if len(args) < 3 {
		err := errors.New("Required arguments are missing: StartDate, EndDate, Type")
		return err
	}

	// Aquire an Api.
	api, err := apiFromConfig()
	if err != nil {
		return err
	}

	// We do things concurrently, so make a channel for returning
	// responses on and one to use as a counting semaphore.
	type recordResponse struct {
		parent *gbdx.Result_Parent
		err   error
	}
	ch := make(chan recordResponse)
	sema := make(chan struct{}, 100) // Max concurrent requests.
	quit := make(chan struct{})
	defer close(quit)

	// Search concurrently based on aoi.
	loop_count := 1
	j := 3
	wkt := getSearchAreaFromArgs(args)
        polygons := getSmallerSearches(wkt)
	for i := 0; i < len(addAreas); i += 1 {
		wkt = getSearchAreaFromArgs([]string{addAreas[i]})
		polygons = append(polygons, getSmallerSearches(wkt)...)
	}
	if len(polygons) >= 1 {
		loop_count = len(polygons)
		j = 4
	}

	filter := getFilterFromArgs(args)
        if len(filter) > 0 {
		args[4] = filter
		j = 5
        }

	var n sync.WaitGroup
	for i := 0; i < loop_count; i += 1 {
		polygon_str := ""
		if len(polygons) >= 1 {
			polygon_str = polygons[i]
		}

		n.Add(1)
		go func(a []string, p_str string) {
			defer n.Done()

			select {
			case sema <- struct{}{}:
			case <-quit:
				return
			}

			if len(a) > 3 {
				a[3] = p_str
			}

			var o recordResponse
			o.parent, o.err = api.GetRecords(a)
			<-sema
			ch <- o
		}(args[0:j],polygon_str)
	}

	// Close ch once all goroutines have returned.
	go func() {
		n.Wait()
		close(ch)
	}()

	var parents Records
	for resp := range ch {

		// On error, send quit signal and drain the channel.
		if resp.err != nil {
                        // TODO: why is it here and why does it hang ??
			// quit <- struct{}{}
			go func() {
			 	for range ch {
					// Drain any stragglers.
				}
			}()
			return resp.err
		}
		parents.Append(resp.parent)

	}

	// Marshall the aggregated result.
	records := dedup_records(parents)
	fmt.Printf("[\n")
	for i := 0; i < len(records); i += 1 {
		result, err := json.Marshal(records[i])
		if err == nil {
			if i < len(records)-1 {
				fmt.Printf("%s,\n", result)
			} else {
				fmt.Printf("%s\n", result)
			}
		}
	}
	fmt.Printf("]\n")
	return cacheToken(api)
}

// recordCmd represents a catalog search
var recordCmd = &cobra.Command{
	Use:   "catalog startDate endDate Type SearchAreaWkt Filter",
	Short: "Search the catalog for records using GBDX",
	Long: `Search the catalog for records using GBDX

Acquisition IDs can be specified as space delimited arguments on the
command line, or if given no arguments, passed in delimited by
newlines via stdin.
`,
	RunE: search,
}

func init() {
	RootCmd.AddCommand(recordCmd)
}
