package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/sashabaranov/go-openai"
)

type Config struct {
	Model   string `toml:"model"`
	Size    string `toml:"size"`
	Quality string `toml:"quality"`
	Style   string `toml:"style"`
	Prompt  string `toml:"prompt"`
}

func LoadFromToml(tomlPath string) (Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(tomlPath, &cfg); err != nil {
		return Config{}, err
	}

	switch cfg.Model {
	case openai.CreateImageModelDallE2:
	case openai.CreateImageModelDallE3:
	default:
		return Config{}, fmt.Errorf("invalid model: %s", cfg.Model)
	}

	switch cfg.Size {
	case openai.CreateImageSize256x256:
	case openai.CreateImageSize512x512:
	case openai.CreateImageSize1024x1024:
	case openai.CreateImageSize1792x1024, openai.CreateImageSize1024x1792:
		if cfg.Model != openai.CreateImageModelDallE3 {
			return Config{}, fmt.Errorf("size %s is only supported by model %s", cfg.Size, openai.CreateImageModelDallE3)
		}
	default:
		return Config{}, fmt.Errorf("invalid size: %s", cfg.Size)
	}

	switch cfg.Quality {
	case openai.CreateImageQualityHD:
	case openai.CreateImageQualityStandard:
	default:
		return Config{}, fmt.Errorf("invalid quality: %s", cfg.Quality)
	}

	switch cfg.Style {
	case openai.CreateImageStyleVivid:
	case openai.CreateImageStyleNatural:
	default:
		return Config{}, fmt.Errorf("invalid style: %s", cfg.Style)
	}

	return cfg, nil
}
