# Comemora

**Comemora** é um backend para gerenciamento de aniversários e datas comemorativas, construído em **Go** seguindo os princípios de **Clean Architecture (Hexagonal)**.

prod url: [comemora.up.railway.app](https://comemora.up.railway.app)

O sistema opera sem interface gráfica ("headless"), focado em gerenciar eventos e orquestrar notificações automáticas via integrações (Teams, WhatsApp, Email).

## Funcionalidades

- **Gerenciamento de Eventos**: CRUD para aniversários, casamentos, datas de trabalho, etc.
- **Motor de Verificação**: Checagem diária para eventos do dia e alertas antecipados (3 dias antes) de eventos importantes.
- **Mensagens Inteligentes**: Templates automáticos de felicitações por tipo de evento. Suporta mensagens personalizadas.
- **Notificações Pluggáveis**: Arquitetura preparada para múltiplos canais (Email, WhatsApp, Teams, SMS). Atualmente implementado com Console/log.
- **Integração em Massa**: Importação e Exportação de dados via CSV.

## Tecnologias

- **Linguagem**: Go 1.21+
- **Arquitetura**: Ports and Adapters (Hexagonal)
- **Web Framework**: Chi
- **Database**: PostgreSQL 16 + GORM (com AutoMigrate)
- **Configuração**: godotenv (arquivo `.env`)
- **Infraestrutura local**: Docker Compose

## Banco de Dados

### Configuração

O projeto usa **PostgreSQL**. As credenciais são lidas via variáveis de ambiente (ou arquivo `.env`):

| Variável      | Padrão      | Descrição                        |
|---------------|-------------|----------------------------------|
| `DB_HOST`     | `localhost` | Host do PostgreSQL               |
| `DB_PORT`     | `5432`      | Porta do PostgreSQL              |
| `DB_USER`     | `user`      | Usuário do banco                 |
| `DB_PASSWORD` | `password`  | Senha do banco                   |
| `DB_NAME`     | `comemora`  | Nome do banco de dados           |
| `DB_SSLMODE`  | `disable`   | Modo SSL (`disable` em dev)      |
| `SERVER_PORT` | `8080`      | Porta do servidor HTTP           |

### Tabela: `events`

O GORM cria e mantém a tabela automaticamente via `AutoMigrate` na inicialização. Não é necessário rodar migrations manualmente.

| Coluna                | Tipo        | Descrição                                              |
|-----------------------|-------------|--------------------------------------------------------|
| `id`                  | `bigint PK` | Identificador único (auto-incremento)                  |
| `name`                | `text`      | Nome da pessoa ou evento                               |
| `day`                 | `integer`   | Dia do evento (1–31)                                   |
| `month`               | `integer`   | Mês do evento (1–12)                                   |
| `year`                | `integer`   | Ano (0 = recorrente/desconhecido)                      |
| `type`                | `text`      | Tipo: `aniversario`, `casamento`, `namoro`, `pet`, `trabalho`, `luto`, `outro` |
| `tags`                | `jsonb`     | Array de tags customizadas                             |
| `preferred_channel`   | `text`      | Canal preferido: `Email`, `WhatsApp`, `Teams`, `SMS`   |
| `contact_destination` | `text`      | Destino da notificação (e-mail, telefone, etc.)        |
| `custom_message`      | `text`      | Mensagem personalizada (opcional)                      |
| `is_important`        | `boolean`   | Se verdadeiro, notifica com 3 dias de antecedência     |
| `created_at`          | `timestamptz` | Data de criação                                      |
| `updated_at`          | `timestamptz` | Data de última atualização                           |

### Configuração em um banco novo

Em um banco PostgreSQL novo, basta:

1. **Criar o banco** com o nome `comemora` (ou o nome definido em `DB_NAME`):
   ```sql
   CREATE DATABASE comemora;
   ```

2. **Iniciar o servidor** — o GORM roda o `AutoMigrate` automaticamente e cria a tabela `events` com todas as colunas na primeira execução.

Não é necessário rodar nenhum script SQL manualmente.

## Como Executar

### Pré-requisitos

- [Go](https://go.dev/dl/) 1.21+
- [Docker](https://docs.docker.com/get-docker/) e Docker Compose

### 1. Subir o banco de dados

```bash
docker compose up -d
```

Isso sobe um PostgreSQL 16 na porta `5432` com o banco `comemora` já criado.

### 2. Configurar variáveis de ambiente

```bash
cp .env.example .env
```

O `.env.example` já vem com os valores padrão compatíveis com o Docker Compose. Edite se necessário.

### 3. Instalar dependências

```bash
go mod tidy
```

### 4. Executar o servidor

```bash
go run cmd/comemora/main.go
```

O servidor inicia na porta `8080`:
```
[Comemora] 2026/01/16 12:00:00 listening on 0.0.0.0:8080
```

Na primeira execução, o GORM cria automaticamente a tabela `events` no banco.

## Guia da API

### 1. Criar um Evento

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ana Souza",
    "day": 20,
    "month": 5,
    "type": "aniversario",
    "is_important": true,
    "preferred_channel": "WhatsApp",
    "contact_destination": "+5511999998888"
  }'
```

**Tipos de evento válidos:** `aniversario`, `casamento`, `namoro`, `pet`, `trabalho`, `luto`, `outro`

**Canais válidos:** `Email`, `WhatsApp`, `Teams`, `SMS`

### 2. Listar Eventos

```bash
curl http://localhost:8080/events
```

### 3. Disparar Verificação Diária

Verifica eventos de hoje e eventos importantes nos próximos 3 dias, e envia notificações.

```bash
curl http://localhost:8080/trigger-check
```

Observe o terminal do servidor para ver os logs das notificações.

### 4. Exportar CSV

```bash
curl -O http://localhost:8080/events/export
```

Baixa um arquivo `events.csv` com todos os eventos.

### 5. Importar CSV

```bash
curl -X POST http://localhost:8080/events/import \
  --data-binary @meus_eventos.csv
```

O formato do CSV deve seguir o mesmo padrão do export: `ID,Name,Day,Month,Year,Type,IsImportant,Channel,Contact`.

## Automação com Cron Job

Configure uma tarefa agendada para chamar o trigger diariamente:

```bash
crontab -e
```

Adicione a linha para rodar todos os dias às 09:00:

```cron
0 9 * * * curl -s http://localhost:8080/trigger-check > /dev/null
```

## Estrutura do Projeto

```
comemora/
├── cmd/comemora/         # Ponto de entrada (main.go) — composição das dependências
├── internal/
│   ├── core/
│   │   ├── domain/       # Entidade Event, tipos e lógica de mensagens
│   │   ├── ports/        # Interfaces: EventRepository, Notifier, Service
│   │   └── services/     # EventService — lógica de aplicação
│   └── adapters/
│       ├── handler/      # HTTP handlers, rotas (Chi)
│       ├── repository/   # PostgresRepository (GORM)
│       └── notifier/     # ConsoleNotifier (log)
├── docs/                 # Decisões de arquitetura (ADRs)
├── docker-compose.yml    # PostgreSQL local
└── .env.example          # Variáveis de ambiente de exemplo
```

## Como Estender

O sistema foi desenhado para ser extensível via interfaces (ports).

**Quer adicionar notificações via Telegram?**

1. Crie `internal/adapters/notifier/telegram.go`
2. Implemente a interface `Notifier` definida em `internal/core/ports/ports.go`:
   ```go
   type Notifier interface {
       Send(ctx context.Context, event domain.Event) error
   }
   ```
3. Em `cmd/comemora/main.go`, substitua `notifier.NewConsoleNotifier(logger)` pelo seu `TelegramNotifier`

O mesmo padrão vale para trocar o banco de dados: implemente `EventRepository` e injete no `main.go`.

---

Desenvolvido com Go.
