B
# 📄 AGENTS.md

## Web Responsive Questionnaire System (IPIP-BFM-50 & SRQ-29)

---

# 1. 📌 Overview

## 1.1 Purpose

Membangun aplikasi web responsif untuk:

* Pengisian kuesioner **IPIP-BFM-50 (kepribadian)**
* Pengisian kuesioner **SRQ-29 (skrining kesehatan mental)**
* Otomatisasi skoring & interpretasi
* Ekspor hasil per partisipan

## 1.2 Target User

* Peneliti
* Psikolog / tenaga kesehatan
* Institusi pendidikan / organisasi

## 1.3 Key Objectives

* UX sederhana & cepat (micro-interaction driven)
* Akurasi skoring sesuai standar instrumen
* Data aman & dapat diekspor

---

# 2. 🧩 Scope & Features

## 2.1 Core Features

| Feature                | Description                         |
| ---------------------- | ----------------------------------- |
| Participant Form       | Input identitas (nama, usia, dll)   |
| Questionnaire Engine   | Render dinamis SRQ-29 & IPIP-BFM-50 |
| Auto Scoring           | Real-time scoring                   |
| Result Dashboard       | Hasil per partisipan                |
| Export                 | CSV / Excel / PDF                   |
| Admin Panel (optional) | Monitoring data                     |

---

# 3. 🧠 Questionnaire Logic

---

## 3.1 SRQ-29 Logic

### Data Structure

* Total: 29 pertanyaan
* Tipe: Boolean (Ya/Tidak)

### Scoring Rule

* YA = 1
* TIDAK = 0 

---

### Interpretation Logic

#### A. Neurotic (Q1–20)

```
score = sum(Q1–Q20)
if score >= 5 → indikasi masalah emosional
if score >= 6 → rekomendasi rujukan
```

#### B. Substance Use (Q21)

```
if Q21 == YA → indikasi penggunaan zat
```

#### C. Psychotic (Q22–24)

```
if any(Q22–Q24 == YA) → indikasi psikotik (urgent)
```

#### D. PTSD (Q25–29)

```
if any(Q25–Q29 == YA) → indikasi PTSD
```

📌 Berdasarkan instrumen resmi (halaman 2 PDF SRQ-29) 

---

## 3.2 IPIP-BFM-50 Logic

### Data Structure

* 50 item
* Skala Likert 1–5

### Scoring

#### Positive Key (+)

| Jawaban | Skor |
| ------- | ---- |
| STS     | 1    |
| TS      | 2    |
| N       | 3    |
| S       | 4    |
| SS      | 5    |

#### Negative Key (-)

| Jawaban | Skor |
| ------- | ---- |
| STS     | 5    |
| TS      | 4    |
| N       | 3    |
| S       | 2    |
| SS      | 1    |

---

### Dimensi

| Dimensi             | + Items             | - Items                |
| ------------------- | ------------------- | ---------------------- |
| Extraversion        | 1,11,21,31,41       | 6,16,26,36,46          |
| Agreeableness       | 7,17,27,37,42,47    | 2,12,22,32             |
| Conscientiousness   | 3,13,23,33,43,48    | 8,18,28,38             |
| Emotional Stability | 9,19                | 4,14,24,29,34,39,44,49 |
| Intellect           | 5,15,25,35,40,45,50 | 10,20,30               |

📌 Berdasarkan dokumen IPIP-BFM-50 

---

### Final Score

```
dimension_score = sum(item_scores)
```

---

# 4. 🧱 System Architecture

## 4.1 High-Level Architecture

```
[ Next.js Frontend ]
        ↓
[ API Gateway (Go Fiber / Gin) ]
        ↓
[ Service Layer ]
        ↓
[ Database (PostgreSQL) ]
```

---

## 4.2 Tech Stack

### Frontend

* Next.js (App Router)
* TypeScript
* Tailwind CSS
* Framer Motion

### Backend

* Golang (Gin / Fiber)
* PostgreSQL


---

# 5. 🎨 UI/UX Design System

## 5.1 Color Palette

| Color   | Usage                     |
| ------- | ------------------------- |
| #1c0f13 | Primary (background dark) |
| #6e7e85 | Secondary                 |
| #b7cece | Accent                    |
| #bbbac6 | Neutral                   |
| #e2e2e2 | Background light          |

---

## 5.2 Micro Interactions

* Progress bar (per questionnaire)
* Animated radio selection (Framer Motion)
* Smooth step transitions
* Auto-save indicator
* Result reveal animation

---

## 5.3 UX Flow

```
Landing
 → Input Data
 → Pilih Kuesioner
 → Isi (step-by-step)
 → Submit
 → Result Page
 → Export
```

---

# 6. 🧮 API Design

## 6.1 Endpoint

### Participant

```
POST /participants
GET /participants/:id
```

### Questionnaire

```
GET /questionnaires/srq29
GET /questionnaires/ipip
```

### Submission

```
POST /responses
```

### Scoring

```
POST /scoring
```

### Export

```
GET /export/:participantId
```

---

## 6.2 Example Response (Scoring)

```json
{
  "participant_id": "uuid",
  "srq29": {
    "neurotic_score": 7,
    "neurotic_status": "indikasi",
    "substance": true,
    "psychotic": false,
    "ptsd": true
  },
  "ipip": {
    "extraversion": 32,
    "agreeableness": 28,
    "conscientiousness": 35,
    "emotional_stability": 20,
    "intellect": 30
  }
}
```

---

# 7. 🗄️ Database Design

## 7.1 Tables

### participants

```
id (uuid)
name
age
gender
created_at
```

### responses

```
id
participant_id
questionnaire_type (srq/ipip)
answers (jsonb)
```

### scores

```
id
participant_id
srq_score (jsonb)
ipip_score (jsonb)
```

---

# 8. 📤 Export Feature

## Format:

* CSV


## Content:

* Data partisipan
* Jawaban
* Skoring
* Interpretasi

---

# 9. 🔐 Security

* HTTPS
* Data anonymization 
* Rate limiting

---

# 10. ⚙️ Performance

* SSR untuk SEO
* Client-side hydration ringan
* Lazy loading questionnaire
* Debounced autosave

---

# 11. 🧪 Testing

* Unit Test (Go + Jest)
* Integration Test
* E2E (Playwright)

---

# 12. 🚀 Deployment

* Frontend: Edgeone pages
* Backend: Heroku 
* DB: Managed PostgreSQL
* Bisa juga run local windows

---

# 13. 📈 Future Enhancements

* Dashboard analytics
* Multi-user role
* AI insight (opsional)
* Longitudinal tracking

---

# 14. 📌 Acceptance Criteria

* ✅ User dapat mengisi 2 kuesioner tanpa reload
* ✅ Skoring otomatis akurat sesuai dokumen
* ✅ Hasil dapat diekspor
* ✅ UI responsif & smooth

---


