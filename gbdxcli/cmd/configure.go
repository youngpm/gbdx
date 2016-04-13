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

	"path"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type secretString string

// String returns secretString types as a string with hidden entries.
func (s secretString) String() (str string) {
	for i, c := range s {
		if i > 3 && len(s)-i < 5 {
			str += string(c)
		} else {
			str += "*"
		}
	}
	return
}

func configure(cmd *cobra.Command, args []string) (err error) {

	var config GBDXConfiguration

	val := viper.GetString("default.gbdx_username")
	fmt.Printf("%q\n", val)
	fmt.Printf("keys are %v\n", viper.AllKeys())
	for _, k := range viper.AllKeys() {
		fmt.Printf("%q: %q\n", k, viper.Get(k))
	}

	return err

	gbdxPath, err := ensureGBDXDir()
	if err != nil {
		return err
	}

	// Get the configuration from the command line.
	var configVars = []struct {
		varName  string
		prompt   string
		val      *string
		isSecret bool
	}{
		{"Username", "GBDX User Name", &config.Username, false},
		{"Password", "GBDX Password", &config.Password, true},
		{"ClientID", "GBDX Client ID", &config.ClientID, false},
		{"ClientPassword", "GBDX Client Password", &config.ClientPassword, true},
	}
	for _, configVar := range configVars {
		// Pretty print the prompt for this variable.
		fmt.Printf(configVar.prompt)
		if val := viper.GetString(configVar.varName); len(val) > 0 {
			if configVar.isSecret {
				fmt.Printf(" [%s]", secretString(val))
			} else {
				fmt.Printf(" [%s]", val)
			}
		}
		fmt.Printf(": ")

		// Read user input into the variable.  If an error is
		// returned, we ignore it and just use the default, which is a
		// nice way to handle the user just hitting enter when wanting to
		// keep the default.
		_, err := fmt.Scanln(configVar.val)
		if err != nil {
			*configVar.val = viper.GetString(configVar.varName)
		}
	}

	// Save to the GBDX config file.
	file, err := os.Create(path.Join(gbdxPath, "configure.toml"))
	if err != nil {
		return err
	}
	defer file.Close()

	enc := toml.NewEncoder(file)
	enc.Indent = ""

	err = enc.Encode(
		&struct {
			Default *GBDXConfiguration `toml:"default"`
		}{&config},
	)
	return err
}

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: configure,
}

func init() {
	RootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
