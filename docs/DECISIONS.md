# Architecture Decisions Records (ADR)

## ADR 001: Linguagem Go
**Decisão:** Utilizar Go 1.21+.
**Motivo:** Performance, baixo footprint de memória, binário estático fácil de deployar e forte tipagem.
**Contexto:** O usuário solicitou explicitamente Go para um backend leve.

## ADR 002: Arquitetura Hexagonal (Ports and Adapters)
**Decisão:** Separar estritamente Core (Domain/Services) de Adapters (Handler/Repository).
**Motivo:** Permitir mudar o banco de dados (SQLite -> Postgres) ou canais de notificação (Console -> WhatsApp) sem alterar uma linha de regra de negócio. Facilita testes unitários.

## ADR 003: SQLite como Banco de Dados Inicial
**Decisão:** SQLite em arquivo local com GORM.
**Motivo:** "Zero config". Não requer Docker ou serviço rodando. GORM abstrai o SQL, facilitando migração futura.

## ADR 004: Gerenciamento de Dependências
**Decisão:** Injeção de Dependência Manual (Pure Go) no `main.go`.
**Motivo:** Evitar frameworks de DI "mágicos" (como Uber Dig ou Wire) para manter o código simples e legível para quem está aprendendo Go.

## ADR 005: Tratamento de Mensagens
**Decisão:** Lógica de `GetContent` no Domain.
**Motivo:** A mensagem é uma propriedade intrínseca do evento na hora do disparo. Manter templates padrões no código (map) por simplicidade, em vez de criar uma tabela no banco para isso agora (MVP).

## ADR 006: HTTP Service Pattern
**Decisão:** Padrão "Grafana/Mat Ryer".
**Motivo:** O uso de struct `NewServer`, rotas centralizadas em `routes.go` e função `run` no main facilita a leitura e manutenção, além de ser um padrão de indústria robusto.
