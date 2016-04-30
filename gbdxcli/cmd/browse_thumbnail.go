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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngpm/gbdx"
)

var validThumbnailOrienations = map[string]bool{
	"standard":  true,
	"maxwidth":  true,
	"maxheight": true,
	"natres":    true,
}

func thumbnail(cmd *cobra.Command, args []string) (err error) {
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

	// Validate the dimension that has been requested.
	dim := viper.GetInt("dim")
	if dim != 256 && dim != 512 && dim != 1024 {
		return fmt.Errorf("requested thumbnail dimension of %d is invalid", dim)
	}

	// Validate the orientation that has been requested.
	orientation := viper.GetString("orient")
	if !validThumbnailOrienations[orientation] {
		return fmt.Errorf("requested thumbnail orientation of %s is invalid", orientation)
	}

	// Get the thumbnail.
	return gbdx.Thumbnail(args[0], dim, orientation, f)
}

// thumbnailCmd represents the thumbnail command
var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail [CID] [OUT]",
	Short: "Download a thumbnail",
	Long: `Download a thumbnail image.

A thumbnail differs from a browse in that you can explicitly control
the width or height of the image.

You must provide [CID]; the catalog id of the thumbnail png you want.
outfile is optional; if not provided, the file will be streamed to
stdout.`,
	RunE: thumbnail,
}

func init() {
	browseCmd.AddCommand(thumbnailCmd)
	thumbnailCmd.Flags().Int("dim", 256, "dimension of the thumbnail (\"256\", \"512\", or \"1024\")")
	thumbnailCmd.Flags().String("orient", "standard",
		"orientation of thumbnail (\"standard\", \"maxwidth\", \"maxheight\", or \"natres\")")
	viper.BindPFlag("dim", thumbnailCmd.Flags().Lookup("dim"))
	viper.BindPFlag("orient", thumbnailCmd.Flags().Lookup("orient"))
}
