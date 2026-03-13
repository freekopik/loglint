package loglint

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Settings struct {
		AllowedSymbols string   `yaml:"allowed_symbols"`
		ExcludeFiles   []string `yaml:"exclude_files"`
	} `yaml:"settings"`

	SensitiveWords map[string]string `yaml:"sensitive_words"`
	TargetMethods  []string          `yaml:"target_methods"`
}

var configPath string

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", ".loglint.yaml", "path to configuration file")
}

func getConfig() Config {
	cfg := Config{}
	cfg.SensitiveWords = map[string]string{
		"password": "user authenticated successfully",
		"token":    "token validated",
		"api_key":  "api request completed",
		"secret":   "sensitive data hidden",
	}
	cfg.TargetMethods = []string{"Info", "Error", "Warn", "Debug", "Fatal", "Panic", "Print", "Printf"}

	if configPath == "" {
		return cfg
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg
	}

	var userCfg Config
	if err := yaml.Unmarshal(data, &userCfg); err != nil {
		return cfg
	}

	if userCfg.Settings.AllowedSymbols != "" {
		cfg.Settings.AllowedSymbols = userCfg.Settings.AllowedSymbols
	}

	if len(userCfg.Settings.ExcludeFiles) > 0 {
		cfg.Settings.ExcludeFiles = userCfg.Settings.ExcludeFiles
	}

	if len(userCfg.TargetMethods) > 0 {
		cfg.TargetMethods = userCfg.TargetMethods
	}

	for k, v := range userCfg.SensitiveWords {
		if v == "" {
			v = "sensitive data redacted"
		}
		cfg.SensitiveWords[k] = v
	}

	return cfg
}
