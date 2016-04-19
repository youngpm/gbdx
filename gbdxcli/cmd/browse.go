// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngpm/gbdx"
)

var validBrowseDims = map[string]bool{
	"small":  true,
	"medium": true,
	"large":  true,
	"natres": true,
}

func browse(cmd *cobra.Command, args []string) (err error) {
	// Deal with positional arguments.
	var f *os.File
	switch len(args) {
	case 1:
		f = os.Stdout
	case 2:
		f, err = os.Create(args[1])
		if err != nil {
			return fmt.Errorf("Can't create %s: %v", args[1], err)
		}
		defer f.Close()
	default:
		return fmt.Errorf("one or two positional arguments must be provided")
	}

	// If metadata requested, return that.
	if viper.GetBool("metadata") {
		return gbdx.BrowseMetadata(args[0], f)
	}

	// Validate the dimension that has been requested.
	dim := viper.GetString("dim")
	if !validBrowseDims[dim] {
		return fmt.Errorf("requested browse dimension of %q is invalid", dim)
	}

	// Get the browse.
	return gbdx.Browse(args[0], dim, viper.GetBool("json"), f)
}

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse [cid] [outfile]",
	Short: "Download a browse",
	Long: `Download browse imagery.

You must provide [cid]; the catalog id of the browse png you want.
outfile is optional; if not provided, the file will be streamed to
stdout.`,
	RunE: browse,
}

func init() {
	RootCmd.AddCommand(browseCmd)

	browseCmd.Flags().String("dim", "small", "dimension of browse (\"small\", \"medium\", \"large\", or \"natres\")")
	browseCmd.Flags().Bool("json", false, "Return JSON containing paths to both the image and metadata of the browse")
	browseCmd.Flags().Bool("metadata", false, "Return GeoJSON metadata describing the browse rather than the browse itself")
	viper.BindPFlag("dim", browseCmd.Flags().Lookup("dim"))
	viper.BindPFlag("json", browseCmd.Flags().Lookup("json"))
	viper.BindPFlag("metadata", browseCmd.Flags().Lookup("metadata"))

}
