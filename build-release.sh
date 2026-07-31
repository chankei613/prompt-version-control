#!/bin/bash
# Prompt Version Control release build: wails build → codesign → notarize → staple → zip
# 個人情報はスクリプトに直書きせず、全て環境変数で渡す。
set -e

APPLE_ID="${APPLE_ID:-}"
TEAM_ID="${TEAM_ID:-}"
APP_PASSWORD="${APP_PASSWORD:-}"
DEVELOPER_NAME="${DEVELOPER_NAME:-}"

APP_NAME="prompt-version-control"
VERSION=$(grep 'AppVersion' version.go | sed 's/.*"\(.*\)".*/\1/')
APP_PATH="build/bin/${APP_NAME}.app"
ZIP_PATH="${APP_NAME}-${VERSION}.zip"
ENTITLEMENTS="build/darwin/entitlements.plist"

if [ -z "$TEAM_ID" ] || [ -z "$APPLE_ID" ] || [ -z "$DEVELOPER_NAME" ] || [ -z "$APP_PASSWORD" ]; then
  echo "ERROR: APPLE_ID / TEAM_ID / DEVELOPER_NAME / APP_PASSWORD が未設定です。"
  echo "以下のコマンドで Developer ID 一覧を確認してください:"
  echo "  security find-identity -v -p codesigning"
  echo ""
  echo "実行例:"
  echo "  APPLE_ID=you@example.com TEAM_ID=ABC1234567 DEVELOPER_NAME=\"Your Name\" APP_PASSWORD=xxxx-xxxx-xxxx-xxxx ./build-release.sh"
  exit 1
fi

IDENTITY="Developer ID Application: ${DEVELOPER_NAME} (${TEAM_ID})"

echo "==> Building Prompt Version Control v${VERSION}..."

xattr -cr build/ 2>/dev/null || true
export PATH="$PATH:$HOME/go/bin"
wails build -platform darwin/universal -o "${APP_NAME}"

echo "==> Code signing..."
codesign \
  --deep \
  --force \
  --verify \
  --verbose \
  --sign "${IDENTITY}" \
  --options runtime \
  --entitlements "${ENTITLEMENTS}" \
  "${APP_PATH}"

codesign --verify --deep --strict "${APP_PATH}"
echo "    Signature OK"

echo "==> Creating zip for notarization..."
ditto -c -k --keepParent "${APP_PATH}" "${ZIP_PATH}"

echo "==> Submitting for notarization (this takes a few minutes)..."
xcrun notarytool submit "${ZIP_PATH}" \
  --apple-id "${APPLE_ID}" \
  --team-id "${TEAM_ID}" \
  --password "${APP_PASSWORD}" \
  --wait

echo "==> Stapling notarization ticket..."
xcrun stapler staple "${APP_PATH}"
xcrun stapler validate "${APP_PATH}"
echo "    Staple OK"

echo "==> Re-zipping stapled app for distribution..."
rm -f "${ZIP_PATH}"
ditto -c -k --keepParent "${APP_PATH}" "${ZIP_PATH}"

echo ""
echo "Done: ${ZIP_PATH}"
echo ""
echo "次のステップ:"
echo "  1. gh release create v${VERSION} ${ZIP_PATH} でGitHub Releaseを作成"
echo "  2. landing/ にzipを同梱してVercelへデプロイ"
