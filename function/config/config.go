package config

import "github.com/caarlos0/env"

type Config struct {
	ClientID     string   `env:"CLIENT_ID,required"`
	ClientSecret string   `env:"CLIENT_SECRET,required"`
	Scope        string   `env:"SCOPE"`
	StreamerIds  []string `env:"STREAMER_IDS,required" envSeparator:","`
	Region       string   `env:"AWS_REGION,required"`
}

func Parse() (*Config, error) {
	var conf Config
	if err := env.Parse(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
