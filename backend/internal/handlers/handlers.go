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

// Handler holds the application-level dependencies for all HTTP handlers.
type Handler struct {
	ParticipantService *services.ParticipantService
	ResponseService    *services.ResponseService
	ScoreService       *services.ScoreService
	ScoringService     *services.ScoringService
}

// NewHandler wires up the handler with all required services.
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

// ── Participant handlers ───────────────────────────────────────────────────────

// CreateParticipant handles POST /api/participants
func (h *Handler) CreateParticipant(c *fiber.Ctx) error {
	var req models.ParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}
	if req.Age < 15 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Age must be at least 15"})
	}
	if req.Gender != "male" && req.Gender != "female" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Gender must be 'male' or 'female'"})
	}

	participant, err := h.ParticipantService.CreateParticipant(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create participant"})
	}
	return c.Status(http.StatusCreated).JSON(participant)
}

// GetParticipant handles GET /api/participants/:id
func (h *Handler) GetParticipant(c *fiber.Ctx) error {
	id := c.Params("id")
	participant, err := h.ParticipantService.GetParticipant(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Participant not found"})
	}
	return c.JSON(participant)
}

// DeleteParticipant handles DELETE /api/participants/:id
//
// Performs a soft-delete cascade:
//  1. Validates that the participant exists.
//  2. Soft-deletes scores, responses, and the participant in one DB transaction.
//  3. Returns 200 with a confirmation message.
//
// Deleted data is retained in the database for research audit purposes;
// it is simply excluded from all active queries via `WHERE deleted_at IS NULL`.
func (h *Handler) DeleteParticipant(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Participant ID is required"})
	}

	if err := h.ParticipantService.DeleteParticipant(id); err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "already deleted") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Participant not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete participant"})
	}

	return c.JSON(fiber.Map{"message": "Participant deleted successfully"})
}

// ── Response / Scoring handlers ───────────────────────────────────────────────

// SubmitResponse handles POST /api/responses
func (h *Handler) SubmitResponse(c *fiber.Ctx) error {
	var req models.SubmissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.ParticipantID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Participant ID is required"})
	}
	if req.QuestionnaireType == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Questionnaire type is required"})
	}
	if len(req.Answers) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Answers are required"})
	}

	qType := strings.ToLower(strings.TrimSpace(req.QuestionnaireType))
	if qType != "srq29" && qType != "ipip-bfm-50" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Questionnaire type must be 'srq29' or 'ipip-bfm-50'",
		})
	}

	response, err := h.ResponseService.SaveResponse(req.ParticipantID, qType, req.Answers)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save response"})
	}

	var srqScore *models.SRQScore
	var ipipScore *models.IPIPScore

	switch qType {
	case "srq29":
		srqScore, err = h.ScoringService.CalculateSRQScore(req.Answers)
	case "ipip-bfm-50":
		ipipScore, err = h.ScoringService.CalculateIPIPScore(req.Answers)
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate score"})
	}

	if _, err = h.ScoreService.SaveScore(req.ParticipantID, srqScore, ipipScore); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save score"})
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire type"})
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate score"})
	}
	return c.JSON(result)
}

// ── Score handlers ────────────────────────────────────────────────────────────

// GetScores handles GET /api/scores/:participantId
func (h *Handler) GetScores(c *fiber.Ctx) error {
	participantID := c.Params("participantId")
	score, err := h.ScoreService.GetScores(participantID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Scores not found"})
	}
	return c.JSON(score)
}

// ── Export handlers ───────────────────────────────────────────────────────────

// ExportCSV handles GET /api/export/:participantId
func (h *Handler) ExportCSV(c *fiber.Ctx) error {
	participantID := c.Params("participantId")

	participant, err := h.ParticipantService.GetParticipant(participantID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Participant not found"})
	}

	score, err := h.ScoreService.GetScores(participantID)
	if err != nil {
		score = nil // export even if scores are missing
	}

	csvData, err := utils.ExportToCSV(participant, score)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export data"})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=participant_"+participantID+".csv")
	c.Set("Content-Length", fmt.Sprintf("%d", len(csvData)))
	return c.SendString(csvData)
}

// ExportAllCSV handles GET /api/export
func (h *Handler) ExportAllCSV(c *fiber.Ctx) error {
	participants, err := h.ParticipantService.GetAllParticipants()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get participants"})
	}

	scoresMap, err := h.ScoreService.GetAllScores()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get scores"})
	}

	scoresList := make([]*models.Score, len(participants))
	for i, p := range participants {
		if score, ok := scoresMap[p.ID]; ok {
			scoresList[i] = score
		}
	}

	csvData, err := utils.ExportMultipleToCSV(participants, scoresList)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export data"})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="all_participants_results.csv"`)
	c.Set("Content-Length", fmt.Sprintf("%d", len(csvData)))
	return c.SendString(csvData)
}

// ── Dashboard handler ─────────────────────────────────────────────────────────

// GetDashboardData handles GET /api/dashboard
func (h *Handler) GetDashboardData(c *fiber.Ctx) error {
	participants, err := h.ParticipantService.GetAllParticipants()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get participants"})
	}

	scoresMap, err := h.ScoreService.GetAllScores()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get scores"})
	}

	type DashboardItem struct {
		Participant *models.Participant `json:"participant"`
		Score       *models.Score       `json:"score"`
	}

	data := make([]DashboardItem, 0, len(participants))
	for _, p := range participants {
		data = append(data, DashboardItem{
			Participant: p,
			Score:       scoresMap[p.ID],
		})
	}
	return c.JSON(data)
}
