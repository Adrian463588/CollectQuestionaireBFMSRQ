package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/database"
	"backend/internal/models"
	"backend/pkg/scoring"

	"github.com/google/uuid"
)

// ── ParticipantService ────────────────────────────────────────────────────────

// ParticipantService handles all participant-related database operations.
type ParticipantService struct {
	DB *database.DB
}

// NewParticipantService creates a new ParticipantService.
func NewParticipantService(db *database.DB) *ParticipantService {
	return &ParticipantService{DB: db}
}

// CreateParticipant persists a new participant.
func (s *ParticipantService) CreateParticipant(req models.ParticipantRequest) (*models.Participant, error) {
	participant := &models.Participant{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Age:       req.Age,
		Gender:    req.Gender,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO participants (id, name, age, gender, created_at)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := s.DB.Connection.Exec(query,
		participant.ID, participant.Name, participant.Age,
		participant.Gender, participant.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create participant: %w", err)
	}
	return participant, nil
}

// GetParticipant returns a non-deleted participant by ID.
func (s *ParticipantService) GetParticipant(id string) (*models.Participant, error) {
	p := &models.Participant{}
	query := `SELECT id, name, age, gender, created_at
	          FROM participants
	          WHERE id = $1 AND deleted_at IS NULL`

	err := s.DB.Connection.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Age, &p.Gender, &p.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get participant: %w", err)
	}
	return p, nil
}

// GetAllParticipants returns all non-deleted participants ordered by creation date.
func (s *ParticipantService) GetAllParticipants() ([]*models.Participant, error) {
	query := `SELECT id, name, age, gender, created_at
	          FROM participants
	          WHERE deleted_at IS NULL
	          ORDER BY created_at DESC`

	rows, err := s.DB.Connection.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all participants: %w", err)
	}
	defer rows.Close()

	var participants []*models.Participant
	for rows.Next() {
		p := &models.Participant{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Gender, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan participant: %w", err)
		}
		participants = append(participants, p)
	}
	return participants, nil
}

// DeleteParticipant performs a cascading soft-delete for a participant and all
// associated responses and scores within a single database transaction.
// This preserves data integrity for research audit while removing it from
// active queries (which all filter WHERE deleted_at IS NULL).
func (s *ParticipantService) DeleteParticipant(id string) error {
	// Verify participant exists and is not already deleted.
	_, err := s.GetParticipant(id)
	if err != nil {
		return fmt.Errorf("participant not found or already deleted: %w", err)
	}

	tx, err := s.DB.Connection.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	now := time.Now()

	// 1. Soft-delete scores
	if _, err = tx.Exec(
		`UPDATE scores SET deleted_at = $1 WHERE participant_id = $2 AND deleted_at IS NULL`,
		now, id,
	); err != nil {
		return fmt.Errorf("soft-delete scores: %w", err)
	}

	// 2. Soft-delete responses
	if _, err = tx.Exec(
		`UPDATE responses SET deleted_at = $1 WHERE participant_id = $2 AND deleted_at IS NULL`,
		now, id,
	); err != nil {
		return fmt.Errorf("soft-delete responses: %w", err)
	}

	// 3. Soft-delete participant
	if _, err = tx.Exec(
		`UPDATE participants SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`,
		now, id,
	); err != nil {
		return fmt.Errorf("soft-delete participant: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// ── ResponseService ───────────────────────────────────────────────────────────

// ResponseService handles questionnaire response persistence.
type ResponseService struct {
	DB *database.DB
}

// NewResponseService creates a new ResponseService.
func NewResponseService(db *database.DB) *ResponseService {
	return &ResponseService{DB: db}
}

// SaveResponse persists questionnaire answers for a participant.
func (s *ResponseService) SaveResponse(participantID, questionnaireType string, answers map[string]interface{}) (*models.Response, error) {
	response := &models.Response{
		ID:                uuid.New().String(),
		ParticipantID:     participantID,
		QuestionnaireType: questionnaireType,
		Answers:           answers,
		CreatedAt:         time.Now(),
	}

	answersJSON, err := json.Marshal(answers)
	if err != nil {
		return nil, fmt.Errorf("marshal answers: %w", err)
	}

	query := `INSERT INTO responses (id, participant_id, questionnaire_type, answers, created_at)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err = s.DB.Connection.Exec(query,
		response.ID, response.ParticipantID, response.QuestionnaireType,
		answersJSON, response.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("save response: %w", err)
	}
	return response, nil
}

// ── ScoreService ──────────────────────────────────────────────────────────────

// ScoreService handles score persistence and retrieval.
type ScoreService struct {
	DB             *database.DB
	ScoringService *ScoringService
}

// NewScoreService creates a new ScoreService.
func NewScoreService(db *database.DB) *ScoreService {
	return &ScoreService{
		DB:             db,
		ScoringService: NewScoringService(),
	}
}

// SaveScore upserts the score record for a participant.
// SRQ and IPIP scores are merged into the same row so a participant
// always has a single consolidated score record.
func (s *ScoreService) SaveScore(participantID string, srqScore *models.SRQScore, ipipScore *models.IPIPScore) (*models.Score, error) {
	score := &models.Score{
		ParticipantID: participantID,
		CreatedAt:     time.Now(),
	}

	var existingID string
	var existingSRQJSON, existingIPIPJSON []byte

	err := s.DB.Connection.QueryRow(
		`SELECT id, srq_score, ipip_score FROM scores
		 WHERE participant_id = $1 AND deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT 1`,
		participantID,
	).Scan(&existingID, &existingSRQJSON, &existingIPIPJSON)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("check existing score: %w", err)
	}

	if err == sql.ErrNoRows {
		// New record
		score.ID = uuid.New().String()
		score.SRQScore = srqScore
		score.IPIPScore = ipipScore

		srqJSON, ipipJSON, marshalErr := marshalScoreJSON(score.SRQScore, score.IPIPScore)
		if marshalErr != nil {
			return nil, marshalErr
		}

		_, execErr := s.DB.Connection.Exec(
			`INSERT INTO scores (id, participant_id, srq_score, ipip_score, created_at)
			 VALUES ($1, $2, $3, $4, $5)`,
			score.ID, score.ParticipantID, srqJSON, ipipJSON, score.CreatedAt,
		)
		if execErr != nil {
			return nil, fmt.Errorf("insert score: %w", execErr)
		}
		return score, nil
	}

	// Merge with existing record
	score.ID = existingID

	if existingSRQJSON != nil {
		existing := &models.SRQScore{}
		if err := json.Unmarshal(existingSRQJSON, existing); err != nil {
			return nil, fmt.Errorf("unmarshal existing SRQ score: %w", err)
		}
		score.SRQScore = existing
	}
	if existingIPIPJSON != nil {
		existing := &models.IPIPScore{}
		if err := json.Unmarshal(existingIPIPJSON, existing); err != nil {
			return nil, fmt.Errorf("unmarshal existing IPIP score: %w", err)
		}
		score.IPIPScore = existing
	}

	if srqScore != nil {
		score.SRQScore = srqScore
	}
	if ipipScore != nil {
		score.IPIPScore = ipipScore
	}

	srqJSON, ipipJSON, marshalErr := marshalScoreJSON(score.SRQScore, score.IPIPScore)
	if marshalErr != nil {
		return nil, marshalErr
	}

	_, execErr := s.DB.Connection.Exec(
		`UPDATE scores SET srq_score = $1, ipip_score = $2, created_at = $3 WHERE id = $4`,
		srqJSON, ipipJSON, score.CreatedAt, score.ID,
	)
	if execErr != nil {
		return nil, fmt.Errorf("update score: %w", execErr)
	}
	return score, nil
}

// GetScores retrieves the latest score record for a participant.
func (s *ScoreService) GetScores(participantID string) (*models.Score, error) {
	score := &models.Score{}
	var srqJSON, ipipJSON []byte

	query := `SELECT id, participant_id, srq_score, ipip_score, created_at
	          FROM scores
	          WHERE participant_id = $1 AND deleted_at IS NULL
	          ORDER BY created_at DESC LIMIT 1`

	err := s.DB.Connection.QueryRow(query, participantID).Scan(
		&score.ID, &score.ParticipantID, &srqJSON, &ipipJSON, &score.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get scores: %w", err)
	}

	if err := unmarshalScores(srqJSON, ipipJSON, score); err != nil {
		return nil, err
	}
	return score, nil
}

// GetAllScores returns the latest score for every active participant, keyed by participant ID.
func (s *ScoreService) GetAllScores() (map[string]*models.Score, error) {
	query := `
		SELECT t.id, t.participant_id, t.srq_score, t.ipip_score, t.created_at
		FROM scores t
		INNER JOIN (
			SELECT participant_id, MAX(created_at) AS max_date
			FROM scores
			WHERE deleted_at IS NULL
			GROUP BY participant_id
		) latest ON t.participant_id = latest.participant_id AND t.created_at = latest.max_date
		WHERE t.deleted_at IS NULL
	`

	rows, err := s.DB.Connection.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all scores: %w", err)
	}
	defer rows.Close()

	scoresMap := make(map[string]*models.Score)
	for rows.Next() {
		score := &models.Score{}
		var srqJSON, ipipJSON []byte

		if err := rows.Scan(
			&score.ID, &score.ParticipantID, &srqJSON, &ipipJSON, &score.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan score: %w", err)
		}

		_ = unmarshalScores(srqJSON, ipipJSON, score)
		scoresMap[score.ParticipantID] = score
	}
	return scoresMap, nil
}

// ── ScoringService ────────────────────────────────────────────────────────────

// ScoringService is a thin facade that delegates to the strongly-typed
// scorers in pkg/scoring. This keeps the handler layer decoupled from
// scoring implementation details (Dependency Inversion Principle).
type ScoringService struct {
	srqScorer  *scoring.SRQScorer
	ipipScorer *scoring.IPIPScorer
}

// NewScoringService creates a new ScoringService with default scorers.
func NewScoringService() *ScoringService {
	return &ScoringService{
		srqScorer:  scoring.NewSRQScorer(),
		ipipScorer: scoring.NewIPIPScorer(),
	}
}

// CalculateSRQScore computes SRQ-29 scores from raw answers.
func (s *ScoringService) CalculateSRQScore(answers map[string]interface{}) (*models.SRQScore, error) {
	return s.srqScorer.Calculate(answers)
}

// CalculateIPIPScore computes IPIP-BFM-50 scores from raw answers.
func (s *ScoringService) CalculateIPIPScore(answers map[string]interface{}) (*models.IPIPScore, error) {
	return s.ipipScorer.Calculate(answers)
}

// ── private helpers ───────────────────────────────────────────────────────────

func marshalScoreJSON(srqScore *models.SRQScore, ipipScore *models.IPIPScore) ([]byte, []byte, error) {
	var srqJSON, ipipJSON []byte
	var err error

	if srqScore != nil {
		srqJSON, err = json.Marshal(srqScore)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal SRQ score: %w", err)
		}
	}
	if ipipScore != nil {
		ipipJSON, err = json.Marshal(ipipScore)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal IPIP score: %w", err)
		}
	}
	return srqJSON, ipipJSON, nil
}

func unmarshalScores(srqJSON, ipipJSON []byte, score *models.Score) error {
	if srqJSON != nil {
		score.SRQScore = &models.SRQScore{}
		if err := json.Unmarshal(srqJSON, score.SRQScore); err != nil {
			return fmt.Errorf("unmarshal SRQ score: %w", err)
		}
	}
	if ipipJSON != nil {
		score.IPIPScore = &models.IPIPScore{}
		if err := json.Unmarshal(ipipJSON, score.IPIPScore); err != nil {
			return fmt.Errorf("unmarshal IPIP score: %w", err)
		}
	}
	return nil
}
