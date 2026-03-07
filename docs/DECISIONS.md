# Architecture Decision Records (ADR)

## ADR 001: Linguagem Go
**Decisão:** Utilizar Go 1.21+.
**Motivo:** Performance, baixo footprint de memória, binário estático fácil de deployar e forte tipagem.

## ADR 002: Arquitetura Hexagonal (Ports and Adapters)
**Decisão:** Separar estritamente Core (Domain/Services) de Adapters (Handler/Repository/Notifier).
**Motivo:** Permite mudar o banco de dados ou canais de notificação sem alterar uma linha de regra de negócio. Facilita testes unitários e a adição de novos adaptadores.

## ADR 003: PostgreSQL como Banco de Dados
**Decisão:** PostgreSQL via GORM com AutoMigrate.
**Motivo:** O projeto foi iniciado pensando em deploy na Railway, que fornece PostgreSQL gerenciado. O GORM abstrai o SQL e gerencia as migrations automaticamente na inicialização, eliminando a necessidade de scripts manuais.
**Nota histórica:** A versão inicial usava SQLite para prototipagem, mas PostgreSQL foi adotado para o ambiente de produção.

## ADR 004: Injeção de Dependência Manual no `main.go`
**Decisão:** Composição de dependências feita manualmente no `main.go`, sem frameworks de DI.
**Motivo:** Evitar frameworks "mágicos" (Uber Dig, Wire). O grafo de dependências do projeto é simples o suficiente para ser legível diretamente no código de bootstrap.

## ADR 005: Lógica de Mensagens no Domain
**Decisão:** Método `GetContent()` vive na entidade `Event` em `domain/event.go`.
**Motivo:** A mensagem é uma propriedade intrínseca do evento. Templates padrão por tipo de evento ficam em um `map` no código (não no banco) por ser uma decisão de MVP — simples, sem overhead operacional.

## ADR 006: HTTP Service Pattern ("Grafana Style")
**Decisão:** Struct `NewServer`, rotas centralizadas em `routes.go`, função `run()` no main.
**Motivo:** Padrão documentado por Mat Ryer (Grafana Labs). Facilita leitura, manutenção e testes de integração ao separar a composição da inicialização do servidor.

## ADR 007: MultiNotifier com Dispatch por Canal
**Decisão:** `MultiNotifier` que despacha para o notificador correto com base em `event.PreferredChannel`, com fallback para `ConsoleNotifier`.
**Motivo:** Permite habilitar múltiplos canais simultaneamente sem alterar a interface `Notifier` nem a lógica do `EventService`. Cada canal é opt-in via variável de ambiente — se a variável não está configurada, o canal simplesmente não é registrado.
