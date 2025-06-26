# 🚀 Human Intelligence - Deployment Guide

> **Architecture**: SvelteKit Frontend + Go Backend + PostgreSQL Database  
> **Recommended Stack**: Vercel + Railway + Supabase  
> **Estimated Cost**: $5-15/month depending on usage

## 📋 **Deployment Architecture Overview**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   (SvelteKit)   │────│   (Go/Gin)     │────│   (PostgreSQL)  │
│   Vercel        │    │   Railway       │    │   Supabase      │
│   FREE          │    │   $5/month      │    │   FREE          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🎯 **Recommended Deployment Strategy**

### **Why This Stack?**
- ✅ **Vercel**: Best SvelteKit deployment experience, free tier
- ✅ **Railway**: Excellent Go backend support, automatic deployments
- ✅ **Supabase**: Managed PostgreSQL with admin dashboard, free 500MB
- ✅ **Total Cost**: ~$5/month (Railway only paid service)
- ✅ **Scalability**: Each service can scale independently

---

## 🗄️ **Database Setup (Supabase)**

### **Step 1: Create Supabase Project**
1. Go to [supabase.com](https://supabase.com)
2. Create account and new project
3. Choose region closest to your users
4. Save the connection details

### **Step 2: Database Configuration**
```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Run your existing migration files
-- Copy SQL from backend/migrations/*.sql
```

### **Step 3: Get Connection String**
```bash
# Found in: Settings > Database > Connection String
postgresql://postgres:[PASSWORD]@[HOST]:5432/postgres
```

---

## 🛠️ **Backend Deployment (Railway)**

### **Step 1: Prepare Repository**
```bash
# Add Dockerfile to backend directory
cd backend
touch Dockerfile
```

**Dockerfile:**
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8081
CMD ["./main"]
```

### **Step 2: Deploy to Railway**
1. Visit [railway.app](https://railway.app)
2. Connect GitHub repository
3. Select `backend` directory as root
4. Railway auto-detects Dockerfile

### **Step 3: Environment Variables**
Set in Railway dashboard:
```bash
# Database
DATABASE_URL=postgresql://postgres:[PASSWORD]@[SUPABASE_HOST]:5432/postgres

# JWT
JWT_SECRET=your-super-secure-jwt-secret-min-32-chars

# Server
PORT=8081
GIN_MODE=release

# CORS (add your frontend domain)
FRONTEND_URL=https://your-app.vercel.app

# Optional: Redis for sessions
REDIS_URL=redis://...
```

### **Step 4: Custom Domain (Optional)**
```bash
# In Railway dashboard:
# Settings > Domains > Add Custom Domain
api.yourdomain.com
```

---

## 🎨 **Frontend Deployment (Vercel)**

### **Step 1: Update Configuration**
**frontend/vite.config.js:**
```javascript
export default defineConfig({
  plugins: [sveltekit()],
  
  // Production API URL
  define: {
    __API_URL__: JSON.stringify(
      process.env.NODE_ENV === 'production' 
        ? 'https://your-backend.railway.app' 
        : 'http://localhost:8081'
    ),
  },

  // ... rest of config
});
```

**frontend/src/lib/stores/auth.js:**
```javascript
// Update API base URL
const API_BASE_URL = process.env.NODE_ENV === 'production'
  ? 'https://your-backend.railway.app/api/v1'
  : 'http://localhost:8081/api/v1';
```

### **Step 2: Deploy to Vercel**
```bash
# Install Vercel CLI
npm i -g vercel

# From frontend directory
cd frontend
vercel

# Follow prompts:
# Framework: SvelteKit
# Build command: npm run build
# Output directory: build
```

### **Step 3: Environment Variables**
Set in Vercel dashboard:
```bash
# API Configuration
VITE_API_URL=https://your-backend.railway.app
VITE_WS_URL=wss://your-backend.railway.app

# Optional: Analytics
VITE_ANALYTICS_ID=your-analytics-id
```

### **Step 4: Custom Domain**
```bash
# In Vercel dashboard:
# Settings > Domains
your-domain.com
www.your-domain.com
```

---

## 🔧 **Alternative Deployment Options**

### **Option A: All-in-One Platform (Render)**
```bash
# Frontend: Render Static Site
# Backend: Render Web Service  
# Database: Render PostgreSQL
# Cost: ~$15/month
# Pros: Single platform, simpler management
```

### **Option B: Traditional VPS (DigitalOcean)**
```bash
# Server: $6/month droplet
# Setup: Docker Compose
# Database: Same server or managed
# Pros: Full control, cost-effective at scale
```

### **Option C: Serverless Backend (Vercel Functions)**
```bash
# Rewrite Go backend as Vercel API routes
# Pros: Both frontend and backend on Vercel
# Cons: Significant code rewrite required
```

---

## 🐳 **Docker Compose (VPS Deployment)**

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  database:
    image: postgres:15
    environment:
      POSTGRES_DB: human_intelligence
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    environment:
      DATABASE_URL: postgres://postgres:${DB_PASSWORD}@database:5432/human_intelligence
      JWT_SECRET: ${JWT_SECRET}
      PORT: 8081
      GIN_MODE: release
    ports:
      - "8081:8081"
    depends_on:
      - database

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      VITE_API_URL: http://localhost:8081
    depends_on:
      - backend

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl
    depends_on:
      - frontend
      - backend

volumes:
  postgres_data:
```

---

## 🌐 **Domain & SSL Configuration**

### **DNS Records**
```bash
# Main domain
A     @           → [Vercel IP]
CNAME www         → your-app.vercel.app

# API subdomain  
CNAME api         → your-backend.railway.app

# Optional: Admin subdomain
CNAME admin       → your-admin.vercel.app
```

### **SSL Certificates**
- **Vercel**: Automatic SSL via Let's Encrypt
- **Railway**: Automatic SSL for custom domains
- **Supabase**: SSL included

---

## 📊 **Cost Breakdown**

### **Recommended Stack**
```bash
Vercel (Frontend):        $0/month    (Free tier)
Railway (Backend):        $5/month    (Hobby plan)
Supabase (Database):      $0/month    (Free 500MB)
Domain (optional):        $12/year    (~$1/month)
TOTAL:                    ~$6/month
```

### **Scaling Costs**
```bash
# When you outgrow free tiers:
Vercel Pro:              $20/month   (Better performance)
Railway Pro:             $20/month   (More resources)
Supabase Pro:            $25/month   (8GB database)
CDN (Cloudflare):        $0/month    (Free tier)
```

---

## 🔍 **Monitoring & Maintenance**

### **Health Checks**
```bash
# Railway: Built-in health checks
# Vercel: Automatic deployment health
# Supabase: Database monitoring dashboard
```

### **Logging**
```bash
# Railway: Real-time logs in dashboard
# Vercel: Function logs and analytics
# Supabase: Database logs and metrics
```

### **Backups**
```bash
# Supabase: Automatic daily backups (free tier)
# Railway: Database backups available
# Code: GitHub repository (version control)
```

---

## 🚨 **Troubleshooting Guide**

### **Common Issues**

#### **CORS Errors**
```go
// backend/internal/config/cors.go
func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:3000",
            "https://your-app.vercel.app",  // Add production URL
        },
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    })
}
```

#### **Database Connection Issues**
```bash
# Check connection string format
# Ensure Supabase allows external connections
# Verify environment variables are set
```

#### **WebSocket Issues**
```bash
# Railway supports WebSockets by default
# Check WS_URL environment variable
# Ensure frontend connects to wss:// (not ws://)
```

### **Deployment Checklist**
- [ ] Database created and migrations run
- [ ] Backend deployed with correct environment variables
- [ ] Frontend built and deployed
- [ ] CORS configured for production domain
- [ ] Custom domains configured (if applicable)
- [ ] SSL certificates working
- [ ] Health checks passing
- [ ] User registration/login tested in production

---

## 🎯 **Quick Deploy Commands**

```bash
# Deploy backend to Railway
git push origin main  # Railway auto-deploys

# Deploy frontend to Vercel
cd frontend
vercel --prod

# Check deployment status
vercel ls
railway status
```

---

## 🔄 **CI/CD Pipeline (GitHub Actions)**

**.github/workflows/deploy.yml:**
```yaml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  deploy-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Deploy to Railway
        run: |
          # Railway CLI deployment
          echo "Backend deploys automatically via Railway GitHub integration"

  deploy-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Deploy to Vercel
        run: |
          cd frontend
          npm ci
          npm run build
          npx vercel --prod --token ${{ secrets.VERCEL_TOKEN }}
```

---

## 🎉 **Go Live Checklist**

### **Pre-Launch**
- [ ] Test all features in production environment
- [ ] Verify database seeding works
- [ ] Test user registration and authentication
- [ ] Check WebSocket connections
- [ ] Validate file uploads (if applicable)
- [ ] Performance testing
- [ ] Security audit

### **Launch**
- [ ] Update DNS records
- [ ] Monitor error logs
- [ ] Test from multiple devices/browsers
- [ ] Set up monitoring alerts
- [ ] Backup production database

### **Post-Launch**
- [ ] Monitor performance metrics
- [ ] Set up regular backups
- [ ] Plan scaling strategy
- [ ] Document any issues and solutions

---

**🚀 Ready to deploy your Human Intelligence platform!**

*Choose your deployment strategy and follow the relevant sections above. The recommended Vercel + Railway + Supabase stack provides the best balance of simplicity, performance, and cost.*