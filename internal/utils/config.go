package utils

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	OutDir               string `mapstructure:"OUT_DIR"`
	FileName             string `mapstructure:"FILE_NAME"`
	Subkill3rWorkerCount string `mapstructure:"SUBKILL3R_WORKER_COUNT"`
	Subkill3rServerAddr  string `mapstructure:"SUBKILL3R_SERVER_ADDR"`
	Subkill3rWordlist    string `mapstructure:"SUBKILL3R_WORDLIST"`
	PurednsWordlist      string `mapstructure:"PUREDNS_WORDLIST"`
	PurednsResolvers     string `mapstructure:"PUREDNS_RESOLVERS"`
	PurednsNumOfThreads  int    `mapstructure:"PUREDNS_NUM_OF_THREADS"`
}

func LoadConfig(path string, docFS embed.FS) (config Config, err error) {
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

	// Set the default path for subkill3r wordlist to
	subkill3r_wordlist, err := ExtractEmbeddedFileToTempDir(docFS, "docs/subdomains-1000.txt", "subdomains-1000.txt")
	if err != nil {
		log.Panic(err)
	}

	// Set the default path for puredns wordlist
	puredns_wordlist, err := ExtractEmbeddedFileToTempDir(docFS, "docs/subdomains-top-110k.txt", "subdomains-top-110k.txt")
	if err != nil {
		log.Panic(err)
	}

	puredns_resolvers, err := ExtractEmbeddedFileToTempDir(docFS, "docs/resolvers.txt", "resolvers.txt")
	if err != nil {
		log.Panic(err)
	}

	// Setting default values
	viper.SetDefault("OUT_DIR", defaultDir)
	viper.SetDefault("FILE_NAME", "subdomains.txt")
	viper.SetDefault("SUBKILL3R_WORDLIST", subkill3r_wordlist)
	viper.SetDefault("PUREDNS_WORDLIST", puredns_wordlist)
	viper.SetDefault("PUREDNS_RESOLVERS", puredns_resolvers)

	if path == "embedded" {
		// Use the passed embedded FS to read the config file
		configData, err := docFS.ReadFile("docs/config.env")
		if err != nil {
			return config, fmt.Errorf("failed to read embedded config file: %v", err)
		}

		// Use viper's ReadConfig method to read from the byte slice
		err = viper.ReadConfig(bytes.NewBuffer(configData))
		if err != nil {
			return config, fmt.Errorf("failed to read configuration from embedded data: %v", err)
		}
	} else {
		// Load configuration from a file path
		viper.AddConfigPath(path)
		err = viper.ReadInConfig()
		if err != nil {
			return config, fmt.Errorf("failed to read configuration from path %s: %v", path, err)
		}
	}

	// Unmarshal the read configuraiton into the Config struct
	err = viper.Unmarshal(&config)
	return config, err
}
