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
  <br/>
  <img src="https://img.shields.io/badge/status-active-success" alt="Status"/>
  <img src="https://img.shields.io/badge/security-audited-brightgreen" alt="Security"/>
  <img src="https://img.shields.io/github/deployments/Bronsan/CyberLab/github-pages?label=pages" alt="Pages"/>
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
  <br/>
  <a href="https://bronsan.github.io/CyberLab/">🌐 Live Demo</a>
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
Stuck on a challenge? The built-in AI assistant (powered by OpenAI) provides **contextual hints** to guide your thinking without giving away the answer. Supports **source code upload for AI security audit**.

### 🏆 CTF Scoring System
Capture flags, earn points, and climb the leaderboard. Supports **global, weekly, and monthly rankings** with Redis-powered real-time updates.

### ♻️ Auto Resource Reclamation
A background worker scans for expired containers every minute — **auto-destroy, release ports, free disk space**. Default 30-minute timeout, configurable per challenge.

### 🌐 i18n + Theme Switching
Built-in **Chinese, English, and Japanese** language support with one-click switching. **Dark/Light/System** theme modes adapt to any environment.

### 🛡️ Security Hardening

| Measure | Description |
|---------|-------------|
| **Flag Hashing** | SHA-256 HMAC with per-challenge salt — no plaintext flags in DB |
| **Rate Limiting** | Login 10/min, Register 5/min — prevents brute force |
| **JWT Auth** | WebSocket + API both require Bearer Token |
| **Docker Allowlist** | Only approved images can be pulled, preventing container escape |
| **CORS Restriction** | Whitelist-only origins, no wildcard |
| **Security Headers** | X-Frame-Options / XSS-Protection / Content-Type-Options |
| **Zero Secrets in Code** | All passwords and keys injected via environment variables |

---

## Live Demo

The frontend is deployed to GitHub Pages:

**https://bronsan.github.io/CyberLab/**

> Note: The frontend is a static site. Full interactivity (login, challenges, containers) requires the backend API running locally.

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
| **Frontend** | Next.js 15 + React 19 + TypeScript + Tailwind CSS + shadcn/ui | SPA + SSR, i18n, theme switching |
| **State** | Zustand | Lightweight global state |
| **Animation** | Framer Motion | Page transitions, card animations |
| **Backend** | Go + Gin + GORM + JWT + Zap | RESTful API, high performance |
| **Database** | MySQL 8.0 (utf8mb4) + Redis 7 | Persistence + cache |
| **Container** | Docker CLI | Dynamic lab orchestration, allowlist |
| **Realtime** | Gorilla WebSocket | Container status push, notifications |
| **Security** | SHA-256 HMAC + rate limiter + CORS whitelist | Anti-leak, anti-brute-force, anti-escape |

---

## 📥 Download

| Platform | File | Size |
|----------|------|------|
| 🪟 Windows | [cyberlab-v1.0.0-windows-amd64.zip](https://github.com/Bronsan/CyberLab/releases/latest) | ~9.6 MB |
| 🐧 Linux | [cyberlab-v1.0.0-linux-amd64.tar.gz](https://github.com/Bronsan/CyberLab/releases/latest) | ~9.3 MB |
| 🍎 macOS Intel | [cyberlab-v1.0.0-darwin-amd64.tar.gz](https://github.com/Bronsan/CyberLab/releases/latest) | ~9.5 MB |
| 🍎 macOS M1/M2 | [cyberlab-v1.0.0-darwin-arm64.tar.gz](https://github.com/Bronsan/CyberLab/releases/latest) | ~8.7 MB |
| 📱 Android / iOS | PWA — open in browser → Add to Home Screen | free |

```bash
# Linux/macOS quick start
tar xzf cyberlab-v1.0.0-linux-amd64.tar.gz
chmod +x run.sh && ./run.sh
# Auto-generates .env on first run — edit DB password and re-run

# Windows: extract and double-click run.bat
```

> All Releases → [GitHub Releases](https://github.com/Bronsan/CyberLab/releases)

---

## Quick Start (from source)

### Prerequisites

- Go 1.23+
- MySQL 8.0+
- Redis 7+
- Docker (optional, for dynamic challenges)
- Node.js 20+

### Required Environment Variables

```bash
export JWT_SECRET=your-random-secret-key
export DB_PASSWORD=your-database-password
export FLAG_HASH_SALT=your-random-64-char-salt
# AI (optional):
export AI_API_KEY=sk-your-openai-api-key
```

### Backend

```bash
cd backend
go run ./cmd/server
```

- API: `http://localhost:17420`
- Swagger: `http://localhost:17420/swagger/index.html`
- WebSocket: `ws://localhost:17420/ws`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Docker Compose

```bash
cd deploy
docker-compose up -d
```

---

## Security Design

| Category | Measure |
|----------|---------|
| **Auth** | JWT stateless tokens, WebSocket requires JWT |
| **Rate Limit** | Per-IP sliding window rate limiter |
| **Password** | bcrypt hashed |
| **Flag** | SHA-256 HMAC + per-challenge salt |
| **Submissions** | Only flag hashes stored, never plaintext |
| **Docker** | Image allowlist, distroless base image |
| **Config** | Zero secrets committed, all via env vars |
| **CORS** | Whitelist origins, no wildcard |
| **TLS** | Terminated at Nginx reverse proxy |

---

## GitHub Pages Deployment

The frontend automatically deploys to GitHub Pages:

```
https://bronsan.github.io/CyberLab/
```

On every push to `master`, GitHub Actions builds the static site and publishes it.

---

## License

**MIT** — Free to use, modify, and distribute.
