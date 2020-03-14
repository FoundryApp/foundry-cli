package main

import (
	"log"
	"foundry/cli/cmd"
	// "foundry/cli/auth"
	// "context"
	"os"

	"github.com/spf13/viper"
)

func init() {
	// Remove timestamp prefix
	log.SetFlags(0)
}

func main() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to get a config dir", err)
	}

	dirPath := configDir + "/foundrycli"
	confName := "config"
	ext := "json"
	fullPath := dirPath + "/" + confName + "." + ext

	viper.SetConfigName(confName)
	viper.SetConfigType(ext)
	viper.AddConfigPath(dirPath)

	log.Println(dirPath);

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)

		f, err := os.Create(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.WriteString("{}")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()

	// a := auth.New()
	// a.SignIn(context.TODO(), "vasek@foundryapp.co", "123456")

	// log.Println(a)
}
