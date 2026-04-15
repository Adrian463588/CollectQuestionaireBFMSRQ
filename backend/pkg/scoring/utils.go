package scoring

import "fmt"

// SafeBool safely extracts a boolean answer from the answers map.
// Returns false if the key is absent or the value is not a bool.
func SafeBool(answers map[string]interface{}, key string) bool {
	val, ok := answers[key]
	if !ok {
		return false
	}
	b, ok := val.(bool)
	return ok && b
}

// SafeFloat safely extracts a numeric answer as float64 from the answers map.
// Handles both float64 (JSON default) and int types.
// Returns 0 if the key is absent or the value cannot be coerced.
func SafeFloat(answers map[string]interface{}, key string) float64 {
	val, ok := answers[key]
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}
	return 0
}

// ReverseScore computes the reversed score for a negatively-keyed item.
// Formula: (max + 1) - val. For a 1-5 Likert scale, max = 5 → 6 - val.
func ReverseScore(val, max int) int {
	return (max + 1) - val
}

// SumBoolRange counts how many questions in the inclusive range [from, to]
// have a "true" (YA) answer. Keys are string representations of integers.
func SumBoolRange(answers map[string]interface{}, from, to int) int {
	count := 0
	for i := from; i <= to; i++ {
		if SafeBool(answers, fmt.Sprintf("%d", i)) {
			count++
		}
	}
	return count
}

// AnyBoolInRange returns true if at least one question in [from, to] is answered YA.
func AnyBoolInRange(answers map[string]interface{}, from, to int) bool {
	for i := from; i <= to; i++ {
		if SafeBool(answers, fmt.Sprintf("%d", i)) {
			return true
		}
	}
	return false
}

// CountBoolInRange counts exact number of YA answers in [from, to].
func CountBoolInRange(answers map[string]interface{}, from, to int) int {
	return SumBoolRange(answers, from, to)
}

// SumDimension computes the raw integer sum for a set of positively-keyed
// and negatively-keyed item numbers. Negative items are reverse-scored.
// All raw answers are on a 1-5 Likert scale.
func SumDimension(answers map[string]interface{}, positiveItems, negativeItems []int) int {
	total := 0
	for _, item := range positiveItems {
		v := int(SafeFloat(answers, fmt.Sprintf("%d", item)))
		if v >= 1 && v <= 5 {
			total += v
		}
	}
	for _, item := range negativeItems {
		v := int(SafeFloat(answers, fmt.Sprintf("%d", item)))
		if v >= 1 && v <= 5 {
			total += ReverseScore(v, 5)
		}
	}
	return total
}

// MeanDimension computes the mean score for a dimension.
// Returns 0 if no valid items are found.
func MeanDimension(answers map[string]interface{}, positiveItems, negativeItems []int) float64 {
	sum := float64(SumDimension(answers, positiveItems, negativeItems))
	count := len(positiveItems) + len(negativeItems)
	if count == 0 {
		return 0
	}
	// Round to 2 decimal places for clean output
	return roundTo2(sum / float64(count))
}

// roundTo2 rounds a float64 to 2 decimal places.
func roundTo2(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}
