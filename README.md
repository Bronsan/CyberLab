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

<!-- Language tabs -->
<div align="center">
  <table>
    <tr>
      <td><a href="#中文"><kbd>🇨🇳 中文</kbd></a></td>
      <td><a href="#english"><kbd>🇺🇸 English</kbd></a></td>
      <td><a href="#日本語"><kbd>🇯🇵 日本語</kbd></a></td>
    </tr>
  </table>
</div>

---

<!-- ==================== 中文 ==================== -->

<a name="中文"></a>

# 🇨🇳 CyberLab

> **像真正的黑客一样练习**

CyberLab 是一个开源的在线网络安全靶场平台，为安全爱好者和渗透测试学习者提供**隔离的、动态的、真实漏洞复现环境**。每个挑战都会为你创建一个独立的 Docker 容器，你可以像在真实环境中一样进行漏洞挖掘和利用。

不同于传统 CTF 平台将所有用户连接到同一个站点，CyberLab 为 **每个用户生成专属环境**，保证互不干扰、数据隔离，并且支持自动回收释放资源。

---

## 📸 预览

<table>
  <tr>
    <td align="center"><b>首页演示</b></td>
    <td align="center"><b>挑战详情</b></td>
    <td align="center"><b>排行榜</b></td>
  </tr>
  <tr>
    <td><img src="https://placehold.co/600x400/0a0a0a/00ff41?text=CyberLab+Landing&font=montserrat" alt="Landing" width="100%"/></td>
    <td><img src="https://placehold.co/600x400/0a0a0a/00ff41?text=Challenge+Detail&font=montserrat" alt="Challenge" width="100%"/></td>
    <td><img src="https://placehold.co/600x400/0a0a0a/00ff41?text=Leaderboard&font=montserrat" alt="Ranking" width="100%"/></td>
  </tr>
</table>

---

## ✨ 核心特性

### 🎯 动态靶场环境
每个挑战创建**独立的 Docker 容器**，分配随机端口，返回专属访问地址。用户之间完全隔离，互不影响。

### 🔓 真实漏洞场景
覆盖主流 Web 安全漏洞类型：**SQL 注入、XSS 跨站脚本、RCE 远程命令执行、SSRF 服务端请求伪造、文件上传绕过、JWT 安全问题**等。基于真实 CVE 构建。

### 🤖 AI 辅助解题
卡住了？内置 AI 助手（接入 OpenAI）提供**上下文相关的提示**，引导思路但不直接给出答案。同时支持**上传源码进行 AI 审计**，自动识别漏洞。

### 🏆 CTF 计分体系
提交 Flag、获取积分、冲击排行榜。支持**总榜、周榜、月榜**三种维度，Redis 缓存实时更新。

### ♻️ 自动资源回收
定时任务每分钟扫描过期容器，**自动销毁 + 归还端口 + 释放磁盘**。默认超时 30 分钟，题目可独立配置。

### 🌐 国际化 + 主题切换
内置**中、英、日三语**支持，一键切换。**深色/浅色/跟随系统**三种主题模式，适配任何环境。

---

## 🏗️ 技术架构

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
│  Tailwind CSS   │ │  60+ RESTful 接口     │
│  shadcn/ui      │ │  WebSocket 实时通信   │
└─────────────────┘ └──────┬────────────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                  │
  ┌──────▼──────┐  ┌──────▼──────┐  ┌────────▼──────┐
  │   MySQL 8.0 │  │   Redis 7   │  │  Docker Engine│
  │   持久化存储 │  │  缓存/排行榜 │  │  动态容器管理 │
  └─────────────┘  └─────────────┘  └───────────────┘
```

| 层级 | 技术 | 用途 |
|------|------|------|
| **前端** | Next.js 15 + React 19 + TypeScript + Tailwind CSS + shadcn/ui | SPA + SSR，响应式 UI |
| **状态管理** | Zustand | 轻量全局状态（用户/认证/WebSocket） |
| **动画** | Framer Motion | 页面过渡、卡片动画 |
| **后端** | Go + Gin + GORM + JWT + Zap | RESTful API，高性能 |
| **数据库** | MySQL 8.0 (utf8mb4) + Redis 7 | 业务存储 + 缓存加速 |
| **容器** | Docker CLI | 动态靶场创建/销毁 |
| **实时通信** | Gorilla WebSocket | 容器状态推送、通知 |
| **部署** | Docker Compose + Nginx | 一键部署，反向代理 |

---

## 🚀 快速开始

### 前置依赖

- Go 1.23+
- MySQL 8.0+
- Redis 7+
- Docker（可选，用于动态靶场）
- Node.js 20+

### 启动后端

```bash
cd backend
export GOPROXY=https://proxy.golang.org,direct  # Linux/macOS
# 或 Windows:
set GOPROXY=https://proxy.golang.org,direct

# 配置数据库（编辑 config/config.yaml）
# 然后启动：
go run ./cmd/server
```

启动后访问：
- API: `http://localhost:8080`
- Swagger 文档: `http://localhost:8080/swagger/index.html`
- WebSocket: `ws://localhost:8080/ws`

### 启动前端

```bash
cd frontend
npm install
npm run dev
```

访问: `http://localhost:3000`

### Docker Compose 一键部署

```bash
cd deploy
cp ../backend/.env.example .env
docker-compose up -d
```

---

## 📖 API 参考

**基础路径**: `/api/v1`

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/auth/register` | 用户注册 |
| `POST` | `/auth/login` | 用户登录 |
| `GET` | `/challenges` | 挑战列表（分页/筛选） |
| `GET` | `/challenges/:id` | 挑战详情 |
| `GET` | `/ranking/global` | 总排行榜 |
| `GET` | `/ranking/week` | 周排行榜 |
| `GET` | `/ranking/month` | 月排行榜 |
| `GET` | `/announcement` | 公告列表 |

### 认证接口（需 Bearer Token）

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/auth/profile` | 个人信息 |
| `PUT` | `/auth/profile` | 修改资料 |
| `POST` | `/container/start` | 启动靶场环境 |
| `GET` | `/container/status/:id` | 容器状态 |
| `POST` | `/container/stop` | 停止容器 |
| `POST` | `/submit` | 提交 Flag |
| `POST` | `/ai/hint` | AI 解题提示 |
| `POST` | `/ai/audit` | AI 代码审计 |

### WebSocket

```
ws://host/ws?userId={id}
```

**事件推送**：
- `container_created` — 容器创建成功
- `container_destroyed` — 容器已销毁
- `container_error` — 容器异常
- `ranking_update` — 排行榜刷新
- `announcement` — 系统公告

---

## 🗂️ 项目结构

```
CyberLab/
├── backend/                    # Go 后端服务
│   ├── cmd/server/main.go      # 入口
│   ├── internal/
│   │   ├── config/             # 配置加载
│   │   ├── handlers/           # HTTP 处理器
│   │   ├── middleware/         # JWT 认证
│   │   ├── models/            # 数据模型（9 表）
│   │   ├── repository/        # 数据访问层
│   │   ├── services/          # 业务逻辑
│   │   ├── scheduler/         # 容器回收
│   │   └── router/            # 路由注册
│   └── pkg/
│       ├── docker/            # Docker CLI 封装
│       ├── ws/                # WebSocket Hub
│       └── utils/             # 工具函数
│
├── frontend/                   # Next.js 前端
│   └── src/
│       ├── app/               # 页面路由
│       ├── components/        # UI 组件
│       ├── locales/           # 多语言（中/英/日）
│       ├── providers/         # 主题 + 语言 Provider
│       ├── services/          # API 客户端
│       └── store/             # Zustand 状态管理
│
├── deploy/                     # 部署配置
│   ├── docker-compose.yml
│   └── nginx/nginx.conf
│
└── docs/                       # 开发文档
```

---

## 📄 License

**MIT** — 可自由使用、修改、商用。

---

<!-- ==================== ENGLISH ==================== -->

<a name="english"></a>

# 🇺🇸 CyberLab

> **Practice Like a Real Hacker**

CyberLab is an open-source cybersecurity training platform that provides **isolated, dynamic, real-vulnerability reproduction environments** for security enthusiasts and penetration testing learners. Each challenge spawns a dedicated Docker container just for you — explore and exploit in a realistic, risk-free environment.

Unlike traditional CTF platforms that connect all users to the same site, CyberLab creates **per-user exclusive environments** ensuring full isolation, data separation, and automatic resource reclamation.

---

## ✨ Key Features

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

## 🚀 Quick Start

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

API: `http://localhost:8080` · Swagger: `http://localhost:8080/swagger/index.html`

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

## 📄 License

**MIT** — Free to use, modify, and distribute.

---

<!-- ==================== 日本語 ==================== -->

<a name="日本語"></a>

# 🇯🇵 CyberLab

> **本物のハッカーのように練習しよう**

CyberLab は、セキュリティ愛好家やペネトレーションテスト学習者のためのオープンソースのサイバーセキュリティトレーニングプラットフォームです。**分離された動的な脆弱性再現環境**を提供し、各チャレンジはあなただけのために専用の Docker コンテナを起動します。

従来の CTF プラットフォームとは異なり、CyberLab は**ユーザーごとに専用環境**を作成するため、相互干渉がなく、データは完全に分離され、リソースは自動的に解放されます。

---

## ✨ 主な機能

### 🎯 動的ラボ環境
各チャレンジは**独立した Docker コンテナ**を作成し、ランダムなポートと専用 URL を割り当てます。ユーザー間の干渉は一切ありません。

### 🔓 実際の脆弱性
**SQL インジェクション、XSS、RCE、SSRF、ファイルアップロードバイパス、JWT セキュリティ問題**など、実際の CVE に基づいた脆弱性を練習できます。

### 🤖 AI アシスト
チャレンジで行き詰まりましたか？ビルトイン AI アシスタント（OpenAI 搭載）が**文脈に応じたヒント**を提供し、答えを直接教えずに考え方を導きます。ソースコードをアップロードして**AI セキュリティ監査**を受けることも可能です。

### ♻️ 自動リソース回収
バックグラウンドワーカーが毎分期限切れコンテナをスキャンし、**自動破棄、ポート解放、ディスク解放**を実施。デフォルトのタイムアウトは 30 分、チャレンジごとに設定可能です。

---

## 🚀 クイックスタート

### バックエンド

```bash
cd backend
export GOPROXY=https://proxy.golang.org,direct
go run ./cmd/server
```

### フロントエンド

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

## 📄 ライセンス

**MIT** — 自由に使用、修正、配布できます。
