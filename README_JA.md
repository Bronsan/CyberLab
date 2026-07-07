<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://img.shields.io/badge/CyberLab-00ff41?style=for-the-badge&logo=hackthebox&logoColor=black">
    <img src="https://img.shields.io/badge/CyberLab-059669?style=for-the-badge&logo=hackthebox&logoColor=white" alt="CyberLab">
  </picture>
</p>

<p align="center">
  <strong><code>本物のハッカーのように練習しよう</code></strong>
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
      <td><a href="README_EN.md"><kbd>🇺🇸 English</kbd></a></td>
      <td><kbd>🇯🇵 日本語</kbd></td>
    </tr>
  </table>
</div>

---

# CyberLab

> **本物のハッカーのように練習しよう**

CyberLab は、セキュリティ愛好家やペネトレーションテスト学習者のためのオープンソースのサイバーセキュリティトレーニングプラットフォームです。**分離された動的な脆弱性再現環境**を提供し、各チャレンジはあなただけのために専用の Docker コンテナを起動します。

従来の CTF プラットフォームとは異なり、CyberLab は**ユーザーごとに専用環境**を作成するため、相互干渉がなく、データは完全に分離され、リソースは自動的に解放されます。

---

## 主な機能

### 🎯 動的ラボ環境
各チャレンジは**独立した Docker コンテナ**を作成し、ランダムなポートと専用 URL を割り当てます。ユーザー間の干渉は一切ありません。

### 🔓 実際の脆弱性
**SQL インジェクション、XSS、RCE、SSRF、ファイルアップロードバイパス、JWT セキュリティ問題**など、実際の CVE に基づいた脆弱性を練習できます。

### 🤖 AI アシスト
チャレンジで行き詰まりましたか？ビルトイン AI アシスタント（OpenAI 搭載）が**文脈に応じたヒント**を提供し、答えを直接教えずに考え方を導きます。ソースコードをアップロードして**AI セキュリティ監査**を受けることも可能です。

### 🏆 CTF スコアリング
フラグをキャプチャしてポイントを獲得し、リーダーボードを駆け上がろう。**グローバル、週間、月間**ランキングに対応。Redis によるリアルタイム更新。

### ♻️ 自動リソース回収
バックグラウンドワーカーが毎分期限切れコンテナをスキャンし、**自動破棄、ポート解放、ディスク解放**を実施。デフォルトのタイムアウトは 30 分、チャレンジごとに設定可能です。

### 🌐 i18n + テーマ切替
**中国語、英語、日本語**に対応。ワンクリックで切替可能。**ダーク/ライト/システム**のテーマモードでどんな環境にも適応。

---

## アーキテクチャ

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
│  Tailwind CSS   │ │  60+ RESTful API      │
│  shadcn/ui      │ │  WebSocket リアルタイム│
└─────────────────┘ └──────┬────────────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                  │
  ┌──────▼──────┐  ┌──────▼──────┐  ┌────────▼──────┐
  │   MySQL 8.0 │  │   Redis 7   │  │  Docker Engine│
  │   ストレージ │  │  キャッシュ │  │  コンテナ管理 │
  └─────────────┘  └─────────────┘  └───────────────┘
```

| レイヤー | 技術 | 用途 |
|----------|------|------|
| **フロントエンド** | Next.js 15 + React 19 + TypeScript + Tailwind CSS + shadcn/ui | SPA + SSR、レスポンシブUI |
| **状態管理** | Zustand | 軽量グローバル状態 |
| **アニメーション** | Framer Motion | ページ遷移、カードアニメーション |
| **バックエンド** | Go + Gin + GORM + JWT + Zap | RESTful API、高パフォーマンス |
| **データベース** | MySQL 8.0 (utf8mb4) + Redis 7 | 永続化 + キャッシュ |
| **コンテナ** | Docker CLI | 動的ラボ環境の管理 |
| **リアルタイム** | Gorilla WebSocket | コンテナ状態プッシュ、通知 |
| **デプロイ** | Docker Compose + Nginx | ワンクリックデプロイ、リバースプロキシ |

---

## クイックスタート

### 前提条件

- Go 1.23+
- MySQL 8.0+
- Redis 7+
- Docker（動的チャレンジ用、オプション）
- Node.js 20+

### バックエンド

```bash
cd backend
export GOPROXY=https://proxy.golang.org,direct
go run ./cmd/server
```

- API: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- WebSocket: `ws://localhost:8080/ws`

### フロントエンド

```bash
cd frontend
npm install
npm run dev
```

アクセス: `http://localhost:3000`

### Docker Compose

```bash
cd deploy
docker-compose up -d
```

---

## ライセンス

**MIT** — 自由に使用、修正、配布できます。
