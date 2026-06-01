#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

OLD_DASH="new""-api"
OLD_SPACED="New"" API"
OLD_COMPACT="New""API"
OLD_COMPACT_UPPER="NEW""API"
OLD_HYPHEN_TITLE="New""-API"
OLD_WEB_PACKAGE="newapi""-web"
OLD_DESKTOP_ID="com.""newapi"".desktop"
OLD_PLAIN="newapi"

NEW_DASH="fast-api"
NEW_SPACED="Fast API"
NEW_COMPACT="FastAPI"
NEW_COMPACT_UPPER="FASTAPI"
NEW_HYPHEN_TITLE="Fast-API"
NEW_WEB_PACKAGE="fastapi-web"
NEW_DESKTOP_ID="com.fastapi.desktop"
NEW_PLAIN="fastapi"

export OLD_DASH OLD_SPACED OLD_COMPACT OLD_COMPACT_UPPER OLD_HYPHEN_TITLE OLD_WEB_PACKAGE OLD_DESKTOP_ID OLD_PLAIN
export NEW_DASH NEW_SPACED NEW_COMPACT NEW_COMPACT_UPPER NEW_HYPHEN_TITLE NEW_WEB_PACKAGE NEW_DESKTOP_ID NEW_PLAIN

# 1. Rename the systemd service file if it exists
if [ -f "${OLD_DASH}.service" ] && [ ! -e "${NEW_DASH}.service" ]; then
  if git ls-files --error-unmatch "${OLD_DASH}.service" >/dev/null 2>&1; then
    git mv "${OLD_DASH}.service" "${NEW_DASH}.service"
  else
    mv "${OLD_DASH}.service" "${NEW_DASH}.service"
  fi
fi

# 2. Main replacement pass
while IFS= read -r -d '' file; do
  perl -0pi -e '
    s/\Q$ENV{OLD_DASH}\E/$ENV{NEW_DASH}/g;
    s/\Q$ENV{OLD_SPACED}\E/$ENV{NEW_SPACED}/g;
    s/\Q$ENV{OLD_COMPACT}\E/$ENV{NEW_COMPACT}/g;
    s/\Q$ENV{OLD_COMPACT_UPPER}\E/$ENV{NEW_COMPACT_UPPER}/g;
    s/\Q$ENV{OLD_HYPHEN_TITLE}\E/$ENV{NEW_HYPHEN_TITLE}/g;
    s/\Q$ENV{OLD_WEB_PACKAGE}\E/$ENV{NEW_WEB_PACKAGE}/g;
    s/\Q$ENV{OLD_DESKTOP_ID}\E/$ENV{NEW_DESKTOP_ID}/g;
    s/\Q$ENV{OLD_PLAIN}\E/$ENV{NEW_PLAIN}/g;
  ' "$file"
done < <(git grep -Il -z \
  -e "$OLD_DASH" \
  -e "$OLD_SPACED" \
  -e "$OLD_COMPACT" \
  -e "$OLD_COMPACT_UPPER" \
  -e "$OLD_HYPHEN_TITLE" \
  -e "$OLD_WEB_PACKAGE" \
  -e "$OLD_DESKTOP_ID" \
  -e "$OLD_PLAIN" \
  -- . || true)

# 3. Revert external upstream documentation URLs that should NOT be rebranded.
#    These domains and paths belong to the upstream project and must stay intact.
PROTECTED_MAP=(
  "docs.fastapi.ai|docs.newapi.ai"
  "docs.fastapi.pro|docs.newapi.pro"
  "doc.fastapi.pro|doc.newapi.pro"
  "www.fastapi.ai|www.newapi.ai"
  "fastapi.ai|newapi.ai"
  "fastapi.pro|newapi.pro"
  "fastapi.com|newapi.com"
  "/api/fastapi/models.json|/api/newapi/models.json"
  "/api/fastapi/vendors.json|/api/newapi/vendors.json"
  "/llm-metadata/api/fastapi/|/llm-metadata/api/newapi/"
)

while IFS= read -r -d '' file; do
  for pair in "${PROTECTED_MAP[@]}"; do
    CHANGED="${pair%%|*}"
    ORIGINAL="${pair##*|}"
    perl -0pi -e "s/\Q${CHANGED}\E/${ORIGINAL}/g" "$file" 2>/dev/null || true
  done
done < <(git grep -FI -l -z \
  -e "docs.fastapi.ai" \
  -e "docs.fastapi.pro" \
  -e "doc.fastapi.pro" \
  -e "www.fastapi.ai" \
  -e "fastapi.ai" \
  -e "fastapi.pro" \
  -e "fastapi.com" \
  -e "/api/fastapi/models.json" \
  -e "/api/fastapi/vendors.json" \
  -e "/llm-metadata/api/fastapi/" \
  -- . || true)

echo "Rebrand complete."
