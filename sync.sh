#!/bin/bash
set -euo pipefail

# ================================================
# Script: sync.sh
# Description:
#   将当前仓库的 main 分支和指定 tag 推送到另一个仓库。
#   支持自动识别目标仓库（GitHub 或 CNB）。
#
# Env Vars:
#   TOKEN   - 目标仓库的访问令牌（必填）
#   TO_GIT  - 目标仓库标识字符串（如 github / cnb）（必填）
#   TAG     - 可选，指定要推送的 tag 名称
#
# Usage:
#   TOKEN=xxx TO_GIT=github TAG=v1.0.0 bash sync.sh
#
# Author: 你自己，宇宙第一骚
# ================================================
# ==== Validate parameters ==== #


: "${TOKEN:?TOKEN is required}"
: "${TO_GIT:?TO_GIT is required}"
TAG="${TAG:-}"

# ==== Determine remote URL ==== #
if [[ "$TO_GIT" == *github* ]]; then
  TO_GIT_URL="https://chihqiang:${TOKEN}@github.com/chihqiang/tlsctl.git"
elif [[ "$TO_GIT" == *cnb* ]]; then
  TO_GIT_URL="https://cnb:${TOKEN}@cnb.cool/zhiqiangwang/tlsctl.git"
else
  echo "❌ Unknown TO_GIT source: $TO_GIT"
  exit 1
fi

# ==== Add or update remote ==== #
git remote add togit "$TO_GIT_URL" 2>/dev/null || git remote set-url togit "$TO_GIT_URL"

# ==== Push main branch ==== #
echo "🚀 Pushing main branch..."
git push -u togit HEAD:main -f

# ==== Push tag if specified ==== #
if [[ -n "$TAG" ]]; then
  if ! git rev-parse "refs/tags/$TAG" >/dev/null 2>&1; then
    echo "⚠️  Local tag '$TAG' does not exist. Skipping."
    exit 0
  fi

  echo "🔍 Checking if remote tag '$TAG' exists..."
  if git ls-remote --tags togit | grep -q "refs/tags/$TAG"; then
    echo "🚫 Remote tag '$TAG' already exists. Skipping push."
  else
    echo "🏷️  Pushing tag: $TAG"
    git push togit "refs/tags/$TAG"
  fi
fi

echo "✅ Sync complete"
