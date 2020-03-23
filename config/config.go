package config

import (
	"os"
	"github.com/spf13/viper"
)

func Init() error {
	// /Users/vasekmlejnsky/Library/Application Support/foundrycli
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}


	// configDir = "/Users/vasekmlejnsky/Developer"

	dirPath := configDir + "/foundrycli"
	confName := "config"
	ext := "json"
	fullPath := dirPath + "/" + confName + "." + ext

	viper.SetConfigName(confName)
	viper.SetConfigType(ext)
	viper.AddConfigPath(dirPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)

		f, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer f.Close()
		f.WriteString("{}")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func Set(key string, val interface{}) {
	viper.Set(key, val)
}

func Get(key string) interface{} {
	val := viper.Get(key)
	return val
}

func Write() error {
	return viper.WriteConfig()
}
