package cmd

import (
	"os"
	"path"
	"runtime"
)

// ensureGBDXDir will create the gbdx directory if it doesn't already exist.
func ensureGBDXDir() (string, error) {
	gbdxPath := path.Join(userHomeDir(), ".gbdx")
	err := os.MkdirAll(gbdxPath, 0600)
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
