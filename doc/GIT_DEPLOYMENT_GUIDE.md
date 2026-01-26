# 🚀 Git-Based Deployment Guide

Panduan deployment menggunakan Git untuk update aplikasi dengan mudah tanpa perlu zip/upload manual.

---

## 📋 Prerequisites

1. **Git terinstall di server**
2. **Repository Git** (GitHub, GitLab, atau Bitbucket)
3. **SSH access ke server**

---

## 🔧 Setup Awal (One-Time Setup)

### 1. Push Project ke Git Repository

**Di local machine:**

```bash
cd d:\Stack\werk-ticketing

# Initialize git (jika belum)
git init

# Add remote repository
git remote add origin https://github.com/username/werk-ticketing.git
# atau SSH: git remote add origin git@github.com:username/werk-ticketing.git

# Add files
git add .
git commit -m "Initial commit"

# Push to repository
git push -u origin main
```

### 2. Clone Repository di Server

**Di server production:**

```bash
# Clone repository ke /var/helpdeskgolang
cd /var
sudo git clone https://github.com/username/werk-ticketing.git helpdeskgolang

# Set ownership
sudo chown -R $USER:$USER /var/helpdeskgolang
```

### 3. Setup Environment Files di Server

```bash
cd /var/helpdeskgolang

# Backend - create .env.production
cd backend
nano .env.production
```

**Isi `.env.production` untuk backend:**

```env
APP_ENV=production
SERVER_PORT=8080
LOG_LEVEL=info
GIN_MODE=release
LOG_FORMAT=json

DB_HOST=localhost
DB_PORT=3306
DB_USER=armmada
DB_PASSWORD=your_secure_password
DB_NAME=armmada

JWT_SECRET=your_very_secure_random_jwt_secret

ARMMADA_BASE_URL=https://support.armmada.id/api/v1/
ARMMADA_USERNAME=armmadaweb
ARMMADA_PASSWORD=j8f2yDzuVhYI4eG67hbsbck0
ARMMADA_PAGE_KEY=eyJsYXN0X2lkIjoxMDAwfQ==
ARMMADA_COMPANY_ID=135
ARMMADA_GROUP_ID=134
ARMMADA_LOCATION_ID=136
```

```bash
# Frontend - create .env.production
cd ../frontend
nano .env.production
```

**Isi `.env.production` untuk frontend:**

```env
VITE_API_BASE_URL=/api
VITE_USE_DUMMY=false
```

### 4. Update .gitignore

Pastikan `.gitignore` tidak commit file environment:

```bash
cd /var/helpdeskgolang
nano .gitignore
```

**Tambahkan:**

```
# Environment files
.env
.env.local
.env.production
.env.development

# Build outputs
backend/werk-ticketing-backend
frontend/dist/
frontend/node_modules/
backend/tmp/
```

---

## 🔄 Deployment Workflow

### Cara Deploy Update (Setiap Ada Perubahan)

**1. Di Local Machine - Push Changes:**

```bash
cd d:\Stack\werk-ticketing

# Make your changes...
# Edit files, fix bugs, add features, etc.

# Commit and push
git add .
git commit -m "Your commit message"
git push origin main
```

**2. Di Server - Pull & Deploy:**

```bash
# SSH ke server
ssh user@your-server

# Navigate to project
cd /var/helpdeskgolang

# Pull latest changes
git pull origin main

# Deploy backend
cd backend
cp .env.production .env
go build -o werk-ticketing-backend main.go
sudo systemctl restart werk-ticketing-backend

# Deploy frontend
cd ../frontend
cp .env.production .env
bun install
bun run build
sudo cp -r dist/* /var/www/werk-ticketing/

# Reload nginx
sudo systemctl reload nginx

echo "✅ Deployment complete!"
```

---

## 🤖 Automated Deploy Script

Buat script untuk otomasi deployment:

**`/var/helpdeskgolang/deploy.sh`:**

```bash
#!/bin/bash

set -e

echo "🚀 Starting deployment..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Pull latest code
echo -e "${BLUE}📥 Pulling latest code...${NC}"
git pull origin main

# Deploy Backend
echo -e "${BLUE}🔨 Building backend...${NC}"
cd backend
cp .env.production .env
go build -ldflags="-s -w" -o werk-ticketing-backend main.go
echo -e "${GREEN}✅ Backend built${NC}"

echo -e "${BLUE}🔄 Restarting backend service...${NC}"
sudo systemctl restart werk-ticketing-backend
echo -e "${GREEN}✅ Backend restarted${NC}"

# Deploy Frontend
echo -e "${BLUE}🔨 Building frontend...${NC}"
cd ../frontend
cp .env.production .env
bun install --production
bun run build
echo -e "${GREEN}✅ Frontend built${NC}"

echo -e "${BLUE}📦 Deploying frontend...${NC}"
sudo cp -r dist/* /var/www/werk-ticketing/
echo -e "${GREEN}✅ Frontend deployed${NC}"

# Reload Nginx
echo -e "${BLUE}🔄 Reloading Nginx...${NC}"
sudo systemctl reload nginx
echo -e "${GREEN}✅ Nginx reloaded${NC}"

# Check status
echo -e "${BLUE}📊 Checking services...${NC}"
sudo systemctl status werk-ticketing-backend --no-pager -l
echo ""

echo -e "${GREEN}✅ Deployment complete!${NC}"
echo -e "${BLUE}🌐 Application is live at: https://bantuan.werk.co.id${NC}"
```

**Make it executable:**

```bash
chmod +x /var/helpdeskgolang/deploy.sh
```

**Usage:**

```bash
cd /var/helpdeskgolang
./deploy.sh
```

---

## 🎯 Quick Deploy Commands

### Deploy Everything (Backend + Frontend)

```bash
cd /var/helpdeskgolang && ./deploy.sh
```

### Deploy Backend Only

```bash
cd /var/helpdeskgolang/backend
git pull origin main
cp .env.production .env
go build -o werk-ticketing-backend main.go
sudo systemctl restart werk-ticketing-backend
```

### Deploy Frontend Only

```bash
cd /var/helpdeskgolang/frontend
git pull origin main
cp .env.production .env
bun install
bun run build
sudo cp -r dist/* /var/www/werk-ticketing/
sudo systemctl reload nginx
```

---

## 🔐 Setup SSH Key (Recommended)

Untuk avoid password prompt saat git pull:

**Di server:**

```bash
# Generate SSH key
ssh-keygen -t ed25519 -C "server@werk.co.id"

# Copy public key
cat ~/.ssh/id_ed25519.pub
```

**Di GitHub/GitLab:**

1. Go to Settings → SSH Keys
2. Add new SSH key
3. Paste public key

**Update remote URL:**

```bash
cd /var/helpdeskgolang
git remote set-url origin git@github.com:username/werk-ticketing.git
```

---

## 📝 Best Practices

### 1. Branching Strategy

```bash
# Development branch
git checkout -b development

# Feature branches
git checkout -b feature/new-feature

# Merge to main for production
git checkout main
git merge development
git push origin main
```

### 2. Tag Releases

```bash
# Tag version
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# Deploy specific version
git checkout v1.0.0
./deploy.sh
```

### 3. Rollback

```bash
# View commits
git log --oneline

# Rollback to previous commit
git reset --hard <commit-hash>
./deploy.sh

# Or checkout specific tag
git checkout v1.0.0
./deploy.sh
```

---

## 🚨 Troubleshooting

### Git Pull Conflicts

```bash
# Stash local changes
git stash

# Pull latest
git pull origin main

# Apply stashed changes
git stash pop
```

### Permission Issues

```bash
# Fix ownership
sudo chown -R $USER:$USER /var/helpdeskgolang

# Fix permissions
chmod +x /var/helpdeskgolang/deploy.sh
```

### Service Not Starting

```bash
# Check logs
sudo journalctl -u werk-ticketing-backend -n 50

# Check status
sudo systemctl status werk-ticketing-backend
```

---

## 📊 Comparison: Manual vs Git

| Aspect                 | Manual (Zip/Upload) | Git-Based   |
| ---------------------- | ------------------- | ----------- |
| **Update Time**        | 5-10 minutes        | 1-2 minutes |
| **Steps**              | 6 steps             | 1 command   |
| **Error Prone**        | High                | Low         |
| **Rollback**           | Difficult           | Easy        |
| **Version Control**    | Manual              | Automatic   |
| **Team Collaboration** | Difficult           | Easy        |

---

## ✅ Migration Checklist

- [ ] Push project to Git repository
- [ ] Clone repository to server
- [ ] Create `.env.production` files
- [ ] Update `.gitignore`
- [ ] Create `deploy.sh` script
- [ ] Test deployment with `./deploy.sh`
- [ ] Setup SSH key (optional)
- [ ] Document workflow for team

---

**Sekarang deployment hanya perlu 2 langkah:**

1. **Local:** `git push`
2. **Server:** `./deploy.sh`

Jauh lebih cepat dan mudah! 🚀
