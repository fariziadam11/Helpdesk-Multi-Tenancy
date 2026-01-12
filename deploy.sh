#!/bin/bash

set -e

echo "🚀 Starting deployment..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
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
sleep 2
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
echo -e "${BLUE}📊 Checking backend service...${NC}"
if sudo systemctl is-active --quiet werk-ticketing-backend; then
    echo -e "${GREEN}✅ Backend service is running${NC}"
else
    echo -e "${YELLOW}⚠️  Backend service is not running!${NC}"
    sudo systemctl status werk-ticketing-backend --no-pager -l
fi

echo ""
echo -e "${GREEN}✅ Deployment complete!${NC}"
echo -e "${BLUE}🌐 Application is live at: https://bantuan.werk.co.id${NC}"
