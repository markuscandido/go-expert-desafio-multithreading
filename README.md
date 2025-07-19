# Desafio Go Expert - Multithreading

Este projeto é uma solução para o desafio de multithreading do curso Go Expert. O objetivo é criar uma aplicação que consulta duas APIs de CEP diferentes simultaneamente e retorna a resposta daquela que for mais rápida, com um timeout global de 1 segundo.

O código foi refatorado para seguir os princípios SOLID, utilizando uma arquitetura de provedores para desacoplar a lógica de negócio das implementações concretas de busca de CEP.

## Funcionalidades

- **Arquitetura Extensível (SOLID)**: Utiliza uma interface `Provider` para permitir a fácil adição de novas fontes de CEP sem alterar a lógica principal.
- **Busca Concorrente**: Utiliza goroutines para consultar múltiplos provedores (atualmente [BrasilAPI](https://brasilapi.com.br/) e [ViaCEP](http://viacep.com.br/)) ao mesmo tempo.
- **Timeout Global**: Implementa um `context.WithTimeout` para garantir que a aplicação não espere mais de 1 segundo por uma resposta.
- **Validação de CEP**: O CEP fornecido como entrada é validado para garantir que contém exatamente 8 dígitos numéricos.
- **Testes com Mocks**: Cobertura de testes de 100% para a lógica de negócio, utilizando um `MockProvider` para simular as respostas das APIs de forma rápida e confiável.
- **Erros Centralizados**: Erros padrão da aplicação são definidos como constantes para consistência e manutenibilidade.

## Tecnologias e Conceitos Utilizados

- **Go (Golang)**
- **Princípios SOLID**: Foco em Inversão de Dependência e Responsabilidade Única.
- **Interfaces** para abstração da lógica de busca.
- **Goroutines** para concorrência.
- **Channels** para comunicação entre goroutines.
- **`context`** para controle de timeout e cancelamento.
- **`sync.WaitGroup`** para sincronização.
- **`testify/assert`** para asserções nos testes.
- **Módulos Go** para gerenciamento de dependências.

## Arquitetura (C4 Model)

A arquitetura foi modelada usando C4 para ilustrar as decisões de design.

### Nível 1: Diagrama de Contexto do Sistema

Este diagrama mostra a interação do sistema com seus usuários e com as APIs externas de CEP.

```mermaid
C4Context
  title Diagrama de Contexto: Busca de CEP

  Person(user, "Usuário", "Pessoa que deseja consultar um CEP.")
  System(cep_app, "Aplicação de Busca de CEP", "CLI que orquestra a busca concorrente de CEP em múltiplas APIs.")
  
  System_Ext(brasilApi, "BrasilAPI", "API externa para consulta de CEP.")
  System_Ext(viaCep, "ViaCEP", "API externa para consulta de CEP.")


  Rel(user, cep_app, "Consulta um CEP via CLI")
  Rel(cep_app, brasilApi, "Busca dados do CEP", "HTTPS/JSON")
  Rel(cep_app, viaCep, "Busca dados do CEP", "HTTPS/JSON")
```

### Nível 3: Diagrama de Componentes

Este diagrama detalha os principais componentes (pacotes e structs) dentro da aplicação, mostrando como a lógica foi desacoplada através da interface `Provider`.

```mermaid
flowchart TD
    subgraph "Aplicação CLI"
        direction LR
        subgraph "Pacote `main`"
            main["main.go<br/>(Ponto de entrada)"]
        end

        subgraph "Pacote `cep`"
            cep_logic["cep.go<br/>(Orquestração)"]
            provider_interface["provider.go<br/>(Interface Provider)"]
        end

        subgraph "Pacote `provider`"
            brasil_api["brasil_api.go<br/>(BrasilAPIProvider)"]
            via_cep["via_cep.go<br/>(ViaCEPProvider)"]
        end
    end

    subgraph "APIs Externas"
        direction LR
        brasil_api_ext["BrasilAPI"]
        via_cep_ext["ViaCEP"]
    end

    main --> cep_logic
    cep_logic --> provider_interface
    brasil_api -- implementa --> provider_interface
    via_cep -- implementa --> provider_interface
    main --> brasil_api
    main --> via_cep
    brasil_api --> brasil_api_ext
    via_cep --> via_cep_ext

    style main fill:#D1E8FF,stroke:#367dD3,stroke-width:2px
    style cep_logic fill:#D1E8FF,stroke:#367dD3,stroke-width:2px
    style provider_interface fill:#FFF2CC,stroke:#d6b656,stroke-width:2px
    style brasil_api fill:#D5E8D4,stroke:#82b366,stroke-width:2px
    style via_cep fill:#D5E8D4,stroke:#82b366,stroke-width:2px
    style brasil_api_ext fill:#F8CECC,stroke:#b85450,stroke-width:2px
    style via_cep_ext fill:#F8CECC,stroke:#b85450,stroke-width:2px
```

## Como Executar

### Pré-requisitos

- Go 1.18 ou superior.

### Execução

1.  Clone o repositório:
    ```sh
    git clone https://github.com/markuscandido/go-expert-desafio-multithreading.git
    cd go-expert-desafio-multithreading
    ```

2.  Execute a aplicação passando um CEP como argumento:
    ```sh
    go run main.go 01001000
    ```

    Você pode usar um CEP com ou sem formatação:
    ```sh
    go run main.go 01001-000
    ```

### Executando os Testes

Para rodar a suíte de testes e ver a cobertura, execute o comando:

```sh
go test ./... -cover
```
