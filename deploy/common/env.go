package common

import "github.com/caarlos0/env/v11"

func ParseConfig[T any]() (*T, error) {
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
