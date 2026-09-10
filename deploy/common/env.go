package common

import "github.com/chihqiang/tlsctl/pkg/envconfig"

func ParseConfig[T any]() (*T, error) {
	return envconfig.ParseConfig[T]()
}
