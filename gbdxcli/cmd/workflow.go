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
	"fmt"

	"github.com/spf13/cobra"
)

func get_workflow(cmd *cobra.Command, args []string) (err error) {
	fmt.Printf("found workflow")

	// Aquire an Api.
	api, err := apiFromConfig()
	if err != nil {
		return err
	}

	return cacheToken(api)
}

// workflowCmd represents a workflow search
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Search for workflows using GBDX",
	Long: `Search for workflows using GBDX

Workflow IDs can be specified as space delimited arguments on the
command line, or if given no arguments, passed in delimited by
newlines via stdin.
`,
	RunE: get_workflow,
}

func init() {
	RootCmd.AddCommand(workflowCmd)
}
