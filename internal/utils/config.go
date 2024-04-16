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
	OutDir                         string `mapstructure:"OUT_DIR"`
	Subkill3rWorkerCount           int    `mapstructure:"SUBKILL3R_WORKER_COUNT"`
	Subkill3rServerAddr            string `mapstructure:"SUBKILL3R_SERVER_ADDR"`
	Subkill3rWordlist              string `mapstructure:"SUBKILL3R_WORDLIST"`
	PurednsWordlist                string `mapstructure:"PUREDNS_WORDLIST"`
	PurednsResolvers               string `mapstructure:"PUREDNS_RESOLVERS"`
	PurednsNumOfThreads            int    `mapstructure:"PUREDNS_NUM_OF_THREADS"`
	GotatorPermlist                string `mapstructure:"GOTATOR_PERMLIST"`
	GotatorDepth                   int    `mapstructure:"GOTATOR_DEPTH"`
	GotatorNumbers                 int    `mapstructure:"GOTATOR_NUMBERS"`
	GotatorNumOfThreads            int    `mapstructure:"GOTATOR_NUM_OF_THREADS"`
	GotatorMindup                  bool   `mapstructure:"GOTATOR_MINDUP"`
	GotatorAdv                     bool   `mapstructure:"GOTATOR_ADV"`
	GotatorMd                      bool   `mapstructure:"GOTATOR_MD"`
	SubfinderNumOfThreads          int    `mapstructure:"SUBFINDER_NUM_OF_THREADS"`
	AmassTimeout                   int    `mapstructure:"AMASS_TIMEOUT"`
	GowitnessTimeout               int    `mapstructure:"GOWITNESS_TIMEOUT"`
	GowitnessResolutionX           int    `mapstructure:"GOWITNESS_RESOLUTION_X"`
	GowitnessResolutionY           int    `mapstructure:"GOWITNESS_RESOLUTION_Y"`
	GowitnessNumOfThreads          int    `mapstructure:"GOWITNESS_NUM_OF_THREADS"`
	GowitnessFullpage              bool   `mapstructure:"GOWITNESS_FULLPAGE"`
	GowitnessScreenshotFilter      bool   `mapstructure:"GOWITNESS_SCREENSHOT_FILTER"`
	GowitnessScreenshotFilterCodes string `mapstructure:"GOWITNESS_SCREENSHOT_FILTER_CODES"`
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
	puredns_wordlist, err := ExtractEmbeddedFileToTempDir(docFS, "docs/subdomains-top-20k.txt", "subdomains-top-20k.txt")
	if err != nil {
		log.Panic(err)
	}

	// Set the default resolvers for puredns
	puredns_resolvers, err := ExtractEmbeddedFileToTempDir(docFS, "docs/resolvers.txt", "resolvers.txt")
	if err != nil {
		log.Panic(err)
	}

	// Set the default permutation list for gotator
	gotator_permlist, err := ExtractEmbeddedFileToTempDir(docFS, "docs/permlist.txt", "permlist.txt")
	if err != nil {
		log.Panic(err)
	}

	// Setting default values

	// main configs
	viper.SetDefault("OUT_DIR", defaultDir)

	// PASSIVE_ENUM configs

	// subfinder configs
	viper.SetDefault("SUBFINDER_NUM_OF_THREADS", 100)

	// amass configs
	viper.SetDefault("AMASS_TIMEOUT", 1)

	// subkill3r configs
	viper.SetDefault("SUBKILL3R_WORDLIST", subkill3r_wordlist)
	viper.SetDefault("SUBKILL3R_WORKER_COUNT", 1000)
	viper.SetDefault("SUBKILL3R_SERVER_ADDR", "8.8.8.8:53")

	// ACTIVE_ENUM configs

	// puredns configs
	viper.SetDefault("PUREDNS_WORDLIST", puredns_wordlist)
	viper.SetDefault("PUREDNS_RESOLVERS", puredns_resolvers)
	viper.SetDefault("PUREDNS_NUM_OF_THREADS", 100)

	// gotator configs
	viper.SetDefault("GOTATOR_PERMLIST", gotator_permlist)
	viper.SetDefault("GOTATOR_DEPTH", 1)
	viper.SetDefault("GOTATOR_NUMBERS", 3)
	viper.SetDefault("GOTATOR_NUM_OF_THREADS", 100)
	viper.SetDefault("GOTATOR_MINDUP", false)
	viper.SetDefault("GOTATOR_ADV", false)
	viper.SetDefault("GOTATOR_MD", false)

	// WEB_OPS configs

	// gowitness configs
	viper.SetDefault("GOWITNESS_TIMEOUT", 10)
	viper.SetDefault("GOWITNESS_RESOLUTION_X", 1440)
	viper.SetDefault("GOWITNESS_RESOLUTION_Y", 900)
	viper.SetDefault("GOWITNESS_NUM_OF_THREADS", 4)
	viper.SetDefault("GOWITNESS_FULLPAGE", false)
	viper.SetDefault("GOWITNESS_SCREENSHOT_FILTER", false)

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
