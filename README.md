# Comemora

**Comemora** é um backend para gerenciamento de aniversários e datas comemorativas, construído em **Go** seguindo os princípios de **Clean Architecture (Hexagonal)**.

Produção: [comemora.up.railway.app](https://comemora.up.railway.app)

O sistema opera sem interface gráfica ("headless"), focado em gerenciar eventos e orquestrar notificações automáticas via múltiplos canais (Email, WhatsApp, Teams, Telegram, Discord).

## Funcionalidades

- **Gerenciamento de Eventos**: CRUD para aniversários, casamentos, datas de trabalho, etc.
- **Motor de Verificação**: Checagem diária para eventos do dia e alertas antecipados (3 dias antes) de eventos importantes.
- **Mensagens Inteligentes**: Templates automáticos de felicitações por tipo de evento. Suporta mensagens personalizadas por evento.
- **Notificações Multi-Canal**: Cada evento define seu canal preferido. Canais disponíveis: Email (Resend), WhatsApp (Infobip), Teams, Telegram e Discord. Console como fallback.
- **Integração em Massa**: Importação e exportação de dados via CSV.

## Tecnologias

- **Linguagem**: Go 1.21+
- **Arquitetura**: Ports and Adapters (Hexagonal)
- **Web Framework**: Chi
- **Database**: PostgreSQL 18 + GORM (com AutoMigrate)
- **Configuração**: godotenv (arquivo `.env`)
- **Infraestrutura local**: Docker Compose

## Variáveis de Ambiente

Copie `.env.example` para `.env` e ajuste os valores.

### Banco de Dados

| Variável       | Exemplo                                                           | Descrição                             |
|----------------|-------------------------------------------------------------------|---------------------------------------|
| `DATABASE_URL` | `postgres://user:pass@localhost:5432/comemora?sslmode=disable`    | URL completa do PostgreSQL (padrão Railway) |

### Servidor

| Variável      | Padrão | Descrição              |
|---------------|--------|------------------------|
| `SERVER_PORT` | `8080` | Porta do servidor HTTP |

### Notificadores (todos opcionais)

Os notificadores são ativados automaticamente quando as variáveis correspondentes estão definidas. Se nenhum estiver configurado, as notificações ficam apenas no console (log).

| Canal    | Variáveis necessárias                                        | Provedor                             |
|----------|--------------------------------------------------------------|--------------------------------------|
| Email    | `RESEND_API_KEY`, `RESEND_FROM`                              | [Resend](https://resend.com)         |
| WhatsApp | `INFOBIP_API_KEY`, `INFOBIP_BASE_URL`, `INFOBIP_FROM_NUMBER` | [Infobip](https://infobip.com)       |
| Teams    | `TEAMS_WEBHOOK_URL`                                          | Microsoft Teams Incoming Webhook     |
| Telegram | `TELEGRAM_BOT_TOKEN`                                         | [BotFather](https://t.me/BotFather) |
| Discord  | `DISCORD_WEBHOOK_URL`                                        | Discord Webhook                      |

## Banco de Dados

### Tabela: `events`

O GORM cria e mantém a tabela automaticamente via `AutoMigrate` na inicialização. Não é necessário rodar migrations manualmente.

| Coluna                | Tipo          | Descrição                                                                          |
|-----------------------|---------------|------------------------------------------------------------------------------------|
| `id`                  | `bigint PK`   | Identificador único (auto-incremento)                                              |
| `name`                | `text`        | Nome da pessoa ou evento                                                           |
| `day`                 | `integer`     | Dia do evento (1–31)                                                               |
| `month`               | `integer`     | Mês do evento (1–12)                                                               |
| `year`                | `integer`     | Ano (0 = recorrente ou desconhecido)                                               |
| `type`                | `text`        | Tipo: `aniversario`, `casamento`, `namoro`, `pet`, `trabalho`, `luto`, `outro`     |
| `tags`                | `jsonb`       | Array de tags customizadas                                                         |
| `preferred_channel`   | `text`        | Canal de notificação: `Email`, `WhatsApp`, `Teams`, `Telegram`, `Discord`, `SMS`  |
| `contact_destination` | `text`        | Destino da notificação (e-mail, telefone, chat_id do Telegram, etc.)              |
| `custom_message`      | `text`        | Mensagem personalizada (opcional; se vazio usa o template padrão do tipo)         |
| `is_important`        | `boolean`     | Se verdadeiro, notifica com 3 dias de antecedência                                |
| `created_at`          | `timestamptz` | Data de criação                                                                    |
| `updated_at`          | `timestamptz` | Data de última atualização                                                         |

## Como Executar

### Pré-requisitos

- [Go](https://go.dev/dl/) 1.21+
- [Docker](https://docs.docker.com/get-docker/) e Docker Compose

### 1. Subir o banco de dados

```bash
docker compose up -d
```

Sobe um PostgreSQL 18 na porta `5432` com o banco `comemora` já criado.

### 2. Configurar variáveis de ambiente

```bash
cp .env.example .env
```

O `.env.example` já vem configurado para funcionar com o Docker Compose local. Para habilitar um notificador, adicione as variáveis correspondentes ao seu `.env`.

### 3. Executar o servidor

```bash
go run cmd/comemora/main.go
```

O servidor inicia na porta `8080`. Na inicialização, o log mostra quais notificadores foram habilitados:

```
[Comemora] 2026/01/16 12:00:00 listening on 0.0.0.0:8080
[Comemora] 2026/01/16 12:00:00 notifier: Email (Resend) enabled
[Comemora] 2026/01/16 12:00:00 notifier: WhatsApp (Infobip) enabled
```

Na primeira execução, o GORM cria automaticamente a tabela `events`.

## Guia da API

### Criar um Evento

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

**`type` válidos:** `aniversario` · `casamento` · `namoro` · `pet` · `trabalho` · `luto` · `outro`

**`preferred_channel` válidos:** `Email` · `WhatsApp` · `Teams` · `Telegram` · `Discord` · `SMS`

**`contact_destination` por canal:**

| Canal    | Formato esperado                        |
|----------|-----------------------------------------|
| Email    | `ana@exemplo.com`                       |
| WhatsApp | `+5511999998888`                        |
| Teams    | *(não usa — mensagem vai para o webhook)* |
| Telegram | Chat ID numérico, ex: `123456789`       |
| Discord  | *(não usa — mensagem vai para o webhook)* |

> Se o campo `custom_message` for omitido, o sistema usa um template automático baseado no `type`.

### Listar Eventos

```bash
curl http://localhost:8080/events
```

### Disparar Verificação

Verifica eventos de hoje e eventos com `is_important: true` nos próximos 3 dias, e envia notificações pelo canal de cada evento.

```bash
curl http://localhost:8080/trigger-check
```

### Exportar CSV

```bash
curl -O http://localhost:8080/events/export
```

Baixa um arquivo `events.csv` com todos os eventos.

### Importar CSV

```bash
curl -X POST http://localhost:8080/events/import \
  --data-binary @meus_eventos.csv
```

O formato do CSV deve seguir o mesmo padrão do export: `ID,Name,Day,Month,Year,Type,IsImportant,Channel,Contact`.

### Health Check

```bash
curl http://localhost:8080/health
```

Retorna `200 OK` quando o serviço está no ar. Útil para monitoramento e Railway health checks.

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
│   │   ├── domain/       # Entidade Event, tipos, canais e lógica de mensagens
│   │   ├── ports/        # Interfaces: EventRepository, Notifier, Service
│   │   └── services/     # EventService — lógica de aplicação
│   └── adapters/
│       ├── handler/      # HTTP handlers e rotas (Chi)
│       ├── repository/   # PostgresRepository (GORM)
│       └── notifier/     # MultiNotifier + Email, WhatsApp, Teams, Telegram, Discord, Console
├── docs/                 # Decisões de arquitetura (ADRs) e contexto do projeto
├── docker-compose.yml    # PostgreSQL local
└── .env.example          # Variáveis de ambiente de exemplo
```

## Como Estender

O sistema foi desenhado para ser extensível via interfaces.

**Adicionar um novo canal de notificação** (ex.: Slack):

1. Crie `internal/adapters/notifier/slack_notifier.go` e implemente a interface `Notifier`:
   ```go
   // internal/core/ports/ports.go
   type Notifier interface {
       Send(ctx context.Context, event domain.Event) error
   }
   ```
2. Adicione a constante do canal em `internal/core/domain/event.go`:
   ```go
   ChannelSlack NotificationChannel = "Slack"
   ```
3. Em `cmd/comemora/main.go`, registre no `MultiNotifier`:
   ```go
   if url := os.Getenv("SLACK_WEBHOOK_URL"); url != "" {
       multi.Register(domain.ChannelSlack, notifier.NewSlackNotifier(url))
   }
   ```

**Trocar o banco de dados:** implemente a interface `EventRepository` e injete no `main.go`.

---

Desenvolvido com Go.
