#!/usr/bin/env bash
# Build the tl-bin pacman package from a published GitHub release.
#
# Unlike build-rpm.sh, which compiles from source, this packages the release
# tarball the release job has already uploaded — the same artefact AUR users
# would download. It therefore has to run after the release exists.
#
# packaging/aur/PKGBUILD is the single source of packaging truth. This script
# copies it, retargets it at VERSION and refreshes its checksums; it never
# defines the package a second time, because two definitions would drift.
set -euo pipefail

VERSION="${VERSION:-$(git describe --tags --abbrev=0 2>/dev/null || echo dev)}"
PKGREL="${PKGREL:-1}"
PKGBUILD_SRC="${PKGBUILD_SRC:-packaging/aur/PKGBUILD}"
OUT_DIR="${OUT_DIR:-dist/pacman}"
WORK_DIR="${WORK_DIR:-build/pacman}"

if [ ! -f "${PKGBUILD_SRC}" ]; then
  echo "PKGBUILD not found at ${PKGBUILD_SRC}" >&2
  exit 1
fi

rm -rf "${WORK_DIR}"
mkdir -p "${WORK_DIR}" "${OUT_DIR}"
cp "${PKGBUILD_SRC}" "${WORK_DIR}/PKGBUILD"

sed -i \
  -e "s/^pkgver=.*/pkgver=${VERSION}/" \
  -e "s/^pkgrel=.*/pkgrel=${PKGREL}/" \
  "${WORK_DIR}/PKGBUILD"

# updpkgsums refreshes every sha256sums_* array, both architectures, not just
# the build host's. Sums for the arch we are not building still ship to AUR
# users, so a partial refresh would leave them stale and unverifiable.
( cd "${WORK_DIR}" && updpkgsums )

# makepkg refuses to run as root; CI runs it as an unprivileged build user.
( cd "${WORK_DIR}" && makepkg --force --clean --noconfirm )

# Debug packages are deliberately not published: they roughly double the
# repository size and rpm-repo does not carry them either. makepkg still
# produces one when the build host enables the debug option, so filter it out
# here rather than depending on the runner's makepkg.conf.
find "${WORK_DIR}" -maxdepth 1 -name '*.pkg.tar.*' \
  ! -name '*-debug-*' -exec cp {} "${OUT_DIR}/" \;

echo "Built into ${OUT_DIR}:"
ls -1 "${OUT_DIR}"
