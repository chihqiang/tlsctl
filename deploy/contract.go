package deploy

import (
	"context"
	"github.com/go-acme/lego/v4/certificate"
)

type IDeploy interface {
	WithEnvConfig() error
	Deploy(ctx context.Context, certificate *certificate.Resource) error
}
