package models

import (
	"time"
)

// Participant represents a person who takes the questionnaire
type Participant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Age       int       `json:"age" db:"age"`
	Gender    string    `json:"gender" db:"gender"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Response represents answers to a questionnaire
type Response struct {
	ID                string                 `json:"id" db:"id"`
	ParticipantID     string                 `json:"participant_id" db:"participant_id"`
	QuestionnaireType string                 `json:"questionnaire_type" db:"questionnaire_type"`
	Answers           map[string]interface{} `json:"answers" db:"answers"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
}

// Score represents calculated scores from questionnaire responses
type Score struct {
	ID            string     `json:"id" db:"id"`
	ParticipantID string     `json:"participant_id" db:"participant_id"`
	SRQScore      *SRQScore  `json:"srq_score,omitempty" db:"srq_score"`
	IPIPScore     *IPIPScore `json:"ipip_score,omitempty" db:"ipip_score"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// SRQScore represents SRQ-29 scoring results
type SRQScore struct {
	NeuroticScore  int    `json:"neurotic_score"`
	NeuroticStatus string `json:"neurotic_status"`
	SubstanceUse   bool   `json:"substance_use"`
	Psychotic      bool   `json:"psychotic"`
	PTSD           bool   `json:"ptsd"`
}

// IPIPScore represents IPIP-BFM-50 scoring results
type IPIPScore struct {
	Extraversion       int `json:"extraversion"`
	Agreeableness      int `json:"agreeableness"`
	Conscientiousness  int `json:"conscientiousness"`
	EmotionalStability int `json:"emotional_stability"`
	Intellect          int `json:"intellect"`
}

// SRQAnswer represents a single SRQ-29 answer
type SRQAnswer struct {
	QuestionNumber int  `json:"question_number"`
	Answer         bool `json:"answer"` // true = Ya, false = Tidak
}

// IPIPAnswer represents a single IPIP-BFM-50 answer
type IPIPAnswer struct {
	QuestionNumber int `json:"question_number"`
	Score          int `json:"score"` // 1-5 Likert scale
}

// SubmissionRequest represents the request body for submitting a questionnaire
type SubmissionRequest struct {
	ParticipantID     string                 `json:"participant_id"`
	QuestionnaireType string                 `json:"questionnaire_type"`
	Answers           map[string]interface{} `json:"answers"`
}

// ParticipantRequest represents the request body for creating a participant
type ParticipantRequest struct {
	Name   string `json:"name" validate:"required"`
	Age    int    `json:"age" validate:"required,min=15"`
	Gender string `json:"gender" validate:"required,oneof=male female"`
}
