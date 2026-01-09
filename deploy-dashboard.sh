#!/bin/bash

# エラーが発生したら停止
set -e

# プロジェクトルートに移動
cd "$(dirname "$0")"
ROOT_DIR=$(pwd)

echo "--- Building Dashboard-V2 ---"
cd "$ROOT_DIR/dashboard-v2"

# 依存関係のインストール（必要に応じて）
# npm install

# ビルド実行
npm run build

echo "--- Deploying to Auth Backend ---"
DEST_DIR="$ROOT_DIR/auth/src/dashboard"

# 配布先ディレクトリのクリーンアップと作成
rm -rf "$DEST_DIR"
mkdir -p "$DEST_DIR"

# ビルド成果物 (dist) をコピー
cp -r dist/* "$DEST_DIR/"

echo "--- Done! ---"
echo "Dashboard has been deployed to $DEST_DIR"
