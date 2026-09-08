#!/usr/bin/env bash
# Add the built package to a checkout of the pacman repository and refresh its
# database. Does not commit or push; the caller decides that.
#
# Usage: REPO_DIR=path/to/pacman-repo .github/scripts/update-pacman-repo.sh
set -euo pipefail

REPO_DIR="${REPO_DIR:?REPO_DIR must point at a checkout of the pacman repository}"
PKG_DIR="${PKG_DIR:-dist/pacman}"
DB_NAME="${DB_NAME:-holbreich}"
ARCH="${ARCH:-x86_64}"

ARCH_DIR="${REPO_DIR}/${ARCH}"
mkdir -p "${ARCH_DIR}"

shopt -s nullglob
packages=("${PKG_DIR}"/*.pkg.tar.*)
shopt -u nullglob
if [ ${#packages[@]} -eq 0 ]; then
  echo "no packages found in ${PKG_DIR}" >&2
  exit 1
fi

# makepkg stamps a build date into the package, so rebuilding the same tag
# yields a byte-different file for identical software. Publishing it would
# change the database checksum and make every user re-download for nothing,
# and would defeat the caller's "already current" check. Re-running a release
# job is therefore a no-op unless the version actually moved; a genuine
# rebuild bumps pkgrel, which changes the entry name and passes this guard.
db="${ARCH_DIR}/${DB_NAME}.db.tar.gz"
if [ -f "${db}" ]; then
  for pkg in "${packages[@]}"; do
    entry="$(basename "${pkg}")"
    entry="${entry%-*.pkg.tar.*}"  # strip -<arch>.pkg.tar.<ext>, keeping pkgrel
    if bsdtar -tf "${db}" 2>/dev/null | grep -qx "${entry}/"; then
      echo "${entry} is already in the repository; nothing to publish"
      exit 0
    fi
  done
fi

cp "${packages[@]}" "${ARCH_DIR}/"

# repo-add updates the existing database rather than replacing it, so this
# repository can host more than one project — rpm-repo already carries two.
# --prevent-downgrade refuses to move an entry backwards, which turns a
# mistaken re-run of an older tag into an error instead of a silent
# downgrade for every user.
(
  cd "${ARCH_DIR}"
  repo-add --new --prevent-downgrade "${DB_NAME}.db.tar.gz" "$(basename "${packages[0]}")"
)

# repo-add leaves .db and .files as symlinks to the .tar.gz files. GitHub
# Pages does not resolve symlinks: it serves the git blob, which holds the
# target's *name*. pacman fetches "<db>.db" by default, so it would receive
# the string "holbreich.db.tar.gz" where a gzip archive is expected and fail
# with a corrupt-database error. Replace them with real copies.
for link in "${DB_NAME}.db" "${DB_NAME}.files"; do
  target="${ARCH_DIR}/${link}"
  if [ -L "${target}" ]; then
    resolved="$(readlink -f "${target}")"
    rm "${target}"
    cp "${resolved}" "${target}"
  fi
done

# repo-add keeps a .old backup of each database it rewrites. Useful locally,
# but the caller commits with `git add -A`, so leaving them would accumulate
# a stale copy of every database in the published repository.
rm -f "${ARCH_DIR}"/*.old

echo "Repository ${ARCH_DIR} now contains:"
ls -1 "${ARCH_DIR}"
