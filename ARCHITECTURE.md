# System Architecture

## High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        USER BROWSER                         │
│                    http://localhost:3000                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ HTTP Requests
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   FRONTEND (Next.js)                        │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  App Router (React Server Components)                │   │
│  │                                                       │   │
│  │  Pages:                                               │   │
│  │  • / (Landing - Participant Form)                    │   │
│  │  • /questionnaire-select                             │   │
│  │  • /questionnaire/srq29                              │   │
│  │  • /questionnaire/ipip                               │   │
│  │  • /results/[type]                                   │   │
│  │                                                       │   │
│  │  Technologies:                                        │   │
│  │  • TypeScript                                         │   │
│  │  • Tailwind CSS                                       │   │
│  │  • Framer Motion                                      │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ Fetch API Calls
                         │ http://localhost:8080/api/*
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  BACKEND (Go Fiber)                         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Router                                               │   │
│  │                                                       │   │
│  │  Routes:                                              │   │
│  │  • POST /api/participants                            │   │
│  │  • GET  /api/participants/:id                        │   │
│  │  • POST /api/responses                               │   │
│  │  • POST /api/scoring                                 │   │
│  │  • GET  /api/scores/:participantId                   │   │
│  │  • GET  /api/export/:participantId                   │   │
│  └────────┬─────────────────────────────────────────────┘   │
           │
           │
           ▼
┌─────────────────────────────────────────────────────────────┐
│                    SERVICE LAYER                            │
│  ┌──────────────────┐  ┌───────────────┐  ┌──────────────┐  │
│  │ParticipantService│  │ResponseService│  │ScoreService  │  │
│  │                  │  │               │  │              │  │
│  │• Create          │  │• Save         │  │• Save        │  │
│  │• Get             │  │               │  │• Get         │  │
│  └──────────────────┘  └───────────────┘  └──────┬───────┘  │
│                                                   │          │
│                                    ┌──────────────┴───────┐  │
│                                    │  ScoringService     │  │
│                                    │                     │  │
│                                    │• CalculateSRQScore  │  │
│                                    │• CalculateIPIPScore │  │
│                                    └─────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ SQL Queries
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  DATABASE (PostgreSQL)                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Tables:                                              │   │
│  │                                                       │   │
│  │  participants                                         │   │
│  │  ├─ id (UUID, PK)                                    │   │
│  │  ├─ name (VARCHAR)                                   │   │
│  │  ├─ age (INTEGER)                                    │   │
│  │  ├─ gender (VARCHAR)                                 │   │
│  │  └─ created_at (TIMESTAMP)                           │   │
│  │                                                       │   │
│  │  responses                                            │   │
│  │  ├─ id (UUID, PK)                                    │   │
│  │  ├─ participant_id (UUID, FK)                        │   │
│  │  ├─ questionnaire_type (VARCHAR)                     │   │
│  │  ├─ answers (JSONB)                                  │   │
│  │  └─ created_at (TIMESTAMP)                           │   │
│  │                                                       │   │
│  │  scores                                               │   │
│  │  ├─ id (UUID, PK)                                    │   │
│  │  ├─ participant_id (UUID, FK)                        │   │
│  │  ├─ srq_score (JSONB)                                │   │
│  │  ├─ ipip_score (JSONB)                               │   │
│  │  └─ created_at (TIMESTAMP)                           │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## Data Flow

### Participant Registration Flow
```
User fills form → Frontend validates → POST /api/participants
→ ParticipantService.CreateParticipant → INSERT into participants table
→ Return participant UUID → Store in sessionStorage
```

### Questionnaire Submission Flow
```
User answers questions → Frontend collects answers → Submit button clicked
→ POST /api/responses with {participant_id, answers}
→ ResponseService.SaveResponse → INSERT into responses table
→ ScoringService.Calculate[SRQ/IPIP]Score → Calculate scores
→ ScoreService.SaveScore → INSERT into scores table
→ Return complete results → Navigate to results page
```

### Export Flow
```
User clicks "Export to CSV" → GET /api/export/:participantId
→ ParticipantService.GetParticipant → Query participants table
→ ScoreService.GetScores → Query scores table
→ utils.ExportToCSV → Generate CSV string
→ Return CSV file with download headers
```

---

## Component Architecture (Frontend)

```
app/
├── layout.tsx (Root Layout)
│   └── Metadata configuration
│   └── Global font loading
│
├── page.tsx (Landing)
│   └── Participant form
│   └── Session storage
│   └── Navigation to /questionnaire-select
│
├── questionnaire-select/page.tsx
│   └── Card selection UI
│   └── Navigation to /questionnaire/[type]
│
├── questionnaire/
│   ├── srq29/page.tsx
│   │   └── Step-by-step questionnaire
│   │   └── Yes/No buttons
│   │   └── Progress bar
│   │   └── Submit to API
│   │
│   └── ipip/page.tsx
│       └── Step-by-step questionnaire
│       └── Likert scale buttons
│       └── Progress bar
│       └── Submit to API
│
└── results/[type]/page.tsx
    └── Fetch scores from API
    └── Display results
    └── Export to CSV button
```

---

## Service Architecture (Backend)

```
internal/
├── config/
│   └── config.go
│       └── Load environment variables
│       └── Build database URL
│
├── database/
│   └── database.go
│       └── PostgreSQL connection
│       └── Auto-migration
│
├── models/
│   └── models.go
│       └── Data structures
│       └── JSON tags
│
├── services/
│   ├── scoring_service.go
│   │   └── CalculateSRQScore()
│   │   └── CalculateIPIPScore()
│   │   └── calculateDimension()
│   │
│   └── repository.go
│       └── ParticipantService
│       └── ResponseService
│       └── ScoreService
│
├── handlers/
│   └── handlers.go
│       └── HTTP handlers
│       └── Request parsing
│       └── Response formatting
│
└── middleware/
    └── cors.go
        └── CORS configuration

pkg/
└── utils/
    └── export.go
        └── CSV generation
```

---

## Technology Stack Layers

```
┌─────────────────────────────────────────┐
│          Presentation Layer             │
│  • Next.js 16 (App Router)             │
│  • React 19                            │
│  • TypeScript                          │
│  • Tailwind CSS 4                      │
│  • Framer Motion                       │
└─────────────────────────────────────────┘
                  ↕
┌─────────────────────────────────────────┐
│          Application Layer              │
│  • Go 1.21+                            │
│  • Fiber v2 (Web Framework)            │
│  • RESTful API                         │
│  • Business Logic                      │
└─────────────────────────────────────────┘
                  ↕
┌─────────────────────────────────────────┐
│           Data Layer                    │
│  • PostgreSQL 12+                      │
│  • JSONB for flexible data             │
│  • UUID primary keys                   │
│  • Foreign key constraints             │
└─────────────────────────────────────────┘
```

---

## Scoring Algorithm Flow

### SRQ-29 Scoring
```
Input: map[string]bool{"1": true, "2": false, ...}
    ↓
Calculate Neurotic Score (Q1-Q20)
    ↓
Sum "true" answers
    ↓
Determine Status:
  • >= 6 → "rekomendasi_rujukan"
  • >= 5 → "indikasi_masalah_emosional"
  • < 5  → "normal"
    ↓
Check Substance Use (Q21)
    ↓
Check Psychotic (Q22-Q24) - ANY true
    ↓
Check PTSD (Q25-Q29) - ANY true
    ↓
Output: SRQScore struct
```

### IPIP-BFM-50 Scoring
```
Input: map[string]float64{"1": 4, "2": 2, ...}
    ↓
For each dimension:
  • Extraversion
  • Agreeableness
  • Conscientiousness
  • Emotional Stability
  • Intellect
    ↓
Sum positive items (normal scoring)
    ↓
Sum negative items (reverse scoring: 6 - value)
    ↓
Total = positive_sum + negative_sum
    ↓
Output: IPIPScore struct with 5 dimension scores
```

---

## Security Architecture

```
┌─────────────────────────────────────────┐
│           Client Browser                │
└──────────────┬──────────────────────────┘
               │
               │ 1. CORS Preflight Check
               ▼
┌─────────────────────────────────────────┐
│          Fiber CORS Middleware          │
│  • Allowed Origins: * (configure)       │
│  • Allowed Methods: GET,POST,PUT,DELETE │
│  • Allowed Headers: Content-Type, etc.  │
└──────────────┬──────────────────────────┘
               │
               │ 2. Request Validation
               ▼
┌─────────────────────────────────────────┐
│         Handler Validation              │
│  • Required fields check               │
│  • Age >= 15 constraint                │
│  • Gender enum validation              │
│  • Type assertions                     │
└──────────────┬──────────────────────────┘
               │
               │ 3. Parameterized Queries
               ▼
┌─────────────────────────────────────────┐
│       PostgreSQL Driver (lib/pq)        │
│  • SQL injection prevention             │
│  • Type-safe queries                   │
│  • Connection pooling                  │
└─────────────────────────────────────────┘
```

---

**This architecture follows SOLID principles, DRY practices, and clean code standards as specified in AGENTS.md** ✅
