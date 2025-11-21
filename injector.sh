#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(pwd)"
OUT_DIR="${ROOT_DIR}/provider"
OUT_FILE="${OUT_DIR}/generated_providers.go"
SEARCH_DIRS=(config middleware controller services repository)

TMP_IMPORTS="$(mktemp)"
TMP_CONSTRUCTORS="$(mktemp)"

cleanup() { rm -f "$TMP_IMPORTS" "$TMP_CONSTRUCTORS"; }
trap cleanup EXIT

# module name from go.mod
MODULE="$(sed -n 's/^module //p' go.mod | tr -d '\r' || true)"

mkdir -p "$OUT_DIR"

echo "// Code generated; DO NOT EDIT." > "$OUT_FILE"
echo "// $(date -u +"%Y-%m-%dT%H:%M:%SZ")" >> "$OUT_FILE"
echo >> "$OUT_FILE"
echo "package provider" >> "$OUT_FILE"
echo >> "$OUT_FILE"

echo "import (" > "$TMP_IMPORTS"
echo "    \"fmt\"" >> "$TMP_IMPORTS"
echo >> "$TMP_IMPORTS"

found_any=0

for d in "${SEARCH_DIRS[@]}"; do
  if [ ! -d "${ROOT_DIR}/${d}" ]; then continue; fi

  pkgname=$(grep -h '^package ' "${ROOT_DIR}/${d}"/*.go 2>/dev/null | head -n1 | awk '{print $2}')
  [ -z "$pkgname" ] && pkgname="${d}"

  alias="${pkgname}_${d}"
  [ -n "$MODULE" ] && import_path="${MODULE}/${d}" || import_path="./${d}"

  for file in $(find "${ROOT_DIR}/${d}" -maxdepth 1 -name '*.go'); do
    grep -E "^func[[:space:]]+New" "$file" | while read -r line; do
      fname=$(echo "$line" | sed -n 's/^func[[:space:]]\+\(New[A-Za-z0-9_]*\).*/\1/p')
      [ -z "$fname" ] && continue
      params=$(echo "$line" | sed -n 's/^[^(]*(\([^)]*\)).*$/\1/p' | tr -d '[:space:]')

      echo "${d}|${pkgname}|${alias}|${import_path}|${fname}|${params}|${file}" >> "$TMP_CONSTRUCTORS"
      found_any=1
    done
  done
done

if [ "$found_any" -eq 1 ]; then
  awk -F'|' '{print $3" \""$4"\""}' "$TMP_CONSTRUCTORS" | sort -u | while read -r imp; do
    echo "    ${imp}" >> "$TMP_IMPORTS"
  done
fi

echo ")" >> "$TMP_IMPORTS"
cat "$TMP_IMPORTS" >> "$OUT_FILE"
echo >> "$OUT_FILE"

cat >> "$OUT_FILE" <<EOF
func ProvideAll() map[string]interface{} {
	out := make(map[string]interface{})
EOF

if [ "$found_any" -eq 1 ]; then
  awk -F'|' '{print $2 "|" $3 "|" $5 "|" $6 "|" $7}' "$TMP_CONSTRUCTORS" |
  while IFS='|' read -r pkg alias fname params file; do
    key="${pkg}.${fname}"
    call="${alias}.${fname}"
    if [ -z "$params" ]; then
      echo "	// $file" >> "$OUT_FILE"
      echo "	__${fname} := ${call}()" >> "$OUT_FILE"
      echo "	out[\"${key}\"] = __${fname}" >> "$OUT_FILE"
      echo >> "$OUT_FILE"
    else
      echo "	// TODO: ${call}(${params}) requires params" >> "$OUT_FILE"
      echo "	// out[\"${key}\"] = ${call}(...)" >> "$OUT_FILE"
      echo >> "$OUT_FILE"
    fi
  done
else
  echo "	// No constructors found" >> "$OUT_FILE"
fi

echo "	return out" >> "$OUT_FILE"
echo "}" >> "$OUT_FILE"

command -v gofmt >/dev/null && gofmt -w "$OUT_FILE"
echo "✅ Generated at ${OUT_FILE}"
