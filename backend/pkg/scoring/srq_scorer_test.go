package scoring

import (
	"fmt"
	"testing"
)

// helpers for building answer maps in tests
func boolAnswers(trueItems []int, total int) map[string]interface{} {
	m := make(map[string]interface{}, total)
	for i := 1; i <= total; i++ {
		m[itoa(i)] = false
	}
	for _, n := range trueItems {
		m[itoa(n)] = true
	}
	return m
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

// ── neuroticStatus / emotional classification ─────────────────────────────────

func TestSRQ_EmotionalNormal(t *testing.T) {
	scorer := NewSRQScorer()
	// 4 YA on Q1-20 → normal
	answers := boolAnswers([]int{1, 2, 3, 4}, 29)
	res, err := scorer.Calculate(answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "neurotic_score", 4, res.NeuroticScore)
	assertEqual(t, "neurotic_status", "normal", res.NeuroticStatus)
	assertEqual(t, "emotional_disorder", 0, res.EmotionalDisorder)
	assertEqual(t, "overall_risk", "rendah", res.OverallRisk)
}

func TestSRQ_EmotionalIndikasi(t *testing.T) {
	scorer := NewSRQScorer()
	// 5 YA on Q1-20 → indikasi
	answers := boolAnswers([]int{1, 2, 3, 4, 5}, 29)
	res, err := scorer.Calculate(answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "neurotic_score", 5, res.NeuroticScore)
	assertEqual(t, "neurotic_status", "indikasi_gme", res.NeuroticStatus)
	assertEqual(t, "emotional_disorder", 1, res.EmotionalDisorder)
	assertEqual(t, "overall_risk", "sedang", res.OverallRisk)
}

func TestSRQ_EmotionalRujukan(t *testing.T) {
	scorer := NewSRQScorer()
	// 6 YA on Q1-20 → rekomendasi rujukan
	answers := boolAnswers([]int{1, 2, 3, 4, 5, 6}, 29)
	res, err := scorer.Calculate(answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "neurotic_status", "rekomendasi_rujukan", res.NeuroticStatus)
}

// ── Domain 2: Substance ───────────────────────────────────────────────────────

func TestSRQ_SubstanceUse(t *testing.T) {
	scorer := NewSRQScorer()
	answers := boolAnswers([]int{21}, 29)
	res, _ := scorer.Calculate(answers)
	assertBool(t, "substance_use", true, res.SubstanceUse)
	assertEqual(t, "substance_dummy", 1, res.SubstanceDummy)
}

// ── Domain 3: Psychotic ───────────────────────────────────────────────────────

func TestSRQ_Psychotic_OneItem(t *testing.T) {
	scorer := NewSRQScorer()
	answers := boolAnswers([]int{22}, 29)
	res, _ := scorer.Calculate(answers)
	assertBool(t, "psychotic", true, res.Psychotic)
	assertEqual(t, "psychotic_count", 1, res.PsychoticCount)
	assertEqual(t, "psychotic_dummy", 1, res.PsychoticDummy)
	assertEqual(t, "overall_risk", "kritis", res.OverallRisk)
}

func TestSRQ_Psychotic_AllItems(t *testing.T) {
	scorer := NewSRQScorer()
	answers := boolAnswers([]int{22, 23, 24}, 29)
	res, _ := scorer.Calculate(answers)
	assertEqual(t, "psychotic_count", 3, res.PsychoticCount)
}

func TestSRQ_NoPsychotic(t *testing.T) {
	scorer := NewSRQScorer()
	answers := boolAnswers(nil, 29)
	res, _ := scorer.Calculate(answers)
	assertBool(t, "psychotic", false, res.Psychotic)
}

// ── Domain 4: PTSD ────────────────────────────────────────────────────────────

func TestSRQ_PTSD(t *testing.T) {
	scorer := NewSRQScorer()
	answers := boolAnswers([]int{25, 27}, 29)
	res, _ := scorer.Calculate(answers)
	assertBool(t, "ptsd", true, res.PTSD)
	assertEqual(t, "ptsd_count", 2, res.PTSDCount)
	assertEqual(t, "ptsd_dummy", 1, res.PTSDDummy)
	assertEqual(t, "overall_risk", "kritis", res.OverallRisk)
}

// ── TotalScore ────────────────────────────────────────────────────────────────

func TestSRQ_TotalScore_AllYes(t *testing.T) {
	scorer := NewSRQScorer()
	all := make([]int, 29)
	for i := range all {
		all[i] = i + 1
	}
	answers := boolAnswers(all, 29)
	res, _ := scorer.Calculate(answers)
	assertEqual(t, "total_score", 29, res.TotalScore)
	assertEqual(t, "neurotic_score", 20, res.NeuroticScore)
}

func TestSRQ_TotalScore_AllNo(t *testing.T) {
	scorer := NewSRQScorer()
	answers := boolAnswers(nil, 29)
	res, _ := scorer.Calculate(answers)
	assertEqual(t, "total_score", 0, res.TotalScore)
}

// ── Overall Risk ──────────────────────────────────────────────────────────────

func TestSRQ_OverallRisk_Tinggi(t *testing.T) {
	scorer := NewSRQScorer()
	// emotional (≥5) + substance → tinggi
	answers := boolAnswers([]int{1, 2, 3, 4, 5, 21}, 29)
	res, _ := scorer.Calculate(answers)
	assertEqual(t, "overall_risk", "tinggi", res.OverallRisk)
}

// ── test helpers ──────────────────────────────────────────────────────────────

func assertEqual[T comparable](t *testing.T, field string, want, got T) {
	t.Helper()
	if want != got {
		t.Errorf("field %q: want %v, got %v", field, want, got)
	}
}

func assertBool(t *testing.T, field string, want, got bool) {
	t.Helper()
	if want != got {
		t.Errorf("field %q: want %v, got %v", field, want, got)
	}
}
