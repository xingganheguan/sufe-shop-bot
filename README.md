
1.后台增加了2FA验证
---
# Telegram Shop Bot - 电报商城机器人

[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org/doc/go1.22)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![CI](https://github.com/Shannon-x/sufe-shop-bot/actions/workflows/ci.yml/badge.svg)](https://github.com/Shannon-x/sufe-shop-bot/actions/workflows/ci.yml)
[![Docker](https://github.com/Shannon-x/sufe-shop-bot/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/Shannon-x/sufe-shop-bot/actions/workflows/docker-publish.yml)
[![Docker Image](https://img.shields.io/badge/Docker-ghcr.io%2Fshannon--x%2Fsufe--shop--bot-blue.svg)](https://ghcr.io/shannon-x/sufe-shop-bot)
[![Telegram Bot API](https://img.shields.io/badge/Telegram%20Bot%20API-v5.5.1-blue.svg)](https://core.telegram.org/bots/api)

> [<img src="https://img.sufe.me/file/1760201136357_isufe_icon.webp" alt="苏菲家宽" height="20" />](https://sufe.pro)
> **苏菲家宽** — 提供极度纯净的家宽代理： [sufe.pro](https://sufe.pro)

一个功能完善的 Telegram 电商机器人系统，专门用于数字商品（充值卡、会员卡、激活码等）的自动化销售。支持多语言、多支付方式、完整的后台管理系统。


<!-- Sponsor / Ad -->
<div align="center">
  <a href="https://sufe.pro">
    <img src="https://img.sufe.me/file/1760201136357_isufe_icon.webp" alt="苏菲家宽" width="140">
  </a>

  <h3>苏菲家宽 · 极度纯净的家宽代理</h3>
  <p>
    为 AI使用、跨境业务等场景提供
    <b>稳定、干净、低延迟</b> 的家宽出口。<br/>
    多地区线路可选  · 可选静态/动态 IP · 覆盖 IPv4/IPv6
  </p>

  <p><b>👉 立即访问：</b> <a href="https://sufe.pro">sufe.pro</a></p>
</div>



## 🌟 核心功能

### 商城功能
- 📦 **商品管理** - 支持多种商品类型，灵活的库存管理
- 💳 **支付集成** - 支持支付宝、微信支付（通过易支付）
- 🎫 **自动发货** - 支付成功后自动发送卡密
- 💰 **余额系统** - 用户充值、余额支付、混合支付
- 📱 **充值卡** - 生成和管理充值卡
- 📊 **订单管理** - 完整的订单生命周期管理

### Bot 功能
- 🌐 **多语言支持** - 中文、英文（可扩展）
- ⌨️ **交互式界面** - 优雅的键盘布局和内联按钮
- 🎯 **用户引导** - 新手友好的操作流程
- 📢 **广播系统** - 向用户批量发送通知
- 🎫 **工单系统** - 完整的客服支持
- ❓ **FAQ 管理** - 常见问题自动回答

### 管理功能
- 🖥️ **Web 管理后台** - 功能完善的管理界面
- 📈 **数据统计** - 实时销售数据和用户分析
- 👥 **用户管理** - 用户信息、余额、订单历史
- 🔔 **通知系统** - 新订单、缺货等实时通知
- 🛡️ **安全管理** - JWT 认证、访问控制、操作审计
- ⚙️ **系统配置** - 在线修改系统设置

## 🚀 快速开始

### 环境要求
- Go 1.22+
- PostgreSQL 15+ 或 SQLite
- Redis 7+ （可选，用于缓存）
- Docker & Docker Compose （推荐）

### 1. 获取代码
```bash
git clone https://github.com/Shannon-x/sufe-shop-bot.git
cd sufe-shop-bot
```

### 2. 快速部署
```bash
# 运行交互式部署脚本
chmod +x deploy.sh
./deploy.sh
```

部署脚本会引导你：
- 选择部署方式（完整部署/简单部署/外部数据库）
- 配置必要的环境变量
- 自动初始化数据库
- 启动所有服务

### 3. 手动部署

#### 使用 Docker Compose（推荐）

**快速部署**（使用预构建镜像）：
```bash
# 拉取最新镜像
docker pull ghcr.io/shannon-x/sufe-shop-bot:latest

# 复制环境变量模板
cp .env.production .env

# 编辑配置
vim .env

# 启动服务
docker-compose up -d
```

**完整部署**（包含 PostgreSQL 和 Redis）：
```bash
# 复制环境变量模板
cp .env.production .env

# 编辑配置
vim .env

# 启动服务（包含数据库）
docker-compose -f docker-compose.full.yml up -d
```

**本地构建部署**：
```bash
# 如需本地构建，编辑 docker-compose.yml
# 注释 image 行，取消注释 build 行
docker-compose up -d --build
```

#### 本地开发
```bash
# 安装依赖
go mod download

# 设置环境变量
export BOT_TOKEN=your_telegram_bot_token
export ADMIN_TOKEN=your_admin_token
export DB_TYPE=sqlite
export DB_NAME=./data/shop.db

# 运行服务
go run cmd/server/main.go
```

## 📋 配置说明

### 必需配置
```env
# Telegram Bot 配置
BOT_TOKEN=your_bot_token_here

# 管理员配置
ADMIN_TOKEN=your_secure_admin_token
ADMIN_TELEGRAM_IDS=123456789,987654321  # 管理员 Telegram ID（逗号分隔）

# 应用配置
BASE_URL=https://your-domain.com
PORT=9147
```

### 数据库配置
```env
# SQLite（默认）
DB_TYPE=sqlite
DB_NAME=./data/shop.db

# PostgreSQL
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=shopbot
DB_USER=shopbot
DB_PASSWORD=your_password
```

### 支付配置
```env
# 易支付配置
EPAY_PID=your_merchant_id
EPAY_KEY=your_merchant_key
EPAY_GATEWAY=https://pay.example.com
EPAY_RETURN_URL=https://your-domain.com/payment/return
EPAY_NOTIFY_URL=https://your-domain.com/payment/notify
```

### 高级配置
```env
# Redis 缓存（可选）
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Webhook 模式（可选）
USE_WEBHOOK=false
WEBHOOK_URL=https://your-domain.com
WEBHOOK_PORT=8443

# 安全配置
JWT_SECRET=your_jwt_secret_key
ENABLE_RATE_LIMIT=true
ENABLE_SECURITY_HEADERS=true
```

完整配置说明请查看 [.env.production](.env.production)

## 🏗️ 项目结构

```
sufe-shop-bot/
├── cmd/
│   └── server/         # 应用入口
├── internal/
│   ├── app/           # 应用初始化
│   ├── bot/           # Telegram Bot 逻辑
│   ├── httpadmin/     # Web 管理后台
│   ├── store/         # 数据存储层
│   ├── payment/       # 支付集成
│   ├── notification/  # 通知系统
│   └── ...           # 其他模块
├── templates/         # HTML 模板
├── static/           # 静态资源
├── scripts/          # 工具脚本
├── docker-compose*.yml # Docker 配置
└── deploy.sh         # 部署脚本
```

## 📚 使用指南

### Bot 命令
- `/start` - 开始使用商城
- `/help` - 获取帮助
- `/language` - 切换语言
- `/balance` - 查看余额
- `/orders` - 我的订单
- `/ticket` - 创建工单

### 管理后台
访问 `https://your-domain.com/admin` 使用配置的 `ADMIN_TOKEN` 登录。

主要功能：
- **商品管理** - 添加商品、管理库存、批量上传卡密
- **订单管理** - 查看订单详情、处理退款
- **用户管理** - 查看用户信息、调整余额
- **系统设置** - 配置支付、通知等参数
- **数据统计** - 销售报表、用户分析

## 🔧 开发指南

### 本地开发
```bash
# 安装开发工具
go install github.com/cosmtrek/air@latest

# 热重载开发
air
```

### 运行测试
```bash
go test ./...
```

### 构建生产版本
```bash
go build -ldflags "-s -w" -o shopbot cmd/server/main.go
```

## 🚀 部署选项

### 1. Docker 部署（推荐）
提供三种预配置的部署方式：
- `docker-compose.full.yml` - 完整部署（应用+PostgreSQL+Redis）
- `docker-compose.simple.yml` - 简单部署（应用+SQLite）
- `docker-compose.yml` - 自定义部署（使用外部数据库）

### 2. 1Panel 部署
支持在 1Panel 面板中一键部署，使用原始的 `docker-compose.yml`。

### 3. 手动部署
- 支持 Systemd 服务
- 支持反向代理（Nginx/Caddy）
- 支持负载均衡和高可用

详细部署文档请查看 [DEPLOY.md](DEPLOY.md)

## 🔒 安全特性

- **JWT 认证** - 安全的 API 访问控制
- **密码策略** - 可配置的密码复杂度
- **速率限制** - 防止暴力破解和 DDoS
- **会话管理** - 并发控制和超时设置
- **数据加密** - 敏感数据加密存储
- **审计日志** - 完整的操作记录
- **CSRF 保护** - 防止跨站请求伪造

## 📊 监控和日志

- **Prometheus 指标** - 导出系统指标
- **结构化日志** - 使用 zap 高性能日志
- **错误追踪** - 详细的错误上下文
- **性能监控** - 请求延迟、数据库查询分析

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)

## 📞 联系方式

- GitHub Issues: [提交问题](https://github.com/Shannon-x/sufe-shop-bot/issues)

---

**注意**：使用本项目前，请确保遵守 Telegram 的服务条款和当地法律法规。

<table>
<tr>
<td width="160" valign="middle" align="center">
  <a href="https://sufe.pro">
    <img src="https://img.sufe.me/file/1760201136357_isufe_icon.webp" alt="苏菲家宽" width="140">
  </a>
</td>
<td valign="middle">

**苏菲家宽** — 极度纯净的家宽代理  
提供 **稳定、干净、低延迟** 的家宽网络出口。  
• 多地区线路可选　• 独享/共享灵活选择　• 静态/动态 IP 可选　• IPv4/IPv6 支持  

**👉 立刻访问： [sufe.pro](https://sufe.pro)**
</td>
</tr>
</table>
