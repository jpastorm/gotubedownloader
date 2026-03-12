package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	DownloadDir        string `json:"downloadDir"`
	MaxConcurrent      int    `json:"maxConcurrent"`
	FragmentConcurrent int    `json:"fragmentConcurrent"`
	DefaultMode        string `json:"defaultMode"`
	SubtitleLang       string `json:"subtitleLang"`
	SpeedLimit         string `json:"speedLimit"`
	Proxy              string `json:"proxy"`
	SkipDuplicates     bool   `json:"skipDuplicates"`
	ContinueDownloads  bool   `json:"continueDownloads"`
	EmbedMetadata      bool   `json:"embedMetadata"`
	DownloadSubtitles  bool   `json:"downloadSubtitles"`
	EmbedSubtitles     bool   `json:"embedSubtitles"`
	NotifyOnComplete   bool   `json:"notifyOnComplete"`
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
	cfgPath  string
)

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".gotubedownloader")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func GetConfigDir() (string, error) {
	return getConfigDir()
}

func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	defaultDir := filepath.Join(homeDir, "Downloads", "GotubeDownloader")
	return &Config{
		DownloadDir:        defaultDir,
		MaxConcurrent:      2,
		FragmentConcurrent: 4,
		DefaultMode:        "video",
		SubtitleLang:       "en",
		SkipDuplicates:     true,
		ContinueDownloads:  true,
		NotifyOnComplete:   true,
	}
}

func Load() (*Config, error) {
	var loadErr error
	once.Do(func() {
		dir, err := getConfigDir()
		if err != nil {
			loadErr = err
			return
		}
		cfgPath = filepath.Join(dir, "config.json")
		instance = DefaultConfig()
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			if os.IsNotExist(err) {
				loadErr = Save(instance)
				return
			}
			loadErr = err
			return
		}
		if err := json.Unmarshal(data, instance); err != nil {
			loadErr = err
		}
	})
	return instance, loadErr
}

func Save(cfg *Config) error {
	mu.Lock()
	defer mu.Unlock()
	if cfgPath == "" {
		dir, err := getConfigDir()
		if err != nil {
			return err
		}
		cfgPath = filepath.Join(dir, "config.json")
	}
	instance = cfg
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, data, 0644)
}

func Get() *Config {
	mu.RLock()
	defer mu.RUnlock()
	if instance == nil {
		instance = DefaultConfig()
	}
	return instance
}
