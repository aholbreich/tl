#!/usr/bin/env sh
set -eu

REPO="${TL_INSTALL_REPO:-aholbreich/tl}"
VERSION="${TL_INSTALL_VERSION:-latest}"
BIN_DIR="${TL_INSTALL_BIN_DIR:-}"

usage() {
  cat <<EOF
Install tl from GitHub release binaries.

Usage:
  install.sh [--version VERSION] [--bin-dir DIR]

Options:
  --version VERSION  Release version to install, for example 0.4.4 (default: latest)
  --bin-dir DIR      Directory to install tl into (default: /usr/local/bin if writable, otherwise \$HOME/.local/bin)
  -h, --help         Show this help

Environment:
  TL_INSTALL_VERSION  Same as --version
  TL_INSTALL_BIN_DIR  Same as --bin-dir
  TL_INSTALL_REPO     GitHub repo to install from (default: aholbreich/tl)
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      if [ "$#" -lt 2 ]; then
        echo "error: --version requires a value" >&2
        exit 2
      fi
      VERSION="$2"
      shift 2
      ;;
    --version=*)
      VERSION="${1#--version=}"
      shift
      ;;
    --bin-dir)
      if [ "$#" -lt 2 ]; then
        echo "error: --bin-dir requires a value" >&2
        exit 2
      fi
      BIN_DIR="$2"
      shift 2
      ;;
    --bin-dir=*)
      BIN_DIR="${1#--bin-dir=}"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "error: unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

case "$(uname -s)" in
  Linux) os="linux" ;;
  Darwin) os="darwin" ;;
  *)
    echo "error: unsupported OS: $(uname -s)" >&2
    exit 1
    ;;
esac

case "$(uname -m)" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "error: unsupported architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

if [ -z "${BIN_DIR}" ]; then
  if [ -d /usr/local/bin ] && [ -w /usr/local/bin ]; then
    BIN_DIR="/usr/local/bin"
  else
    BIN_DIR="${HOME}/.local/bin"
  fi
fi

asset="tl-${os}-${arch}.tar.gz"
if [ "${VERSION}" = "latest" ]; then
  url="https://github.com/${REPO}/releases/latest/download/${asset}"
else
  url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"
fi

if command -v mktemp >/dev/null 2>&1; then
  tmp_dir="$(mktemp -d)"
else
  tmp_dir="/tmp/tl-install.$$"
  mkdir -p "${tmp_dir}"
fi
trap 'rm -rf "${tmp_dir}"' EXIT HUP INT TERM

archive="${tmp_dir}/${asset}"

echo "Downloading ${url}"
if command -v curl >/dev/null 2>&1; then
  curl -fsSL -o "${archive}" "${url}"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "${archive}" "${url}"
else
  echo "error: curl or wget is required" >&2
  exit 1
fi

if ! command -v tar >/dev/null 2>&1; then
  echo "error: tar is required" >&2
  exit 1
fi

tar -xzf "${archive}" -C "${tmp_dir}"
if [ ! -f "${tmp_dir}/tl" ]; then
  echo "error: archive did not contain tl binary" >&2
  exit 1
fi

mkdir -p "${BIN_DIR}"
if command -v install >/dev/null 2>&1; then
  install -m 0755 "${tmp_dir}/tl" "${BIN_DIR}/tl"
else
  cp "${tmp_dir}/tl" "${BIN_DIR}/tl"
  chmod 0755 "${BIN_DIR}/tl"
fi

echo "Installed tl to ${BIN_DIR}/tl"
"${BIN_DIR}/tl" --version

case ":${PATH}:" in
  *":${BIN_DIR}:"*) ;;
  *) echo "Note: ${BIN_DIR} is not on PATH" ;;
esac
