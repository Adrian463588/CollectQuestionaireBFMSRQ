package services

import (
	"fmt"

	"backend/internal/models"
)

// ScoringService handles scoring logic for questionnaires
type ScoringService struct{}

// NewScoringService creates a new scoring service
func NewScoringService() *ScoringService {
	return &ScoringService{}
}

// getBoolAnswer safely extracts a boolean answer from the answers map.
func getBoolAnswer(answers map[string]interface{}, key string) bool {
	val, ok := answers[key]
	if !ok {
		return false
	}
	boolVal, ok := val.(bool)
	return ok && boolVal
}

// sumBoolAnswers counts how many questions in [from, to] have a true answer.
func sumBoolAnswers(answers map[string]interface{}, from, to int) int {
	count := 0
	for i := from; i <= to; i++ {
		if getBoolAnswer(answers, fmt.Sprintf("%d", i)) {
			count++
		}
	}
	return count
}

// anyBoolTrue returns true if any question in [from, to] has a true answer.
func anyBoolTrue(answers map[string]interface{}, from, to int) bool {
	for i := from; i <= to; i++ {
		if getBoolAnswer(answers, fmt.Sprintf("%d", i)) {
			return true
		}
	}
	return false
}

// neuroticStatus interprets a neurotic score per SRQ-29 rules from AGENTS.md.
func neuroticStatus(score int) string {
	switch {
	case score >= 6:
		return "rekomendasi_rujukan"
	case score >= 5:
		return "indikasi_masalah_emosional"
	default:
		return "normal"
	}
}

// CalculateSRQScore calculates SRQ-29 scores based on answers.
// Answers must be keyed by question number ("1" … "29") with boolean values.
func (s *ScoringService) CalculateSRQScore(answers map[string]interface{}) (*models.SRQScore, error) {
	neuro := sumBoolAnswers(answers, 1, 20)

	score := &models.SRQScore{
		NeuroticScore:  neuro,
		NeuroticStatus: neuroticStatus(neuro),
		SubstanceUse:   getBoolAnswer(answers, "21"),
		Psychotic:      anyBoolTrue(answers, 22, 24),
		PTSD:           anyBoolTrue(answers, 25, 29),
	}

	return score, nil
}

// CalculateIPIPScore calculates IPIP-BFM-50 scores based on answers.
// Answers must be keyed by question number ("1" … "50") with float64/int values (1-5 Likert).
func (s *ScoringService) CalculateIPIPScore(answers map[string]interface{}) (*models.IPIPScore, error) {
	score := &models.IPIPScore{
		// Extraversion: + items (1,11,21,31,41), - items (6,16,26,36,46)
		Extraversion: s.calculateDimension(answers,
			[]int{1, 11, 21, 31, 41},
			[]int{6, 16, 26, 36, 46},
		),
		// Agreeableness: + items (7,17,27,37,42,47), - items (2,12,22,32)
		Agreeableness: s.calculateDimension(answers,
			[]int{7, 17, 27, 37, 42, 47},
			[]int{2, 12, 22, 32},
		),
		// Conscientiousness: + items (3,13,23,33,43,48), - items (8,18,28,38)
		Conscientiousness: s.calculateDimension(answers,
			[]int{3, 13, 23, 33, 43, 48},
			[]int{8, 18, 28, 38},
		),
		// Emotional Stability: + items (9,19), - items (4,14,24,29,34,39,44,49)
		EmotionalStability: s.calculateDimension(answers,
			[]int{9, 19},
			[]int{4, 14, 24, 29, 34, 39, 44, 49},
		),
		// Intellect: + items (5,15,25,35,40,45,50), - items (10,20,30)
		Intellect: s.calculateDimension(answers,
			[]int{5, 15, 25, 35, 40, 45, 50},
			[]int{10, 20, 30},
		),
	}

	return score, nil
}

// calculateDimension sums scores for a single IPIP dimension.
// Positive items use direct scoring (1-5); negative items use reverse scoring (6 - value).
func (s *ScoringService) calculateDimension(answers map[string]interface{}, positiveItems, negativeItems []int) int {
	total := 0

	for _, item := range positiveItems {
		if val, ok := answers[fmt.Sprintf("%d", item)]; ok {
			if v, ok := val.(float64); ok {
				total += int(v)
			}
		}
	}

	for _, item := range negativeItems {
		if val, ok := answers[fmt.Sprintf("%d", item)]; ok {
			if v, ok := val.(float64); ok {
				total += 6 - int(v)
			}
		}
	}

	return total
}
