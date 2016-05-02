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
	"github.com/youngpm/gbdx"
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
	// Make sure the GBDX directory exists.
	gbdxPath, err := ensureGBDXDir()
	if err != nil {
		return fmt.Errorf("configure failed to create the gbdx directory: %v", err)
	}

	// Load the existing profile, if there is one.
	var profile GBDXProfile
	if err = viper.Unmarshal(&profile); err != nil {
		return fmt.Errorf("configure failed to parse the configuration: %v", err)
	}

	// Get the configuration from the command line.
	var configVars = []struct {
		prompt   string
		val      *string
		isSecret bool
	}{
		{"GBDX User Name", &profile.ActiveConfig.Username, false},
		{"GBDX Password", &profile.ActiveConfig.Password, true},
		{"GBDX Client ID", &profile.ActiveConfig.ClientID, true},
		{"GBDX Client Secret", &profile.ActiveConfig.ClientSecret, true},
	}
	for _, configVar := range configVars {
		// Pretty print the prompt for this variable.
		fmt.Printf(configVar.prompt)
		if val := *configVar.val; len(val) > 0 {
			if configVar.isSecret {
				fmt.Printf(" [%s]", secretString(val[max(0, len(val)-10):]))
			} else {
				fmt.Printf(" [%s]", val)
			}
		}
		fmt.Printf(": ")

		// Get user input for this value.
		var s string
		if n, err := fmt.Scanln(&s); err != nil && n > 0 {
			// Gobble up remaining tokens if any.
			for n, err := fmt.Scanln(&s); err != nil && n > 0; {
			}
			return fmt.Errorf("your input is bogus: %v", err)

		}
		if len(s) > 0 {
			*configVar.val = s
		}
	}

	// Read in configuration file if it exists.
	var confFile string
	profilesOut := make(map[string]gbdx.Config)
	if confFile = viper.ConfigFileUsed(); len(confFile) > 0 {
		_, err = toml.DecodeFile(confFile, &profilesOut)
		if err != nil {
			return fmt.Errorf("configure failed to parse the configurtion: %v", err)
		}
	} else {
		confFile = path.Join(gbdxPath, "credentials.toml")
	}

	// Update/add the new profile.
	profilesOut[viper.GetString("profile")] = profile.ActiveConfig

	// Save to the GBDX config file.
	file, err := os.Create(confFile)
	if err != nil {
		return fmt.Errorf("configure failed to create on disk the updated configuration: %v", err)
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(profilesOut)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Stores your GBDX configuration",
	Long:  `Store your GBDX configuration in the .gbdx home directory.`,
	RunE:  configure,
}

func init() {
	RootCmd.AddCommand(configureCmd)
}
