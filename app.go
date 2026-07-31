package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chankei613/prompt-version-control/internal/api"
	"github.com/chankei613/prompt-version-control/internal/db"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const apiAddr = "127.0.0.1:8425"

// App はWailsのバインディング。実処理は internal/api.Server が持っている。
type App struct {
	ctx    context.Context
	server *api.Server
	srv    *http.Server
	ready  bool
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir := appDataDir()
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		runtime.LogErrorf(ctx, "data dir error: %s", err)
		return
	}

	conn, err := db.Init(filepath.Join(dataDir, "prompt-version-control.db"))
	if err != nil {
		runtime.LogErrorf(ctx, "db init error: %s", err)
		return
	}
	a.server = api.New(conn)

	a.srv = &http.Server{Addr: apiAddr, Handler: a.server.Router()}
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			runtime.LogErrorf(ctx, "api server error: %s", err)
		}
	}()

	a.ready = true
	runtime.LogInfof(ctx, "Prompt Version Control ready (api: http://%s, data: %s)", apiAddr, dataDir)
}

func (a *App) shutdown(ctx context.Context) {
	if a.srv != nil {
		_ = a.srv.Close()
	}
	if a.server != nil {
		if sqlDB, err := a.server.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

var errNotReady = &notReadyError{}

type notReadyError struct{}

func (e *notReadyError) Error() string { return "app not ready — check startup logs" }

// ─── フロントエンドへ公開するメソッド ──────────────────────────────────────────

func (a *App) GetAppVersion() string {
	return AppVersion
}

func (a *App) GetAPIURL() string {
	return "http://" + apiAddr
}

func (a *App) ListPrompts() ([]db.PromptDoc, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListPrompts()
}

func (a *App) GetPrompt(id string) (db.PromptDoc, error) {
	if !a.ready {
		return db.PromptDoc{}, errNotReady
	}
	return a.server.GetPrompt(id)
}

func (a *App) CreatePrompt(key, description, content, message string) (db.PromptDoc, error) {
	if !a.ready {
		return db.PromptDoc{}, errNotReady
	}
	return a.server.CreatePrompt(api.CreatePromptInput{Key: key, Description: description, Content: content, Message: message})
}

func (a *App) DeletePrompt(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.DeletePrompt(id)
}

func (a *App) CreateVersion(promptID, content, message string) (db.PromptVersion, error) {
	if !a.ready {
		return db.PromptVersion{}, errNotReady
	}
	return a.server.CreateVersion(promptID, api.CreateVersionInput{Content: content, Message: message})
}

func (a *App) ListVersions(promptID string) ([]db.PromptVersion, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListVersions(promptID)
}

func (a *App) Rollback(promptID, versionID string) (db.PromptDoc, error) {
	if !a.ready {
		return db.PromptDoc{}, errNotReady
	}
	return a.server.Rollback(promptID, versionID)
}

func (a *App) Diff(versionID, againstID string) (api.DiffResult, error) {
	if !a.ready {
		return api.DiffResult{}, errNotReady
	}
	return a.server.Diff(versionID, againstID)
}

func (a *App) QualityTrend(promptID string) ([]api.VersionQuality, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.QualityTrend(promptID)
}

func (a *App) ListKeys() ([]db.AgentKey, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListKeys()
}

func (a *App) IssueKey(name string) (api.IssueKeyResult, error) {
	if !a.ready {
		return api.IssueKeyResult{}, errNotReady
	}
	return a.server.IssueKey(name)
}

func (a *App) RevokeKey(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.RevokeKey(id)
}

// Quit はアプリを完全終了する（Settings 画面から呼ぶ）。
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

func appDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".prompt-version-control")
}
