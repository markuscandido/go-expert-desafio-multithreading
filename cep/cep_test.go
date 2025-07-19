package cep_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/markuscandido/go-expert-desafio-multithreading/cep"
	"github.com/markuscandido/go-expert-desafio-multithreading/cep/provider"
	"github.com/stretchr/testify/assert"
)

func TestGetCepData_FastestWins(t *testing.T) {
	fastProvider := provider.NewMockProvider("FastProvider", func(ctx context.Context, cep string) (*provider.APIResponse, error) {
		return &provider.APIResponse{Source: "FastProvider", Payload: "fast response"}, nil
	}, 10*time.Millisecond)

	slowProvider := provider.NewMockProvider("SlowProvider", func(ctx context.Context, cep string) (*provider.APIResponse, error) {
		return &provider.APIResponse{Source: "SlowProvider", Payload: "slow response"}, nil
	}, 100*time.Millisecond)

	result, err := cep.GetCepData("12345678", slowProvider, fastProvider)

	assert.NoError(t, err)
	assert.Contains(t, result, "FastProvider")
	assert.Contains(t, result, "fast response")
}

func TestGetCepData_Timeout(t *testing.T) {
	slowProvider := provider.NewMockProvider("SlowProvider", func(ctx context.Context, cep string) (*provider.APIResponse, error) {
		return &provider.APIResponse{Source: "SlowProvider", Payload: "slow response"}, nil
	}, 1500*time.Millisecond) // Exceeds the 1-second timeout

	_, err := cep.GetCepData("12345678", slowProvider)

	assert.Error(t, err)
	assert.Equal(t, cep.ErrTimeout, err)
}

func TestGetCepData_InvalidCep(t *testing.T) {
	_, err := cep.GetCepData("12345")
	assert.Error(t, err)
	assert.Equal(t, cep.ErrInvalidCEP, err)
}

func TestGetCepData_NoProviders(t *testing.T) {
	result, err := cep.GetCepData("12345678") // No providers given
	assert.Error(t, err)
	assert.Equal(t, cep.ErrNoValidResponse, err)
	assert.Equal(t, "", result)
}

func TestGetCepData_AllProvidersFail(t *testing.T) {
	failingProvider1 := provider.NewMockProvider("Failing1", func(ctx context.Context, cep string) (*provider.APIResponse, error) {
		return nil, fmt.Errorf("provider 1 failed")
	}, 10*time.Millisecond)

	failingProvider2 := provider.NewMockProvider("Failing2", func(ctx context.Context, cep string) (*provider.APIResponse, error) {
		return nil, fmt.Errorf("provider 2 failed")
	}, 20*time.Millisecond)

	_, err := cep.GetCepData("12345678", failingProvider1, failingProvider2)

	assert.Error(t, err)
	assert.Equal(t, cep.ErrNoValidResponse, err)
}

func TestGetCepData_OneProviderFails(t *testing.T) {
	successProvider := provider.NewMockProvider("SuccessProvider", func(ctx context.Context, cep string) (*provider.APIResponse, error) {
		return &provider.APIResponse{Source: "SuccessProvider", Payload: "success"}, nil
	}, 20*time.Millisecond)

	failingProvider := provider.NewMockProvider("FailingProvider", func(ctx context.Context, cep string) (*provider.APIResponse, error) {
		return nil, fmt.Errorf("provider failed")
	}, 10*time.Millisecond)

	result, err := cep.GetCepData("12345678", failingProvider, successProvider)

	assert.NoError(t, err)
	assert.Contains(t, result, "SuccessProvider")
}

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		name    string
		cep     string
		want    string
		wantErr bool
	}{
		{"Valid CEP", "12345-678", "12345678", false},
		{"Valid CEP with spaces", " 12345678 ", "12345678", false},
		{"Invalid CEP - too short", "12345", "", true},
		{"Invalid CEP - too long", "123456789", "", true},
		{"Invalid CEP - with letters", "12345-abc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cep.ValidateCEP(tt.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCEP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateCEP() = %v, want %v", got, tt.want)
			}
		})
	}
}
