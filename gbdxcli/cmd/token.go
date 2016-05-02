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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func printToken(cmd *cobra.Command, args []string) (err error) {

	api, err := apiFromConfig()
	if err != nil {
		return err
	}

	token, err := api.Token()
	if err != nil {
		return err
	}

	result, err := json.Marshal(token)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", result)

	return cacheToken(api)
}

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get a GBDX token",
	Long:  "Returns a new GBDX token.",
	RunE:  printToken,
}

func init() {
	RootCmd.AddCommand(tokenCmd)
}
