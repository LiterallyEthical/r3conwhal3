package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)


type Config struct {
	OutDir string `mapstructure:"OUT_DIR"`
	FileName string `mapstructure:"FILE_NAME"`
	Subkill3rWorkerCount string `mapstructure:"SUBKILL3R_WORKER_COUNT"`
	Subkill3rServerAddr string `mapstructure:"SUBKILL3R_SERVER_ADDR"`
	Subkill3rWordlist string `mapstructure:"SUBKILL3R_WORDLIST"`	
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	
	viper.AutomaticEnv()

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Couldn't find user's home directory, %v", err))
		return
	}
	
	// Set the default value for outDir
	defaultDir := filepath.Join(homeDir, "r3conwhal3", "results")
	
	// Setting default values
	viper.SetDefault("OUT_DIR", defaultDir)
	viper.SetDefault("FILE_NAME", "subdomains.txt")	
	viper.SetDefault("SUBKILL3R_WORDLIST", "")



	err = viper.ReadInConfig()
	if err != nil {
		return
	}


	err = viper.Unmarshal(&config)

	return
}