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

	"github.com/spf13/cobra"
	"github.com/youngpm/gbdx"
)

func orderStatus(cmd *cobra.Command, args []string) (err error) {

	// Read order ids from stdin (line seperated)  if given no arguments.
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
	ch := make(chan orderResponse, len(args))
	sema := make(chan struct{}, 10) // Max concurrent requests.

	// Aquire an Api.
	api, err := apiFromConfig()
	if err != nil {
		return err
	}

	// Concurrently fetch the order status.
	for _, id := range args {
		go func(id string) {
			sema <- struct{}{}
			var o orderResponse
			o.order, o.err = api.OrderStatus(id)
			<-sema
			ch <- o
		}(id)
	}

	// Aggregate the results.
	var orders Orders
	for range args {
		o := <-ch
		if o.err != nil {
			return o.err
		}
		orders.Append(o.order)
	}

	// Marshall the aggregated result.
	result, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", result)
	return cacheToken(api)
}

// orderCmd represents the order command
var orderStatusCmd = &cobra.Command{
	Use:   "status [ORDERIDs]",
	Short: "Check the status of GBDX orders",
	Long: `Check the status of GBDX orders.

Pass the order ids in space delimited arguments to check multiple
orders, or if given no arguments, passed in delimited by newlines via
stdin.

The requests are done concurrently for speed.`,
	RunE: orderStatus,
}

func init() {
	orderCmd.AddCommand(orderStatusCmd)
}
