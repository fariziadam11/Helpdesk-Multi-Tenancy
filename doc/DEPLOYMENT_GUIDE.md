# 🚀 Deployment Guide - Werk Ticketing

Panduan lengkap untuk deploy **Backend Go** dan **Frontend Vue (Bun runtime)** ke server production tanpa Docker.

---

## 📋 Prerequisites

### Di Server Production

1. **Go 1.21+** - untuk menjalankan backend
2. **Bun** - untuk menjalankan frontend
3. **MySQL 8** - database
4. **Nginx** - reverse proxy (opsional tapi direkomendasikan)
5. **systemd** - untuk process management (Linux)

---

## 🔧 Persiapan Server

### 1. Install Go

```bash
# Download Go (sesuaikan versi)
wget https://go.dev/dl/go1.25.4.linux-amd64.tar.gz

# Extract
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz

# Add to PATH (tambahkan ke ~/.bashrc atau ~/.profile)
export PATH=$PATH:/usr/local/go/bin

# Verify
go version
```

### 2. Install Bun

```bash
# Install Bun
curl -fsSL https://bun.sh/install | bash

# Verify
bun --version
```

### 3. Setup MySQL

```bash
# Install MySQL (Ubuntu/Debian)
sudo apt update
sudo apt install mysql-server

# Secure installation
sudo mysql_secure_installation

# Login ke MySQL
sudo mysql -u root -p

# Buat database dan user
CREATE DATABASE armmada CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'armmada'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT ALL PRIVILEGES ON armmada.* TO 'armmada'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 4. Install Nginx (Opsional)

```bash
sudo apt install nginx
sudo systemctl enable nginx
sudo systemctl start nginx
```

---

## 📦 Deploy Backend (Go)

### 1. Upload Backend ke Server

```bash
# Di local machine, zip backend folder
cd d:\Stack\werk-ticketing
tar -czf backend.tar.gz backend/

# Upload ke server (gunakan scp, rsync, atau FTP)
scp backend.tar.gz user@your-server:/var/

# Di server, extract
cd /var
tar -xzf backend.tar.gz
mv backend helpdeskgolang/backend
# Atau jika sudah ada folder helpdeskgolang:
# tar -xzf backend.tar.gz -C helpdeskgolang/
cd helpdeskgolang/backend
```

### 2. Setup Environment Variables

```bash
# Copy dan edit .env
cp .env.example .env
nano .env
```

**Konfigurasi `.env` untuk production:**

```env
# Application Configuration
APP_ENV=production

# Server Configuration
SERVER_PORT=8080

# Logging Configuration
LOG_LEVEL=info
GIN_MODE=release
LOG_FORMAT=json

# Database Connection
DB_HOST=localhost
DB_PORT=3306
DB_USER=armmada
DB_PASSWORD=your_secure_password
DB_NAME=armmada

# JWT Configuration
JWT_SECRET=your_very_secure_random_jwt_secret_here

# InvGate Armmada API Configuration
ARMMADA_BASE_URL=https://support.armmada.id/api/v1/
ARMMADA_USERNAME=armmadaweb
ARMMADA_PASSWORD=j8f2yDzuVhYI4eG67hbsbck0
ARMMADA_PAGE_KEY=eyJsYXN0X2lkIjoxMDAwfQ==
ARMMADA_COMPANY_ID=135
ARMMADA_GROUP_ID=134
ARMMADA_LOCATION_ID=136
```

### 3. Build Backend

```bash
# Build binary
go build -o werk-ticketing-backend main.go

# Test run
./werk-ticketing-backend
```

### 4. Setup Systemd Service

Buat file service untuk auto-start backend:

```bash
sudo nano /etc/systemd/system/werk-ticketing-backend.service
```

**Isi file:**

```ini
[Unit]
Description=Werk Ticketing Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/helpdeskgolang/backend
ExecStart=/var/helpdeskgolang/backend/werk-ticketing-backend
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=werk-ticketing-backend

# Environment
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
```

**Enable dan start service:**

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service (auto-start on boot)
sudo systemctl enable werk-ticketing-backend

# Start service
sudo systemctl start werk-ticketing-backend

# Check status
sudo systemctl status werk-ticketing-backend

# View logs
sudo journalctl -u werk-ticketing-backend -f
```

---

## 🎨 Deploy Frontend (Vue + Bun)

### 1. Upload Frontend ke Server

```bash
# Di local machine
cd d:\Stack\werk-ticketing
tar -czf frontend.tar.gz frontend/

# Upload ke server
scp frontend.tar.gz user@your-server:/var/

# Di server
cd /var
tar -xzf frontend.tar.gz
mv frontend helpdeskgolang/frontend
# Atau jika sudah ada folder helpdeskgolang:
# tar -xzf frontend.tar.gz -C helpdeskgolang/
cd helpdeskgolang/frontend
```

### 2. Setup Environment Variables

```bash
# Buat file .env untuk production
nano .env
```

**Isi `.env`:**

```env
# API Base URL - sesuaikan dengan domain/IP server Anda
VITE_API_BASE_URL=http://your-server-ip:8080/api

# Atau jika menggunakan domain dengan Nginx:
# VITE_API_BASE_URL=https://api.yourdomain.com/api

# Disable dummy data
VITE_USE_DUMMY=false
```

### 3. Build Frontend

```bash
# Install dependencies
bun install

# Build untuk production
bun run build

# Hasil build ada di folder dist/
```

### 4. Opsi Deployment Frontend

#### **Opsi A: Serve dengan Bun (Preview Mode)**

```bash
# Preview production build
bun run preview --host 0.0.0.0 --port 3000
```

**Setup systemd service:**

```bash
sudo nano /etc/systemd/system/werk-ticketing-frontend.service
```

```ini
[Unit]
Description=Werk Ticketing Frontend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/helpdeskgolang/frontend
ExecStart=/root/.bun/bin/bun run preview --host 0.0.0.0 --port 3000
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=werk-ticketing-frontend

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable werk-ticketing-frontend
sudo systemctl start werk-ticketing-frontend
sudo systemctl status werk-ticketing-frontend
```

#### **Opsi B: Serve dengan Nginx (Recommended)**

```bash
# Copy build files ke nginx directory
sudo mkdir -p /var/www/werk-ticketing
sudo cp -r dist/* /var/www/werk-ticketing/
sudo chown -R www-data:www-data /var/www/werk-ticketing
```

**Konfigurasi Nginx:**

```bash
sudo nano /etc/nginx/sites-available/werk-ticketing
```

```nginx
server {
    listen 80;
    server_name your-domain.com;  # Ganti dengan domain Anda

    # Frontend
    root /var/www/werk-ticketing;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # Backend API Proxy
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/javascript application/json;
}
```

**Enable site:**

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/werk-ticketing /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

---

## 🔒 Setup SSL dengan Let's Encrypt (Opsional)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Obtain SSL certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal sudah di-setup otomatis
# Test renewal
sudo certbot renew --dry-run
```

### ⚠️ Troubleshooting Certbot Issues

#### Error: "Certificate Authority failed to verify"

Jika Anda mendapat error seperti:

```
Domain: bantuan.werk.co.id
Type:   unauthorized
Detail: Invalid response from http://bantuan.werk.co.id/.well-known/acme-challenge/...
```

**Penyebab**: Nginx tidak dapat diakses dari internet karena listening pada localhost saja.

**Solusi**:

1. **Periksa konfigurasi Nginx** - pastikan listening pada port 80, bukan `127.0.0.1:8000`:

```nginx
# ❌ SALAH - hanya localhost
listen 127.0.0.1:8000;

# ✅ BENAR - accessible dari internet
listen 80;
server_name bantuan.werk.co.id;
```

2. **Update konfigurasi**:

```bash
sudo nano /etc/nginx/sites-available/werk-ticketing
# Ubah listen 127.0.0.1:8000; menjadi listen 80;

# Test konfigurasi
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

3. **Verifikasi domain dapat diakses**:

```bash
# Test dari server
curl -I http://bantuan.werk.co.id

# Atau dari browser, buka:
# http://bantuan.werk.co.id
```

4. **Pastikan firewall membuka port 80 dan 443**:

```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw status
```

5. **Jalankan Certbot lagi**:

```bash
sudo certbot --nginx -d bantuan.werk.co.id
```

---

## 📊 Monitoring & Maintenance

### Check Service Status

```bash
# Backend
sudo systemctl status werk-ticketing-backend
sudo journalctl -u werk-ticketing-backend -f

# Frontend (jika pakai systemd)
sudo systemctl status werk-ticketing-frontend
sudo journalctl -u werk-ticketing-frontend -f

# Nginx
sudo systemctl status nginx
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Restart Services

```bash
# Backend
sudo systemctl restart werk-ticketing-backend

# Frontend
sudo systemctl restart werk-ticketing-frontend

# Nginx
sudo systemctl reload nginx
```

### Update Application

**Backend:**

```bash
cd /var/helpdeskgolang/backend
git pull  # atau upload file baru
go build -o werk-ticketing-backend main.go
sudo systemctl restart werk-ticketing-backend
```

**Frontend:**

```bash
cd /var/helpdeskgolang/frontend
git pull  # atau upload file baru
bun install
bun run build
sudo cp -r dist/* /var/www/werk-ticketing/
sudo systemctl reload nginx
```

---

## 🔥 Firewall Configuration

```bash
# Allow HTTP, HTTPS, SSH
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Jika backend langsung exposed (tidak direkomendasikan)
# sudo ufw allow 8080/tcp

# Enable firewall
sudo ufw enable
sudo ufw status
```

---

## ✅ Checklist Deployment

### Backend

- [ ] Go terinstall di server
- [ ] MySQL database dibuat
- [ ] File backend di-upload
- [ ] `.env` dikonfigurasi dengan benar
- [ ] Binary di-build (`go build`)
- [ ] Systemd service dibuat dan running
- [ ] Backend accessible di `http://localhost:8080`

### Frontend

- [ ] Bun terinstall di server
- [ ] File frontend di-upload
- [ ] `.env` dikonfigurasi (VITE_API_BASE_URL)
- [ ] Dependencies terinstall (`bun install`)
- [ ] Production build dibuat (`bun run build`)
- [ ] Nginx dikonfigurasi atau Bun preview running
- [ ] Frontend accessible di browser

### Security

- [ ] Firewall dikonfigurasi
- [ ] SSL certificate terinstall (jika ada domain)
- [ ] JWT_SECRET diganti dengan random string
- [ ] Database password aman
- [ ] `.env` file tidak accessible dari web

### Testing

- [ ] Buka frontend di browser
- [ ] Test register user
- [ ] Test login
- [ ] Test create ticket
- [ ] Test list tickets
- [ ] Test ticket detail
- [ ] Check logs untuk errors

---

## 🆘 Troubleshooting

### Backend tidak start

```bash
# Check logs
sudo journalctl -u werk-ticketing-backend -n 50

# Check database connection
mysql -u armmada -p armmada

# Check port
sudo netstat -tulpn | grep 8080
```

### Frontend tidak accessible

```bash
# Check nginx config
sudo nginx -t

# Check nginx logs
sudo tail -f /var/log/nginx/error.log

# Check file permissions
ls -la /var/www/werk-ticketing
```

### CORS errors

Pastikan backend CORS middleware sudah dikonfigurasi untuk allow origin dari frontend domain.

### Database connection failed

- Check MySQL service: `sudo systemctl status mysql`
- Check credentials di `.env`
- Check MySQL user permissions

---

## 📝 Notes

1. **Backup Database**: Setup automated MySQL backup
2. **Log Rotation**: Configure logrotate untuk backend logs
3. **Monitoring**: Consider setup monitoring (Prometheus, Grafana)
4. **Auto-deployment**: Setup CI/CD dengan GitHub Actions (opsional)

---

**Selamat! Aplikasi Werk Ticketing sudah deployed! 🎉**
