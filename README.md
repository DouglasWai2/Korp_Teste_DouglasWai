# NF System

Sistema de estudo composto por:
- frontend em Angular
- microserviço `estoque` em Go
- microserviço `faturamento` em Go
- banco PostgreSQL

O projeto permite:
- cadastrar produtos
- adicionar saldo a produtos
- remover produtos
- cadastrar notas fiscais com itens
- listar notas fiscais
- visualizar detalhes de uma nota fiscal
- imprimir uma nota fiscal, baixando saldo no serviço de estoque

## Arquitetura

Estrutura principal:

```text
.
├── db/
│   ├── init.sql
│   └── migrations/
├── services/
│   ├── estoque/
│   └── faturamento/
├── web/
├── docker-compose.yaml
└── .env
```

Serviços:
- `estoque`: roda em `http://localhost:8080`
- `faturamento`: roda em `http://localhost:8081`
- `web`: roda em `http://localhost:4200`
- `postgres`: roda em `localhost:5432`

## Tecnologias

Frontend:
- Angular 21
- TypeScript
- RxJS
- FormsModule
- Angular Router

Backend:
- Go 1.26.2
- Gin
- PostgreSQL
- `lib/pq`
- `godotenv`

Banco:
- PostgreSQL 16 via Docker Compose

## Pré-requisitos

Instale antes de iniciar:
- Docker e Docker Compose
- Go `1.26.2` ou compatível
- Node.js e npm

## Variáveis de ambiente

O projeto utiliza `.env` na raiz.

Exemplo atual:

```env
DB_URL=postgres://wai_user:admin@localhost:5432/nf_system_db?sslmode=disable
```

Observações:
- `estoque` usa `DB_URL`
- `faturamento` usa `DB_URL`
- `faturamento` também aceita `ESTOQUE_API_URL`
- se `ESTOQUE_API_URL` não for informado, o default é `http://localhost:8080`

## Como inicializar o projeto

### 1. Subir o banco

Na raiz do projeto:

```bash
docker compose up -d postgres
```

Verificar se o banco está pronto:

```bash
docker compose ps
docker compose exec postgres pg_isready -U wai_user -d nf_system_db
```

## 2. Aplicar schema e migrations

Se estiver iniciando um banco novo, o `init.sql` já é executado automaticamente pelo container.

Se quiser aplicar manualmente:

```bash
docker compose exec postgres psql -U wai_user -d nf_system_db -f /docker-entrypoint-initdb.d/init.sql
```

Migration de quantidade em `nf_produtos`:

```bash
docker compose exec postgres psql -U wai_user -d nf_system_db -f /migrations/001_add_quantidade_to_nf_produtos.sql
```

## 3. Subir o microserviço de estoque

Em outro terminal:

```bash
cd services/estoque
go run ./cmd/main.go
```

API disponível em:

```text
http://localhost:8080
```

Health check:

```bash
curl http://localhost:8080/api
```

## 4. Subir o microserviço de faturamento

Em outro terminal:

```bash
cd services/faturamento
go run ./cmd/main.go
```

API disponível em:

```text
http://localhost:8081
```

Health check:

```bash
curl http://localhost:8081/api/faturamento
```

## 5. Subir o frontend

Em outro terminal:

```bash
cd web
npm install
npm start
```

Aplicação disponível em:

```text
http://localhost:4200
```

## Ordem recomendada de subida

```bash
docker compose up -d postgres
cd services/estoque && go run ./cmd/main.go
cd services/faturamento && go run ./cmd/main.go
cd web && npm install && npm start
```

## Comandos úteis

### Banco

Subir banco:

```bash
docker compose up -d postgres
```

Reiniciar banco:

```bash
docker compose restart postgres
```

Resetar banco local:

```bash
docker compose down -v
docker compose up -d postgres
```

### Backend Go

Build do serviço `estoque`:

```bash
cd services/estoque
env GOCACHE=/tmp/go-build-estoque go build ./...
```

Build do serviço `faturamento`:

```bash
cd services/faturamento
env GOCACHE=/tmp/go-build-faturamento go build ./...
```

### Frontend

Instalar dependências:

```bash
cd web
npm install
```

Rodar em desenvolvimento:

```bash
npm start
```

Build:

```bash
npm run build
```

Testes:

```bash
npm test
```

## Resumo das APIs

### Estoque

Base URL:

```text
http://localhost:8080
```

Rotas principais:
- `GET /api`
- `POST /api/products`
- `GET /api/products`
- `GET /api/products/:codigo`
- `PATCH /api/products/:codigo/increment`
- `PATCH /api/products/:codigo/decrement`
- `DELETE /api/products/:codigo`

Exemplo de criação:

```json
{
  "codigo": "ABC123",
  "descricao": "Produto exemplo",
  "saldo": 10
}
```

Exemplo para adicionar saldo:

```json
{
  "quantidade": 5
}
```

### Faturamento

Base URL:

```text
http://localhost:8081
```

Rotas principais:
- `GET /api/faturamento`
- `POST /api/faturamento/notas-fiscais`
- `GET /api/faturamento/notas-fiscais`
- `GET /api/faturamento/notas-fiscais/:numero`
- `PATCH /api/faturamento/notas-fiscais/:numero/imprimir`

Exemplo de criação de nota fiscal:

```json
{
  "itens": [
    { "codigo_produto": "ABC123", "quantidade": 2 },
    { "codigo_produto": "XYZ999", "quantidade": 1 }
  ]
}
```

## Frontend

Telas disponíveis:
- `Produtos`
- `Notas fiscais`
- `Detalhes da nota fiscal`

Funcionalidades da interface:
- cadastro de produto
- reposição de saldo
- remoção de produto
- criação de nota com múltiplos itens
- impressão de nota
- detalhamento dos itens da nota

## Tratamento de erros

Frontend:
- mensagens amigáveis por serviço
- exemplos:
  - `Problema no servico de produtos, tente novamente mais tarde.`
  - `Problema no servico de faturamento, tente novamente mais tarde.`

Backend:
- uso de status HTTP adequados
- `400` para requisição inválida
- `404` para recurso não encontrado
- `409` para conflitos de negócio
- `500` para erro interno

## Observações de desenvolvimento

- o frontend chama as APIs por URL absoluta em `localhost`
- os serviços Go habilitam CORS para `http://localhost:4200`
- o faturamento depende do serviço de estoque para:
  - validar produtos na criação da nota
  - baixar saldo ao imprimir a nota

## Problemas comuns

### Frontend não carrega produtos

Verifique:

```bash
curl http://localhost:8080/api/products
```

### Frontend não carrega notas fiscais

Verifique:

```bash
curl http://localhost:8081/api/faturamento/notas-fiscais
```

### Banco não sobe corretamente

Ver logs:

```bash
docker compose logs postgres
```

### Reiniciar tudo do zero

```bash
docker compose down -v
docker compose up -d postgres

cd services/estoque && go run ./cmd/main.go
cd services/faturamento && go run ./cmd/main.go
cd web && npm install && npm start
```
