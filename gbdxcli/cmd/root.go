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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gbdxcli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Ensure $HOME/.gbdx is created.  We write config as well as store persistent goods here, so it must exist.
	if err := ensureGBDXDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed ensuring existence of $HOME/.gbdx directory\n")
		os.Exit(1)
	}

	// Set configuration variables.
	viper.SetEnvPrefix("gbdx")
	viper.SetDefault("username", "")
	viper.SetDefault("password", "")
	viper.SetDefault("clientID", "")
	viper.SetDefault("clientSecret", "")

	// Where to find the configuration.
	viper.SetConfigName("config")                    // name of gbdx config file (without extension)
	viper.AddConfigPath(path.Join("$HOME", ".gbdx")) // adding gbdx directory as first search path
	viper.AutomaticEnv()                             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// GBDXConfiguration holds the various configuation items we need to interact with GBDX.
type GBDXConfiguration struct {
	username       string
	password       string
	clientID       string
	clientPassword string
}
