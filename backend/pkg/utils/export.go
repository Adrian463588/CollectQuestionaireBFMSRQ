package utils

import (
	"encoding/csv"
	"fmt"
	"strings"

	"backend/internal/models"
)

// ExportToCSV exports participant data, responses, and scores to CSV format
func ExportToCSV(
	participant *models.Participant,
	score *models.Score,
) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Write header
	header := []string{
		"Participant ID",
		"Name",
		"Age",
		"Gender",
		"Created At",
	}

	// Add SRQ-29 headers if exists
	if score != nil && score.SRQScore != nil {
		header = append(header,
			"SRQ Neurotic Score",
			"SRQ Neurotic Status",
			"SRQ Substance Use",
			"SRQ Psychotic",
			"SRQ PTSD",
			"SRQ Classification",
		)
	}

	// Add IPIP-BFM-50 headers if exists
	if score != nil && score.IPIPScore != nil {
		header = append(header,
			"IPIP Extraversion",
			"IPIP Agreeableness",
			"IPIP Conscientiousness",
			"IPIP Emotional Stability",
			"IPIP Intellect",
			"IPIP Dominant Trait",
		)
	}

	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	record := []string{
		participant.ID,
		participant.Name,
		fmt.Sprintf("%d", participant.Age),
		participant.Gender,
		participant.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if score != nil && score.SRQScore != nil {
		substanceUse := "No"
		if score.SRQScore.SubstanceUse {
			substanceUse = "Yes"
		}
		psychotic := "No"
		if score.SRQScore.Psychotic {
			psychotic = "Yes"
		}
		ptsd := "No"
		if score.SRQScore.PTSD {
			ptsd = "Yes"
		}

		record = append(record,
			fmt.Sprintf("%d", score.SRQScore.NeuroticScore),
			score.SRQScore.NeuroticStatus,
			substanceUse,
			psychotic,
			ptsd,
			getSRQClass(score.SRQScore),
		)
	}

	if score != nil && score.IPIPScore != nil {
		record = append(record,
			fmt.Sprintf("%d", score.IPIPScore.Extraversion),
			fmt.Sprintf("%d", score.IPIPScore.Agreeableness),
			fmt.Sprintf("%d", score.IPIPScore.Conscientiousness),
			fmt.Sprintf("%d", score.IPIPScore.EmotionalStability),
			fmt.Sprintf("%d", score.IPIPScore.Intellect),
			getIPIPClass(score.IPIPScore),
		)
	}

	if err := writer.Write(record); err != nil {
		return "", fmt.Errorf("failed to write CSV record: %w", err)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return builder.String(), nil
}

// ExportMultipleToCSV exports multiple participants data to CSV
func ExportMultipleToCSV(
	participants []*models.Participant,
	scores []*models.Score,
) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Write header
	header := []string{
		"Participant ID",
		"Name",
		"Age",
		"Gender",
		"Created At",
		"SRQ Neurotic Score",
		"SRQ Neurotic Status",
		"SRQ Substance Use",
		"SRQ Psychotic",
		"SRQ PTSD",
		"SRQ Classification",
		"IPIP Extraversion",
		"IPIP Agreeableness",
		"IPIP Conscientiousness",
		"IPIP Emotional Stability",
		"IPIP Intellect",
		"IPIP Dominant Trait",
	}

	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write records
	for i, participant := range participants {
		record := []string{
			participant.ID,
			participant.Name,
			fmt.Sprintf("%d", participant.Age),
			participant.Gender,
			participant.CreatedAt.Format("2006-01-02 15:04:05"),
			"", "", "No", "No", "No", "-",
			"0", "0", "0", "0", "0", "-",
		}

		if i < len(scores) && scores[i] != nil {
			score := scores[i]
			if score.SRQScore != nil {
				substanceUse := "No"
				if score.SRQScore.SubstanceUse {
					substanceUse = "Yes"
				}
				psychotic := "No"
				if score.SRQScore.Psychotic {
					psychotic = "Yes"
				}
				ptsd := "No"
				if score.SRQScore.PTSD {
					ptsd = "Yes"
				}

				record[5] = fmt.Sprintf("%d", score.SRQScore.NeuroticScore)
				record[6] = score.SRQScore.NeuroticStatus
				record[7] = substanceUse
				record[8] = psychotic
				record[9] = ptsd
				record[10] = getSRQClass(score.SRQScore)
			}

			if score.IPIPScore != nil {
				record[11] = fmt.Sprintf("%d", score.IPIPScore.Extraversion)
				record[12] = fmt.Sprintf("%d", score.IPIPScore.Agreeableness)
				record[13] = fmt.Sprintf("%d", score.IPIPScore.Conscientiousness)
				record[14] = fmt.Sprintf("%d", score.IPIPScore.EmotionalStability)
				record[15] = fmt.Sprintf("%d", score.IPIPScore.Intellect)
				record[16] = getIPIPClass(score.IPIPScore)
			}
		}

		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return builder.String(), nil
}

func getSRQClass(s *models.SRQScore) string {
	if s == nil {
		return "-"
	}
	var flags []string
	if s.NeuroticScore >= 5 {
		flags = append(flags, "Indikasi GME")
	}
	if s.SubstanceUse {
		flags = append(flags, "Penggunaan Zat")
	}
	if s.Psychotic {
		flags = append(flags, "Gejala Psikotik")
	}
	if s.PTSD {
		flags = append(flags, "Gejala PTSD")
	}
	if len(flags) > 0 {
		return strings.Join(flags, " | ")
	}
	return "Normal"
}

func getIPIPClass(s *models.IPIPScore) string {
	if s == nil {
		return "-"
	}
	maxVal := s.Extraversion
	dominant := "Extraversion"
	
	if s.Agreeableness > maxVal {
		maxVal = s.Agreeableness
		dominant = "Agreeableness"
	}
	if s.Conscientiousness > maxVal {
		maxVal = s.Conscientiousness
		dominant = "Conscientiousness"
	}
	if s.EmotionalStability > maxVal {
		maxVal = s.EmotionalStability
		dominant = "Emotional Stability"
	}
	if s.Intellect > maxVal {
		maxVal = s.Intellect
		dominant = "Intellect"
	}
	return dominant
}
