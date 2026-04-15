package scoring

import (
	"fmt"
	"math"
	"testing"
)

// floatAnswers builds an answer map for IPIP (Likert 1-5 values).
func floatAnswers(vals map[int]float64) map[string]interface{} {
	m := make(map[string]interface{}, len(vals))
	for k, v := range vals {
		m[fmt.Sprintf("%d", k)] = v
	}
	return m
}

// allMid returns answers where all 50 items are set to mid value (3).
func allMid() map[string]interface{} {
	m := make(map[string]interface{}, 50)
	for i := 1; i <= 50; i++ {
		m[fmt.Sprintf("%d", i)] = float64(3)
	}
	return m
}

// ── Mean calculation ──────────────────────────────────────────────────────────

// When all items = 3 (neutral), every dimension mean must equal 3.0.
// Positive items: 3. Negative items: 6-3 = 3. Mean = 30/10 = 3.0.
func TestIPIP_AllMid_MeanEqualThree(t *testing.T) {
	scorer := NewIPIPScorer()
	res, err := scorer.Calculate(allMid())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertFloat(t, "extraversion", 3.0, res.Extraversion)
	assertFloat(t, "agreeableness", 3.0, res.Agreeableness)
	assertFloat(t, "conscientiousness", 3.0, res.Conscientiousness)
	assertFloat(t, "emotional_stability", 3.0, res.EmotionalStability)
	assertFloat(t, "intellect", 3.0, res.Intellect)
}

// ── Reverse scoring ───────────────────────────────────────────────────────────

// Extraversion negative items: 6,16,26,36,46. Setting those to 5 should yield
// reversed score of 1 for each. Setting positive items to 5 → 5 each.
// Mean = (5*5 + 1*5) / 10 = (25+5)/10 = 3.0
func TestIPIP_ReverseScoring_Extraversion(t *testing.T) {
	scorer := NewIPIPScorer()
	m := make(map[string]interface{}, 50)
	for i := 1; i <= 50; i++ {
		m[fmt.Sprintf("%d", i)] = float64(1) // base 1
	}
	// Extra positive = 5
	for _, item := range []int{1, 11, 21, 31, 41} {
		m[fmt.Sprintf("%d", item)] = float64(5)
	}
	// Extra negative = 5 → reversed = 1
	for _, item := range []int{6, 16, 26, 36, 46} {
		m[fmt.Sprintf("%d", item)] = float64(5)
	}
	res, _ := scorer.Calculate(m)
	// pos: 5*5=25, neg: 6-5=1 → 1*5=5, total=30, mean=3.0
	assertFloat(t, "extraversion", 3.0, res.Extraversion)
}

// All answers = 1 → positive items score 1, negative items reversed = 5.
// Extraversion: (1*5 + 5*5) / 10 = 30/10 = 3.0 (same as all-5 case by symmetry)
func TestIPIP_AllOne_Symmetric(t *testing.T) {
	scorer := NewIPIPScorer()
	m := make(map[string]interface{}, 50)
	for i := 1; i <= 50; i++ {
		m[fmt.Sprintf("%d", i)] = float64(1)
	}
	res, _ := scorer.Calculate(m)
	assertFloat(t, "extraversion (all 1)", 3.0, res.Extraversion)
}

// All positive items = 5, all negative = 1 → max Extraversion = 5.0
func TestIPIP_MaxExtraversion(t *testing.T) {
	scorer := NewIPIPScorer()
	m := make(map[string]interface{}, 50)
	for i := 1; i <= 50; i++ {
		m[fmt.Sprintf("%d", i)] = float64(1) // base
	}
	for _, item := range []int{1, 11, 21, 31, 41} {
		m[fmt.Sprintf("%d", item)] = float64(5)
	}
	for _, item := range []int{6, 16, 26, 36, 46} {
		m[fmt.Sprintf("%d", item)] = float64(1) // reversed = 5
	}
	res, _ := scorer.Calculate(m)
	assertFloat(t, "extraversion max", 5.0, res.Extraversion)
	assertEqualStr(t, "extra_label", "sangat_tinggi", res.ExtraLabel)
}

// ── Interpretation labels ─────────────────────────────────────────────────────

func TestIPIP_Labels(t *testing.T) {
	cases := []struct {
		mean  float64
		label string
	}{
		{4.5, "sangat_tinggi"},
		{4.0, "tinggi"},
		{3.5, "tinggi"},
		{3.0, "rata_rata"},
		{2.5, "rata_rata"},
		{2.0, "rendah"},
		{1.5, "rendah"},
		{1.0, "sangat_rendah"},
	}
	for _, tc := range cases {
		got := interpretMean(tc.mean)
		if got != tc.label {
			t.Errorf("interpretMean(%.1f): want %q, got %q", tc.mean, tc.label, got)
		}
	}
}

// ── Sum preservation ──────────────────────────────────────────────────────────

func TestIPIP_SumField_AllFive(t *testing.T) {
	scorer := NewIPIPScorer()
	m := make(map[string]interface{}, 50)
	for i := 1; i <= 50; i++ {
		m[fmt.Sprintf("%d", i)] = float64(5)
	}
	res, _ := scorer.Calculate(m)
	// Extraversion: pos=5*5=25, neg=(6-5)*5=5 → sum=30
	if res.ExtraversionSum != 30 {
		t.Errorf("extraversion_sum: want 30, got %d", res.ExtraversionSum)
	}
}

// ── helper ────────────────────────────────────────────────────────────────────

func assertFloat(t *testing.T, field string, want, got float64) {
	t.Helper()
	if math.Abs(want-got) > 0.01 {
		t.Errorf("field %q: want %.2f, got %.2f", field, want, got)
	}
}

func assertEqualStr(t *testing.T, field, want, got string) {
	t.Helper()
	if want != got {
		t.Errorf("field %q: want %q, got %q", field, want, got)
	}
}
