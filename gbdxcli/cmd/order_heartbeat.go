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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/youngpm/gbdx"
)

// heartbeatCmd represents the heartbeat command
var heartbeatCmd = &cobra.Command{
	Use:   "heartbeat",
	Short: "Check if the GBDX ordering endpoint is alive",
	Long:  `Check if the GBDX ordering endpoint is alive.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := gbdx.OrderingHeartbeat()
		if err != nil {
			fmt.Printf("GBDX ordering is down, reason: %v\n", err)
		}
		fmt.Printf("GBDX ordering is up.\n")
	},
}

func init() {
	orderCmd.AddCommand(heartbeatCmd)
}
