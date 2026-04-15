package scoring

import "backend/internal/models"

// ipipDimension defines the item mapping for a single Big Five dimension.
// Item numbers correspond to official IPIP-BFM-50 scoring key
// (Akhtar & Azwar, 2019; ipip.ori.org/new_ipip-50-item-scale.htm).
type ipipDimension struct {
	Name          string
	PositiveItems []int
	NegativeItems []int
}

// ipipDimensions holds the official item-to-dimension mapping for all five traits.
// Verified against IPIP-BFM-50.doc.md (Indonesian adaptation, UGM 2019).
var ipipDimensions = []ipipDimension{
	{
		Name:          "Extraversion",
		PositiveItems: []int{1, 11, 21, 31, 41},
		NegativeItems: []int{6, 16, 26, 36, 46},
	},
	{
		Name:          "Agreeableness",
		PositiveItems: []int{7, 17, 27, 37, 42, 47},
		NegativeItems: []int{2, 12, 22, 32},
	},
	{
		Name:          "Conscientiousness",
		PositiveItems: []int{3, 13, 23, 33, 43, 48},
		NegativeItems: []int{8, 18, 28, 38},
	},
	{
		Name:          "EmotionalStability",
		PositiveItems: []int{9, 19},
		NegativeItems: []int{4, 14, 24, 29, 34, 39, 44, 49},
	},
	{
		Name:          "Intellect",
		PositiveItems: []int{5, 15, 25, 35, 40, 45, 50},
		NegativeItems: []int{10, 20, 30},
	},
}

// IPIPScorer implements Scorer[*models.IPIPScore] for the IPIP-BFM-50 instrument.
// Scoring method: MEAN of 10 reverse-scored items per dimension (range 1.0–5.0).
// This is the recommended method for psychometric research publications (Akhtar & Azwar, 2019).
type IPIPScorer struct{}

// NewIPIPScorer creates a new IPIPScorer.
func NewIPIPScorer() *IPIPScorer {
	return &IPIPScorer{}
}

// Calculate computes IPIP-BFM-50 scores from a map of Likert answers (1–5)
// keyed by question number string (e.g. "1"…"50").
//
// Steps:
//  1. Apply reverse scoring to negatively-keyed items (6 - original)
//  2. Compute mean per dimension (sum / 10 items)
//  3. Attach interpretation label based on mean value
//  4. Preserve raw sum for optional alternative analysis
func (s *IPIPScorer) Calculate(answers map[string]interface{}) (*models.IPIPScore, error) {
	score := &models.IPIPScore{}

	for _, dim := range ipipDimensions {
		mean := MeanDimension(answers, dim.PositiveItems, dim.NegativeItems)
		sum := SumDimension(answers, dim.PositiveItems, dim.NegativeItems)
		label := interpretMean(mean)

		switch dim.Name {
		case "Extraversion":
			score.Extraversion = mean
			score.ExtraversionSum = sum
			score.ExtraLabel = label
		case "Agreeableness":
			score.Agreeableness = mean
			score.AgreeablenessSum = sum
			score.AgreLabel = label
		case "Conscientiousness":
			score.Conscientiousness = mean
			score.ConscientiousnessSum = sum
			score.ConsLabel = label
		case "EmotionalStability":
			score.EmotionalStability = mean
			score.EmotionalStabilitySum = sum
			score.StabLabel = label
		case "Intellect":
			score.Intellect = mean
			score.IntellectSum = sum
			score.IntellLabel = label
		}
	}

	return score, nil
}

// interpretMean converts a mean score (1.0–5.0) into a descriptive label.
// Thresholds are based on the 5-point Likert scale midpoints and standard
// psychometric practice for absolute interpretation (not norm-referenced).
//
// For norm-referenced (Z-score) interpretation, the researcher should
// standardize scores against the collected sample distribution using SPSS/R.
//
//	Mean ≥ 4.5 → "sangat_tinggi"
//	Mean ≥ 3.5 → "tinggi"
//	Mean ≥ 2.5 → "rata_rata"
//	Mean ≥ 1.5 → "rendah"
//	Mean  < 1.5 → "sangat_rendah"
func interpretMean(mean float64) string {
	switch {
	case mean >= 4.5:
		return "sangat_tinggi"
	case mean >= 3.5:
		return "tinggi"
	case mean >= 2.5:
		return "rata_rata"
	case mean >= 1.5:
		return "rendah"
	default:
		return "sangat_rendah"
	}
}
