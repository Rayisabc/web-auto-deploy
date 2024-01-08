// Author: Rui
// Date: 2023/01/05 16:20
// Description: Configure Reader

package util

import (
	"os"
	"path/filepath"
	"web-auto-deploy/model"

	"github.com/BurntSushi/toml"
)

type CfgUtil struct {
}

func (c *CfgUtil) ReadConfig() model.Config {
	filename := "config.toml"
	execDir, err := os.Executable()
	if err != nil {
		Error.Println("execDir error")
	}

	execDir = filepath.Dir(execDir)
	configFile := filepath.Join(execDir, filename)
	data, err := os.ReadFile(configFile)
	if err != nil {
		Warn.Println("Custom config.toml not found")
		data, err = os.ReadFile("config.toml")
		if err != nil {
			panic(err)
		}
		Warn.Println("Use the default config")
	}

	var config model.Config
	if _, err := toml.Decode(string(data), &config); err != nil {
		Error.Fatalln("Unable to read the config,please check the configuration")
		panic(err)
	}

	return config
}
