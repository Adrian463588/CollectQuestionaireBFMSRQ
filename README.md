# CollectBFMSRQ - Questionnaire Scoring System

A web-based responsive questionnaire system for IPIP-BFM-50 (personality) and SRQ-29 (mental health screening) with automated scoring and interpretation.

## 📋 Features

- **Participant Management**: Collect participant identity information
- **Dual Questionnaires**: 
  - SRQ-29: Mental health screening (29 yes/no questions)
  - IPIP-BFM-50: Personality assessment (50 Likert scale questions)
- **Automated Scoring**: Real-time scoring with interpretation
- **Results Dashboard**: Visual display of scores with recommendations
- **Export Functionality**: CSV export for data analysis
- **Responsive Design**: Works on desktop, tablet, and mobile
- **Smooth UX**: Micro-interactions with Framer Motion

## 🏗️ Architecture

```
CollectBFMSRQ/
├── frontend/              # Next.js frontend application
│   ├── app/              # App Router pages
│   ├── components/       # Reusable React components
│   ├── data/            # Questionnaire data files
│   ├── types/           # TypeScript type definitions
│   └── package.json
├── backend/              # Go backend API
│   ├── cmd/server/      # Main server entry point
│   ├── internal/        # Internal packages
│   │   ├── config/      # Configuration management
│   │   ├── database/    # Database connection & migrations
│   │   ├── handlers/    # HTTP handlers
│   │   ├── models/      # Data models
│   │   └── services/    # Business logic
│   ├── pkg/utils/       # Public utilities
│   └── go.mod
└── AGENTS.md            # Project specification
```

## 🚀 Quick Start

### Prerequisites

- Node.js 18+ and npm
- Go 1.21+
- PostgreSQL 12+

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Create `.env` file from example:
```bash
copy .env.example .env
```

3. Edit `.env` with your PostgreSQL credentials:
```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=questionnaire_db
```

4. Create PostgreSQL database:
```sql
CREATE DATABASE questionnaire_db;
```

5. Run the backend server:
```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Run the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## 📡 API Endpoints

### Participants
- `POST /api/participants` - Create a new participant
- `GET /api/participants/:id` - Get participant by ID

### Responses
- `POST /api/responses` - Submit questionnaire responses

### Scoring
- `POST /api/scoring` - Calculate scores from answers
- `GET /api/scores/:participantId` - Get scores for a participant

### Export
- `GET /api/export/:participantId` - Export participant data to CSV

## 🧮 Scoring Logic

### SRQ-29

**Neurotic (Q1-Q20)**:
- Score = sum of "Yes" answers
- Score >= 5: Indikasi masalah emosional
- Score >= 6: Rekomendasi rujukan

**Substance Use (Q21)**:
- Yes: Indikasi penggunaan zat

**Psychotic (Q22-Q24)**:
- Any Yes: Indikasi psikotik (urgent)

**PTSD (Q25-Q29)**:
- Any Yes: Indikasi PTSD

### IPIP-BFM-50

**Positive Key (+)**: STS=1, TS=2, N=3, S=4, SS=5
**Negative Key (-)**: STS=5, TS=4, N=3, S=2, SS=1

**Dimensions**:
- Extraversion: + items (1,11,21,31,41), - items (6,16,26,36,46)
- Agreeableness: + items (7,17,27,37,42,47), - items (2,12,22,32)
- Conscientiousness: + items (3,13,23,33,43,48), - items (8,18,28,38)
- Emotional Stability: + items (9,19), - items (4,14,24,29,34,39,44,49)
- Intellect: + items (5,15,25,35,40,45,50), - items (10,20,30)

## 🎨 Design System

### Color Palette
- Primary: `#1c0f13` (dark background)
- Secondary: `#6e7e85`
- Accent: `#b7cece`
- Neutral: `#bbbac6`
- Light: `#e2e2e2`

### Micro-interactions
- Progress bar per questionnaire
- Animated radio selection
- Smooth step transitions
- Auto-save indicator
- Result reveal animation

## 🛠️ Tech Stack

### Frontend
- **Next.js 14** (App Router)
- **TypeScript**
- **Tailwind CSS**
- **Framer Motion**

### Backend
- **Go** with Fiber framework
- **PostgreSQL**
- **RESTful API**

## 📤 Export Format

CSV export includes:
- Participant data
- Answers
- Scores
- Interpretations

## 🔐 Security Considerations

- HTTPS in production
- Data anonymization
- Rate limiting
- Input validation
- CORS configuration

## 🧪 Testing

Backend tests (to be added):
```bash
cd backend
go test ./...
```

Frontend tests (to be added):
```bash
cd frontend
npm test
```

## 🚀 Production Deployment

### Frontend (Vercel/Edge)
```bash
npm run build
```

### Backend (Heroku/Cloud)
```bash
# Build binary
cd backend
go build -o server cmd/server/main.go

# Run server
./server
```

## 📈 Future Enhancements

- Dashboard analytics
- Multi-user role system
- AI insights
- Longitudinal tracking
- PDF export
- Admin panel

## 📄 License

This project is built based on the specification in AGENTS.md

## 👥 Contributors

Built following SOLID principles, DRY, and clean code practices.

## 📞 Support

For issues or questions, please refer to the project documentation in AGENTS.md
