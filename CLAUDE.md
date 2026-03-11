# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PicoClaw is an ultra-lightweight personal AI assistant written in Go, designed to run on resource-constrained hardware ($10 devices, <10MB RAM). It connects messaging channels (Telegram, Discord, Slack, WeChat, etc.) to LLM providers (OpenAI, Anthropic, etc.) via an agent loop that supports tool use, memory, and scheduled tasks. Forked from [nanobot](https://github.com/HKUDS/nanobot) and rewritten in Go.

## Build & Development Commands

```bash
make build          # Build binary (runs go generate first)
make test           # Run all tests
make fmt            # Format code (uses golangci-lint formatters)
make lint           # Run full linter suite
make vet            # Static analysis
make check          # Full pre-commit check: deps + fmt + vet + test
make generate       # Run go generate only
make clean          # Remove build artifacts
```

Run a single test:
```bash
CGO_ENABLED=0 go test -run TestName -v ./pkg/session/
```

Build uses `CGO_ENABLED=0` by default and the `-tags stdjson` build tag. Binary output goes to `build/`.

## Architecture

### Core Flow: Channel → Bus → Agent Loop → Provider

1. **Channels** (`pkg/channels/`) — Adapters for messaging platforms. Each channel implements a common interface and converts platform messages into internal `bus.InboundMessage`. Channels include: Telegram, Discord, Slack, Feishu/Lark, DingTalk, WeChat (wecom), WhatsApp, IRC, Matrix, LINE, QQ, OneBot, webhook, and hardware channels (MaixCam, Pico).

2. **Message Bus** (`pkg/bus/`) — Decouples channels from the agent loop. Channels publish inbound messages; the agent loop subscribes.

3. **Agent Loop** (`pkg/agent/`) — The core orchestration layer. `AgentLoop` receives messages from the bus, resolves which agent instance to use (via `AgentRegistry`), builds conversation context (system prompt + session history + tool definitions), calls the LLM provider, processes tool calls in a loop, and sends the final response back through the channel manager. Supports multi-agent routing via `pkg/routing/`.

4. **Providers** (`pkg/providers/`) — LLM backend adapters. Implements `LLMProvider` interface (`Chat` + `GetDefaultModel`). Includes OpenAI-compatible, Anthropic (native SDK), Claude CLI, Codex CLI, GitHub Copilot, and a factory pattern for generic HTTP providers. Provider fallback chains (`FallbackChain`) handle automatic failover with classification (auth, rate limit, billing, timeout, overloaded).

5. **Tools** (`pkg/tools/`) — Agent capabilities exposed as function calls. Built-in tools: shell execution, filesystem operations, file editing, web fetch, message send, cron scheduling, MCP tool proxy, sub-agent spawning, I2C/SPI hardware interfaces. Tools are registered in `pkg/tools/registry.go`.

6. **Skills** (`pkg/skills/`) — Higher-level composable behaviors built on top of tools, stored in `~/.picoclaw/workspace/skills/`.

### Key Supporting Packages

- **`pkg/config/`** — JSON config loading with env var override support (via `caarlos0/env`). Config struct in `config.go`, defaults in `defaults.go`. Supports config migration across versions.
- **`pkg/session/`** — Conversation session/history management with per-chat isolation.
- **`pkg/memory/`** — Persistent agent memory across sessions.
- **`pkg/state/`** — Runtime state management (active sessions, typing indicators).
- **`pkg/mcp/`** — Model Context Protocol client for connecting to external MCP servers as tool sources.
- **`pkg/media/`** — Media file handling (images, audio, documents) across channels.
- **`pkg/voice/`** — Speech-to-text transcription.
- **`pkg/commands/`** — Slash-command registry for channel-level commands.

### CLI Structure (`cmd/`)

- **`cmd/picoclaw/`** — Main CLI binary. Uses `spf13/cobra`. Subcommands: `agent`, `gateway`, `auth`, `cron`, `status`, `skills`, `onboard`, `migrate`, `version`.
- **`cmd/picoclaw-launcher/`** — Web-based launcher with auth (serves a management UI).
- **`cmd/picoclaw-launcher-tui/`** — Terminal UI launcher using `tview`.

### Channel Interface Patterns

Channels can optionally implement capability interfaces beyond the base channel:
- `TypingCapable` — show typing indicators
- `MessageEditor` — edit sent messages
- `ReactionCapable` — add reactions to messages
- `PlaceholderCapable` — send "thinking..." placeholders
- `CommandRegistrarCapable` — register slash commands on the platform

## Code Style & Linting

- **Formatters**: gofmt, gofumpt, goimports, gci (import ordering: stdlib → external → local module), golines (max 120 chars)
- `interface{}` → `any` rewrite rule is enforced
- `golangci-lint` v2 with version 2 config format (`.golangci.yaml`)
- Import order: standard library, then external packages, then `github.com/sipeed/picoclaw/...`
- Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/) in imperative mood

## CI Pipeline (PR checks)

PRs must pass three jobs: **Lint** (golangci-lint), **Security Check** (govulncheck), and **Tests** (go test ./...). All run `go generate ./...` first.