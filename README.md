<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://img.shields.io/badge/CyberLab-00ff41?style=for-the-badge&logo=hackthebox&logoColor=black">
    <img src="https://img.shields.io/badge/CyberLab-059669?style=for-the-badge&logo=hackthebox&logoColor=white" alt="CyberLab">
  </picture>
</p>

<p align="center">
  <strong><code>像真正的黑客一样练习</code></strong>
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
      <td><kbd>🇨🇳 中文</kbd></td>
      <td><a href="README_EN.md"><kbd>🇺🇸 English</kbd></a></td>
      <td><a href="README_JA.md"><kbd>🇯🇵 日本語</kbd></a></td>
    </tr>
  </table>
  <br/>
  <a href="https://bronsan.github.io/CyberLab/">🌐 在线预览</a>
</div>

---

# CyberLab

> **像真正的黑客一样练习**

CyberLab 是一个开源的在线网络安全靶场平台，为安全爱好者和渗透测试学习者提供 **隔离的、动态的、真实漏洞复现环境**。每个挑战都会为你创建一个独立的 Docker 容器，你可以像在真实环境中一样进行漏洞挖掘和利用。

不同于传统 CTF 平台将所有用户连接到同一个站点，CyberLab 为 **每个用户生成专属环境**，保证互不干扰、数据隔离，并且支持自动回收释放资源。

---

## 核心特性

### 🎯 动态靶场环境
每个挑战创建**独立的 Docker 容器**，分配随机端口，返回专属访问地址。用户之间完全隔离，互不影响。

### 🔓 真实漏洞场景
覆盖主流 Web 安全漏洞类型：**SQL 注入、XSS 跨站脚本、RCE 远程命令执行、SSRF 服务端请求伪造、文件上传绕过、JWT 安全问题**等。基于真实 CVE 构建。

### 🤖 AI 辅助解题
卡住了？内置 AI 助手（接入 OpenAI）提供**上下文相关的提示**，引导思路但不直接给出答案。支持**上传源码进行 AI 审计**，自动识别漏洞。

### 🏆 CTF 计分体系
提交 Flag、获取积分、冲击排行榜。支持**总榜、周榜、月榜**三种维度，Redis 缓存实时更新。

### ♻️ 自动资源回收
定时任务每分钟扫描过期容器，**自动销毁 + 归还端口 + 释放磁盘**。默认超时 30 分钟，题目可独立配置。

### 🌐 国际化 + 主题切换
内置**中、英、日三语**支持，一键切换。**深色/浅色/跟随系统**三种主题模式，适配任何环境。

### 🛡️ 安全加固

| 防护措施 | 说明 |
|---------|------|
| **Flag 哈希存储** | SHA-256 HMAC + 每挑战独立盐值，数据库泄露不暴露 Flag |
| **速率限制** | 登录 10 次/分钟、注册 5 次/分钟，防暴力破解 |
| **JWT 认证** | WebSocket 和 API 均需 Bearer Token，无匿名访问 |
| **Docker 白名单** | 镜像拉取限白名单，防止恶意镜像逃逸 |
| **CORS 限定** | 只允许白名单域名跨域访问 |
| **安全响应头** | X-Frame-Options / XSS-Protection / Content-Type-Options |
| **无密钥提交** | 配置文件中不含任何密码或密钥，全部通过环境变量注入 |

---

## 在线演示

项目前端已部署到 GitHub Pages：

**https://bronsan.github.io/CyberLab/**

> 说明：前端为静态站点，后端 API 需本地启动才能体验完整交互（登录、挑战、容器等）。

---

## 技术架构

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
| **前端** | Next.js 15 + React 19 + TypeScript + Tailwind CSS + shadcn/ui | SPA + SSR，响应式 UI，三语 i18n |
| **状态管理** | Zustand | 轻量全局状态 |
| **动画** | Framer Motion | 页面过渡、卡片动画 |
| **后端** | Go + Gin + GORM + JWT + Zap | RESTful API，高性能 |
| **数据库** | MySQL 8.0 (utf8mb4) + Redis 7 | 业务存储 + 缓存加速 |
| **容器** | Docker CLI | 动态靶场创建/销毁，镜像白名单 |
| **实时通信** | Gorilla WebSocket | 容器状态推送、通知 |
| **部署** | Docker Compose + Nginx | 一键部署，反向代理 |
| **安全** | SHA-256 HMAC + 限流中间件 + 白名单 CORS | 防泄露、防爆破、防逃逸 |

---

## 快速开始

### 前置依赖

- Go 1.23+
- MySQL 8.0+
- Redis 7+
- Docker（可选，用于动态靶场）
- Node.js 20+

### 环境变量（必须设置）

```bash
export JWT_SECRET=your-random-secret-key
export DB_PASSWORD=your-database-password
export FLAG_HASH_SALT=your-random-64-char-salt
# AI 功能可选：
export AI_API_KEY=sk-your-openai-api-key
```

### 启动后端

```bash
cd backend
go run ./cmd/server
```

- API: `http://localhost:17420`
- Swagger 文档: `http://localhost:17420/swagger/index.html`
- WebSocket: `ws://localhost:17420/ws`

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
# 编辑 .env 填入实际密钥
docker-compose up -d
```

---

## API 参考

基础路径: `/api/v1`

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/auth/register` | 用户注册（限流 5次/分钟） |
| `POST` | `/auth/login` | 用户登录（限流 10次/分钟） |
| `GET` | `/challenges` | 挑战列表（分页/筛选） |
| `GET` | `/challenges/:id` | 挑战详情 |
| `GET` | `/ranking/global` | 总排行榜 |
| `GET` | `/ranking/week` | 周排行榜 |
| `GET` | `/ranking/month` | 月排行榜 |

### 认证接口（需 Bearer Token）

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/auth/profile` | 个人信息 |
| `PUT` | `/auth/profile` | 修改资料 |
| `POST` | `/container/start` | 启动靶场环境 |
| `POST` | `/container/stop` | 停止容器 |
| `POST` | `/submit` | 提交 Flag（哈希对比） |
| `POST` | `/ai/hint` | AI 解题提示 |
| `POST` | `/ai/audit` | AI 代码审计 |

### WebSocket

```
ws://host/ws?token={jwt}
```

---

## 项目结构

```
CyberLab/
├── README.md               # 中文文档
├── README_EN.md            # English
├── README_JA.md            # 日本語
├── backend/                # Go 后端服务
│   ├── cmd/server/         # 入口
│   ├── internal/           # 业务逻辑
│   │   ├── config/         # 配置加载（环境变量覆盖）
│   │   ├── handlers/       # HTTP 处理器
│   │   ├── middleware/     # JWT 认证 + 速率限制
│   │   ├── models/         # 9 张数据表
│   │   ├── repository/     # 数据访问层
│   │   ├── services/       # 业务逻辑
│   │   ├── scheduler/      # 容器回收（1 分钟间隔）
│   │   └── router/         # 路由注册
│   └── pkg/
│       ├── docker/         # Docker CLI 封装 + 镜像白名单
│       ├── ws/             # WebSocket Hub（JWT 认证）
│       └── utils/          # JWT / 响应格式 / 日志 / 哈希
│
├── frontend/               # Next.js 前端
│   └── src/
│       ├── app/            # 7 个页面路由
│       ├── components/     # UI 组件
│       ├── locales/        # 中/英/日 翻译文件
│       ├── providers/      # 主题 + 语言 Provider
│       ├── services/       # API 客户端
│       └── store/          # Zustand 状态管理
│
├── deploy/                 # Docker Compose + Nginx
└── .github/workflows/      # GitHub Pages 自动部署
```

---

## 安全设计

| 类别 | 措施 |
|------|------|
| **认证** | JWT 无状态 Token，WebSocket 强制 JWT 验证 |
| **限流** | 基于 IP 的滑动窗口限流器 |
| **密码** | bcrypt 哈希存储 |
| **Flag** | SHA-256 HMAC + 每挑战独立盐值 |
| **提交记录** | Flag 仅存哈希，不存明文 |
| **Docker** | 镜像白名单，distroless 基础镜像 |
| **配置** | 密钥零提交，全部环境变量注入 |
| **CORS** | 白名单域名，禁止通配符 |
| **TLS** | 生产环境需 Nginx 反代终结 SSL |

---

## 部署到 GitHub Pages

项目前端自动部署到 GitHub Pages：

```
https://bronsan.github.io/CyberLab/
```

每次推送 `master` 分支，GitHub Actions 自动构建静态站点并发布。

---

## License

**MIT** — 可自由使用、修改、商用。
