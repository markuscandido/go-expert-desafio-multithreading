package provider

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

// BrasilAPIProvider implementa a interface Provider para a BrasilAPI.
type BrasilAPIProvider struct{}

// NewBrasilAPIProvider cria uma nova instância de BrasilAPIProvider.
func NewBrasilAPIProvider() *BrasilAPIProvider {
	return &BrasilAPIProvider{}
}

// Fetch busca dados de CEP na BrasilAPI.
func (p *BrasilAPIProvider) Fetch(ctx context.Context, cep string) (*APIResponse, error) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição para BrasilAPI: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição para BrasilAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BrasilAPI retornou status inesperado: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta da BrasilAPI: %w", err)
	}

	log.Printf("%s: Resposta recebida com sucesso", p.GetSourceName())
	return &APIResponse{Source: p.GetSourceName(), Payload: string(body)}, nil
}

// GetSourceName retorna o nome do provedor.
func (p *BrasilAPIProvider) GetSourceName() string {
	return "BrasilAPI"
}
