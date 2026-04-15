// Package scoring provides psychometric scoring logic for SRQ-29 and IPIP-BFM-50.
// Each instrument is implemented as a separate type satisfying the Scorer interface,
// following SOLID principles (Single Responsibility + Open/Closed).
package scoring

// Scorer is a generic interface for any questionnaire scoring engine.
// To add a new instrument, implement this interface — no existing code needs to change.
type Scorer[T any] interface {
	Calculate(answers map[string]interface{}) (T, error)
}
