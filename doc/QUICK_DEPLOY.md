# 🚀 Quick Start - Git Deployment

## One-Time Setup

### 1. Push to Git

```bash
cd d:\Stack\werk-ticketing
git init
git remote add origin https://github.com/username/werk-ticketing.git
git add .
git commit -m "Initial commit"
git push -u origin main
```

### 2. Clone di Server

```bash
cd /var
sudo git clone https://github.com/username/werk-ticketing.git helpdeskgolang
sudo chown -R $USER:$USER /var/helpdeskgolang
```

### 3. Setup Environment

```bash
cd /var/helpdeskgolang/backend
nano .env.production  # Edit dengan config production

cd /var/helpdeskgolang/frontend
nano .env.production  # VITE_API_BASE_URL=/api
```

### 4. Make Deploy Script Executable

```bash
cd /var/helpdeskgolang
chmod +x deploy.sh
```

---

## Daily Workflow

### Update & Deploy

**Local:**

```bash
# Edit code...
git add .
git commit -m "Your changes"
git push
```

**Server:**

```bash
cd /var/helpdeskgolang
./deploy.sh
```

Done! ✅

---

## Common Commands

```bash
# Full deployment
cd /var/helpdeskgolang && ./deploy.sh

# Backend only
cd /var/helpdeskgolang/backend
git pull && cp .env.production .env && go build -o werk-ticketing-backend main.go && sudo systemctl restart werk-ticketing-backend

# Frontend only
cd /var/helpdeskgolang/frontend
git pull && cp .env.production .env && bun run build && sudo cp -r dist/* /var/www/werk-ticketing/

# Check status
sudo systemctl status werk-ticketing-backend
sudo journalctl -u werk-ticketing-backend -f

# Rollback
git log --oneline
git reset --hard <commit-hash>
./deploy.sh
```

---

**Full Guide:** See `GIT_DEPLOYMENT_GUIDE.md`
