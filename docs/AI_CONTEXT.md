# Comemora - AI Context & Documentation

Este documento serve como guia para agentes de IA e desenvolvedores que venham a manter este projeto.

## Propósito do Projeto

**Comemora** é um sistema backend "headless" para gerenciamento de datas comemorativas (aniversários, casamentos, etc.).
O objetivo é ser **simples, leve e desacoplado**. Não possui frontend — funciona via API REST e Cron Jobs externos.

## Arquitetura: Clean Architecture (Hexagonal)

O projeto segue estritamente a arquitetura de **Ports and Adapters**.

### Estrutura de Diretórios

- `cmd/comemora/main.go`: Ponto de entrada. Faz a composição de todas as dependências (DI manual).
- `internal/core/domain/`: **O Coração**. Entidade `Event`, tipos (`EventType`, `NotificationChannel`) e lógica pura (`GetContent`). Sem dependências externas.
- `internal/core/ports/ports.go`: **Contratos**. Interfaces `EventRepository`, `Notifier` e `Service`.
- `internal/core/services/`: **Regras de Negócio**. `EventService` orquestra repositório e notificadores.
- `internal/adapters/handler/`: Camada HTTP (Chi). Handlers inline, rotas centralizadas em `routes.go`.
- `internal/adapters/repository/`: `PostgresRepository` via GORM.
- `internal/adapters/notifier/`: `MultiNotifier` + adaptadores por canal (Email, WhatsApp, Teams, Telegram, Discord, Console).

## Endpoints

| Método | Rota              | Descrição                                                      |
|--------|-------------------|----------------------------------------------------------------|
| GET    | `/`               | Lista os endpoints disponíveis (JSON)                         |
| GET    | `/health`         | Health check (retorna 200 OK)                                 |
| POST   | `/events`         | Cria um evento                                                |
| GET    | `/events`         | Lista todos os eventos                                        |
| GET    | `/events/export`  | Exporta eventos em CSV                                        |
| POST   | `/events/import`  | Importa eventos via CSV                                       |
| GET    | `/trigger-check`  | Dispara verificação e notificações do dia                     |

## Sistema de Notificação

### Como funciona

1. `EventService.CheckAndNotify` busca eventos do dia e eventos importantes nos próximos 3 dias.
2. Para cada evento, chama `Notifier.Send(ctx, event)`.
3. O `MultiNotifier` despacha para o adaptador registrado para `event.PreferredChannel`.
4. Se o canal não estiver configurado (env var ausente), cai no `ConsoleNotifier` (fallback).

### Canais implementados

| Canal    | Adaptador              | Credencial mínima                           |
|----------|------------------------|---------------------------------------------|
| Email    | `EmailNotifier`        | `RESEND_API_KEY`                            |
| WhatsApp | `WhatsAppNotifier`     | `INFOBIP_API_KEY` + `INFOBIP_BASE_URL` + `INFOBIP_FROM_NUMBER` |
| Teams    | `TeamsNotifier`        | `TEAMS_WEBHOOK_URL`                         |
| Telegram | `TelegramNotifier`     | `TELEGRAM_BOT_TOKEN`                        |
| Discord  | `DiscordNotifier`      | `DISCORD_WEBHOOK_URL`                       |
| Console  | `ConsoleNotifier`      | nenhuma (sempre ativo como fallback)        |

### Adicionar um novo canal

1. Criar `internal/adapters/notifier/novo_canal.go` implementando `ports.Notifier`.
2. Adicionar a constante em `internal/core/domain/event.go`.
3. Registrar no `MultiNotifier` em `cmd/comemora/main.go` com a lógica de env var.

## Lógica de Negócio Importante

- **Mensagens**: `Event.GetContent()` em `domain/event.go`. Se `CustomMessage != ""`, usa ela. Caso contrário, usa template do `EventType`.
- **Antecipação**: Eventos com `IsImportant: true` são notificados também 3 dias antes da data.
- **Banco**: PostgreSQL via GORM. `AutoMigrate` cria/atualiza a tabela `events` na inicialização. `DATABASE_URL` é a variável esperada (padrão Railway).
- **Opt-in de canais**: Notificadores só são registrados se as variáveis de ambiente correspondentes estiverem presentes. Facilita deploys parciais.

## Padrões de Código

- **HTTP**: "Grafana Style" — `NewServer`, `routes.go`, handlers como funções que retornam `http.HandlerFunc`, tipos request/response inline.
- **Erros**: Sempre wrapped com `fmt.Errorf("prefixo: %w", err)`. Prefixo indica o canal ou camada.
- **HTTP Client**: Todos os notificadores usam `http.NewRequestWithContext` para propagar cancelamento. Response body sempre drenado antes de fechar (`io.Copy(io.Discard, resp.Body)`) para reutilização de conexão TCP.
- **Configuração**: Lida via `os.Getenv`. `godotenv.Load()` carrega `.env` em dev; em produção usa variáveis reais do ambiente.
