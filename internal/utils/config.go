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
	EnableSubkill3r                bool   `mapstructure:"ENABLE_SUBKILL3R"`
	EnableAssetfinder              bool   `mapstructure:"ENABLE_ASSETFINDER"`
	EnableAmass                    bool   `mapstructure:"ENABLE_AMASS"`
	EnableGowitness                bool   `mapstructure:"ENABLE_GOWITNESS"`
	EnableFFUF                     bool   `mapstructure:"ENABLE_FFUF"`
	EnableWebGalery                bool   `mapstructure:"ENABLE_WEB_GALERY"`
	EnableSubzy                    bool   `mapstructure:"ENABLE_SUBZY"`
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
	FFUFNumOfThreads               int    `mapstructure:"FFUF_NUM_OF_THREADS"`
	FFUFMaxtime                    int    `mapstructure:"FFUF_MAXTIME"`
	FFUFRate                       int    `mapstructure:"FFUF_RATE"`
	FFUFTimeout                    int    `mapstructure:"FFUF_TIMEOUT"`
	FFUFWordlist                   string `mapstructure:"FFUF_WORDLIST"`
	FFUFMatchHTTPCode              string `mapstructure:"FFUF_MATCH_HTTP_CODE"`
	FFUFFilterResponseSize         string `mapstructure:"FFUF_FILTER_RESPONSE_SIZE"`
	FFUFOutputFormat               string `mapstructure:"FFUF_OUTPUT_FORMAT"`
	FFUFOutput                     string `mapstructure:"FFUF_OUTPUT"`
	FFUFSF                         bool   `mapstructure:"FFUF_SF"`
	FFUFSE                         bool   `mapstructure:"FFUF_SE"`
	SUBZYConcurrency               int    `mapstructure:"SUBZY_CONCURRENCY"`
	SUBZYTimeout                   int    `mapstructure:"SUBZY_TIMEOUT"`
	SUBZYHideFails                 bool   `mapstructure:"SUBZY_HIDE_FAILS"`
	SUBZYHTTPS                     bool   `mapstructure:"SUBZY_HTTPS"`
	SUBZYVerifySSL                 bool   `mapstructure:"SUBZY_VERIFY_SSL"`
	SUBZYVuln                      bool   `mapstructure:"SUBZY_VULN"`
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

	// Set the default dir fuzzing wordlist for ffuf
	ffuf_wordlist, err := ExtractEmbeddedFileToTempDir(docFS, "docs/common.txt", "commmon.txt")
	if err != nil {
		log.Panic(err)
	}

	// Setting default values

	// main configs
	viper.SetDefault("OUT_DIR", defaultDir)
	viper.SetDefault("ENABLE_WEB_GALERY", true)

	// PASSIVE_ENUM configs
	viper.SetDefault("ENABLE_ASSETFINDER", true)
	viper.SetDefault("ENABLE_AMASS", true)
	viper.SetDefault("ENABLE_SUBKILL3R", true)
	viper.SetDefault("ENABLE_GOWITNESS", true)
	viper.SetDefault("ENABLE_FFUF", true)

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

	// FFUF settings
	viper.SetDefault("FFUF_NUM_OF_THREADS", 40)
	viper.SetDefault("FFUF_MAXTIME", 600)
	viper.SetDefault("FFUF_RATE", 0)
	viper.SetDefault("FFUF_TIMEOUT", 10)
	viper.SetDefault("FFUF_WORDLIST", ffuf_wordlist)
	viper.SetDefault("FFUF_MATCH_HTTP_CODE", "200-299,301,302,307,401,403,405,500")
	viper.SetDefault("FFUF_FILTER_RESPONSE_SIZE", 0)
	viper.SetDefault("FFUF_OUTPUT_FORMAT", "json")
	viper.SetDefault("FFUF_OUTPUT", "ffuf_out")
	viper.SetDefault("FFUF_SF", false)
	viper.SetDefault("FFUF_SE", false)

	// WEB_OPS configs

	// subzy configs
	viper.SetDefault("SUBZY_CONCURRENCY", 10)
	viper.SetDefault("SUBZY_TIMEOUT", 10)
	viper.SetDefault("SUBZY_HIDE_FAILS", false)
	viper.SetDefault("SUBZY_HTTPS", false)
	viper.SetDefault("SUBZY_VERIFY_SSL", false)
	viper.SetDefault("SUBZY_VULN", false)

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
