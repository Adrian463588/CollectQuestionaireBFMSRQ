package scoring

import "backend/internal/models"

// srqConfig holds configurable cut-off thresholds for SRQ-29.
// Using a config struct follows the Open/Closed principle:
// cut-offs can be adjusted without changing scoring logic.
type srqConfig struct {
	// EmotionalCutoffIndikasi is the minimum neurotic score for "indikasi GME" classification.
	// WHO standard: ≥5 (Scribd SRQ-29 doc; psikiatri.fk.undip.ac.id)
	EmotionalCutoffIndikasi int
	// EmotionalCutoffRujukan is the minimum score for "rekomendasi rujukan" classification.
	EmotionalCutoffRujukan int
}

// defaultSRQConfig returns WHO-recommended default thresholds.
func defaultSRQConfig() srqConfig {
	return srqConfig{
		EmotionalCutoffIndikasi: 5,
		EmotionalCutoffRujukan:  6,
	}
}

// SRQScorer implements Scorer[*models.SRQScore] for the SRQ-29 instrument.
// Domain breakdown per WHO guideline and official SRQ-29 instrument:
//   - Q1–20  : Gangguan Mental Emosional (neurotic/GME)
//   - Q21    : Penggunaan Zat Psikoaktif
//   - Q22–24 : Gejala Psikotik
//   - Q25–29 : Gejala PTSD
type SRQScorer struct {
	config srqConfig
}

// NewSRQScorer creates a new SRQScorer with default WHO thresholds.
func NewSRQScorer() *SRQScorer {
	return &SRQScorer{config: defaultSRQConfig()}
}

// Calculate computes SRQ-29 scores from a map of boolean answers keyed by
// question number string (e.g. "1"…"29"). true = YA (1), false = TIDAK (0).
func (s *SRQScorer) Calculate(answers map[string]interface{}) (*models.SRQScore, error) {
	// ── Domain 1: Gangguan Mental Emosional (Q1–20) ──────────────────────────
	emotionalScore := SumBoolRange(answers, 1, 20)
	emotionalStatus := s.classifyEmotional(emotionalScore)

	// ── Domain 2: Penggunaan Zat Psikoaktif (Q21) ────────────────────────────
	substanceUse := SafeBool(answers, "21")

	// ── Domain 3: Gejala Psikotik (Q22–24) ───────────────────────────────────
	psychoticCount := CountBoolInRange(answers, 22, 24)
	psychotic := psychoticCount >= 1

	// ── Domain 4: Gejala PTSD (Q25–29) ───────────────────────────────────────
	ptsdCount := CountBoolInRange(answers, 25, 29)
	ptsd := ptsdCount >= 1

	// ── Aggregate ─────────────────────────────────────────────────────────────
	totalScore := SumBoolRange(answers, 1, 29)
	overallRisk := s.classifyOverallRisk(emotionalScore, substanceUse, psychotic, ptsd)

	// ── Research dummy variables (0/1) ────────────────────────────────────────
	emotionalDisorder := boolToInt(emotionalScore >= s.config.EmotionalCutoffIndikasi)
	substanceDummy := boolToInt(substanceUse)
	psychoticDummy := boolToInt(psychotic)
	ptsdDummy := boolToInt(ptsd)

	return &models.SRQScore{
		// Core fields
		NeuroticScore:  emotionalScore,
		NeuroticStatus: emotionalStatus,
		SubstanceUse:   substanceUse,
		Psychotic:      psychotic,
		PTSD:           ptsd,
		// Extended fields
		PsychoticCount:    psychoticCount,
		PTSDCount:         ptsdCount,
		TotalScore:        totalScore,
		OverallRisk:       overallRisk,
		EmotionalDisorder: emotionalDisorder,
		SubstanceDummy:    substanceDummy,
		PsychoticDummy:    psychoticDummy,
		PTSDDummy:         ptsdDummy,
	}, nil
}

// classifyEmotional returns a clinical interpretation label for the neurotic score.
// Based on WHO SRQ-29 guidelines (two-tier cut-off).
func (s *SRQScorer) classifyEmotional(score int) string {
	switch {
	case score >= s.config.EmotionalCutoffRujukan:
		return "rekomendasi_rujukan" // ≥6 → refer to mental health professional
	case score >= s.config.EmotionalCutoffIndikasi:
		return "indikasi_gme" // ≥5 → emotional disorder indication
	default:
		return "normal"
	}
}

// classifyOverallRisk derives a global risk level from all four SRQ-29 domains.
// Logic follows Sprint2.md global interpretation table:
//   - kritis  : psychotic OR PTSD present (serious, requires urgent referral)
//   - tinggi  : emotional disorder AND (substance OR multi-domain)
//   - sedang  : emotional disorder only (score ≥5)
//   - rendah  : all domains normal
func (s *SRQScorer) classifyOverallRisk(emotionalScore int, substance, psychotic, ptsd bool) string {
	hasEmotional := emotionalScore >= s.config.EmotionalCutoffIndikasi
	hasSevere := psychotic || ptsd

	switch {
	case hasSevere:
		return "kritis"
	case hasEmotional && substance:
		return "tinggi"
	case hasEmotional:
		return "sedang"
	default:
		return "rendah"
	}
}

// boolToInt converts a boolean to 1 (true) or 0 (false).
// Used to create dummy variables suitable for logistic regression / SEM.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
