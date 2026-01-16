# CelebrationHub - AI Context & Documentation

Este documento serve como guia para Agentes de IA e desenvolvedores que venham a manter este projeto.

## 🎯 Propósito do Projeto
**CelebrationHub** é um sistema backend "headless" para gerenciamento de datas comemorativas (aniversários, eventos, etc.).
O objetivo principal é ser **simples, leve e desacoplado**. Ele não possui Frontend acoplado; funciona via API REST e Cron Jobs.

## 🏗 Arquitetura: Clean Architecture (Hexagonal)
O projeto segue estritamente a arquitetura de **Ports and Adapters**.

### Estrutura de Diretórios
*   `cmd/`: Entrypoints. `main.go`.
*   `internal/core/domain`: **O Coração**. Contém as entidades (`Event`) e lógica pura de dados (`GetContent`). NENHUMA dependência externa aqui (sem frameworks, sem banco).
*   `internal/core/ports`: **Contratos**. Interfaces que definem o que o Core precisa (`EventRepository`, `Notifier`) e o que ele oferece (`Service`).
*   `internal/core/services`: **Implementação das Regras**. Orquestra o fluxo. Chama repositórios e notifiers.
*   `internal/adapters`: **O Mundo Externo**. Implementações concretas das Interfaces.
    *   `handler/`: Camada HTTP (Entrada). Framework `Chi`.
    *   `repository/`: Banco de Dados (Saída). Framework `GORM`/SQLite.
    *   `notifier/`: Integrações de Mensagem (Saída). Atualmente Mock/Console.

## 🛠 Design Decisions (Decisões de Projeto)

### 1. Go Pattern: "Grafana Style" HTTP Services
Adotamos o estilo documentado por Mat Ryer (Grafana Labs):
*   `routes.go`: Centraliza todas as rotas.
*   `NewServer`: Construtor que recebe dependências.
*   `run()`: Função principal que facilita testes e injecta dependências.
*   **Inline Types**: Request/Response structs definidos DENTRO dos handlers para evitar acoplamento desnecessário.

### 2. SQLite
Escolha inicial para manter o custo zero e portabilidade (arquivo único). O uso de GORM permite migrar para PostgreSQL trocando apenas o driver no `main.go`.

### 3. Notificações Pluggables
A interface `Notifier` (`internal/core/ports/ports.go`) é crítica.
*   Hoje: `ConsoleNotifier` (apenas logs).
*   Futuro: Criar `WhatsAppNotifier` ou `TeamsNotifier` implementando a mesma interface.

## 🧠 Lógica de Negócio Importante
*   **Mensagens Padrão**: Definidas em `domain/event.go`. Se `CustomMessage` for vazio, o sistema escolhe um template baseado no `EventType`.
*   **Check Engine**: O endpoint `GET /trigger-check` é o coração da automação. Ele deve ser chamado via Cron externo.

## 🤖 Como Estender
1.  **Novo Canal de Notificação**: Crie `internal/adapters/notifier/seu_canal.go`. Implemente `Send`. Injete no `main.go`.
2.  **Novo Filtro de Evento**: Adicione métodos na interface `EventRepository` e a implementação em `sqlite_repo.go`.
3.  **Alterar Mensagens Padrão**: Edite `internal/core/domain/event.go`.

---
*Gerado por Antigravity*
