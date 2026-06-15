# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

### Go Backend (api/)
- **Development**: `cd api && go run main.go` (uses config.toml)
- **Build**: `cd api && make` (builds both amd64 and arm64 binaries)
- **Individual builds**: `make amd64` or `make arm64`
- **Clean**: `make clean`
- **Config**: Copy `config.sample.toml` to `config.toml` and configure

### Web Frontend (web/)
- **Development**: `cd web && npm run dev` (runs on Vite dev server with --host)
- **Build**: `cd web && npm run build`
- **Lint**: `cd web && npm run lint` (ESLint with auto-fix)

### Testing
- Backend tests: `cd api/test && bash run_crawler_test.sh`
- No specific frontend test configuration found

## Project Architecture

### Backend (Go)
- **Framework**: Gin web framework with dependency injection via uber-go/fx
- **Database**: GORM with MySQL, Redis for caching, LevelDB for local storage
- **Authentication**: JWT tokens with Redis session storage
- **Middleware**: CORS, authorization, parameter handling, static resource serving
- **Structure**:
  - `handler/`: HTTP request handlers (REST API endpoints)
  - `service/`: Business logic services (AI integrations, payments, etc.)
  - `store/`: Database models and data access layer
  - `core/`: Application server and middleware configuration
  - `utils/`: Utility functions and helpers

### Frontend (Vue.js)
- **Framework**: Vue 3 with Composition API
- **UI Components**: Element Plus + Vant (mobile components)
- **State Management**: Pinia
- **Routing**: Vue Router with nested routes
- **Build Tool**: Vite
- **CSS**: Stylus preprocessor with Tailwind CSS utilities
- **Features**: Responsive design (desktop/mobile views), theme switching (dark/light)

### Key Features
- **AI Chat**: Multiple chat models and conversation management
- **Image Generation**: MidJourney, Stable Diffusion, DALL-E integration
- **Audio/Video**: Suno music creation, Luma/KeLing video generation
- **User Management**: Authentication, payments, power logs, invitations
- **Admin Panel**: Comprehensive management interface

### Database Models
Key entities: User, ChatItem, ChatMessage, ChatRole, ChatModel, Order, Product, AdminUser, and various job types for AI services.

### API Structure
- User APIs: `/api/user/*` (auth, profile, settings)
- Chat APIs: `/api/chat/*` (conversations, messages)
- AI Service APIs: `/api/mj/*`, `/api/sd/*`, `/api/dall/*`, `/api/suno/*`, `/api/video/*`
- Admin APIs: `/api/admin/*` (management functions)

### Configuration
- Backend: TOML configuration file (`config.toml`)
- Database: MySQL with automatic migrations
- Services: Redis, various AI API integrations
- File Storage: Local, Aliyun OSS, MinIO, Qiniu options

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **geekai** (8484 symbols, 23360 relationships, 240 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> If any GitNexus tool warns the index is stale, run `npx gitnexus analyze` in terminal first.

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `gitnexus_impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `gitnexus_detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `gitnexus_query({query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `gitnexus_context({name: "symbolName"})`.

## Never Do

- NEVER edit a function, class, or method without first running `gitnexus_impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `gitnexus_rename` which understands the call graph.
- NEVER commit changes without running `gitnexus_detect_changes()` to check affected scope.

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/geekai/context` | Codebase overview, check index freshness |
| `gitnexus://repo/geekai/clusters` | All functional areas |
| `gitnexus://repo/geekai/processes` | All execution flows |
| `gitnexus://repo/geekai/process/{name}` | Step-by-step execution trace |

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->
