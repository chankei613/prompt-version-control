# 技術選定

**決定日:** 2026-08-01

| レイヤー | 採用 |
|---|---|
| Desktop | Wails v2（v2.12.0固定） |
| Backend | Go 1.22+ |
| Frontend | Vue 3 + Vite + Pinia |
| DB | SQLite + GORM |
| CI | GitHub Actions（Go 1.23 / macos-14） |
| 配布 | `.app`/`.exe` + コード署名・公証（環境変数ベースのbuild-release.sh） |

## 前製品群から必ず引き継ぐこと

1. `.golangci.yml` はv2形式で最初から書く
2. `frontend/wailsjs/` は `.gitignore` に入れない
3. CIのgo-testジョブは `./internal/... ./cmd/...` にスコープ限定
4. `build-release.sh`・`wails.json` に個人情報を直書きしない
5. Phase 3のUI実装に、アプリ内Help画面を最初から含める
6. **テストファイルを新規作成・変更したら、コミット前に必ず `npm run typecheck` を再実行する**。
   wailsjsのクラス型は `api.XxxResult.createFrom({...})` で構築する
7. release.ymlはharness-managerのv0.1.1構成をそのままコピー
8. Vercelデプロイ後は実際に割り当てられたURLを確認してからog:url/og:imageを設定する
   （クリーンな短縮ドメインになるか"-xxx"サフィックス付きになるかは事前に分からない）

## Apple Developer認証情報
TEAM_ID: `D5R956CRBE` / Developer ID Application: `keisuke haraguchi`。
Apple app-specific passwordは既存製品で発行済みのものを使い回せる可能性が高い。
