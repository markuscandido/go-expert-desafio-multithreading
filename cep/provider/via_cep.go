package provider

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ViaCEPProvider implementa a interface Provider para a ViaCEP.
type ViaCEPProvider struct{}

// NewViaCEPProvider cria uma nova instância de ViaCEPProvider.
func NewViaCEPProvider() *ViaCEPProvider {
	return &ViaCEPProvider{}
}

// Fetch busca dados de CEP na ViaCEP.
func (p *ViaCEPProvider) Fetch(ctx context.Context, cep string) (*APIResponse, error) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição para ViaCEP: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição para ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP retornou status inesperado: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta da ViaCEP: %w", err)
	}

	log.Printf("%s: Resposta recebida com sucesso", p.GetSourceName())
	return &APIResponse{Source: p.GetSourceName(), Payload: string(body)}, nil
}

// GetSourceName retorna o nome do provedor.
func (p *ViaCEPProvider) GetSourceName() string {
	return "ViaCEP"
}
