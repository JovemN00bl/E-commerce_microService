# 🛍️ Go Microservices E-commerce

Um sistema de e-commerce robusto e escalável construído com **Go (Golang)**, utilizando arquitetura de microsserviços. Este projeto demonstra a implementação de padrões modernos de desenvolvimento backend, incluindo comunicação síncrona de alta performance e processamento assíncrono de eventos.

## 🏗️ Arquitetura do Sistema

O sistema é dividido em serviços independentes que se comunicam para realizar operações de negócio complexas:

1.  **API Gateway:** A porta de entrada. Gerencia o roteamento de requisições e a validação de segurança via Middleware JWT.
2.  **Auth Service:** Gerencia o ciclo de vida do usuário, criptografia de senhas (bcrypt) e emissão de tokens de acesso.
3.  **Products Service:** Gerencia o catálogo. Expõe API REST para o cliente e um servidor **gRPC** para consulta interna de estoque.
4.  **Orders Service:** O orquestrador. Gerencia transações atômicas de pedidos no banco, comunica-se via gRPC com o serviço de Produtos e publica eventos no RabbitMQ.
5.  **Infraestrutura:** Bancos de dados PostgreSQL isolados por serviço e RabbitMQ para mensageria.

## 🚀 Tecnologias Utilizadas

* **Linguagem:** Go (Golang) 1.21+
* **Framework Web:** Gin Gonic
* **Banco de Dados:** PostgreSQL (Driver `pgx/v5`)
* **Comunicação Síncrona:** gRPC & Protocol Buffers
* **Comunicação Assíncrona:** RabbitMQ (AMQP 0-9-1)
* **Autenticação:** JWT (JSON Web Tokens)
* **Containerização:** Docker & Docker Compose

## 📂 Estrutura de Pastas

```bash
E-commerce_micro/
├── api-gateway/       # Proxy Reverso e Autenticação (Porta 8000)
├── auth-service/      # Login e Registro (Porta 8080)
├── products-service/  # Catálogo + Server gRPC (Porta 8081 / 50051)
├── orders-service/    # Pedidos + Client gRPC (Porta 8082)
├── proto/             # Contratos gRPC (.proto) compartilhados
└── docker-compose.yml # Orquestração dos containers (DBs e Broker)
