// cmd/smoketest はPrompt Version ControlのAPIを一時DBで自前起動し、
// ブートストラップ鍵発行 → プロンプト作成 → バージョン追加 → diff → 品質スコア →
// ロールバック（履歴が消えないこと）、の一連が通しで動くことを確認する。
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/chankei613/prompt-version-control/internal/api"
	"github.com/chankei613/prompt-version-control/internal/db"
)

func main() {
	dbPath := "smoketest.db"
	_ = os.Remove(dbPath)
	defer func() { _ = os.Remove(dbPath) }()

	conn, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	srv := httptest.NewServer(api.NewRouter(conn))
	defer srv.Close()

	issueBody, _ := json.Marshal(map[string]string{"name": "smoketest"})
	resp, err := http.Post(srv.URL+"/api/v1/keys", "application/json", bytes.NewReader(issueBody))
	if err != nil {
		log.Fatal(err)
	}
	var issued api.IssueKeyResult
	if err := json.NewDecoder(resp.Body).Decode(&issued); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if issued.APIKey == "" {
		log.Fatal("FAIL: bootstrap key issuance returned empty key")
	}
	fmt.Println("PASS: bootstrap key issued")

	authed := func(method, path string, body []byte) *http.Response {
		req, _ := http.NewRequest(method, srv.URL+path, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+issued.APIKey)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		return resp
	}

	// create prompt with v1
	createBody, _ := json.Marshal(api.CreatePromptInput{
		Key: "claude_rules", Content: "Be concise.\nAlways cite sources.", Message: "initial",
	})
	resp = authed(http.MethodPost, "/api/v1/prompts", createBody)
	var doc db.PromptDoc
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if doc.ID == "" || doc.CurrentVersionID == "" {
		log.Fatalf("FAIL: prompt creation missing id/current_version_id: %+v", doc)
	}
	v1ID := doc.CurrentVersionID
	fmt.Println("PASS: prompt created with initial version")

	// add v2
	v2Body, _ := json.Marshal(api.CreateVersionInput{Content: "Be concise.\nAlways cite sources.\nNever guess.", Message: "add guessing rule"})
	resp = authed(http.MethodPost, "/api/v1/prompts/"+doc.ID+"/versions", v2Body)
	var v2 db.PromptVersion
	if err := json.NewDecoder(resp.Body).Decode(&v2); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if v2.VersionNo != 2 {
		log.Fatalf("FAIL: expected version_no=2, got %d", v2.VersionNo)
	}
	fmt.Println("PASS: second version added, version_no incremented")

	// diff v1 vs v2
	resp = authed(http.MethodGet, "/api/v1/versions/"+v1ID+"/diff?against="+v2.ID, nil)
	var diff api.DiffResult
	if err := json.NewDecoder(resp.Body).Decode(&diff); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	addCount := 0
	for _, l := range diff.Lines {
		if l.Type == "add" {
			addCount++
		}
	}
	if addCount != 1 {
		log.Fatalf("FAIL: expected 1 added line in diff, got %d (%+v)", addCount, diff.Lines)
	}
	fmt.Println("PASS: diff correctly shows 1 added line")

	// rate v1 low, v2 high
	r1, _ := json.Marshal(api.CreateRatingInput{PromptVersionID: v1ID, Score: 0.4})
	_ = authed(http.MethodPost, "/api/v1/ratings", r1).Body.Close()
	r2, _ := json.Marshal(api.CreateRatingInput{PromptVersionID: v2.ID, Score: 0.9})
	_ = authed(http.MethodPost, "/api/v1/ratings", r2).Body.Close()

	resp = authed(http.MethodGet, "/api/v1/prompts/"+doc.ID+"/quality", nil)
	var trend []api.VersionQuality
	if err := json.NewDecoder(resp.Body).Decode(&trend); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(trend) != 2 || trend[0].AvgScore != 0.4 || trend[1].AvgScore != 0.9 {
		log.Fatalf("FAIL: unexpected quality trend: %+v", trend)
	}
	fmt.Println("PASS: quality trend shows v2 scoring higher than v1")

	// rollback to v1 — history must remain (2 versions), only CurrentVersionID changes
	rollbackBody, _ := json.Marshal(map[string]string{"version_id": v1ID})
	resp = authed(http.MethodPost, "/api/v1/prompts/"+doc.ID+"/rollback", rollbackBody)
	var rolledBack db.PromptDoc
	if err := json.NewDecoder(resp.Body).Decode(&rolledBack); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if rolledBack.CurrentVersionID != v1ID {
		log.Fatalf("FAIL: expected current_version_id=%s after rollback, got %s", v1ID, rolledBack.CurrentVersionID)
	}

	resp = authed(http.MethodGet, "/api/v1/prompts/"+doc.ID+"/versions", nil)
	var versions []db.PromptVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(versions) != 2 {
		log.Fatalf("FAIL: rollback must not delete history, expected 2 versions, got %d", len(versions))
	}
	fmt.Println("PASS: rollback moved current pointer without deleting history")

	fmt.Println("SMOKE TEST OK")
}
