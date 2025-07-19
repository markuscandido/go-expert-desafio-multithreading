# Desafio Go Expert - Multithreading

Este projeto é uma solução para o desafio de multithreading do curso Go Expert. O objetivo é criar uma aplicação que consulta duas APIs de CEP diferentes simultaneamente e retorna a resposta daquela que for mais rápida, com um timeout global de 1 segundo.

## Funcionalidades

- **Busca Concorrente**: Utiliza goroutines para consultar a [BrasilAPI](https://brasilapi.com.br/) e a [ViaCEP](http://viacep.com.br/) ao mesmo tempo.
- **Timeout Global**: Implementa um `context.WithTimeout` para garantir que a aplicação não espere mais de 1 segundo por uma resposta.
- **Validação de CEP**: O CEP fornecido como entrada é validado para garantir que contém exatamente 8 dígitos numéricos.
- **Testes Unitários**: Cobertura de testes para a lógica de validação e para as funções de busca, utilizando servidores de teste (`httptest`) para simular as APIs.

## Tecnologias e Conceitos Utilizados

- **Go (Golang)**
- **Goroutines** para concorrência.
- **Channels** para comunicação entre goroutines.
- **`context`** para controle de timeout e cancelamento de requisições.
- **`sync.WaitGroup`** para sincronização.
- **`net/http`** para realizar as requisições HTTP.
- **`httptest`** para criar mocks de servidores nos testes.
- **Módulos Go** para gerenciamento de dependências.

## Como Executar

### Pré-requisitos

- Go 1.18 ou superior.

### Execução

1. Clone o repositório:
   ```sh
   git clone https://github.com/markuscandido/go-expert-desafio-multithreading.git
   cd go-expert-desafio-multithreading
   ```

2. Execute a aplicação passando um CEP como argumento:
   ```sh
   go run main.go 01001000
   ```

   Você pode usar um CEP com ou sem formatação:
   ```sh
   go run main.go 01001-000
   ```

### Executando os Testes

Para rodar a suíte de testes, execute o comando:

```sh
go test -v
```
