# 🎉 Comemora

**Comemora** é um backend leve e eficiente para gerenciamento de aniversários e datas comemorativas, construído em **Go** seguindo os princípios de **Clean Architecture (Hexagonal)**.

O sistema opera sem interface gráfica ("headless"), focado em gerenciar eventos e orquestrar notificações automáticas via integrações (Teams, WhatsApp, Email).

## ✨ Funcionalidades

*   **Gerenciamento de Eventos**: CRUD completo para aniversários, casamentos, datas de trabalho, etc.
*   **Motor de Verificação**: Checagem diária automática para eventos do dia e alertas antecipados de eventos importantes.
*   **Mensagens Inteligentes**: Templates automáticos de felicitações (ex: "Feliz Aniversário", "Bodas de Casamento") caso nenhuma mensagem personalizada seja fornecida.
*   **Notificações Pluggáveis**: Arquitetura pronta para múltiplos canais (atualmente com Mock/Console log, fácil de estender para APIs reais).
*   **Portabilidade**: Banco de dados SQLite em arquivo local (sem necessidade de configurar servidores de banco complexos).
*   **Integração em Massa**: Importação e Exportação de dados via CSV.

## 🛠 Tecnologias

*   **Linguagem**: Go (Golang) 1.21+
*   **Arquitetura**: Ports and Adapters (Hexagonal)
*   **Web Framework**: Chi (Router leve e idiomático)
*   **Database**: SQLite + GORM
*   **Padrões**: Service Layer, Repository Pattern, Dependency Injection (Pure Go).

## 🚀 Como Executar

### Pré-requisitos
*   [Go](https://go.dev/dl/) instalado na máquina.

### Instalação

1.  Clone o repositório (ou baixe os arquivos):
    ```bash
    git clone https://github.com/seu-usuario/comemora.git
    cd comemora
    ```

2.  Instale as dependências:
    ```bash
    go mod tidy
    ```

3.  Execute o servidor:
    ```bash
    go run cmd/comemora/main.go
    ```

    O servidor iniciará na porta `8080`. Você verá logs como:
    ```
    [Comemora] 2026/01/16 12:00:00 listening on [::]:8080
    ```

## 🔌 Guia da API

Aqui estão exemplos de como interagir com a API usando `curl`.

### 1. Criar um Evento
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ana Souza",
    "day": 20,
    "month": 5,
    "type": "Aniversário",
    "is_important": true,
    "preferred_channel": "WhatsApp",
    "contact_destination": "+5511999998888"
  }'
```

### 2. Listar Eventos
```bash
curl http://localhost:8080/events
```

### 3. Disparar Verificação Diária (Trigger)
Este é o endpoint que o **Cron Job** deve chamar. Ele verifica se há eventos hoje e dispara as notificações.
```bash
curl http://localhost:8080/trigger-check
```
*Observe o terminal onde o servidor está rodando para ver os logs das notificações enviadas.*

### 4. Exportar CSV
```bash
curl -O http://localhost:8080/events/export
```
Isso baixará um arquivo `events.csv`.

### 5. Importar CSV
```bash
curl -X POST http://localhost:8080/events/import \
  --data-binary @meus_eventos.csv
```

## ⏰ Configuração de Automação (Cron Job)

Para que o sistema funcione sozinho, você deve configurar uma tarefa agendada no seu sistema operacional para chamar o endpoint de verificação uma vez por dia.

### No Linux/Mac (Crontab)

1.  Abra o editor do cron:
    ```bash
    crontab -e
    ```

2.  Adicione a seguinte linha para rodar todos os dias às 09:00 da manhã:
    ```cron
    0 9 * * * curl -s http://localhost:8080/trigger-check > /dev/null
    ```

## 🏗 Estrutura do Projeto

Para desenvolvedores que desejam manter ou estender o projeto:

*   `cmd/comemora/`: Ponto de entrada (`main.go`). Onde tudo é conectado.
*   `internal/core/`:
    *   `domain/`: Onde vivem as Entidades (`Event`) e regras de negócio puras.
    *   `ports/`: Interfaces (contratos) para Repositórios e Notificadores.
    *   `services/`: Implementação da lógica de aplicação (EventService).
*   `internal/adapters/`:
    *   `handler/`: Camada HTTP (Handlers, Rotas, JSON decode/encode).
    *   `repository/`: Implementação do banco de dados (SQLite).
    *   `notifier/`: Integração com sistemas de mensagem (Console, e futuramente outros).
*   `docs/`: Documentação técnica detalhada e decisões de arquitetura (ADRs).

## 🤝 Como Contribuir / Estender

O sistema foi desenhado para ser extensível.

*   **Quer adicionar notificações reais via Telegram?**
    1.  Crie `internal/adapters/notifier/telegram.go`.
    2.  Implemente a interface `Notifier` definida em `ports`.
    3.  Vá em `main.go` e troque o `ConsoleNotifier` pelo seu `TelegramNotifier`.

---
Desenvolvido com ❤️ e Go.
