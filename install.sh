#!/bin/sh

set -eu

REPO_OWNER="drackthor"
REPO_NAME="ysort"
BINARY_NAME="ysort"
API_BASE_URL="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}"

fatal() {
  printf '[ysort-install] ERROR: %s\n' "$*" >&2
  exit 1
}

is_true() {
  value=$(printf '%s' "${1:-}" | tr '[:upper:]' '[:lower:]')
  case "$value" in
    1|true|yes|y|on)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

QUIET_VALUE="${YSORT_QUIET:-0}"
ASSUME_YES_VALUE="${YSORT_YES:-0}"

log() {
  if ! is_true "$QUIET_VALUE"; then
    printf '[ysort-install] %s\n' "$*"
  fi
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fatal "required command '$1' is not available"
}

expand_path() {
  input_path="$1"
  case "$input_path" in
    "~")
      printf '%s' "$HOME"
      ;;
    "~/"*)
      printf '%s/%s' "$HOME" "${input_path#~/}"
      ;;
    *)
      printf '%s' "$input_path"
      ;;
  esac
}

path_contains() {
  directory="$1"
  case ":$PATH:" in
    *":$directory:"*|*":$directory/:"*)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

normalize_os() {
  raw_os=$(printf '%s' "$1" | tr '[:upper:]' '[:lower:]')
  case "$raw_os" in
    linux)
      printf 'linux'
      ;;
    darwin|mac|macos|osx)
      printf 'darwin'
      ;;
    windows|win)
      printf 'windows'
      ;;
    *)
      printf ''
      ;;
  esac
}

normalize_arch() {
  raw_arch=$(printf '%s' "$1" | tr '[:upper:]' '[:lower:]')
  case "$raw_arch" in
    x86_64|x64|amd64)
      printf 'amd64'
      ;;
    aarch64|arm64)
      printf 'arm64'
      ;;
    *)
      printf ''
      ;;
  esac
}

resolve_os() {
  if [ -n "${YSORT_OS:-}" ]; then
    os_value=$(normalize_os "$YSORT_OS")
    [ -n "$os_value" ] || fatal "unsupported YSORT_OS='$YSORT_OS'"
    printf '%s' "$os_value"
    return 0
  fi

  detected=$(normalize_os "$(uname -s 2>/dev/null || true)")
  if [ -n "$detected" ]; then
    printf '%s' "$detected"
    return 0
  fi

  printf 'linux'
}

resolve_arch() {
  if [ -n "${YSORT_ARCH:-}" ]; then
    arch_value=$(normalize_arch "$YSORT_ARCH")
    [ -n "$arch_value" ] || fatal "unsupported YSORT_ARCH='$YSORT_ARCH'"
    printf '%s' "$arch_value"
    return 0
  fi

  detected=$(normalize_arch "$(uname -m 2>/dev/null || true)")
  if [ -n "$detected" ]; then
    printf '%s' "$detected"
    return 0
  fi

  printf 'amd64'
}

resolve_install_dir() {
  if [ -n "${YSORT_INSTALL_DIR:-}" ]; then
    printf '%s' "$(expand_path "$YSORT_INSTALL_DIR")"
    return 0
  fi

  for candidate in "$HOME/.local/bin" "$HOME/bin" "/usr/local/bin"; do
    if [ -d "$candidate" ] && path_contains "$candidate"; then
      printf '%s' "$candidate"
      return 0
    fi
  done

  return 1
}

api_get() {
  curl -fsSL \
    -H 'Accept: application/vnd.github+json' \
    -H 'X-GitHub-Api-Version: 2022-11-28' \
    "$1"
}

resolve_release_json() {
  if [ -n "${YSORT_VERSION:-}" ]; then
    requested="$YSORT_VERSION"

    second_candidate=""
    case "$requested" in
      v*)
        second_candidate="${requested#v}"
        ;;
      *)
        second_candidate="v${requested}"
        ;;
    esac

    if json=$(api_get "$API_BASE_URL/releases/tags/$requested" 2>/dev/null); then
      printf '%s' "$json"
      return 0
    fi

    if [ -n "$second_candidate" ] && [ "$second_candidate" != "$requested" ]; then
      if json=$(api_get "$API_BASE_URL/releases/tags/$second_candidate" 2>/dev/null); then
        printf '%s' "$json"
        return 0
      fi
    fi

    fatal "version '$YSORT_VERSION' was not found in GitHub releases"
  fi

  api_get "$API_BASE_URL/releases/latest"
}

extract_tag_name() {
  printf '%s' "$1" | jq -r '.tag_name // empty'
}

extract_asset_url() {
  printf '%s' "$1" | jq -r \
    --arg os "$2" \
    --arg arch "$3" \
    --arg bin "$BINARY_NAME" \
    'first(
      .assets[]?.browser_download_url
      | select(test("/" + $bin + "_[^/]*_" + $os + "_" + $arch + "\\.tar\\.gz$"))
    ) // empty'
}

confirm_install() {
  if is_true "$ASSUME_YES_VALUE"; then
    log "Skipping approval prompt because YSORT_YES is enabled."
    return 0
  fi

  if [ ! -r /dev/tty ]; then
    fatal "no interactive terminal for approval. Set YSORT_YES=1 to pre-approve"
  fi

  printf 'Proceed with installation? [y/N]: '
  if ! IFS= read -r answer < /dev/tty; then
    fatal "could not read confirmation input"
  fi

  normalized=$(printf '%s' "$answer" | tr '[:upper:]' '[:lower:]')
  case "$normalized" in
    y|yes)
      return 0
      ;;
    *)
      fatal "installation aborted by user"
      ;;
  esac
}

main() {
  require_cmd curl
  require_cmd tar
  require_cmd uname
  require_cmd mktemp
  require_cmd grep
  require_cmd jq

  resolved_os=$(resolve_os)
  resolved_arch=$(resolve_arch)

  if ! install_dir=$(resolve_install_dir); then
    fatal "no suitable install directory found. Create one of '$HOME/.local/bin', '$HOME/bin', '/usr/local/bin' and add it to PATH, or set YSORT_INSTALL_DIR"
  fi

  log "Resolving release metadata from GitHub API..."
  release_json=$(resolve_release_json)
  release_tag=$(extract_tag_name "$release_json")
  [ -n "$release_tag" ] || fatal "failed to parse release tag from GitHub API response: $release_json"

  release_version="$release_tag"
  case "$release_version" in
    v*)
      release_version="${release_version#v}"
      ;;
  esac

  asset_url=$(extract_asset_url "$release_json" "$resolved_os" "$resolved_arch")
  [ -n "$asset_url" ] || fatal "no release archive found for os='$resolved_os' arch='$resolved_arch' in release '$release_tag'"

  log "Calculated installation settings:"
  log "  Repository:   ${REPO_OWNER}/${REPO_NAME}"
  log "  Version tag:  $release_tag"
  log "  Version:      $release_version"
  log "  OS:           $resolved_os"
  log "  Arch:         $resolved_arch"
  log "  Install dir:  $install_dir"
  log "  Archive URL:  $asset_url"

  confirm_install

  if [ ! -d "$install_dir" ]; then
    log "Install directory does not exist. Creating: $install_dir"
    mkdir -p "$install_dir" || fatal "failed to create install directory '$install_dir'"
  fi

  [ -w "$install_dir" ] || fatal "install directory '$install_dir' is not writable"

  tmp_dir=$(mktemp -d "${TMPDIR:-/tmp}/ysort-install.XXXXXX")
  archive_path="$tmp_dir/ysort.tar.gz"

  cleanup() {
    rm -rf "$tmp_dir"
  }
  trap cleanup EXIT INT TERM

  log "Downloading release archive..."
  curl -fL "$asset_url" -o "$archive_path"

  log "Locating '$BINARY_NAME' inside archive..."
  binary_in_archive=$(tar -tzf "$archive_path" | grep -E '/ysort$|^ysort$' | head -n 1 || true)
  [ -n "$binary_in_archive" ] || fatal "could not find '$BINARY_NAME' in downloaded archive"

  log "Extracting binary from archive path: $binary_in_archive"
  tar -xzf "$archive_path" -C "$tmp_dir" "$binary_in_archive"

  source_binary="$tmp_dir/$binary_in_archive"
  [ -f "$source_binary" ] || fatal "expected extracted binary at '$source_binary'"

  target_binary="$install_dir/$BINARY_NAME"
  log "Installing binary to: $target_binary"
  mv "$source_binary" "$target_binary"
  chmod 0755 "$target_binary"

  log "Installation complete."
  log "Run '$BINARY_NAME --help' to verify the installation."
  log "Installed binary path: $target_binary"
}

main "$@"
