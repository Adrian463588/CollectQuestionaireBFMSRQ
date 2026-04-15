# Quick Start Guide - CollectBFMSRQ

This guide will help you get the application running in 5 minutes.

## Prerequisites

Before you begin, ensure you have:

- **Node.js 18+** and npm installed
- **Go 1.21+** installed  
- **PostgreSQL 12+** installed and running

## Step-by-Step Setup

### 1. Verify Prerequisites

Check if you have the required software:

```bash
node --version
go version
psql --version
```

If any are missing, download and install them:
- Node.js: https://nodejs.org/
- Go: https://golang.org/dl/
- PostgreSQL: https://www.postgresql.org/download/windows/

### 2. Create PostgreSQL Database

Open PostgreSQL command line:

```bash
psql -U postgres
```

Create the database:

```sql
CREATE DATABASE questionnaire_db;
\q
```

### 3. Configure Backend

Navigate to the backend folder:

```bash
cd backend
```

Create a `.env` file (copy from `.env.example`):

```bash
copy .env.example .env
```

Edit `.env` with your database credentials:

```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=questionnaire_db
```

### 4. Start Backend Server

In the first terminal window, run:

```bash
cd backend
go run cmd\server\main.go
```

You should see:
```
Server starting on port 8080
Successfully connected to database
Database migrations completed successfully
```

### 5. Start Frontend Server

Open a **new terminal window** and run:

```bash
cd frontend
npm run dev
```

You should see:
```
✓ Ready in X ms
○ Local:   http://localhost:3000
```

### 6. Open the Application

Open your web browser and go to:

```
http://localhost:3000
```

## Using the Application

### Step 1: Enter Participant Information
- Fill in your name, age (must be 15+), and gender
- Click "Continue to Questionnaire Selection"

### Step 2: Choose a Questionnaire
- **SRQ-29**: Mental health screening (29 yes/no questions)
- **IPIP-BFM-50**: Personality assessment (50 Likert scale questions)

### Step 3: Complete the Questionnaire
- Answer each question
- Use Previous/Next buttons to navigate
- Click Submit when finished

### Step 4: View Results
- See your scores immediately
- Click "Export to CSV" to download your data
- Click "Take Another Questionnaire" to start over

## Troubleshooting

### Backend won't start

**Error: "failed to connect to database"**
- Check that PostgreSQL is running
- Verify your `.env` file credentials
- Ensure the database `questionnaire_db` exists

**Error: "port 8080 already in use"**
- Change `SERVER_PORT` in `.env` to another port (e.g., 8081)
- Update the API URLs in frontend accordingly

### Frontend won't start

**Error: "npm is not recognized"**
- Install Node.js from https://nodejs.org/
- Restart your terminal after installation

**Error: "Cannot find module"**
- Run `npm install` in the frontend directory
- Delete `node_modules` folder and run `npm install` again

### Database errors

**Error: "relation does not exist"**
- The tables haven't been created
- Restart the backend - it auto-creates tables on startup

**Error: "permission denied"**
- Check your PostgreSQL user permissions
- Verify the password in `.env` is correct

## API Testing

You can test the API directly using curl or Postman:

### Create a participant:
```bash
curl -X POST http://localhost:8080/api/participants \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test User\",\"age\":25,\"gender\":\"male\"}"
```

### Submit SRQ-29 responses:
```bash
curl -X POST http://localhost:8080/api/responses \
  -H "Content-Type: application/json" \
  -d "{\"participant_id\":\"YOUR_ID\",\"answers\":{\"1\":true,\"2\":false,\"3\":true}}"
```

### Get scores:
```bash
curl http://localhost:8080/api/scores/YOUR_PARTICIPANT_ID
```

### Export to CSV:
```bash
curl http://localhost:8080/api/export/YOUR_PARTICIPANT_ID -o results.csv
```

## Production Deployment

### Frontend (Vercel)
```bash
cd frontend
npm run build
# Deploy the .next folder to Vercel
```

### Backend (Heroku)
```bash
# Add Procfile to backend folder
echo "web: ./server" > Procfile
go build -o server cmd/server/main.go
# Deploy to Heroku
```

## Need Help?

- Check the main README.md for detailed documentation
- Review AGENTS.md for the complete specification
- Check backend logs in the terminal running the backend
- Check frontend logs in the terminal running the frontend

## Default Configuration

- **Backend URL**: http://localhost:8080
- **Frontend URL**: http://localhost:3000
- **Database**: PostgreSQL on localhost:5432
- **Database Name**: questionnaire_db

---

**You're all set! Start collecting questionnaire data today.** 🚀
