package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	once sync.Once
	cfg  *Config // singleton 
)

var (
	DEBUG      = false
	configPath = ""
)

type Config struct {
	TermoTunePath     string        `json:"termo_tune_path"`     // path to termo tune
	PathYTDownloaded      string        `json:"path_yt_downloaded"`    // path to Youtube downloaded files
	PathFFmpeg    string        `json:"path_ffmpeg"`    // path to ffmpeg
	PathFFprobe   string        `json:"path_ffprobe"`   // path to ffprobe
	SearchTimeout time.Duration `json:"search_timeout"` // search timeout
	Theme         string        `json:"theme"`          // UI theme
	DBPath        string        `json:"db_path"`        // path to the database
	DiscordRPC    bool          `json:"discord_rpc"`    // Discord Rich Presence
	LogFile       string        `json:"log_file"`       // path to the log file
	ServerPort    string        `json:"server_port"`    // port to run the server on
}

func MergeConfig(configs, defaultConfig * Config) *Config {
	if configs.TermoTunePath == "" {
		configs.TermoTunePath = defaultConfig.TermoTunePath
	}
	if configs.PathYTDownloaded == "" {
		configs.PathYTDownloaded = defaultConfig.PathYTDownloaded
	}
	if configs.PathFFmpeg == "" {
		configs.PathFFmpeg = defaultConfig.PathFFmpeg
	}
	if configs.PathFFprobe == "" {
		configs.PathFFprobe = defaultConfig.PathFFprobe
	}					
	if configs.SearchTimeout == 0 {
		configs.SearchTimeout = defaultConfig.SearchTimeout
	}
	if configs.Theme == "" {
		configs.Theme = defaultConfig.Theme
	}
	if configs.DBPath == "" {
		configs.DBPath = defaultConfig.DBPath
	}
	if !configs.DiscordRPC {
		configs.DiscordRPC = defaultConfig.DiscordRPC
	}
	if configs.LogFile == "" {
		configs.LogFile = defaultConfig.LogFile
	}
	if configs.ServerPort == "" {
		configs.ServerPort = defaultConfig.ServerPort
	}
	
	return configs
}

func InitConfig() (*Config, error) {
	termotune_path := os.Getenv("TERMOTUNE_PATH")
	config_path := os.Getenv("CONFIG_PATH") 
	
	var config *Config 


	var defaultConfig = &Config{
		TermoTunePath:     termotune_path,
		PathYTDownloaded:  "youtube_downloaded",
		PathFFmpeg:        "ffmpeg",
		PathFFprobe:       "ffprobe",
		SearchTimeout:     60 * time.Second,
		Theme:             "default",    	
		DBPath:            filepath.Join(termotune_path, "db.json"),
		DiscordRPC:        true,
		LogFile:           filepath.Join(termotune_path, "log.txt"),
		ServerPort:        "8080",	   			

	}

	if jsonFile, err := os.ReadFile(config_path); err == nil {
		config = &Config{}
		if err = json.Unmarshal(jsonFile, config); err == nil {
			config = MergeConfig(config, defaultConfig)
			return config, nil
		} else {
			fmt.Println("Error loading config file:", err)
			return nil, errors.New("error loading config file: " + err.Error())
		}
	}
	return defaultConfig, nil
}

func GetConfig() (*Config){
	once.Do(func() {
		cfg, _ = InitConfig()
		})
		return cfg
}


func EditConfigField(field, value string) error {
	config := GetConfig()
	switch field {
	case "termotune_path":
		config.TermoTunePath = value
	case "path_yt_downloaded":
		config.PathYTDownloaded = value
	case "path_ffmpeg":
		config.PathFFmpeg = value
	case "path_ffprobe":
		config.PathFFprobe = value
	case "search_timeout":
		if duration, err := time.ParseDuration(value); err == nil {
			config.SearchTimeout = duration
		} else {
			return err
		}
	case "theme":
		config.Theme = value
	case "db_path":
		config.DBPath = value
	case "discord_rpc":
		if value == "true" {
			config.DiscordRPC = true
		} else if value == "false" {
			config.DiscordRPC = false
		}
	case "log_file":
		config.LogFile = value
	case "server_port":
		config.ServerPort = value
	default:
		return errors.New("unknown field: " + field)
	}

	return saveConfig(config)
}

func saveConfig(config *Config) error {
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, jsonData, 0o644)
}


