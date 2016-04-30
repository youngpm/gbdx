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
	"encoding/json"
	"fmt"

	"bufio"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngpm/gbdx"
)

func order(cmd *cobra.Command, args []string) (err error) {

	// Aquire an Api.
	var profile GBDXProfile
	if err = viper.Unmarshal(&profile); err != nil {
		return err
	}
	api, err := gbdx.NewApi(profile.ActiveConfig)
	if err != nil {
		return err
	}

	// Read catalog ids from stdin (line seperated)  if given no arguments.
	if len(args) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			args = append(args, scanner.Text())
		}
	}

	// Order them.
	order, err := api.NewOrder(args)
	if err != nil {
		return err
	}

	result, err := json.Marshal(order)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", result)
	return err
}

// orderCmd represents the order command
var orderCmd = &cobra.Command{
	Use:   "order [ACQIDS]",
	Short: "Submit orders to GBDX",
	Long: `Order acquisitions from GBDX.

Acquisition IDs can be specified as space delimited arguments on the
command line, or if given no arguments, passed in delimited by
newlines via stdin.
`,
	RunE: order,
}

func init() {
	RootCmd.AddCommand(orderCmd)
}