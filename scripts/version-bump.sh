#!/usr/bin/env bash
# version-bump.sh — Bumps the version in flake.nix and metainfo.xml,
# then creates a signed git tag ready to trigger the release workflow.
#
# Usage: bash scripts/version-bump.sh <new-version>
#   e.g. bash scripts/version-bump.sh 1.5.0

set -euo pipefail

##############################################################################
# Helpers
##############################################################################
die() { echo "ERROR: $*" >&2; exit 1; }
confirm() {
    read -r -p "$* [y/N] " answer
    [[ "${answer,,}" == "y" ]] || die "Aborted."
}

##############################################################################
# Validate input
##############################################################################
NEW_VERSION="${1:-}"
[[ -n "$NEW_VERSION" ]] || die "Usage: $0 <version>  (e.g. 1.5.0)"

# Strip leading 'v' if supplied
NEW_VERSION="${NEW_VERSION#v}"

# Must look like X.Y.Z or X.Y.Z-suffix
[[ "$NEW_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9._-]+)?$ ]] \
    || die "Version must follow X.Y.Z format (got: ${NEW_VERSION})"

FLAKE="flake.nix"
METAINFO="build/io.github.rayz3r0.SonamiGtk.metainfo.xml"
TODAY="$(date +%Y-%m-%d)"

cd "$(git rev-parse --show-toplevel)"

##############################################################################
# Show current state
##############################################################################
CURRENT_FLAKE=$(grep -oP 'version = "\K[^"]+' "$FLAKE" | head -1)
echo "──────────────────────────────────────────────"
echo "  Current version (flake.nix):  ${CURRENT_FLAKE}"
echo "  New version:                  ${NEW_VERSION}"
echo "  Release date:                 ${TODAY}"
echo "──────────────────────────────────────────────"
confirm "Proceed?"

##############################################################################
# 1. Update flake.nix
##############################################################################
echo "==> Updating ${FLAKE}"
sed -i "s/version = \"${CURRENT_FLAKE}\"/version = \"${NEW_VERSION}\"/" "$FLAKE"

# Verify the change landed
grep -q "version = \"${NEW_VERSION}\"" "$FLAKE" \
    || die "Failed to update version in ${FLAKE}"

##############################################################################
# 2. Update metainfo.xml — prepend a new <release> entry
##############################################################################
echo "==> Updating ${METAINFO}"

NEW_ENTRY="    <release version=\"${NEW_VERSION}\" date=\"${TODAY}\">
        <url type=\"details\">https://github.com/RayZ3R0/sonami-gtk/releases/tag/v${NEW_VERSION}</url>
        <description>
            <p>Changes in this release.</p>
        </description>
    </release>"

# Insert right after the <releases> opening tag
python3 - <<PYEOF
import re, sys

path = "${METAINFO}"
new_entry = """${NEW_ENTRY}"""

with open(path, "r") as f:
    content = f.read()

# Check that this version doesn't already exist
if 'version="${NEW_VERSION}"' in content:
    print("WARNING: version ${NEW_VERSION} already exists in metainfo — skipping insertion.")
    sys.exit(0)

# Insert after <releases>
updated = content.replace(
    "<releases>",
    "<releases>\n" + new_entry,
    1
)

if updated == content:
    print("ERROR: Could not find <releases> tag in metainfo.xml", file=sys.stderr)
    sys.exit(1)

with open(path, "w") as f:
    f.write(updated)

print("Inserted release entry for ${NEW_VERSION} dated ${TODAY}")
PYEOF

##############################################################################
# 3. Show diff
##############################################################################
echo ""
echo "==> Diff:"
git diff "$FLAKE" "$METAINFO"

##############################################################################
# 4. Commit + tag
##############################################################################
echo ""
confirm "Commit these changes and create tag v${NEW_VERSION}?"

git add "$FLAKE" "$METAINFO"
git commit -m "chore: release v${NEW_VERSION}"
git tag -a "v${NEW_VERSION}" -m "Release v${NEW_VERSION}"

echo ""
echo "✅  Done! Next steps:"
echo "    git push origin main"
echo "    git push origin v${NEW_VERSION}"
echo ""
echo "    Pushing the tag will trigger the GitHub Actions release workflow."
