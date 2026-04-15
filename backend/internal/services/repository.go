package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/google/uuid"
)

// ParticipantService handles participant-related operations
type ParticipantService struct {
	DB *database.DB
}

// NewParticipantService creates a new participant service
func NewParticipantService(db *database.DB) *ParticipantService {
	return &ParticipantService{DB: db}
}

// CreateParticipant creates a new participant
func (s *ParticipantService) CreateParticipant(req models.ParticipantRequest) (*models.Participant, error) {
	participant := &models.Participant{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Age:       req.Age,
		Gender:    req.Gender,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO participants (id, name, age, gender, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.DB.Connection.Exec(query, participant.ID, participant.Name, participant.Age, participant.Gender, participant.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create participant: %w", err)
	}

	return participant, nil
}

// GetParticipant retrieves a participant by ID
func (s *ParticipantService) GetParticipant(id string) (*models.Participant, error) {
	participant := &models.Participant{}
	query := `SELECT id, name, age, gender, created_at FROM participants WHERE id = $1`

	err := s.DB.Connection.QueryRow(query, id).Scan(
		&participant.ID,
		&participant.Name,
		&participant.Age,
		&participant.Gender,
		&participant.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant: %w", err)
	}

	return participant, nil
}

// GetAllParticipants retrieves all participants
func (s *ParticipantService) GetAllParticipants() ([]*models.Participant, error) {
	var participants []*models.Participant
	query := `SELECT id, name, age, gender, created_at FROM participants ORDER BY created_at DESC`

	rows, err := s.DB.Connection.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		p := &models.Participant{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Gender, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, p)
	}
	return participants, nil
}

// ResponseService handles response-related operations
type ResponseService struct {
	DB *database.DB
}

// NewResponseService creates a new response service
func NewResponseService(db *database.DB) *ResponseService {
	return &ResponseService{DB: db}
}

// SaveResponse saves questionnaire responses
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
		return nil, fmt.Errorf("failed to marshal answers: %w", err)
	}

	query := `INSERT INTO responses (id, participant_id, questionnaire_type, answers, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = s.DB.Connection.Exec(query, response.ID, response.ParticipantID, response.QuestionnaireType, answersJSON, response.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to save response: %w", err)
	}

	return response, nil
}

// ScoreService handles score-related operations
type ScoreService struct {
	DB             *database.DB
	ScoringService *ScoringService
}

// NewScoreService creates a new score service
func NewScoreService(db *database.DB) *ScoreService {
	return &ScoreService{
		DB:             db,
		ScoringService: NewScoringService(),
	}
}

// SaveScore saves calculated scores
func (s *ScoreService) SaveScore(participantID string, srqScore *models.SRQScore, ipipScore *models.IPIPScore) (*models.Score, error) {
	score := &models.Score{
		ParticipantID: participantID,
		CreatedAt:     time.Now(),
	}

	// Merge with an existing row so SRQ and IPIP results are preserved for the same participant.
	var existingID string
	var existingSRQJSON, existingIPIPJSON []byte

	err := s.DB.Connection.QueryRow(
		`SELECT id, srq_score, ipip_score FROM scores WHERE participant_id = $1 ORDER BY created_at DESC LIMIT 1`,
		participantID,
	).Scan(&existingID, &existingSRQJSON, &existingIPIPJSON)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing score: %w", err)
	}

	if err == sql.ErrNoRows {
		score.ID = uuid.New().String()
		score.SRQScore = srqScore
		score.IPIPScore = ipipScore

		srqJSON, ipipJSON, marshalErr := marshalScoreJSON(score.SRQScore, score.IPIPScore)
		if marshalErr != nil {
			return nil, marshalErr
		}

		_, execErr := s.DB.Connection.Exec(
			`INSERT INTO scores (id, participant_id, srq_score, ipip_score, created_at) VALUES ($1, $2, $3, $4, $5)`,
			score.ID,
			score.ParticipantID,
			srqJSON,
			ipipJSON,
			score.CreatedAt,
		)
		if execErr != nil {
			return nil, fmt.Errorf("failed to save score: %w", execErr)
		}

		return score, nil
	}

	score.ID = existingID

	if existingSRQJSON != nil {
		existingSRQ := &models.SRQScore{}
		if unmarshalErr := json.Unmarshal(existingSRQJSON, existingSRQ); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal existing SRQ score: %w", unmarshalErr)
		}
		score.SRQScore = existingSRQ
	}
	if existingIPIPJSON != nil {
		existingIPIP := &models.IPIPScore{}
		if unmarshalErr := json.Unmarshal(existingIPIPJSON, existingIPIP); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal existing IPIP score: %w", unmarshalErr)
		}
		score.IPIPScore = existingIPIP
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
		srqJSON,
		ipipJSON,
		score.CreatedAt,
		score.ID,
	)
	if execErr != nil {
		return nil, fmt.Errorf("failed to update score: %w", execErr)
	}

	return score, nil
}

// GetScores retrieves scores for a participant
func (s *ScoreService) GetScores(participantID string) (*models.Score, error) {
	score := &models.Score{}
	query := `SELECT id, participant_id, srq_score, ipip_score, created_at FROM scores WHERE participant_id = $1 ORDER BY created_at DESC LIMIT 1`

	var srqJSON, ipipJSON []byte

	err := s.DB.Connection.QueryRow(query, participantID).Scan(
		&score.ID,
		&score.ParticipantID,
		&srqJSON,
		&ipipJSON,
		&score.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get scores: %w", err)
	}

	if srqJSON != nil {
		score.SRQScore = &models.SRQScore{}
		if err := json.Unmarshal(srqJSON, score.SRQScore); err != nil {
			return nil, fmt.Errorf("failed to unmarshal SRQ score: %w", err)
		}
	}

	if ipipJSON != nil {
		score.IPIPScore = &models.IPIPScore{}
		if err := json.Unmarshal(ipipJSON, score.IPIPScore); err != nil {
			return nil, fmt.Errorf("failed to unmarshal IPIP score: %w", err)
		}
	}

	return score, nil
}

// GetAllScores returns all latest scores mapped by participant ID
func (s *ScoreService) GetAllScores() (map[string]*models.Score, error) {
	query := `
		SELECT t.id, t.participant_id, t.srq_score, t.ipip_score, t.created_at
		FROM scores t
		INNER JOIN (
			SELECT participant_id, MAX(created_at) as MaxDate
			FROM scores
			GROUP BY participant_id
		) tm ON t.participant_id = tm.participant_id AND t.created_at = tm.MaxDate
	`

	rows, err := s.DB.Connection.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all scores: %w", err)
	}
	defer rows.Close()

	scoresMap := make(map[string]*models.Score)
	for rows.Next() {
		score := &models.Score{}
		var srqJSON, ipipJSON []byte

		err := rows.Scan(
			&score.ID,
			&score.ParticipantID,
			&srqJSON,
			&ipipJSON,
			&score.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan score: %w", err)
		}

		if srqJSON != nil {
			score.SRQScore = &models.SRQScore{}
			json.Unmarshal(srqJSON, score.SRQScore)
		}
		if ipipJSON != nil {
			score.IPIPScore = &models.IPIPScore{}
			json.Unmarshal(ipipJSON, score.IPIPScore)
		}

		scoresMap[score.ParticipantID] = score
	}

	return scoresMap, nil
}

func marshalScoreJSON(srqScore *models.SRQScore, ipipScore *models.IPIPScore) ([]byte, []byte, error) {
	var srqJSON, ipipJSON []byte
	var err error

	if srqScore != nil {
		srqJSON, err = json.Marshal(srqScore)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal SRQ score: %w", err)
		}
	}

	if ipipScore != nil {
		ipipJSON, err = json.Marshal(ipipScore)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal IPIP score: %w", err)
		}
	}

	return srqJSON, ipipJSON, nil
}
