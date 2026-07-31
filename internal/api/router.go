package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Server は全ロジックの実体。HTTPハンドラーとWailsネイティブバインディングの
// 両方がこの同じ Server のメソッドを呼ぶことで、UIとAPIの挙動がズレないようにする。
type Server struct {
	DB *gorm.DB
}

func New(conn *gorm.DB) *Server {
	return &Server{DB: conn}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1/keys", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB, "/api/v1/keys"))
		r.Post("/", s.httpIssueKey)
		r.Get("/", s.httpListKeys)
		r.Delete("/{id}", s.httpRevokeKey)
	})

	r.Route("/api/v1/prompts", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Post("/", s.httpCreatePrompt)
		r.Get("/", s.httpListPrompts)
		r.Get("/{id}", s.httpGetPrompt)
		r.Delete("/{id}", s.httpDeletePrompt)
		r.Post("/{id}/versions", s.httpCreateVersion)
		r.Get("/{id}/versions", s.httpListVersions)
		r.Post("/{id}/rollback", s.httpRollback)
		r.Get("/{id}/quality", s.httpQualityTrend)
	})

	r.Route("/api/v1/versions", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Get("/{id}/diff", s.httpDiff)
	})

	r.Route("/api/v1/ratings", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Post("/", s.httpAddRating)
	})

	return r
}

// NewRouter はcmd/pvcserve（単体HTTPサーバー）向けの簡易コンストラクタ。
func NewRouter(conn *gorm.DB) http.Handler {
	return New(conn).Router()
}
