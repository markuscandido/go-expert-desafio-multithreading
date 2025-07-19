package provider

import (
	"context"
)

// APIResponse é uma estrutura unificada para as respostas das APIs.
type APIResponse struct {
	Source  string
	Payload string
}

// Provider define a interface para um provedor de dados de CEP.
type Provider interface {
	Fetch(ctx context.Context, cep string) (*APIResponse, error)
	GetSourceName() string
}
