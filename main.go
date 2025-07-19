package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

// APIResponse é uma estrutura unificada para as respostas das APIs.
type APIResponse struct {
	Source  string
	Payload string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Uso: go run main.go <cep>")
	}
	cep, err := validateCEP(os.Args[1])
	if err != nil {
		log.Fatalf("Erro de validação: %v", err)
	}

	// Cria um contexto com timeout de 1 segundo.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Garante que o contexto seja cancelado e recursos liberados.

	ch := make(chan APIResponse, 2) // Canal bufferizado para não bloquear goroutines.
	var wg sync.WaitGroup

	wg.Add(2)
	go fetchBrasilAPI(ctx, &wg, "https://brasilapi.com.br/api/cep/v1/"+cep, ch)
	go fetchViaCEP(ctx, &wg, "http://viacep.com.br/ws/"+cep+"/json/", ch)

	// Espera as goroutines terminarem e fecha o canal.
	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case res, ok := <-ch:
		if ok {
			fmt.Printf("Resposta recebida de %s (mais rápida):\n%s\n", res.Source, res.Payload)
		} else {
			// Isso não deve acontecer devido ao timeout, mas é uma boa prática.
			log.Println("Nenhuma resposta recebida e o canal foi fechado.")
		}
	case <-ctx.Done():
		log.Println("Erro: Timeout. Nenhuma API respondeu em 1 segundo.")
	}
}

func fetchBrasilAPI(ctx context.Context, wg *sync.WaitGroup, url string, ch chan<- APIResponse) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("BrasilAPI: Erro ao criar requisição: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Verifica se o erro foi causado pelo cancelamento do contexto.
		if ctx.Err() != nil {
			log.Println("BrasilAPI: Requisição cancelada.")
		} else {
			log.Printf("BrasilAPI: Erro na requisição: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("BrasilAPI: Erro ao ler resposta: %v", err)
		return
	}

	// Envia a resposta para o canal.
	ch <- APIResponse{Source: "BrasilAPI", Payload: string(body)}
}

func fetchViaCEP(ctx context.Context, wg *sync.WaitGroup, url string, ch chan<- APIResponse) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("ViaCEP: Erro ao criar requisição: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			log.Println("ViaCEP: Requisição cancelada.")
		} else {
			log.Printf("ViaCEP: Erro na requisição: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ViaCEP: Erro ao ler resposta: %v", err)
		return
	}

	ch <- APIResponse{Source: "ViaCEP", Payload: string(body)}
}

func validateCEP(cep string) (string, error) {
	// Remove qualquer caractere não numérico.
	re := regexp.MustCompile(`[^0-9]`)
	cleanCEP := re.ReplaceAllString(cep, "")

	if len(cleanCEP) != 8 {
		return "", fmt.Errorf("CEP inválido. Deve conter 8 dígitos numéricos")
	}
	return cleanCEP, nil
}
