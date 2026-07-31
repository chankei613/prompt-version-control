// Package db はPrompt Version ControlのGORMモデルとSQLite初期化を提供する。
// docs/spec.md参照。
package db

import "time"

// PromptDoc は管理対象の1つのプロンプト/worker定義ファイル。
// CurrentVersionID は git の HEAD に相当し、ロールバックはこのポインタを動かすだけで
// 履歴（PromptVersion）は破壊しない。
type PromptDoc struct {
	ID               string `gorm:"primaryKey" json:"id"`
	Key              string `gorm:"uniqueIndex" json:"key"`
	Description      string `json:"description"`
	CurrentVersionID string `json:"current_version_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PromptVersion はイミュータブルなスナップショット（追記専用）。PromptDoc内で連番を持つ。
type PromptVersion struct {
	ID          string `gorm:"primaryKey" json:"id"`
	PromptDocID string `gorm:"index" json:"prompt_doc_id"`
	VersionNo   int    `json:"version_no"`
	Content     string `json:"content"`
	Message     string `json:"message"`

	CreatedAt time.Time `json:"created_at"`
}

// QualityRating は外部システムから届く、特定バージョンへの品質スコア。
type QualityRating struct {
	ID              string  `gorm:"primaryKey" json:"id"`
	PromptVersionID string  `gorm:"index" json:"prompt_version_id"`
	Source          string  `json:"source"`
	Score           float64 `json:"score"`
	Note            string  `json:"note"`

	RatedAt time.Time `json:"rated_at"`
}

// AgentKey — CRUD/Ingestion APIを叩くためのAPIキー。ハッシュのみ保存する。
type AgentKey struct {
	ID         string     `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	APIKeyHash string     `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}
