package cep

import "errors"

var (
	// ErrInvalidCEP é retornado quando o CEP fornecido é inválido.
	ErrInvalidCEP = errors.New("CEP inválido. Deve conter 8 dígitos numéricos")
	// ErrTimeout é retornado quando nenhuma API responde dentro do tempo limite.
	ErrTimeout = errors.New("timeout. Nenhum provedor respondeu em 1 segundo")
	// ErrNoValidResponse é retornado quando nenhum dos provedores consegue retornar uma resposta válida.
	ErrNoValidResponse = errors.New("nenhum provedor retornou um resultado válido")
)
