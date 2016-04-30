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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngpm/gbdx"
)

func orderStatus(cmd *cobra.Command, args []string) (err error) {

	// Aquire an Api.
	var profile GBDXProfile
	if err = viper.Unmarshal(&profile); err != nil {
		return err
	}
	api, err := gbdx.NewApi(profile.ActiveConfig)
	if err != nil {
		return err
	}

	order, err := api.OrderStatus(args[0])
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
var orderStatusCmd = &cobra.Command{
	Use:   "status [ORDERID]",
	Short: "Check the status of a GBDX order",
	Long:  `Check the status of a GBDX order.`,
	RunE:  orderStatus,
}

func init() {
	orderCmd.AddCommand(orderStatusCmd)
}
