// Copyright © 2016 Patrick Young <patrick.mckendree.young@gmail.com>
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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/youngpm/gbdx"
)

func orderLocation(cmd *cobra.Command, args []string) (err error) {
	// Read catalog ids from stdin (line seperated)  if given no arguments.
	if len(args) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			args = append(args, scanner.Text())
		}
	}

	// We do things concurrently, so make a channel for returning
	// responses on and one to use as a counting semaphore.
	type orderResponse struct {
		order *gbdx.Order
		err   error
	}
	ch := make(chan orderResponse)
	sema := make(chan struct{}, 10) // Max concurrent requests.
	quit := make(chan struct{})
	defer close(quit)

	// Aquire an Api.
	api, err := apiFromConfig()
	if err != nil {
		return err
	}

	// Place the order concurrently in batches of 100.
	var n sync.WaitGroup
	for i := 0; i < len(args); i += 100 {
		j := i + 100
		if j > len(args) {
			j = len(args)
		}

		n.Add(1)
		go func(a []string) {
			defer n.Done()

			select {
			case sema <- struct{}{}:
			case <-quit:
				return
			}

			var o orderResponse
			o.order, o.err = api.OrderLocation(a...)
			<-sema
			ch <- o
		}(args[i:j])
	}

	// Close ch once all goroutines have returned.
	go func() {
		n.Wait()
		close(ch)
	}()

	var orders Orders
	for resp := range ch {

		// On error, send quit signal and drain the channel.
		if resp.err != nil {
			quit <- struct{}{}
			go func() {
				for range ch {
					// Drain any stragglers.
				}
			}()
			return resp.err
		}

		orders.Append(resp.order)
	}

	result, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", result)

	return nil
}

var orderLocationCmd = &cobra.Command{
	Use:   "location [CATIDs]",
	Short: "Check the location of GBDX catalog IDs",
	Long: `Check the location of GBDX catalog IDs

Pass the catalog ids in space delimited arguments to check multiple
catalog ids, or if given no arguments, passed in delimited by newlines via
stdin.

A useful command to run to see how many strips might be ordered is
something like

  cat file-of-cids.txt | gbdxcli order location | jq '.acquisitions[].state' | sort | uniq -c

which will return a count of "delivered" and "not_found", the latter
would be ordered if you were to submit those catalog ids.
`,
	RunE: orderLocation}

func init() {
	orderCmd.AddCommand(orderLocationCmd)
}
