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
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/youngpm/gbdx"
)

// Records is an aggregation of the content from catalog searches.
// Currently, the GBDX api will only let you search 100 items at a time,
// but this is inconvienent to deal with.
type Records struct {
	IDs           []string         `json:"search_ids"`
	Collection    []gbdx.Record    `json:"search_records"`
	Collection_1B []gbdx.Record_1B `json:"search_1b_records"`
	Collection_LS []gbdx.Record_Landsat `json:"search_landsat_records"`
	Collection_IH []gbdx.Record_Idaho `json:"search_idaho_records"`
}

func (r *Records) Append(parent *gbdx.Result_Parent, parent_1b *gbdx.Result_1B_Parent, parent_ls *gbdx.Result_Landsat_Parent, parent_ih *gbdx.Result_Idaho_Parent) {
	r.IDs = append(r.IDs, string(len(r.IDs)))
        // if len(parent.SearchTag) == 0 {
 	// 	r.IDs = append(r.IDs, string(len(r.IDs)))
        // } else {
	// 	r.IDs = append(r.IDs, parent.SearchTag)
        // }
        if parent != nil {
		r.Collection = append(r.Collection, parent.Records...)
	}
	if parent_1b != nil {
		r.Collection_1B = append(r.Collection_1B, parent_1b.Records...)
	}
	if parent_ls != nil {
		r.Collection_LS = append(r.Collection_LS, parent_ls.Records...)
	}
	if parent_ih != nil {
                r.Collection_IH = append(r.Collection_IH, parent_ih.Records...)
        }
}

func dedup_records(r Records, most_recent bool, aoi_in_image bool) ([]gbdx.Record) {
	if len(r.IDs) == 0 || len(r.Collection) == 0 {
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
	var recent_rec gbdx.Record
        var recent_timestamp string = "1970-01-01T00:00:00.000Z"
	for _, value := range mRecords {
		retArr = append(retArr, value)
		if most_recent == true {
			if strings.Compare(value.Props.Timestamp,recent_timestamp) > 0 {
				recent_rec = value
				recent_timestamp = value.Props.Timestamp
			} 
		}
	}
	if most_recent == true {
		return []gbdx.Record{recent_rec}
	}
	return retArr
}

func dedup_1b_records(r Records, most_recent bool, aoi_in_image bool) ([]gbdx.Record_1B) {
        if len(r.IDs) == 0 || len(r.Collection_1B) == 0 {
                return []gbdx.Record_1B{}
        }

        mRecords := make(map[string]gbdx.Record_1B)
        for i := 0; i < len(r.Collection_1B); i++ {
                rec := r.Collection_1B[i]
                _, exist := mRecords[rec.ID]
                if exist == false {
                        mRecords[rec.ID] = rec
                }
        }

        retArr := make([]gbdx.Record_1B, 0, len(mRecords))
	var recent_rec gbdx.Record_1B
        var recent_timestamp string = "1970-01-01T00:00:00.000Z"
        for _, value := range mRecords {
                retArr = append(retArr, value)
                if most_recent == true {
                        if strings.Compare(value.Props.Timestamp,recent_timestamp) > 0 {
                                recent_rec = value
                                recent_timestamp = value.Props.Timestamp
                        }
                }
        }
        if most_recent == true {
                return []gbdx.Record_1B{recent_rec}
        }
        return retArr
}

func dedup_ls_records(r Records) ([]gbdx.Record_Landsat) {
        if len(r.IDs) == 0 || len(r.Collection_LS) == 0 {
                return []gbdx.Record_Landsat{}
        }

        mRecords := make(map[string]gbdx.Record_Landsat)
        for i := 0; i < len(r.Collection_LS); i++ {
                rec := r.Collection_LS[i]
                _, exist := mRecords[rec.ID]
                if exist == false {
                        mRecords[rec.ID] = rec
                }
        }

        retArr := make([]gbdx.Record_Landsat, 0, len(mRecords))
        for _, value := range mRecords {
                retArr = append(retArr, value)
        }
        return retArr
}

func dedup_idaho_records(r Records) ([]gbdx.Record_Idaho) {
        if len(r.IDs) == 0 || len(r.Collection_IH) == 0 {
                return []gbdx.Record_Idaho{}
        }

        mRecords := make(map[string]gbdx.Record_Idaho)
        for i := 0; i < len(r.Collection_IH); i++ {
                rec := r.Collection_IH[i]
                _, exist := mRecords[rec.ID]
                if exist == false {
                        mRecords[rec.ID] = rec
                }
        }

        retArr := make([]gbdx.Record_Idaho, 0, len(mRecords))
        for _, value := range mRecords {
                retArr = append(retArr, value)
        }
        return retArr
}

func search(cmd *cobra.Command, args []string) (err error) {

	// Read record ids from stdin (line separated) if given no arguments.
	addAreas := []string{}
	foundAddAreas := false

	if len(args) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 {
				break
			} else if len(line) <= 8 {
				if strings.HasPrefix(line, "#") {
					foundAddAreas = true
				} else {
					args = append(args, line)
				}
			} else if foundAddAreas == false {
				args = append(args, line)
			} else {
				addAreas = append(addAreas, line)
			}
		}
	}
	if len(args) < 3 {
		err := errors.New("Required arguments are missing: StartDate, EndDate, Type(s)")
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
		parent    *gbdx.Result_Parent
		parent_1b *gbdx.Result_1B_Parent
		parent_ls *gbdx.Result_Landsat_Parent
		parent_ih *gbdx.Result_Idaho_Parent
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
		if len(polygons) >= 1 {
			args[4] = filter
			j = 5
		} else {
			args[3] = filter
			j = 4
		}
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

			if len(a) > 3 && len(polygons) >= 1 {
				a[3] = p_str
			}

			var o recordResponse
			o.parent, o.parent_1b, o.parent_ls, o.parent_ih, o.err = api.GetRecords(a)
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
		parents.Append(resp.parent, resp.parent_1b, resp.parent_ls, resp.parent_ih)
	}

	// Marshall the aggregated result.
	specialOp := getSpecialOpFromLimit(filter)
	records := dedup_records(parents, strings.HasPrefix(specialOp,"recen"), strings.HasPrefix(specialOp,"bestfi"))
	record_1bs := dedup_1b_records(parents, strings.HasPrefix(specialOp,"recen"), strings.HasPrefix(specialOp,"bestfi"))
	record_landsats := dedup_ls_records(parents)
	record_idahos := dedup_idaho_records(parents)
	fmt.Printf("[\n")
	for i := 0; i < len(records); i += 1 {
		result, err := json.Marshal(records[i])
		if err == nil {
			if i < len(records)-1 {
				fmt.Printf("%s,\n", result)
			} else if len(record_1bs) == 0 && len(record_landsats) == 0 && len(record_idahos) == 0 {
				fmt.Printf("%s\n", result)
			} else {
				fmt.Printf("%s,\n", result)
			}
		}
	}
        for i := 0; i < len(record_1bs); i += 1 {
                result, err := json.Marshal(record_1bs[i])
                if err == nil {
                        if i < len(record_1bs)-1 {
                                fmt.Printf("%s,\n", result)
                        } else if len(record_landsats) == 0 && len(record_idahos) == 0 {
                                fmt.Printf("%s\n", result)
			} else {
                                fmt.Printf("%s,\n", result)
                        }
		}
        }
	for i := 0; i < len(record_landsats); i += 1 {
                result, err := json.Marshal(record_landsats[i])
                if err == nil {
                        if i < len(record_landsats)-1 {
                                fmt.Printf("%s,\n", result)
                        } else if len(record_idahos) == 0 {
				fmt.Printf("%s\n", result)
			} else {
                                fmt.Printf("%s,\n", result)
                        }
                }
        }
	for i := 0; i < len(record_idahos); i += 1 {
                result, err := json.Marshal(record_idahos[i])
                if err == nil {
                        if i < len(record_idahos)-1 {
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
