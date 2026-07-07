<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://img.shields.io/badge/CyberLab-00ff41?style=for-the-badge&logo=hackthebox&logoColor=black">
    <img src="https://img.shields.io/badge/CyberLab-059669?style=for-the-badge&logo=hackthebox&logoColor=white" alt="CyberLab">
  </picture>
</p>

<p align="center">
  <strong><code>Practice Like a Real Hacker</code></strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/Next.js-15-black?logo=next.js" alt="Next.js"/>
  <img src="https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white" alt="MySQL"/>
  <img src="https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white" alt="Redis"/>
  <img src="https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white" alt="Docker"/>
  <img src="https://img.shields.io/badge/TypeScript-3178C6?logo=typescript&logoColor=white" alt="TypeScript"/>
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License"/>
</p>

---

<div align="center">
  <table>
    <tr>
      <td><a href="README.md"><kbd>🇨🇳 中文</kbd></a></td>
      <td><kbd>🇺🇸 English</kbd></td>
      <td><a href="README_JA.md"><kbd>🇯🇵 日本語</kbd></a></td>
    </tr>
  </table>
</div>

---

# CyberLab

> **Practice Like a Real Hacker**

CyberLab is an open-source cybersecurity training platform that provides **isolated, dynamic, real-vulnerability reproduction environments** for security enthusiasts and penetration testing learners. Each challenge spawns a dedicated Docker container just for you — explore and exploit in a realistic, risk-free environment.

Unlike traditional CTF platforms that connect all users to the same site, CyberLab creates **per-user exclusive environments** ensuring full isolation, data separation, and automatic resource reclamation.

---

## Key Features

### 🎯 Dynamic Lab Environments
Every challenge creates an **isolated Docker container** with a randomly assigned port and unique access URL. Users are completely isolated from each other.

### 🔓 Real-World Vulnerabilities
Covers mainstream web security vulnerabilities: **SQL Injection, XSS, RCE, SSRF, File Upload Bypass, JWT Security Issues**, and more. Built on real CVEs.

### 🤖 AI-Powered Assistance
Stuck on a challenge? The built-in AI assistant (powered by OpenAI) provides **contextual hints** to guide your thinking without giving away the answer. Also supports **source code upload for AI security audit**.

### 🏆 CTF Scoring System
Capture flags, earn points, and climb the leaderboard. Supports **global, weekly, and monthly rankings** with Redis-powered real-time updates.

### ♻️ Auto Resource Reclamation
A background worker scans for expired containers every minute — **auto-destroy, release ports, free disk space**. Default 30-minute timeout, configurable per challenge.

### 🌐 i18n + Theme Switching
Built-in **Chinese, English, and Japanese** language support with one-click switching. **Dark/Light/System** theme modes adapt to any environment.

---

## Architecture

```
┌─────────────────────────────────────────────┐
│                Browser / Client              │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│           Nginx Reverse Proxy                │
└────────┬─────────────────┬───────────────────┘
         │                 │
┌────────▼────────┐ ┌──────▼────────────────┐
│  Next.js 15     │ │  Go Gin API           │
│  Frontend       │ │  Backend              │
│  React 19       │ │  GORM + JWT + Zap     │
│  Tailwind CSS   │ │  60+ RESTful APIs     │
│  shadcn/ui      │ │  WebSocket Realtime   │
└─────────────────┘ └──────┬────────────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                  │
  ┌──────▼──────┐  ┌──────▼──────┐  ┌────────▼──────┐
  │   MySQL 8.0 │  │   Redis 7   │  │  Docker Engine│
  │   Storage   │  │  Cache/Rank │  │  Containers   │
  └─────────────┘  └─────────────┘  └───────────────┘
```

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Frontend** | Next.js 15 + React 19 + TypeScript + Tailwind CSS + shadcn/ui | SPA + SSR, responsive UI |
| **State** | Zustand | Lightweight global state |
| **Animation** | Framer Motion | Page transitions, card animations |
| **Backend** | Go + Gin + GORM + JWT + Zap | RESTful API, high performance |
| **Database** | MySQL 8.0 (utf8mb4) + Redis 7 | Persistence + cache |
| **Container** | Docker CLI | Dynamic lab orchestration |
| **Realtime** | Gorilla WebSocket | Container status push, notifications |
| **Deploy** | Docker Compose + Nginx | One-click deployment, reverse proxy |

---

## Quick Start

### Prerequisites

- Go 1.23+
- MySQL 8.0+
- Redis 7+
- Docker (optional, for dynamic challenges)
- Node.js 20+

### Backend

```bash
cd backend
export GOPROXY=https://proxy.golang.org,direct
go run ./cmd/server
```

- API: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- WebSocket: `ws://localhost:8080/ws`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Visit: `http://localhost:3000`

### Docker Compose

```bash
cd deploy
docker-compose up -d
```

---

## API Reference

Base path: `/api/v1`

### Public Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/auth/register` | User registration |
| `POST` | `/auth/login` | User login |
| `GET` | `/challenges` | Challenge list (paginated) |
| `GET` | `/challenges/:id` | Challenge detail |
| `GET` | `/ranking/global` | Global leaderboard |
| `GET` | `/ranking/week` | Weekly leaderboard |
| `GET` | `/ranking/month` | Monthly leaderboard |
| `GET` | `/announcement` | Announcements |

### Authenticated Endpoints (Bearer Token)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/auth/profile` | User profile |
| `PUT` | `/auth/profile` | Update profile |
| `POST` | `/container/start` | Start challenge environment |
| `GET` | `/container/status/:id` | Container status |
| `POST` | `/container/stop` | Stop container |
| `POST` | `/submit` | Submit flag |
| `POST` | `/ai/hint` | AI hint |
| `POST` | `/ai/audit` | AI code audit |

### WebSocket

```
ws://host/ws?userId={id}
```

Events:
- `container_created` — Container created
- `container_destroyed` — Container destroyed
- `container_error` — Container error
- `ranking_update` — Leaderboard refresh
- `announcement` — System announcement

---

## Project Structure

```
CyberLab/
├── backend/                    # Go backend service
│   ├── cmd/server/main.go      # Entry point
│   ├── internal/               # Business logic
│   │   ├── config/             # Configuration
│   │   ├── handlers/           # HTTP handlers
│   │   ├── middleware/         # JWT auth
│   │   ├── models/             # Data models (9 tables)
│   │   ├── repository/         # Data access layer
│   │   ├── services/           # Business services
│   │   ├── scheduler/          # Container reclaimer
│   │   └── router/             # Route registration
│   └── pkg/                    # Packages
│       ├── docker/             # Docker CLI wrapper
│       ├── ws/                 # WebSocket hub
│       └── utils/              # Utilities
│
├── frontend/                   # Next.js frontend
│   └── src/
│       ├── app/                # Pages & routes
│       ├── components/         # UI components
│       ├── locales/            # i18n (ZH/EN/JA)
│       ├── providers/          # Theme & language providers
│       ├── services/           # API client
│       └── store/              # Zustand state
│
├── deploy/                     # Deployment config
│   ├── docker-compose.yml
│   └── nginx/nginx.conf
│
└── docs/                       # Documentation
```

---

## License

**MIT** — Free to use, modify, and distribute.
