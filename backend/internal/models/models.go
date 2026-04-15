package models

import "time"

// Participant represents a person who takes the questionnaire.
type Participant struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Age       int        `json:"age" db:"age"`
	Gender    string     `json:"gender" db:"gender"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Response represents answers to a questionnaire.
type Response struct {
	ID                string                 `json:"id" db:"id"`
	ParticipantID     string                 `json:"participant_id" db:"participant_id"`
	QuestionnaireType string                 `json:"questionnaire_type" db:"questionnaire_type"`
	Answers           map[string]interface{} `json:"answers" db:"answers"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
}

// Score represents calculated scores from questionnaire responses.
type Score struct {
	ID            string     `json:"id" db:"id"`
	ParticipantID string     `json:"participant_id" db:"participant_id"`
	SRQScore      *SRQScore  `json:"srq_score,omitempty" db:"srq_score"`
	IPIPScore     *IPIPScore `json:"ipip_score,omitempty" db:"ipip_score"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// SRQScore represents SRQ-29 scoring results.
// All existing JSON keys are preserved for backward compatibility.
// New fields are additive and safe for existing consumers.
//
// Scoring is domain-based per WHO SRQ-29 guidelines:
//
//	Q1–20  : Gangguan Mental Emosional (GME)
//	Q21    : Penggunaan Zat Psikoaktif
//	Q22–24 : Gejala Psikotik
//	Q25–29 : Gejala PTSD
type SRQScore struct {
	// ── Core domain fields (backward compatible) ──────────────────────────────
	NeuroticScore  int    `json:"neurotic_score"`  // sum of YA on Q1–20
	NeuroticStatus string `json:"neurotic_status"` // "normal" | "indikasi_gme" | "rekomendasi_rujukan"
	SubstanceUse   bool   `json:"substance_use"`   // Q21 = YA
	Psychotic      bool   `json:"psychotic"`       // any of Q22–24 = YA
	PTSD           bool   `json:"ptsd"`            // any of Q25–29 = YA

	// ── Extended domain detail ────────────────────────────────────────────────
	PsychoticCount int `json:"psychotic_count"` // count of YA in Q22–24 (0–3)
	PTSDCount      int `json:"ptsd_count"`      // count of YA in Q25–29 (0–5)

	// ── Aggregate scores ──────────────────────────────────────────────────────
	TotalScore  int    `json:"total_score"`  // sum of YA Q1–29 (for descriptive / regression)
	OverallRisk string `json:"overall_risk"` // "rendah" | "sedang" | "tinggi" | "kritis"

	// ── Research dummy variables (0 = tidak, 1 = ya) ──────────────────────────
	// These binary variables are ready for use in logistic regression or SEM.
	EmotionalDisorder int `json:"emotional_disorder"` // 1 if neurotic_score ≥ 5
	SubstanceDummy    int `json:"substance_dummy"`    // 1 if substance_use
	PsychoticDummy    int `json:"psychotic_dummy"`    // 1 if psychotic
	PTSDDummy         int `json:"ptsd_dummy"`         // 1 if ptsd
}

// IPIPScore represents IPIP-BFM-50 scoring results.
// Scores are computed as MEAN of 10 items per dimension (range 1.0–5.0)
// after reverse-scoring negatively-keyed items.
// This follows the recommended method for psychometric research publications
// (Akhtar & Azwar, 2019; ipip.ori.org).
type IPIPScore struct {
	// ── Dimension mean scores (1.0–5.0) — primary scoring ────────────────────
	Extraversion       float64 `json:"extraversion"`
	Agreeableness      float64 `json:"agreeableness"`
	Conscientiousness  float64 `json:"conscientiousness"`
	EmotionalStability float64 `json:"emotional_stability"`
	Intellect          float64 `json:"intellect"`

	// ── Interpretation labels per dimension ───────────────────────────────────
	// "sangat_tinggi" | "tinggi" | "rata_rata" | "rendah" | "sangat_rendah"
	// Based on absolute mean thresholds on the 1–5 Likert scale.
	// For norm-referenced interpretation, standardize against sample distribution
	// using external statistical software (SPSS/R) after data export.
	ExtraLabel  string `json:"extra_label"`
	AgreLabel   string `json:"agre_label"`
	ConsLabel   string `json:"cons_label"`
	StabLabel   string `json:"stab_label"`
	IntellLabel string `json:"intell_label"`

	// ── Raw sum scores (integer, range 10–50) — preserved for alternative analysis ──
	ExtraversionSum       int `json:"extraversion_sum"`
	AgreeablenessSum      int `json:"agreeableness_sum"`
	ConscientiousnessSum  int `json:"conscientiousness_sum"`
	EmotionalStabilitySum int `json:"emotional_stability_sum"`
	IntellectSum          int `json:"intellect_sum"`
}

// SRQAnswer represents a single SRQ-29 answer.
type SRQAnswer struct {
	QuestionNumber int  `json:"question_number"`
	Answer         bool `json:"answer"` // true = Ya, false = Tidak
}

// IPIPAnswer represents a single IPIP-BFM-50 answer.
type IPIPAnswer struct {
	QuestionNumber int `json:"question_number"`
	Score          int `json:"score"` // 1-5 Likert scale
}

// SubmissionRequest represents the request body for submitting a questionnaire.
type SubmissionRequest struct {
	ParticipantID     string                 `json:"participant_id"`
	QuestionnaireType string                 `json:"questionnaire_type"`
	Answers           map[string]interface{} `json:"answers"`
}

// ParticipantRequest represents the request body for creating a participant.
type ParticipantRequest struct {
	Name   string `json:"name" validate:"required"`
	Age    int    `json:"age" validate:"required,min=15"`
	Gender string `json:"gender" validate:"required,oneof=male female"`
}
