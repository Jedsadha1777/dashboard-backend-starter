# Deployment Guide

## Development Environment

### Quick Start
```bash
docker-compose up -d
```

Access:
- API: http://localhost:3000
- PostgreSQL: localhost:5432
- Redis: localhost:6379

### Environment Variables
```env
# Required
JWT_SECRET=your-secret-key-minimum-32-characters
DB_PASSWORD=postgres

# Optional (with defaults)
SERVER_PORT=3000
DB_NAME=dashboard
ENVIRONMENT=development
AUTO_SEED=true
```

---

## Production Deployment

### Prerequisites
- VPS with Docker installed (Ubuntu 20.04+ recommended)
- Domain name pointed to server IP
- SSL certificate (or use Let's Encrypt)

### Step 1: Server Setup

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create app directory
sudo mkdir -p /opt/dashboard
cd /opt/dashboard
```

### Step 2: Deploy Application

```bash
# Clone repository
git clone <repository-url> .

# Create production .env
cat > .env << EOF
ENVIRONMENT=production
JWT_SECRET=$(openssl rand -base64 32)
DB_PASSWORD=$(openssl rand -base64 16)
DB_NAME=dashboard_prod
AUTO_SEED=false
ENABLE_SWAGGER=false
EOF

# Build and start
docker-compose build
docker-compose up -d

# Check status
docker-compose ps
```

### Step 3: Setup HTTPS with Nginx

Create `nginx/default.conf`:
```nginx
upstream api {
    server api:3000;
}

server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    # SSL certificates (use certbot for Let's Encrypt)
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    location / {
        proxy_pass http://api;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Block access to sensitive endpoints in production
    location /swagger {
        return 404;
    }
}
```

Add Nginx to docker-compose:
```yaml
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/default.conf:/etc/nginx/conf.d/default.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - api
    restart: unless-stopped
```

### Step 4: Database Backup

Create `scripts/backup.sh`:
```bash
#!/bin/bash
BACKUP_DIR="/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DB_NAME="dashboard_prod"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
docker-compose exec -T postgres pg_dump -U postgres $DB_NAME | gzip > $BACKUP_DIR/backup_$TIMESTAMP.sql.gz

# Keep only last 7 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed: backup_$TIMESTAMP.sql.gz"
```

Setup cron job:
```bash
# Add to crontab
crontab -e

# Daily backup at 2 AM
0 2 * * * /opt/dashboard/scripts/backup.sh
```

### Step 5: Monitoring

#### Basic Monitoring with Uptime Robot
1. Sign up at https://uptimerobot.com
2. Add monitor for: https://yourdomain.com/health
3. Set check interval: 5 minutes

#### Application Logs
```bash
# View logs
docker-compose logs -f api

# Save logs to file
docker-compose logs api > api.log

# Log rotation (add to docker-compose.yml)
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "5"
```

#### System Monitoring
```bash
# Check disk space
df -h

# Check memory
free -h

# Check Docker resources
docker stats
```

---

## Maintenance

### Update Application
```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker-compose build api
docker-compose up -d api

# Check deployment
docker-compose logs -f api
```

### Database Maintenance
```bash
# Access database
docker-compose exec postgres psql -U postgres -d dashboard

# Vacuum database (monthly)
docker-compose exec postgres vacuumdb -U postgres -d dashboard -z

# Restore from backup
gunzip < backup_20240101_020000.sql.gz | docker-compose exec -T postgres psql -U postgres -d dashboard
```

### SSL Certificate Renewal

Using Let's Encrypt:
```bash
# Install certbot
sudo apt install certbot

# Get certificate
sudo certbot certonly --standalone -d yourdomain.com

# Auto-renewal (add to crontab)
0 0 * * * certbot renew --quiet
```

---

## Troubleshooting

### Container won't start
```bash
# Check logs
docker-compose logs api

# Check config
docker-compose config

# Rebuild
docker-compose build --no-cache api
```

### Database connection issues
```bash
# Test connection
docker-compose exec postgres pg_isready

# Check credentials
docker-compose exec postgres psql -U postgres -l
```

### High memory usage
```bash
# Restart containers
docker-compose restart

# Clear unused resources
docker system prune -a
```

### Port already in use
```bash
# Find process using port
lsof -i :3000

# Kill process
kill -9 <PID>
```

---

## Security Checklist

- [ ] Change default passwords
- [ ] Set strong JWT_SECRET (32+ characters)
- [ ] Enable HTTPS
- [ ] Disable Swagger in production
- [ ] Setup firewall (allow only 80, 443, 22)
- [ ] Regular security updates
- [ ] Backup encryption
- [ ] Monitor access logs
- [ ] Rate limiting configured
- [ ] CORS properly configured