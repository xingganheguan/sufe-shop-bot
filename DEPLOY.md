# Shop Bot 部署指南

## 目录
1. [系统要求](#系统要求)
2. [快速开始](#快速开始)
3. [部署方式](#部署方式)
4. [配置说明](#配置说明)
5. [管理员设置](#管理员设置)
6. [Nginx 反向代理](#nginx-反向代理)
7. [更新与维护](#更新与维护)
8. [故障排除](#故障排除)

---

## 系统要求

| 要求 | 最低配置 | 推荐配置 |
|------|----------|----------|
| CPU | 1 核心 | 2+ 核心 |
| 内存 | 512MB | 2GB+ |
| 磁盘 | 5GB | 10GB+ |
| Docker | 20.10+ | 最新版本 |
| Docker Compose | 2.0+ | 最新版本 |
| 系统 | Linux/macOS/Windows(WSL2) | Linux |

---

## 快速开始

### 方式一：使用预构建镜像（推荐）

最快的部署方式，无需本地构建：

```bash
# 1. 克隆项目
git clone https://github.com/Shannon-x/sufe-shop-bot.git
cd sufe-shop-bot

# 2. 准备配置
cp .env.production .env
nano .env  # 编辑配置

# 3. 拉取并启动
docker pull ghcr.io/shannon-x/sufe-shop-bot:latest
docker-compose up -d

# 4. 查看日志
docker-compose logs -f app
```

### 方式二：完整部署（含数据库）

包含 PostgreSQL 和 Redis 的完整部署：

```bash
# 1. 克隆项目
git clone https://github.com/Shannon-x/sufe-shop-bot.git
cd sufe-shop-bot

# 2. 准备配置
cp .env.production .env
nano .env  # 编辑配置

# 3. 启动完整服务栈
docker-compose -f docker-compose.full.yml up -d

# 4. 验证服务状态
docker-compose -f docker-compose.full.yml ps
```

---

## 部署方式

### 🚀 方式一：标准部署（docker-compose.yml）

使用预构建镜像，适用于：
- 快速部署和测试
- 使用外部数据库（如 1Panel、宝塔等管理的数据库）
- 轻量级生产环境

```bash
docker-compose up -d
```

### 🏗️ 方式二：完整部署（docker-compose.full.yml）

包含完整服务栈，适用于：
- 独立服务器部署
- 新环境快速搭建
- 需要完整隔离的生产环境

```bash
docker-compose -f docker-compose.full.yml up -d
```

### 📦 方式三：简单部署（docker-compose.simple.yml）

使用 SQLite 的最简配置，适用于：
- 本地开发测试
- 低流量场景
- 资源受限环境

```bash
docker-compose -f docker-compose.simple.yml up -d
```

### 🔧 方式四：本地构建部署

如需自定义或修改代码：

```bash
# 编辑 docker-compose.yml，注释 image 行，取消 build 注释
# 然后执行：
docker-compose up -d --build
```

---

## 配置说明

### 必需配置

在 `.env` 文件中配置以下必需项：

```env
# ===== Telegram Bot 配置 =====
BOT_TOKEN=your_bot_token_here          # 从 @BotFather 获取

# ===== 管理员配置 =====
ADMIN_TOKEN=your_secure_password       # 管理面板登录密码（建议16位以上）
ADMIN_TELEGRAM_IDS=123456789           # 管理员 Telegram ID（从 @userinfobot 获取）

# ===== 应用配置 =====
BASE_URL=https://your-domain.com       # 您的域名
PORT=9147                               # 应用端口
```

### 数据库配置

#### SQLite（默认，适合小型部署）
```env
DB_TYPE=sqlite
DB_NAME=shop.db
```

#### PostgreSQL（推荐生产环境）
```env
DB_TYPE=postgres
DB_HOST=localhost      # 使用 docker-compose.full.yml 时填 postgres
DB_PORT=5432
DB_NAME=shopbot
DB_USER=shopbot
DB_PASSWORD=your_secure_password
```

### Redis 配置（可选）
```env
REDIS_HOST=localhost   # 使用 docker-compose.full.yml 时填 redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
```

### 支付配置（可选）
```env
EPAY_PID=your_merchant_id
EPAY_KEY=your_secret_key
EPAY_GATEWAY=https://pay.gateway.com
EPAY_RETURN_URL=${BASE_URL}/payment/return
EPAY_NOTIFY_URL=${BASE_URL}/payment/notify
```

---

## 管理员设置

### 获取 Telegram ID

发送任意消息给 [@userinfobot](https://t.me/userinfobot)，机器人会返回您的 ID。

### 配置管理员

在 `.env` 中设置：
```env
# 单个管理员
ADMIN_TELEGRAM_IDS=123456789

# 多个管理员（逗号分隔）
ADMIN_TELEGRAM_IDS=123456789,987654321,555666777
```

### 访问管理面板

1. 启动服务后访问：`https://your-domain.com/admin/`
2. 使用 `ADMIN_TOKEN` 作为密码登录

---

## Nginx 反向代理

### 基础配置

```nginx
server {
    listen 80;
    server_name bot.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name bot.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://127.0.0.1:9147;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持（如需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

### 限制管理面板访问（可选）

```nginx
location /admin {
    # 限制 IP 访问
    allow 1.2.3.4;      # 您的 IP
    deny all;
    
    proxy_pass http://127.0.0.1:9147;
    # ... 其他 proxy 配置
}
```

---

## 更新与维护

### 更新到最新版本

```bash
# 使用预构建镜像
docker pull ghcr.io/shannon-x/sufe-shop-bot:latest
docker-compose down
docker-compose up -d

# 或者本地构建
git pull
docker-compose down
docker-compose up -d --build
```

### 备份数据

#### PostgreSQL 备份
```bash
# 导出
docker-compose exec postgres pg_dump -U shopbot shopbot > backup_$(date +%Y%m%d).sql

# 恢复
docker-compose exec -T postgres psql -U shopbot shopbot < backup_20240101.sql
```

#### SQLite 备份
```bash
# 导出
cp data/shop.db backup/shop_$(date +%Y%m%d).db

# 恢复
cp backup/shop_20240101.db data/shop.db
docker-compose restart app
```

### 查看日志

```bash
# 实时日志
docker-compose logs -f app

# 最后 100 行
docker-compose logs --tail=100 app

# 所有服务日志
docker-compose logs
```

---

## 故障排除

### 常见问题

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| 机器人不响应 | BOT_TOKEN 错误 | 检查 token 是否正确 |
| 数据库连接失败 | 配置错误或服务未启动 | 检查数据库配置和服务状态 |
| 管理面板 404 | 端口未开放 | 检查防火墙和端口映射 |
| 支付失败 | 回调 URL 不可达 | 确认 BASE_URL 配置正确 |

### 诊断命令

```bash
# 检查服务状态
docker-compose ps

# 检查容器健康状态
docker inspect shopbot-app --format='{{.State.Health.Status}}'

# 进入容器调试
docker-compose exec app sh

# 检查网络连接
docker-compose exec app wget -qO- http://localhost:9147/healthz
```

### 重置服务

```bash
# 重启单个服务
docker-compose restart app

# 完全重建
docker-compose down -v  # 警告：-v 会删除数据卷
docker-compose up -d --build
```

---

## 支持

如遇问题：
1. 查看 [故障排除](#故障排除) 部分
2. 搜索 [GitHub Issues](https://github.com/Shannon-x/sufe-shop-bot/issues)
3. 提交新 Issue 获取帮助

---

**祝您使用愉快！** 🚀