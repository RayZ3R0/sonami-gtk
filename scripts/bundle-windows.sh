#!/usr/bin/env bash
# bundle-windows.sh — Bundles tonearm.exe + all required DLLs/assets into a portable ZIP.
# Must be run inside an MSYS2 UCRT64 shell after `go build` has produced tonearm.exe.
#
# Usage: bash scripts/bundle-windows.sh <TAG>
#   e.g. bash scripts/bundle-windows.sh v1.5.0

set -euo pipefail

TAG="${1:?Usage: bundle-windows.sh <tag>}"
BUNDLE_DIR="tonearm-${TAG}-windows-x86_64"
MSYS2_PREFIX="/ucrt64"
GST_PLUGIN_PATH="${MSYS2_PREFIX}/lib/gstreamer-1.0"

echo "==> Creating bundle directory: ${BUNDLE_DIR}"
rm -rf "${BUNDLE_DIR}"
mkdir -p "${BUNDLE_DIR}/lib/gstreamer-1.0"
mkdir -p "${BUNDLE_DIR}/share/icons/hicolor/scalable/apps"
mkdir -p "${BUNDLE_DIR}/share/icons/hicolor/128x128/apps"
mkdir -p "${BUNDLE_DIR}/share/icons/hicolor/symbolic/apps"
mkdir -p "${BUNDLE_DIR}/share/glib-2.0/schemas"
mkdir -p "${BUNDLE_DIR}/share/applications"
mkdir -p "${BUNDLE_DIR}/share/metainfo"
mkdir -p "${BUNDLE_DIR}/lib/gdk-pixbuf-2.0"

# ── 1. Main executable ────────────────────────────────────────────────────────
echo "==> Copying executable"
cp tonearm.exe "${BUNDLE_DIR}/"

# ── 2. Collect required DLLs via ldd ──────────────────────────────────────────
echo "==> Collecting DLLs (ldd)"
collect_dlls() {
    local exe="$1"
    ldd "$exe" 2>/dev/null \
        | grep "${MSYS2_PREFIX}" \
        | awk '{print $3}' \
        | sort -u
}

DLLS=$(collect_dlls "${BUNDLE_DIR}/tonearm.exe")
for dll in $DLLS; do
    [ -f "$dll" ] && cp -n "$dll" "${BUNDLE_DIR}/"
done

# ── 3. GStreamer plugins ──────────────────────────────────────────────────────
echo "==> Copying GStreamer plugins"
if [ -d "${GST_PLUGIN_PATH}" ]; then
    # Core + base + good + bad plugins needed for HTTP streaming + audio
    for plugin in \
        libgstsouphttpsrc.dll \
        libgstaudioconvert.dll \
        libgstaudioresample.dll \
        libgstcoreelements.dll \
        libgstplayback.dll \
        libgstautodetect.dll \
        libgstvolume.dll \
        libgstopusparse.dll \
        libgstflac.dll \
        libgstmpg123.dll \
        libgstaiff.dll \
        libgstwasapi.dll \
        libgstwasapi2.dll \
        libgstdirectsound.dll \
        libgstisomp4.dll \
        libgstmatroska.dll \
        libgstid3demux.dll \
        libgstreplaygain.dll \
    ; do
        src="${GST_PLUGIN_PATH}/${plugin}"
        [ -f "$src" ] && cp "$src" "${BUNDLE_DIR}/lib/gstreamer-1.0/"
    done
    # Also collect DLL deps of each GStreamer plugin
    for plugin in "${BUNDLE_DIR}/lib/gstreamer-1.0/"*.dll; do
        for dll in $(collect_dlls "$plugin"); do
            [ -f "$dll" ] && cp -n "$dll" "${BUNDLE_DIR}/"
        done
    done
fi

# ── 4. GLib schemas ──────────────────────────────────────────────────────────
echo "==> Compiling GSettings schemas"
cp internal/settings/dev.dergs.Tonearm.gschema.xml \
    "${BUNDLE_DIR}/share/glib-2.0/schemas/"
# Copy upstream compiled schemas for GTK/Adwaita
if [ -d "${MSYS2_PREFIX}/share/glib-2.0/schemas" ]; then
    cp "${MSYS2_PREFIX}/share/glib-2.0/schemas/gschemas.compiled" \
       "${BUNDLE_DIR}/share/glib-2.0/schemas/" 2>/dev/null || true
fi
glib-compile-schemas "${BUNDLE_DIR}/share/glib-2.0/schemas/"

# ── 5. Icons ──────────────────────────────────────────────────────────────────
echo "==> Copying icons"
cp internal/icons/hicolor/scalable/apps/dev.dergs.Tonearm.svg \
    "${BUNDLE_DIR}/share/icons/hicolor/scalable/apps/"
cp internal/icons/hicolor/128x128/apps/dev.dergs.Tonearm.png \
    "${BUNDLE_DIR}/share/icons/hicolor/128x128/apps/"
cp internal/icons/hicolor/symbolic/apps/dev.dergs.Tonearm-symbolic.svg \
    "${BUNDLE_DIR}/share/icons/hicolor/symbolic/apps/"

# Copy the Adwaita icon theme (GTK4 needs it on Windows)
if [ -d "${MSYS2_PREFIX}/share/icons/Adwaita" ]; then
    cp -r "${MSYS2_PREFIX}/share/icons/Adwaita" \
        "${BUNDLE_DIR}/share/icons/"
fi

# ── 6. Desktop / Metainfo ────────────────────────────────────────────────────
echo "==> Copying desktop/metainfo"
cp build/dev.dergs.Tonearm.desktop  "${BUNDLE_DIR}/share/applications/"
cp build/dev.dergs.Tonearm.metainfo.xml "${BUNDLE_DIR}/share/metainfo/"

# ── 7. GDK pixbuf loaders (needed for PNG/SVG rendering) ────────────────────
echo "==> Copying GDK pixbuf loaders"
if [ -d "${MSYS2_PREFIX}/lib/gdk-pixbuf-2.0" ]; then
    cp -r "${MSYS2_PREFIX}/lib/gdk-pixbuf-2.0" "${BUNDLE_DIR}/lib/"
    # Collect DLL deps of the loaders too
    for loader in $(find "${BUNDLE_DIR}/lib/gdk-pixbuf-2.0" -name "*.dll" 2>/dev/null); do
        for dll in $(collect_dlls "$loader"); do
            [ -f "$dll" ] && cp -n "$dll" "${BUNDLE_DIR}/"
        done
    done
fi

# ── 8. GTK4 runtime data ─────────────────────────────────────────────────────
echo "==> Copying GTK4 runtime data"
for share_dir in gtk-4.0 themes; do
    if [ -d "${MSYS2_PREFIX}/share/${share_dir}" ]; then
        cp -r "${MSYS2_PREFIX}/share/${share_dir}" "${BUNDLE_DIR}/share/" 2>/dev/null || true
    fi
done

# ── 9. Environment launcher script ──────────────────────────────────────────
echo "==> Writing launcher"
cat > "${BUNDLE_DIR}/tonearm-launch.bat" << 'BATCH'
@echo off
:: Tonearm launcher — sets required env vars then starts the app
set "SCRIPT_DIR=%~dp0"
set "GST_PLUGIN_PATH=%SCRIPT_DIR%lib\gstreamer-1.0"
set "GDK_PIXBUF_MODULEDIR=%SCRIPT_DIR%lib\gdk-pixbuf-2.0\2.10.0\loaders"
set "GDK_PIXBUF_MODULE_FILE=%SCRIPT_DIR%lib\gdk-pixbuf-2.0\2.10.0\loaders.cache"
set "GSETTINGS_SCHEMA_DIR=%SCRIPT_DIR%share\glib-2.0\schemas"
set "GTK_DATA_PREFIX=%SCRIPT_DIR%"
set "XDG_DATA_DIRS=%SCRIPT_DIR%share"
start "" "%SCRIPT_DIR%tonearm.exe" %*
BATCH

# ── 10. ZIP it up ────────────────────────────────────────────────────────────
echo "==> Creating ZIP archive"
ZIP_NAME="${BUNDLE_DIR}.zip"
if command -v powershell &>/dev/null; then
    powershell -Command "Compress-Archive -Path '${BUNDLE_DIR}' -DestinationPath '${ZIP_NAME}' -Force"
else
    7z a -tzip "${ZIP_NAME}" "${BUNDLE_DIR}" >/dev/null
fi

echo "==> Done: ${ZIP_NAME}"
