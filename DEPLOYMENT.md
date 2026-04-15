# Deployment Guide

This guide covers deploying the CollectBFMSRQ application to production environments.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Configuration](#environment-configuration)
3. [Frontend Deployment](#frontend-deployment)
4. [Backend Deployment](#backend-deployment)
5. [Database Setup](#database-setup)
6. [Production Checklist](#production-checklist)
7. [Monitoring & Maintenance](#monitoring--maintenance)

---

## Prerequisites

Before deploying, ensure you have:

- ✅ A domain name (e.g., `yourapp.com`)
- ✅ SSL certificate (HTTPS)
- ✅ PostgreSQL database (managed or self-hosted)
- ✅ Deployment platform accounts (Vercel, Heroku, etc.)
- ✅ Environment variables configured

---

## Environment Configuration

### Backend `.env` for Production

Create a production `.env` file:

```env
# Server Configuration
SERVER_PORT=8080

# Database Configuration (Use production credentials)
DB_HOST=your-db-host.example.com
DB_PORT=5432
DB_USER=your_production_user
DB_PASSWORD=your_secure_password
DB_NAME=questionnaire_db

# CORS - Restrict to your domain
ALLOWED_ORIGINS=https://yourapp.com,https://www.yourapp.com
```

**Security Notes:**
- Never commit `.env` files to version control
- Use secret management (e.g., GitHub Secrets, AWS Secrets Manager)
- Rotate database passwords regularly

---

## Frontend Deployment

### Option 1: Vercel (Recommended)

**Step 1: Connect Repository**
```bash
# Push your code to GitHub
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/yourusername/collectbfmsrq.git
git push -u origin main
```

**Step 2: Deploy to Vercel**
```bash
# Install Vercel CLI
npm i -g vercel

# Login to Vercel
vercel login

# Deploy
cd frontend
vercel --prod
```

**Step 3: Configure Environment Variables**
In Vercel Dashboard:
- Go to Project Settings → Environment Variables
- Add: `NEXT_PUBLIC_API_URL=https://your-backend-url.com`

**Step 4: Update API URL in Frontend**

Create `frontend/.env.production`:
```env
NEXT_PUBLIC_API_URL=https://your-backend.herokuapp.com
```

Update API calls in components:
```typescript
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Use API_URL instead of hardcoded localhost
fetch(`${API_URL}/api/participants`, {...})
```

---

### Option 2: Netlify

**Step 1: Build the Application**
```bash
cd frontend
npm run build
```

**Step 2: Deploy to Netlify**
```bash
# Install Netlify CLI
npm i -g netlify-cli

# Login
netlify login

# Deploy
netlify deploy --prod --dir=.next
```

**Step 3: Configure Build Settings**
In Netlify Dashboard:
- Build command: `npm run build`
- Publish directory: `.next`
- Node version: `18`

---

### Option 3: Docker Deployment

**Create `frontend/Dockerfile`:**
```dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:18-alpine AS runner

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

EXPOSE 3000

CMD ["npm", "start"]
```

**Build and Run:**
```bash
cd frontend
docker build -t collectbfmsrq-frontend .
docker run -p 3000:3000 collectbfmsrq-frontend
```

---

## Backend Deployment

### Option 1: Heroku (Recommended)

**Step 1: Prepare for Heroku**

Create `backend/Procfile`:
```
web: ./server
```

Create `backend/heroku.yml`:
```yaml
build:
  docker:
    web: Dockerfile
```

**Step 2: Create Dockerfile**

Create `backend/Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
```

**Step 3: Deploy to Heroku**
```bash
# Install Heroku CLI
# https://devcenter.heroku.com/articles/heroku-cli

# Login
heroku login

# Create app
cd backend
heroku create your-app-name

# Add PostgreSQL
heroku addons:create heroku-postgresql:mini

# Set environment variables
heroku config:set SERVER_PORT=8080

# Deploy
git push heroku main
```

**Step 4: Verify Deployment**
```bash
heroku logs --tail
heroku open
```

---

### Option 2: Railway

**Step 1: Connect Repository**
- Go to https://railway.app/
- Click "New Project" → "Deploy from GitHub repo"

**Step 2: Configure**
- Root directory: `/backend`
- Build command: `go build -o server cmd/server/main.go`
- Start command: `./server`

**Step 3: Add PostgreSQL**
- In Railway dashboard: Add Database → PostgreSQL
- Railway automatically sets DATABASE_URL

**Step 4: Deploy**
- Push to main branch
- Railway auto-deploys

---

### Option 3: DigitalOcean App Platform

**Step 1: Create App**
- Go to DigitalOcean App Platform
- Create App from GitHub repository

**Step 2: Configure**
- Source: Your GitHub repo
- Root directory: `/backend`
- Build command: `go build -o server cmd/server/main.go`
- Run command: `./server`

**Step 3: Add Database**
- Add Managed Database (PostgreSQL)
- Connection string auto-injected as `DATABASE_URL`

**Step 4: Deploy**
- Set environment variables
- Deploy

---

### Option 4: Self-Hosted VPS

**Step 1: Prepare Server**
```bash
# SSH into server
ssh user@your-server-ip

# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install golang-go postgresql nginx -y
```

**Step 2: Setup PostgreSQL**
```bash
sudo -u postgres psql
CREATE DATABASE questionnaire_db;
CREATE USER youruser WITH PASSWORD 'yourpassword';
GRANT ALL PRIVILEGES ON DATABASE questionnaire_db TO youruser;
\q
```

**Step 3: Build and Run**
```bash
# Upload code to server (scp or git clone)
scp -r backend/ user@your-server-ip:/opt/collectbfmsrq/

# Build
cd /opt/collectbfmsrq/backend
go build -o server cmd/server/main.go

# Create systemd service
sudo nano /etc/systemd/system/collectbfmsrq.service
```

**Create systemd service file:**
```ini
[Unit]
Description=CollectBFMSRQ Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/collectbfmsrq/backend
EnvironmentFile=/opt/collectbfmsrq/backend/.env
ExecStart=/opt/collectbfmsrq/backend/server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**Step 4: Enable and Start Service**
```bash
sudo systemctl daemon-reload
sudo systemctl enable collectbfmsrq
sudo systemctl start collectbfmsrq
sudo systemctl status collectbfmsrq
```

**Step 5: Configure Nginx Reverse Proxy**
```nginx
server {
    listen 80;
    server_name api.yourapp.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
sudo nginx -t
sudo systemctl reload nginx
```

**Step 6: Add SSL with Let's Encrypt**
```bash
sudo apt install certbot python3-certbot-nginx -y
sudo certbot --nginx -d api.yourapp.com
```

---

## Database Setup

### Option 1: Managed PostgreSQL (Recommended for Production)

**Providers:**
- Heroku Postgres
- Supabase
- Neon
- AWS RDS
- DigitalOcean Managed Database

**Benefits:**
- Automated backups
- High availability
- Automatic updates
- Monitoring

### Option 2: Self-Hosted PostgreSQL

**Setup:**
```bash
# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Start service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
sudo -u postgres psql
CREATE DATABASE questionnaire_db;
CREATE USER app_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE questionnaire_db TO app_user;
\c questionnaire_db
GRANT ALL ON SCHEMA public TO app_user;
\q
```

**Backup Strategy:**
```bash
# Daily backup cron
0 2 * * * pg_dump questionnaire_db > /backups/questionnaire_db_$(date +\%Y\%m\%d).sql

# Restore
psql questionnaire_db < /backups/questionnaire_db_20240101.sql
```

---

## Production Checklist

### Before Deployment

- [ ] All tests passing (frontend + backend)
- [ ] Environment variables configured
- [ ] Database credentials secured
- [ ] CORS configured correctly
- [ ] SSL certificates installed
- [ ] Domain DNS configured
- [ ] Error monitoring setup (Sentry, etc.)
- [ ] Backup strategy implemented
- [ ] Rate limiting configured
- [ ] Input validation tested

### Security

- [ ] HTTPS enabled everywhere
- [ ] CORS restricted to production domain
- [ ] Database passwords rotated
- [ ] API rate limiting enabled
- [ ] SQL injection prevention verified
- [ ] XSS protection (CSP headers)
- [ ] Environment variables not exposed to frontend

### Performance

- [ ] Database indexes created
- [ ] Frontend assets optimized
- [ ] CDN configured for static assets
- [ ] Caching strategy implemented
- [ ] Load testing performed
- [ ] Database connection pooling configured

### Monitoring

- [ ] Application logging setup
- [ ] Error tracking (Sentry, Rollbar)
- [ ] Uptime monitoring (UptimeRobot, Pingdom)
- [ ] Database performance monitoring
- [ ] API response time tracking

---

## Monitoring & Maintenance

### Logging

**Backend Logging:**
```go
// Add structured logging in main.go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Info("Server starting",
    zap.String("port", cfg.ServerPort),
    zap.String("database", cfg.DBName),
)
```

**Frontend Logging:**
```typescript
// Use logging service instead of console.log
const logError = (error: Error) => {
  // Send to Sentry/LogRocket
  fetch('/api/log', {
    method: 'POST',
    body: JSON.stringify({ error: error.message }),
  });
};
```

### Health Checks

**Create `/health` endpoint:**
```go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "ok",
        "timestamp": time.Now(),
        "database": db.Connection.Ping() == nil,
    })
})
```

**Setup Uptime Monitoring:**
- Use UptimeRobot (free)
- Monitor `https://api.yourapp.com/health`
- Alert on downtime

### Database Maintenance

**Weekly Tasks:**
- Check database size
- Review slow query log
- Verify backups completed

**Monthly Tasks:**
- Run `VACUUM ANALYZE` on tables
- Review connection pool stats
- Check for unused indexes

```sql
-- Database maintenance queries
VACUUM ANALYZE;
REINDEX DATABASE questionnaire_db;

-- Check table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public';
```

---

## Rollback Strategy

### Frontend Rollback (Vercel)
```bash
# List deployments
vercel ls

# Rollback to previous
vercel rollback <deployment-url>
```

### Backend Rollback (Heroku)
```bash
# View releases
heroku releases

# Rollback
heroku rollback v123
```

### Database Rollback
```bash
# Restore from backup
psql questionnaire_db < /backups/latest.sql
```

---

## Scaling Considerations

### Horizontal Scaling

**Backend:**
- Deploy multiple instances behind load balancer
- Use stateless architecture (JWT tokens, external session store)
- Configure sticky sessions if needed

**Database:**
- Read replicas for query distribution
- Connection pooling (PgBouncer)
- Caching layer (Redis)

### Vertical Scaling

**When to Scale Up:**
- Database CPU > 70% consistently
- Memory usage > 80%
- Response times degrading

**Upgrade Path:**
1. Increase database instance size
2. Add read replicas
3. Implement caching (Redis)
4. Add CDN for frontend assets

---

## Cost Estimates

### Small Scale (< 1000 users/month)
- **Frontend**: Vercel Hobby (Free)
- **Backend**: Heroku Mini ($5/month)
- **Database**: Heroku Postgres Mini ($5/month)
- **Total**: ~$10/month

### Medium Scale (< 10,000 users/month)
- **Frontend**: Vercel Pro ($20/month)
- **Backend**: Heroku Basic ($25/month)
- **Database**: Heroku Postgres Basic ($50/month)
- **Total**: ~$95/month

### Large Scale (> 10,000 users/month)
- **Frontend**: Vercel Enterprise (Custom)
- **Backend**: AWS ECS/EKS ($100+/month)
- **Database**: AWS RDS ($200+/month)
- **Total**: ~$300+/month

---

## Support Contacts

- **Vercel Support**: https://vercel.com/support
- **Heroku Support**: https://help.heroku.com
- **PostgreSQL Support**: https://www.postgresql.org/support/
- **Go Community**: https://golang.org/doc/community

---

**Deployment is complete! Monitor your application and iterate based on user feedback.** 🚀
