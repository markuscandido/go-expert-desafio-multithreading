package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"CEP válido", "01001000", "01001000", false},
		{"CEP com máscara", "01001-000", "01001000", false},
		{"CEP com caracteres", "abc01001000def", "01001000", false},
		{"CEP curto", "12345", "", true},
		{"CEP longo", "123456789", "", true},
		{"CEP vazio", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateCEP(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCEP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateCEP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchBrasilAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"cep": "01001-000"}`)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch := make(chan APIResponse, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go fetchBrasilAPI(ctx, &wg, server.URL, ch)

	wg.Wait()
	close(ch)

	select {
	case res, ok := <-ch:
		if !ok {
			t.Fatal("O canal foi fechado inesperadamente")
		}
		if res.Source != "BrasilAPI" {
			t.Errorf("Fonte esperada 'BrasilAPI', mas foi '%s'", res.Source)
		}
		if res.Payload != `{"cep": "01001-000"}`+"\n" {
			t.Errorf("Payload inesperado: %s", res.Payload)
		}
	case <-ctx.Done():
		t.Fatal("Timeout não deveria ter ocorrido")
	}
}

func TestFetchViaCEP_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch := make(chan APIResponse, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go fetchViaCEP(ctx, &wg, server.URL, ch)

	// Goroutine para fechar o canal quando o WaitGroup terminar
	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case res, ok := <-ch:
		if ok {
			t.Fatalf("Não deveria receber resposta no canal em caso de timeout, mas recebeu: %+v", res)
		}
		// Se !ok, o canal foi fechado, o que é esperado após o wg.Wait().
	case <-time.After(300 * time.Millisecond): // Um tempo maior que o timeout
		t.Fatal("O teste não terminou como esperado")
	}

	// Verifica se o contexto foi de fato cancelado
	if ctx.Err() == nil {
		t.Error("O contexto deveria ter sido cancelado por timeout, mas não foi")
	}
}
