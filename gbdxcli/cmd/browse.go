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

func browse(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("One positional argument must be provided.")
	}
	var f *os.File
	if len(args) == 2 {
		f, err := os.Create(args[1])
		if err != nil {
			return fmt.Errorf("Can't create %s: %v", args[1], err)
		}
		defer f.Close()
	} else {
		f = os.Stdout
	}

	dim := viper.GetString("dim")
	browseSizes := map[string]bool{
		"small":  true,
		"medium": true,
		"large":  true,
		"natres": true,
	}
	if !browseSizes[dim] {
		return fmt.Errorf("Browse dimension of %q is invalid.", dim)
	}
	err := gbdx.Browse(args[0], dim, f)
	return err
}

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: browse,
}

func init() {
	RootCmd.AddCommand(browseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// browseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	browseCmd.Flags().String("dim", "small", "dimension of browse (small, medium, large, or natres)")
	viper.BindPFlag("dim", browseCmd.Flags().Lookup("dim"))

}
