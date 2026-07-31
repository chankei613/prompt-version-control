package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/prompt-version-control/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateRatingInput struct {
	PromptVersionID string  `json:"prompt_version_id"`
	Source          string  `json:"source"`
	Score           float64 `json:"score"`
	Note            string  `json:"note"`
}

// AddRating は特定バージョンへの品質スコアを追加する（追記専用）。
func (s *Server) AddRating(in CreateRatingInput) (db.QualityRating, error) {
	if in.PromptVersionID == "" {
		return db.QualityRating{}, &apiError{"prompt_version_id is required"}
	}
	if _, err := s.GetVersion(in.PromptVersionID); err != nil {
		return db.QualityRating{}, err
	}

	rating := db.QualityRating{
		ID:              uuid.NewString(),
		PromptVersionID: in.PromptVersionID,
		Source:          in.Source,
		Score:           in.Score,
		Note:            in.Note,
		RatedAt:         time.Now(),
	}
	if err := s.DB.Create(&rating).Error; err != nil {
		return db.QualityRating{}, err
	}
	return rating, nil
}

type VersionQuality struct {
	VersionID   string  `json:"version_id"`
	VersionNo   int     `json:"version_no"`
	AvgScore    float64 `json:"avg_score"`
	RatingCount int     `json:"rating_count"`
}

// QualityTrend はPromptDocの各バージョンごとの平均スコアを、バージョン番号順に返す。
// 「このバージョンに変えたら質が下がった」をグラフ化するための集計。
func (s *Server) QualityTrend(promptID string) ([]VersionQuality, error) {
	versions, err := s.ListVersions(promptID)
	if err != nil {
		return nil, err
	}

	// ListVersionsはversion_no降順なので、表示は昇順に揃える
	for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
		versions[i], versions[j] = versions[j], versions[i]
	}

	result := make([]VersionQuality, 0, len(versions))
	for _, v := range versions {
		var ratings []db.QualityRating
		if err := s.DB.Where("prompt_version_id = ?", v.ID).Find(&ratings).Error; err != nil {
			return nil, err
		}
		vq := VersionQuality{VersionID: v.ID, VersionNo: v.VersionNo, RatingCount: len(ratings)}
		if len(ratings) > 0 {
			var sum float64
			for _, r := range ratings {
				sum += r.Score
			}
			vq.AvgScore = sum / float64(len(ratings))
		}
		result = append(result, vq)
	}
	return result, nil
}

// ─── HTTPハンドラー ────────────────────────────────────────────────────

func (s *Server) httpAddRating(w http.ResponseWriter, r *http.Request) {
	var body CreateRatingInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	rating, err := s.AddRating(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, rating)
}

func (s *Server) httpQualityTrend(w http.ResponseWriter, r *http.Request) {
	result, err := s.QualityTrend(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
