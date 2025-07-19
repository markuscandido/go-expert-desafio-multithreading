package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markuscandido/go-expert-desafio-multithreading/cep"
	"github.com/markuscandido/go-expert-desafio-multithreading/cep/provider"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Uso: go run main.go <cep>")
	}

	providers := []provider.Provider{
		provider.NewBrasilAPIProvider(),
		provider.NewViaCEPProvider(),
	}

	result, err := cep.GetCepData(os.Args[1], providers...)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
	fmt.Println(result)
}
