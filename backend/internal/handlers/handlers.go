package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"backend/internal/models"
	"backend/internal/services"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Handler holds the dependencies for HTTP handlers
type Handler struct {
	ParticipantService *services.ParticipantService
	ResponseService    *services.ResponseService
	ScoreService       *services.ScoreService
	ScoringService     *services.ScoringService
}

// NewHandler creates a new handler
func NewHandler(
	participantSvc *services.ParticipantService,
	responseSvc *services.ResponseService,
	scoreSvc *services.ScoreService,
	scoringSvc *services.ScoringService,
) *Handler {
	return &Handler{
		ParticipantService: participantSvc,
		ResponseService:    responseSvc,
		ScoreService:       scoreSvc,
		ScoringService:     scoringSvc,
	}
}

// CreateParticipant handles POST /api/participants
func (h *Handler) CreateParticipant(c *fiber.Ctx) error {
	var req models.ParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}
	if req.Age < 15 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Age must be at least 15",
		})
	}
	if req.Gender != "male" && req.Gender != "female" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Gender must be 'male' or 'female'",
		})
	}

	participant, err := h.ParticipantService.CreateParticipant(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create participant",
		})
	}

	return c.Status(http.StatusCreated).JSON(participant)
}

// GetParticipant handles GET /api/participants/:id
func (h *Handler) GetParticipant(c *fiber.Ctx) error {
	id := c.Params("id")

	participant, err := h.ParticipantService.GetParticipant(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Participant not found",
		})
	}

	return c.JSON(participant)
}

// SubmitResponse handles POST /api/responses
func (h *Handler) SubmitResponse(c *fiber.Ctx) error {
	var req models.SubmissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ParticipantID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Participant ID is required",
		})
	}
	if req.QuestionnaireType == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Questionnaire type is required",
		})
	}
	if len(req.Answers) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Answers are required",
		})
	}

	questionnaireType := strings.ToLower(strings.TrimSpace(req.QuestionnaireType))
	if questionnaireType != "srq29" && questionnaireType != "ipip-bfm-50" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Questionnaire type must be 'srq29' or 'ipip-bfm-50'",
		})
	}

	// Save response
	response, err := h.ResponseService.SaveResponse(req.ParticipantID, questionnaireType, req.Answers)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save response",
		})
	}

	// Calculate scores
	var srqScore *models.SRQScore
	var ipipScore *models.IPIPScore

	if questionnaireType == "srq29" {
		srqScore, err = h.ScoringService.CalculateSRQScore(req.Answers)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to calculate SRQ score",
			})
		}
	} else {
		ipipScore, err = h.ScoringService.CalculateIPIPScore(req.Answers)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to calculate IPIP score",
			})
		}
	}

	// Save score
	_, err = h.ScoreService.SaveScore(req.ParticipantID, srqScore, ipipScore)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save score",
		})
	}

	result := fiber.Map{
		"response_id":    response.ID,
		"participant_id": req.ParticipantID,
	}

	if srqScore != nil {
		result["srq29"] = srqScore
	}
	if ipipScore != nil {
		result["ipip"] = ipipScore
	}

	return c.Status(http.StatusCreated).JSON(result)
}

// CalculateScore handles POST /api/scoring
func (h *Handler) CalculateScore(c *fiber.Ctx) error {
	var req struct {
		QuestionnaireType string                 `json:"questionnaire_type"`
		Answers           map[string]interface{} `json:"answers"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.QuestionnaireType == "" || req.Answers == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Questionnaire type and answers are required",
		})
	}

	var result interface{}
	var err error

	switch req.QuestionnaireType {
	case "srq29":
		result, err = h.ScoringService.CalculateSRQScore(req.Answers)
	case "ipip-bfm-50":
		result, err = h.ScoringService.CalculateIPIPScore(req.Answers)
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid questionnaire type",
		})
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate score",
		})
	}

	return c.JSON(result)
}

// GetScores handles GET /api/scores/:participantId
func (h *Handler) GetScores(c *fiber.Ctx) error {
	participantID := c.Params("participantId")

	score, err := h.ScoreService.GetScores(participantID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Scores not found",
		})
	}

	return c.JSON(score)
}

// ExportCSV handles GET /api/export/:participantId
func (h *Handler) ExportCSV(c *fiber.Ctx) error {
	participantID := c.Params("participantId")

	// Get participant data
	participant, err := h.ParticipantService.GetParticipant(participantID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Participant not found",
		})
	}

	// Get scores
	score, err := h.ScoreService.GetScores(participantID)
	if err != nil {
		// Continue without scores if not found
		score = nil
	}

	// Export to CSV
	csvData, err := utils.ExportToCSV(participant, score)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to export data",
		})
	}

	// Set headers for file download
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=participant_"+participantID+".csv")
	c.Set("Content-Length", fmt.Sprintf("%d", len(csvData)))

	return c.SendString(csvData)
}

// ExportAllCSV handles GET /api/export
func (h *Handler) ExportAllCSV(c *fiber.Ctx) error {
	participants, err := h.ParticipantService.GetAllParticipants()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get participants",
		})
	}

	scoresMap, err := h.ScoreService.GetAllScores()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get scores",
		})
	}

	scoresList := make([]*models.Score, len(participants))
	for i, p := range participants {
		if score, ok := scoresMap[p.ID]; ok {
			scoresList[i] = score
		}
	}

	csvData, err := utils.ExportMultipleToCSV(participants, scoresList)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to export data",
		})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="all_participants_results.csv"`)
	c.Set("Content-Length", fmt.Sprintf("%d", len(csvData)))

	return c.SendString(csvData)
}

// GetDashboardData handles GET /api/dashboard
func (h *Handler) GetDashboardData(c *fiber.Ctx) error {
	participants, err := h.ParticipantService.GetAllParticipants()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get participants",
		})
	}

	scoresMap, err := h.ScoreService.GetAllScores()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get scores",
		})
	}

	type DashboardItem struct {
		Participant *models.Participant `json:"participant"`
		Score       *models.Score       `json:"score"`
	}

	var data []DashboardItem
	for _, p := range participants {
		item := DashboardItem{
			Participant: p,
			Score:       scoresMap[p.ID],
		}
		data = append(data, item)
	}

	return c.JSON(data)
}
