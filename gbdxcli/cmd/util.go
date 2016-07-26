package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
	"github.com/youngpm/gbdx"
)

// ensureGBDXDir will create the gbdx directory if it doesn't already exist.
func ensureGBDXDir() (string, error) {
	gbdxPath := path.Join(userHomeDir(), ".gbdx")
	err := os.MkdirAll(gbdxPath, 0700)
	return gbdxPath, err
}

// userHomeDir returns the home directory of the user.  I've borrowed
// this from https://github.com/spf13/viper/blob/master/util.go .
func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func apiFromConfig() (*gbdx.Api, error) {
	var profile GBDXProfile
	if err := viper.Unmarshal(&profile); err != nil {
		return nil, err
	}
	return gbdx.NewApi(profile.ActiveConfig)
}

// cacheToken updates an existing configuration file with the
// provided one.  Note that we only update the profile as stored in
// viper.
func cacheToken(api *gbdx.Api) error {

	// Read in configuration file if it exists.
	confFile := viper.ConfigFileUsed()
	profilesOut := make(map[string]gbdx.Config)
	_, err := toml.DecodeFile(confFile, &profilesOut)
	if err != nil {
		return fmt.Errorf("failed to parse the configurtion: %v", err)
	}

	// Update/add the token.
	c := profilesOut[viper.GetString("profile")]
	c.Token, err = api.Token()
	if err != nil {
		return fmt.Errorf("failed to fetch a token to cache: %v", err)
	}
	profilesOut[viper.GetString("profile")] = c

	// Save to the GBDX config file.
	file, err := os.Create(confFile)
	if err != nil {
		return fmt.Errorf("failed to write updated configuration to disk: %v", err)
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(profilesOut)
}
