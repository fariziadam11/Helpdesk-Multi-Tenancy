# 📂 Quick Reference - Deployment Paths

## Server Directory Structure

```
/var/helpdeskgolang/
├── backend/
│   ├── main.go
│   ├── werk-ticketing-backend (binary)
│   ├── .env
│   ├── go.mod
│   ├── go.sum
│   ├── internal/
│   └── migrations/
│
└── frontend/
    ├── dist/ (production build)
    ├── src/
    ├── package.json
    ├── .env
    └── bun.lockb
```

## Important Paths

### Backend

- **Project Directory:** `/var/helpdeskgolang/backend`
- **Binary:** `/var/helpdeskgolang/backend/werk-ticketing-backend`
- **Environment:** `/var/helpdeskgolang/backend/.env`
- **Systemd Service:** `/etc/systemd/system/werk-ticketing-backend.service`

### Frontend

- **Project Directory:** `/var/helpdeskgolang/frontend`
- **Build Output:** `/var/helpdeskgolang/frontend/dist`
- **Nginx Serve From:** `/var/www/werk-ticketing` (copy dari dist/)
- **Environment:** `/var/helpdeskgolang/frontend/.env`

### Nginx

- **Config:** `/etc/nginx/sites-available/werk-ticketing`
- **Symlink:** `/etc/nginx/sites-enabled/werk-ticketing`
- **Static Files:** `/var/www/werk-ticketing`

## Quick Commands

### Setup Directory

```bash
# Create main directory
sudo mkdir -p /var/helpdeskgolang
sudo chown $USER:$USER /var/helpdeskgolang

# Create subdirectories
mkdir -p /var/helpdeskgolang/backend
mkdir -p /var/helpdeskgolang/frontend
```

### Upload Files

```bash
# From local machine
scp -r backend/ user@server:/var/helpdeskgolang/
scp -r frontend/ user@server:/var/helpdeskgolang/
```

### Navigate

```bash
# Backend
cd /var/helpdeskgolang/backend

# Frontend
cd /var/helpdeskgolang/frontend
```

### Logs Location

```bash
# Backend logs (systemd)
sudo journalctl -u werk-ticketing-backend -f

# Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

## Permissions

```bash
# Set ownership
sudo chown -R www-data:www-data /var/helpdeskgolang
sudo chown -R www-data:www-data /var/www/werk-ticketing

# Set permissions
chmod 755 /var/helpdeskgolang/backend/werk-ticketing-backend
chmod 600 /var/helpdeskgolang/backend/.env
chmod 600 /var/helpdeskgolang/frontend/.env
```
