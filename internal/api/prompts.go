package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/prompt-version-control/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	errPromptNotFound  = &apiError{"prompt not found"}
	errVersionNotFound = &apiError{"version not found"}
	errKeyRequired     = &apiError{"key is required"}
	errWrongPromptDoc  = &apiError{"version does not belong to this prompt"}
)

type CreatePromptInput struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Content     string `json:"content"` // 初期バージョン（version 1）として保存する
	Message     string `json:"message"`
}

// CreatePrompt はPromptDocと初期バージョン(1)を同時に作成する（HTTP・ネイティブバインディング共用）。
func (s *Server) CreatePrompt(in CreatePromptInput) (db.PromptDoc, error) {
	if in.Key == "" {
		return db.PromptDoc{}, errKeyRequired
	}

	now := time.Now()
	doc := db.PromptDoc{
		ID:          uuid.NewString(),
		Key:         in.Key,
		Description: in.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.DB.Create(&doc).Error; err != nil {
		return db.PromptDoc{}, err
	}

	version := db.PromptVersion{
		ID:          uuid.NewString(),
		PromptDocID: doc.ID,
		VersionNo:   1,
		Content:     in.Content,
		Message:     in.Message,
		CreatedAt:   now,
	}
	if err := s.DB.Create(&version).Error; err != nil {
		return db.PromptDoc{}, err
	}

	doc.CurrentVersionID = version.ID
	if err := s.DB.Save(&doc).Error; err != nil {
		return db.PromptDoc{}, err
	}
	return doc, nil
}

func (s *Server) ListPrompts() ([]db.PromptDoc, error) {
	var rows []db.PromptDoc
	err := s.DB.Order("key asc").Find(&rows).Error
	return rows, err
}

func (s *Server) GetPrompt(id string) (db.PromptDoc, error) {
	var doc db.PromptDoc
	if err := s.DB.First(&doc, "id = ?", id).Error; err != nil {
		return db.PromptDoc{}, errPromptNotFound
	}
	return doc, nil
}

func (s *Server) DeletePrompt(id string) error {
	res := s.DB.Delete(&db.PromptDoc{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errPromptNotFound
	}
	// バージョン履歴も一緒に削除する（PromptDoc自体を消す以上、孤児レコードを残さない）
	s.DB.Delete(&db.PromptVersion{}, "prompt_doc_id = ?", id)
	return nil
}

type CreateVersionInput struct {
	Content string `json:"content"`
	Message string `json:"message"`
}

// CreateVersion は新しいスナップショットを追記し、CurrentVersionIDをそこへ進める。
func (s *Server) CreateVersion(promptID string, in CreateVersionInput) (db.PromptVersion, error) {
	doc, err := s.GetPrompt(promptID)
	if err != nil {
		return db.PromptVersion{}, err
	}

	var maxNo int
	s.DB.Model(&db.PromptVersion{}).Where("prompt_doc_id = ?", promptID).Select("COALESCE(MAX(version_no), 0)").Scan(&maxNo)

	version := db.PromptVersion{
		ID:          uuid.NewString(),
		PromptDocID: promptID,
		VersionNo:   maxNo + 1,
		Content:     in.Content,
		Message:     in.Message,
		CreatedAt:   time.Now(),
	}
	if err := s.DB.Create(&version).Error; err != nil {
		return db.PromptVersion{}, err
	}

	doc.CurrentVersionID = version.ID
	doc.UpdatedAt = time.Now()
	if err := s.DB.Save(&doc).Error; err != nil {
		return db.PromptVersion{}, err
	}
	return version, nil
}

func (s *Server) ListVersions(promptID string) ([]db.PromptVersion, error) {
	var rows []db.PromptVersion
	err := s.DB.Where("prompt_doc_id = ?", promptID).Order("version_no desc").Find(&rows).Error
	return rows, err
}

func (s *Server) GetVersion(id string) (db.PromptVersion, error) {
	var v db.PromptVersion
	if err := s.DB.First(&v, "id = ?", id).Error; err != nil {
		return db.PromptVersion{}, errVersionNotFound
	}
	return v, nil
}

// Rollback はCurrentVersionIDを既存のバージョンへ戻すだけで、履歴は一切変更しない
// （gitのcheckoutに相当。revertのように新しいコミットは作らない）。
func (s *Server) Rollback(promptID, versionID string) (db.PromptDoc, error) {
	doc, err := s.GetPrompt(promptID)
	if err != nil {
		return db.PromptDoc{}, err
	}
	version, err := s.GetVersion(versionID)
	if err != nil {
		return db.PromptDoc{}, err
	}
	if version.PromptDocID != promptID {
		return db.PromptDoc{}, errWrongPromptDoc
	}

	doc.CurrentVersionID = versionID
	doc.UpdatedAt = time.Now()
	if err := s.DB.Save(&doc).Error; err != nil {
		return db.PromptDoc{}, err
	}
	return doc, nil
}

type DiffResult struct {
	FromVersionID string     `json:"from_version_id"`
	ToVersionID   string     `json:"to_version_id"`
	Lines         []DiffLine `json:"lines"`
}

func (s *Server) Diff(versionID, againstID string) (DiffResult, error) {
	from, err := s.GetVersion(versionID)
	if err != nil {
		return DiffResult{}, err
	}
	to, err := s.GetVersion(againstID)
	if err != nil {
		return DiffResult{}, err
	}
	return DiffResult{
		FromVersionID: from.ID,
		ToVersionID:   to.ID,
		Lines:         diffLines(from.Content, to.Content),
	}, nil
}

// ─── HTTPハンドラー ────────────────────────────────────────────────────

func (s *Server) httpCreatePrompt(w http.ResponseWriter, r *http.Request) {
	var body CreatePromptInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	doc, err := s.CreatePrompt(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, doc)
}

func (s *Server) httpListPrompts(w http.ResponseWriter, r *http.Request) {
	rows, err := s.ListPrompts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) httpGetPrompt(w http.ResponseWriter, r *http.Request) {
	doc, err := s.GetPrompt(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

func (s *Server) httpDeletePrompt(w http.ResponseWriter, r *http.Request) {
	if err := s.DeletePrompt(chi.URLParam(r, "id")); err != nil {
		if err == errPromptNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) httpCreateVersion(w http.ResponseWriter, r *http.Request) {
	var body CreateVersionInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	v, err := s.CreateVersion(chi.URLParam(r, "id"), body)
	if err != nil {
		if err == errPromptNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (s *Server) httpListVersions(w http.ResponseWriter, r *http.Request) {
	rows, err := s.ListVersions(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

type rollbackRequest struct {
	VersionID string `json:"version_id"`
}

func (s *Server) httpRollback(w http.ResponseWriter, r *http.Request) {
	var body rollbackRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	doc, err := s.Rollback(chi.URLParam(r, "id"), body.VersionID)
	if err != nil {
		if err == errPromptNotFound || err == errVersionNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

func (s *Server) httpDiff(w http.ResponseWriter, r *http.Request) {
	result, err := s.Diff(chi.URLParam(r, "id"), r.URL.Query().Get("against"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
