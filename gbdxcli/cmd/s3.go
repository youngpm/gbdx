// Copyright © 2016 Peter Schmitt peter.schmitt@digitalglobe.com
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

	"github.com/spf13/cobra"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Interface to S3 Storage Service",
	Long: `GBDX stores ancillary data and derived products in an Amazon Web Services (AWS) S3 bucket. When a workflow is run on the GBDX platform, a task called "StageDataToS3" is typically the last task run by the workflow. This task takes the processed data and places it in the AWS S3 bucket.

The purpose of the GBDX S3 Storage Service is to allow users to access this data. The service provides the temporary credentials required to access a Prefix, Folder, or Object in the S3 bucket.

For more details, see https://gbdxdocs.digitalglobe.com/docs/s3-storage-service-course`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("s3 called")
	},
}

func init() {
	RootCmd.AddCommand(s3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
