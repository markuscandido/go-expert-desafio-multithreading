package cep

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/markuscandido/go-expert-desafio-multithreading/cep/provider"
)

// GetCepData orquestra a busca concorrente de dados de CEP usando provedores.
func GetCepData(cepInput string, providers ...provider.Provider) (string, error) {
	cep, err := ValidateCEP(cepInput)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan *provider.APIResponse, len(providers))
	var wg sync.WaitGroup

	wg.Add(len(providers))
	for _, p := range providers {
		go func(p provider.Provider) {
			defer wg.Done()
			resp, err := p.Fetch(ctx, cep)
			if err != nil {
				log.Printf("%s: Erro ao buscar CEP: %v", p.GetSourceName(), err)
				return
			}
			ch <- resp
		}(p)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case res, ok := <-ch:
		if ok {
			return fmt.Sprintf("Resposta recebida de %s (mais rápida):\n%s", res.Source, res.Payload), nil
		}
		return "", ErrNoValidResponse
	case <-ctx.Done():
		return "", ErrTimeout
	}
}

// ValidateCEP valida e formata uma string de CEP.
func ValidateCEP(cep string) (string, error) {
	re := regexp.MustCompile(`[^0-9]`)
	cep = re.ReplaceAllString(cep, "")
	if len(cep) != 8 {
		return "", ErrInvalidCEP
	}
	return cep, nil
}
