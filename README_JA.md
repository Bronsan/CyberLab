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
      <td><a href="README_EN.md"><kbd>🇺🇸 English</kbd></a></td>
      <td><kbd>🇯🇵 日本語</kbd></td>
    </tr>
  </table>
  <br/>
  <a href="https://bronsan.github.io/CyberLab/">🌐 ライブデモ</a>
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

### 🛡️ セキュリティ対策

| 対策 | 説明 |
|------|------|
| **フラグのハッシュ化** | SHA-256 HMAC + チャレンジ別ソルト、DB漏洩時にフラグ保護 |
| **レート制限** | ログイン 10回/分、登録 5回/分、ブルートフォース防止 |
| **JWT 認証** | WebSocket と API は Bearer Token 必須 |
| **Docker ホワイトリスト** | 許可されたイメージのみ実行可能、エスケープ防止 |
| **CORS 制限** | ホワイトリストドメインのみ、ワイルドカード禁止 |
| **セキュリティヘッダー** | X-Frame-Options / XSS-Protection / Content-Type-Options |
| **パスワード・鍵はコードなし** | すべて環境変数から注入 |

---

## ライブデモ

フロントエンドは GitHub Pages にデプロイされています：

**https://bronsan.github.io/CyberLab/**

> 注：フロントエンドは静的サイトです。完全な操作（ログイン、チャレンジ、コンテナ等）にはバックエンド API の起動が必要です。

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
| **フロントエンド** | Next.js 15 + React 19 + TypeScript + Tailwind CSS + shadcn/ui | SPA + SSR、i18n、テーマ切替 |
| **状態管理** | Zustand | 軽量グローバル状態 |
| **アニメーション** | Framer Motion | ページ遷移、カードアニメーション |
| **バックエンド** | Go + Gin + GORM + JWT + Zap | RESTful API、高パフォーマンス |
| **データベース** | MySQL 8.0 (utf8mb4) + Redis 7 | 永続化 + キャッシュ |
| **コンテナ** | Docker CLI | 動的ラボ環境、ホワイトリスト |
| **リアルタイム** | Gorilla WebSocket | コンテナ状態プッシュ、通知 |
| **セキュリティ** | SHA-256 HMAC + レート制限 + CORS制限 | 漏洩防止、エスケープ防止 |

---

## 📥 ダウンロード

| プラットフォーム | ファイル | サイズ |
|------------------|----------|--------|
| 🪟 Windows | [cyberlab-v1.0.0-windows-amd64.zip](https://github.com/Bronsan/CyberLab/releases/latest) | ~9.6 MB |
| 🐧 Linux | [cyberlab-v1.0.0-linux-amd64.tar.gz](https://github.com/Bronsan/CyberLab/releases/latest) | ~9.3 MB |
| 🍎 macOS Intel | [cyberlab-v1.0.0-darwin-amd64.tar.gz](https://github.com/Bronsan/CyberLab/releases/latest) | ~9.5 MB |
| 🍎 macOS M1/M2 | [cyberlab-v1.0.0-darwin-arm64.tar.gz](https://github.com/Bronsan/CyberLab/releases/latest) | ~8.7 MB |
| 📱 Android / iOS | PWA — ブラウザで開く → ホーム画面に追加 | 無料 |

```bash
# Linux/macOS クイックスタート
tar xzf cyberlab-v1.0.0-linux-amd64.tar.gz
chmod +x run.sh && ./run.sh
# 初回実行で .env が自動生成 — DBパスワードを編集して再実行

# Windows: 解凍して run.bat をダブルクリック
```

> すべてのリリース → [GitHub Releases](https://github.com/Bronsan/CyberLab/releases)

---

## クイックスタート（ソースから）

### 前提条件

- Go 1.23+
- MySQL 8.0+
- Redis 7+
- Docker（動的チャレンジ用、オプション）
- Node.js 20+

### 必要な環境変数

```bash
export JWT_SECRET=your-random-secret-key
export DB_PASSWORD=your-database-password
export FLAG_HASH_SALT=your-random-64-char-salt
# AI（オプション）：
export AI_API_KEY=sk-your-openai-api-key
```

### バックエンド

```bash
cd backend
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

## GitHub Pages デプロイ

フロントエンドは GitHub Pages に自動デプロイされます：

```
https://bronsan.github.io/CyberLab/
```

`master` ブランチにプッシュするたびに、GitHub Actions が静的サイトをビルドして公開します。

---

## ライセンス

**MIT** — 自由に使用、修正、配布できます。
