# Telegram Shop Bot

[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org/doc/go1.22)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![CI](https://github.com/Shannon-x/sufe-shop-bot/actions/workflows/ci.yml/badge.svg)](https://github.com/Shannon-x/sufe-shop-bot/actions/workflows/ci.yml)
[![Docker](https://github.com/Shannon-x/sufe-shop-bot/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/Shannon-x/sufe-shop-bot/actions/workflows/docker-publish.yml)
[![Docker Image](https://img.shields.io/badge/Docker-ghcr.io%2Fshannon--x%2Fsufe--shop--bot-blue.svg)](https://ghcr.io/shannon-x/sufe-shop-bot)
[![Telegram Bot API](https://img.shields.io/badge/Telegram%20Bot%20API-v5.5.1-blue.svg)](https://core.telegram.org/bots/api)

[中文文档](README.md) | English

A feature-rich Telegram e-commerce bot system designed for automated sales of digital products (gift cards, membership cards, activation codes, etc.). Supports multiple languages, multiple payment methods, and a complete backend management system.

## 🌟 Key Features

### E-commerce Features
- 📦 **Product Management** - Support for various product types with flexible inventory management
- 💳 **Payment Integration** - Alipay, WeChat Pay support (via Epay gateway)
- 🎫 **Auto Delivery** - Automatic code delivery after successful payment
- 💰 **Balance System** - User recharge, balance payment, mixed payment
- 📱 **Recharge Cards** - Generate and manage recharge cards
- 📊 **Order Management** - Complete order lifecycle management

### Bot Features
- 🌐 **Multi-language Support** - Chinese, English (extensible)
- ⌨️ **Interactive UI** - Elegant keyboard layouts and inline buttons
- 🎯 **User Guidance** - Beginner-friendly workflow
- 📢 **Broadcast System** - Bulk notifications to users
- 🎫 **Ticket System** - Complete customer support
- ❓ **FAQ Management** - Automatic Q&A responses

### Admin Features
- 🖥️ **Web Admin Panel** - Full-featured management interface
- 📈 **Analytics** - Real-time sales data and user analysis
- 👥 **User Management** - User info, balance, order history
- 🔔 **Notification System** - Real-time alerts for new orders, low stock, etc.
- 🛡️ **Security Management** - JWT auth, access control, audit logs
- ⚙️ **System Configuration** - Online system settings modification

## 🚀 Quick Start

### Requirements
- Go 1.22+
- PostgreSQL 15+ or SQLite
- Redis 7+ (optional, for caching)
- Docker & Docker Compose (recommended)

### 1. Get the Code
```bash
git clone https://github.com/Shannon-x/sufe-shop-bot.git
cd sufe-shop-bot
```

### 2. Quick Deploy
```bash
# Run the interactive deployment script
chmod +x deploy.sh
./deploy.sh
```

The script will guide you to:
- Choose deployment method (full/simple/external database)
- Configure necessary environment variables
- Auto-initialize the database
- Start all services

### 3. Manual Deployment

#### Using Docker Compose (Recommended)

**Quick Deployment** (using pre-built image):
```bash
# Pull latest image
docker pull ghcr.io/shannon-x/sufe-shop-bot:latest

# Copy environment template
cp .env.production .env

# Edit configuration
vim .env

# Start services
docker-compose up -d
```

**Full Deployment** (with PostgreSQL and Redis):
```bash
# Copy environment template
cp .env.production .env

# Edit configuration
vim .env

# Start services (with database)
docker-compose -f docker-compose.full.yml up -d
```

**Local Build Deployment**:
```bash
# For local build, edit docker-compose.yml
# Comment out the image line, uncomment the build line
docker-compose up -d --build
```

#### Local Development
```bash
# Install dependencies
go mod download

# Set environment variables
export BOT_TOKEN=your_telegram_bot_token
export ADMIN_TOKEN=your_admin_token
export DB_TYPE=sqlite
export DB_NAME=./data/shop.db

# Run the service
go run cmd/server/main.go
```

## 📋 Configuration

### Required Settings
```env
# Telegram Bot Config
BOT_TOKEN=your_bot_token_here

# Admin Config
ADMIN_TOKEN=your_secure_admin_token
ADMIN_TELEGRAM_IDS=123456789,987654321  # Admin Telegram IDs (comma-separated)

# Application Config
BASE_URL=https://your-domain.com
PORT=9147
```

### Database Config
```env
# SQLite (default)
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

### Payment Config
```env
# Epay Configuration
EPAY_PID=your_merchant_id
EPAY_KEY=your_merchant_key
EPAY_GATEWAY=https://pay.example.com
EPAY_RETURN_URL=https://your-domain.com/payment/return
EPAY_NOTIFY_URL=https://your-domain.com/payment/notify
```

### Advanced Config
```env
# Redis Cache (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Webhook Mode (optional)
USE_WEBHOOK=false
WEBHOOK_URL=https://your-domain.com
WEBHOOK_PORT=8443

# Security Config
JWT_SECRET=your_jwt_secret_key
ENABLE_RATE_LIMIT=true
ENABLE_SECURITY_HEADERS=true
```

See [.env.production](.env.production) for complete configuration options.

## 🏗️ Project Structure

```
sufe-shop-bot/
├── cmd/
│   └── server/         # Application entry point
├── internal/
│   ├── app/           # Application initialization
│   ├── bot/           # Telegram Bot logic
│   ├── httpadmin/     # Web admin panel
│   ├── store/         # Data storage layer
│   ├── payment/       # Payment integration
│   ├── notification/  # Notification system
│   └── ...           # Other modules
├── templates/         # HTML templates
├── static/           # Static resources
├── scripts/          # Utility scripts
├── docker-compose*.yml # Docker configurations
└── deploy.sh         # Deployment script
```

## 📚 Usage Guide

### Bot Commands
- `/start` - Start using the shop
- `/help` - Get help
- `/language` - Switch language
- `/balance` - Check balance
- `/orders` - My orders
- `/ticket` - Create support ticket

### Admin Panel
Access `https://your-domain.com/admin` and login with the configured `ADMIN_TOKEN`.

Main features:
- **Product Management** - Add products, manage inventory, bulk upload codes
- **Order Management** - View order details, process refunds
- **User Management** - View user info, adjust balance
- **System Settings** - Configure payment, notifications, etc.
- **Analytics** - Sales reports, user analysis

## 🔧 Development Guide

### Local Development
```bash
# Install development tools
go install github.com/cosmtrek/air@latest

# Hot reload development
air
```

### Run Tests
```bash
go test ./...
```

### Build for Production
```bash
go build -ldflags "-s -w" -o shopbot cmd/server/main.go
```

## 🚀 Deployment Options

### 1. Docker Deployment (Recommended)
Three pre-configured deployment methods:
- `docker-compose.full.yml` - Full deployment (App + PostgreSQL + Redis)
- `docker-compose.simple.yml` - Simple deployment (App + SQLite)
- `docker-compose.yml` - Custom deployment (external database)

### 2. 1Panel Deployment
Supports one-click deployment in 1Panel using the original `docker-compose.yml`.

### 3. Manual Deployment
- Systemd service support
- Reverse proxy support (Nginx/Caddy)
- Load balancing and high availability

See [DEPLOY.md](DEPLOY.md) for detailed deployment documentation.

## 🔒 Security Features

- **JWT Authentication** - Secure API access control
- **Password Policy** - Configurable password complexity
- **Rate Limiting** - Prevent brute force and DDoS
- **Session Management** - Concurrent control and timeout settings
- **Data Encryption** - Encrypted storage for sensitive data
- **Audit Logging** - Complete operation records
- **CSRF Protection** - Prevent cross-site request forgery

## 📊 Monitoring and Logging

- **Prometheus Metrics** - Export system metrics
- **Structured Logging** - High-performance logging with zap
- **Error Tracking** - Detailed error context
- **Performance Monitoring** - Request latency, database query analysis

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)

## 📞 Contact

- GitHub Issues: [Submit Issue](https://github.com/Shannon-x/sufe-shop-bot/issues)

---

**Note**: Before using this project, please ensure compliance with Telegram's terms of service and local laws and regulations.