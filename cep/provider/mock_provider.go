package provider

import (
	"context"
	"time"
)

// MockProvider é uma implementação de Provider para testes.
type MockProvider struct {
	SourceName    string
	FetchFunc     func(ctx context.Context, cep string) (*APIResponse, error)
	ExecutionTime time.Duration
}

// NewMockProvider cria um novo MockProvider.
func NewMockProvider(sourceName string, fetchFunc func(ctx context.Context, cep string) (*APIResponse, error), executionTime time.Duration) *MockProvider {
	return &MockProvider{
		SourceName:    sourceName,
		FetchFunc:     fetchFunc,
		ExecutionTime: executionTime,
	}
}

// Fetch executa a função de busca simulada após um atraso.
func (m *MockProvider) Fetch(ctx context.Context, cep string) (*APIResponse, error) {
	select {
	case <-time.After(m.ExecutionTime):
		return m.FetchFunc(ctx, cep)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetSourceName retorna o nome do provedor mock.
func (m *MockProvider) GetSourceName() string {
	return m.SourceName
}
