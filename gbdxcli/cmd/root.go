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
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:  "gbdxcli",
	Long: "A CLI for GBDX.",
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
	// Define global flags.
	RootCmd.PersistentFlags().String("profile", "default", "GBDX profile to use")
	viper.BindPFlag("profile", RootCmd.PersistentFlags().Lookup("profile"))
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// We map a user defined profile from the cli to the active profile.
	viper.RegisterAlias("ActiveConfig", viper.GetString("profile"))

	// Ensure $HOME/.gbdx is created.  We write config as well as store persistent goods here, so it must exist.
	gbdxPath, err := ensureGBDXDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed ensuring existence of %s directory\n", gbdxPath)
		os.Exit(1)
	}

	// Where to find the configuration file.
	viper.SetConfigName("credentials") // name of gbdx config file (without extension)
	viper.AddConfigPath(gbdxPath)      // adding gbdx directory as first search path

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// GBDXConfig holds the various configuation items we need to interact with GBDX.
type GBDXConfig struct {
	Username     string `mapstructure:"gbdx_username" toml:"gbdx_username"`
	Password     string `mapstructure:"gbdx_password" toml:"gbdx_password"`
	ClientID     string `mapstructure:"gbdx_client_id" toml:"gbdx_client_id"`
	ClientSecret string `mapstructure:"gbdx_client_secret" toml:"gbdx_client_secret"`
}

type GBDXProfile struct {
	ActiveConfig GBDXConfig
}
