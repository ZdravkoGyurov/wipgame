package config

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Matchmaking Matchmaking `yaml:"matchmaking"`
	Redis       Redis       `yaml:"redis"`
}
type Matchmaking struct {
	Timeout                      time.Duration `yaml:"timeout"`
	Interval                     time.Duration `yaml:"interval"`
	RangeIncrement               int           `yaml:"rangeIncrement"`
	BaseRatingRange              float64       `yaml:"baseRatingRange"`
	BaseRatingRangeDuration      float64       `yaml:"baseRatingRangeDuration"`
	RatingRangeIncrementInterval float64       `yaml:"ratingRangeIncrementInterval"`
	RatingRangeMultiplier        float64       `yaml:"ratingRangeMultiplier"`
}

type Redis struct {
	Address       string `yaml:"address"`
	Password      string `yaml:"password"`
	HashSetName   string `yaml:"hashSetName"`
	SortedSetName string `yaml:"sortedSetName"`
}

const (
	configFileName = "config.yaml"
)

var configDir = os.Getenv("CONFIG_DIR")

func Load() (Config, error) {
	cfg := Config{}

	configFile, err := os.ReadFile(path.Join(configDir, configFileName))
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return cfg, err
	}

	slog.Info(fmt.Sprintf("loaded config: %+v", cfg))

	return cfg, nil
}
