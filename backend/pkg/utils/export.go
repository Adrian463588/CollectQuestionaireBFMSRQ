// Package utils provides shared utilities for data export.
package utils

import (
	"encoding/csv"
	"fmt"
	"strings"

	"backend/internal/models"
)

// csvHeaders defines every column in the exported CSV.
// This single source of truth ensures header and row are always in sync (DRY).
var csvHeaders = []string{
	// ── Participant ─────────────────────────────────────────────────────────
	"Participant ID",
	"Name",
	"Age",
	"Gender",
	"Created At",
	// ── SRQ-29: domain scores ───────────────────────────────────────────────
	"SRQ Neurotic Score (Q1-20)",
	"SRQ Neurotic Status",
	"SRQ Substance Use (Q21)",
	"SRQ Psychotic (Q22-24)",
	"SRQ Psychotic Count",
	"SRQ PTSD (Q25-29)",
	"SRQ PTSD Count",
	"SRQ Total Score (Q1-29)",
	"SRQ Overall Risk",
	// ── SRQ-29: dummy variables (0/1) for research ─────────────────────────
	"SRQ GME_Dummy",
	"SRQ Substance_Dummy",
	"SRQ Psychotic_Dummy",
	"SRQ PTSD_Dummy",
	// ── IPIP-BFM-50: mean scores (1.0–5.0) ─────────────────────────────────
	"IPIP Extraversion Mean",
	"IPIP Agreeableness Mean",
	"IPIP Conscientiousness Mean",
	"IPIP Emotional Stability Mean",
	"IPIP Intellect Mean",
	// ── IPIP-BFM-50: interpretation labels ─────────────────────────────────
	"IPIP Extraversion Label",
	"IPIP Agreeableness Label",
	"IPIP Conscientiousness Label",
	"IPIP Emotional Stability Label",
	"IPIP Intellect Label",
	// ── IPIP-BFM-50: raw sum scores (10–50) ────────────────────────────────
	"IPIP Extraversion Sum",
	"IPIP Agreeableness Sum",
	"IPIP Conscientiousness Sum",
	"IPIP Emotional Stability Sum",
	"IPIP Intellect Sum",
}

// ExportToCSV exports a single participant's data and scores to a CSV string.
func ExportToCSV(participant *models.Participant, score *models.Score) (string, error) {
	var b strings.Builder
	w := csv.NewWriter(&b)

	if err := w.Write(csvHeaders); err != nil {
		return "", fmt.Errorf("write CSV header: %w", err)
	}
	if err := w.Write(buildRecord(participant, score)); err != nil {
		return "", fmt.Errorf("write CSV record: %w", err)
	}

	w.Flush()
	return b.String(), w.Error()
}

// ExportMultipleToCSV exports all participants to a single CSV string.
// scores[i] corresponds to participants[i]; a nil score means no data yet.
func ExportMultipleToCSV(participants []*models.Participant, scores []*models.Score) (string, error) {
	var b strings.Builder
	w := csv.NewWriter(&b)

	if err := w.Write(csvHeaders); err != nil {
		return "", fmt.Errorf("write CSV header: %w", err)
	}

	for i, p := range participants {
		var score *models.Score
		if i < len(scores) {
			score = scores[i]
		}
		if err := w.Write(buildRecord(p, score)); err != nil {
			return "", fmt.Errorf("write CSV record for participant %s: %w", p.ID, err)
		}
	}

	w.Flush()
	return b.String(), w.Error()
}

// buildRecord constructs a CSV row matching the order of csvHeaders.
// Empty strings are used for absent data — no dummy/placeholder values.
func buildRecord(p *models.Participant, score *models.Score) []string {
	row := make([]string, len(csvHeaders))

	// ── Participant ─────────────────────────────────────────────────────────
	row[0] = p.ID
	row[1] = p.Name
	row[2] = fmt.Sprintf("%d", p.Age)
	row[3] = p.Gender
	row[4] = p.CreatedAt.Format("2006-01-02 15:04:05")

	// ── SRQ-29 ─────────────────────────────────────────────────────────────
	if score != nil && score.SRQScore != nil {
		s := score.SRQScore
		row[5] = fmt.Sprintf("%d", s.NeuroticScore)
		row[6] = s.NeuroticStatus
		row[7] = boolStr(s.SubstanceUse)
		row[8] = boolStr(s.Psychotic)
		row[9] = fmt.Sprintf("%d", s.PsychoticCount)
		row[10] = boolStr(s.PTSD)
		row[11] = fmt.Sprintf("%d", s.PTSDCount)
		row[12] = fmt.Sprintf("%d", s.TotalScore)
		row[13] = s.OverallRisk
		row[14] = fmt.Sprintf("%d", s.EmotionalDisorder)
		row[15] = fmt.Sprintf("%d", s.SubstanceDummy)
		row[16] = fmt.Sprintf("%d", s.PsychoticDummy)
		row[17] = fmt.Sprintf("%d", s.PTSDDummy)
	}

	// ── IPIP-BFM-50 ─────────────────────────────────────────────────────────
	if score != nil && score.IPIPScore != nil {
		ip := score.IPIPScore
		row[18] = fmt.Sprintf("%.2f", ip.Extraversion)
		row[19] = fmt.Sprintf("%.2f", ip.Agreeableness)
		row[20] = fmt.Sprintf("%.2f", ip.Conscientiousness)
		row[21] = fmt.Sprintf("%.2f", ip.EmotionalStability)
		row[22] = fmt.Sprintf("%.2f", ip.Intellect)
		row[23] = ip.ExtraLabel
		row[24] = ip.AgreLabel
		row[25] = ip.ConsLabel
		row[26] = ip.StabLabel
		row[27] = ip.IntellLabel
		row[28] = fmt.Sprintf("%d", ip.ExtraversionSum)
		row[29] = fmt.Sprintf("%d", ip.AgreeablenessSum)
		row[30] = fmt.Sprintf("%d", ip.ConscientiousnessSum)
		row[31] = fmt.Sprintf("%d", ip.EmotionalStabilitySum)
		row[32] = fmt.Sprintf("%d", ip.IntellectSum)
	}

	return row
}

// boolStr converts a bool to "Ya"/"Tidak" for the Indonesian CSV output.
func boolStr(b bool) string {
	if b {
		return "Ya"
	}
	return "Tidak"
}
