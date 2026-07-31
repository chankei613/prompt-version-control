# Prompt Version Control — 仕様書

> 作成: 2026-08-01

## 1. 製品概要
「プロンプトのGit」。プロンプト/worker定義の変更履歴を管理し、ロールバック・2バージョン間diff・
バージョンごとの品質スコア推移を提供する。

## 2. コアコンセプト
- **PromptDoc**: 管理対象1件（key + CurrentVersionIDというHEADポインタ）
- **PromptVersion**: イミュータブルなスナップショット（追記専用、連番）
- **QualityRating**: 外部から届く、特定バージョンへの品質スコア（0.0-1.0）
- **ロールバック**: 新バージョンを作らずCurrentVersionIDを動かすだけ（git checkoutに相当）
- **A/Bテスト**: v0.1.0では「2バージョンのdiff比較」+「バージョンごとの品質推移グラフ」として実現。
  実トラフィックの振り分けはJ. AI Task Routerの領分でありスコープ外

## 3. 機能一覧
### Phase 1: PromptDoc/PromptVersion CRUD・diff・rollback
### Phase 2: QualityRating・品質推移集計
### Phase 3: UI（一覧・履歴・diff・品質推移・Help）

## 4. データストア
```sql
prompt_docs (id, key, description, current_version_id, created_at, updated_at)
prompt_versions (id, prompt_doc_id, version_no, content, message, created_at)
quality_ratings (id, prompt_version_id, source, score, note, rated_at)
agent_keys (id, name, api_key_hash, created_at, revoked_at)
```
