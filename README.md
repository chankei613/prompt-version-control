# Prompt Version Control

「プロンプトのGit」— comet-taskAI ロードマップ Product H。

systemプロンプト・worker定義ファイルの変更履歴を管理し、ロールバック・2バージョン間diff・
バージョンごとの品質スコア推移を提供する。「このプロンプトを変えたら質が下がった」を可視化する。

詳細は [docs/spec.md](docs/spec.md) を参照。

## 現在のステータス: v0.1.0 リリース済み

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・CRUD API・diff・rollback
- [x] Phase 2: 品質スコア・集計
- [x] Phase 3: Wails + Vue3 UI
- [x] Phase 4: 仕上げ・署名・配布・LP

macOSアプリ（署名・公証済み）は [GitHub Releases](https://github.com/chankei613/prompt-version-control/releases) から、
ランディングページは https://prompt-version-control-nine.vercel.app/ から入手できる。
アプリ内のHelpタブに使い方の説明がある。

## 使い方（開発用ヘッドレスサーバー）

```bash
go mod tidy
go run .      # :8425 でAPIサーバー起動
go run ./cmd/smoketest
```

### プロンプトの作成・更新・ロールバック

```bash
curl -X POST localhost:8425/api/v1/prompts \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"key":"claude_rules","content":"Be concise.","message":"initial"}'

curl -X POST localhost:8425/api/v1/prompts/{id}/versions \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"content":"Be concise.\nNever guess.","message":"add guessing rule"}'

curl -X POST localhost:8425/api/v1/prompts/{id}/rollback \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"version_id":"{旧バージョンID}"}'
```

**ロールバックは履歴を消さない。** `current_version_id` というポインタを動かすだけで、
過去のバージョンは全て残り続ける（gitのcheckoutに相当）。

## API

| メソッド | パス | 用途 |
|---|---|---|
| POST/GET/DELETE | `/api/v1/keys` | APIキー管理 |
| POST | `/api/v1/prompts` | プロンプト作成（初期バージョンも同時作成） |
| GET | `/api/v1/prompts` | 一覧 |
| GET | `/api/v1/prompts/{id}` | 単体取得 |
| DELETE | `/api/v1/prompts/{id}` | 削除（履歴ごと） |
| POST | `/api/v1/prompts/{id}/versions` | 新バージョン追加（current_version_idを更新） |
| GET | `/api/v1/prompts/{id}/versions` | 履歴一覧 |
| POST | `/api/v1/prompts/{id}/rollback` | current_version_idを指定バージョンへ戻す |
| GET | `/api/v1/prompts/{id}/quality` | バージョンごとの平均品質スコア推移 |
| GET | `/api/v1/versions/{id}/diff?against={id}` | 2バージョン間のdiff |
| POST | `/api/v1/ratings` | 品質スコア追加（外部システムから） |

## ディレクトリ構成

```
internal/db/       GORMモデル（PromptDoc/PromptVersion/QualityRating/AgentKey）
internal/api/       REST API（prompts/versions/ratings/keys）+ diffアルゴリズム + 認証ミドルウェア
cmd/smoketest/      通しスモークテスト
docs/                設計ドキュメント
```
