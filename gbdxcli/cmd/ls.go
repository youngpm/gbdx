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

func listCustomerBucket(cmd *cobra.Command, args []string) (err error) {
	var profile GBDXProfile
	if err = viper.Unmarshal(&profile); err != nil {
		return err
	}

	api, err := gbdx.NewApi(profile.ActiveConfig)
	if err != nil {
		return err
	}

	// Determine optional PREFIX positional argument
	var prefix string
	if len(args) == 0 {
		prefix = ""
	} else if len(args) == 1 {
		prefix = args[0]
	} else {
		return fmt.Errorf("Expected 0 or 1 arguments, got %d.  Args are %v", len(args), args)
	}

	// Invoke the command!
	err = api.ListBucket(prefix, viper.GetBool("recursive"), os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [user_prefix or NONE]",
	Short: "List objects (files) and common prefixes (sub-directories) under your customer prefix",
	Long:  `List objects (files) and common prefixes (sub-directories) under your customer prefix. Pass in optional user_prefix to restrict the listing to files in sub directories.`,
	RunE:  listCustomerBucket,
}

func init() {
	s3Cmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	lsCmd.Flags().BoolP("recursive", "r", false, "Command is performed  on  all  files  or  objects under the specified directory or prefix.")
	viper.BindPFlag("recursive", lsCmd.Flags().Lookup("recursive"))
}
