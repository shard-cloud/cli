#!/usr/bin/env sh
set -u

if [ "${OS:-}" = "Windows_NT" ]; then
    echo 'error: this script is not supported on Windows use `npm i -g shard-cloud-cli`'
    exit 1
fi

BOLD="$(tput bold 2>/dev/null || printf '')"
GREY="$(tput setaf 0 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
GREEN="$(tput setaf 2 2>/dev/null || printf '')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '')"
RESET="$(tput sgr0 2>/dev/null || printf '')"

REPOSITORY='shard-cloud/cli'

SYSTEM_TYPE="$(uname -s | tr '[:upper:]' '[:lower:]')"
CPU_ARCH="$(uname -m | tr '[:upper:]' '[:lower:]')"
case "${CPU_ARCH}" in
  x86_64) CPU_ARCH="amd64" ;;
  aarch64) CPU_ARCH="arm64" ;;
esac

command_available() {
  command -v "$1" 1>/dev/null 2>&1
}

print_info() {
  printf '%s\n' "${BOLD}${GREY}>${RESET} $*"
}

print_warning() {
  printf '%s\n' "${YELLOW}! $*${RESET}"
}

print_error() {
  printf '%s\n' "${RED}x $*${RESET}" >&2
}

print_success() {
  printf '%s\n' "${GREEN}$@ ${RESET}"
}

print_completed() {
  printf '%s\n' "${GREEN}✓${RESET} $*"
}

fetch_version(){
  api_endpoint="https://api.github.com/repos/${REPOSITORY}/releases/latest"

  if command_available curl; then
    curl -sL $api_endpoint | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
  elif command_available wget; then
    wget -q -O- $api_endpoint | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
  else
    print_error "Unable to find a latest release in github"
    return 1
  fi
}

LATEST_VERSION=$(fetch_version)

fetch_release() {
  output_file=$1
  release_file=$2

  download_url="https://github.com/${REPOSITORY}/releases/download/$LATEST_VERSION/${release_file}"
  
  if command_available curl; then
    curl -sL -o "${output_file}" "${download_url}"
  elif command_available wget; then
    wget -q -O "${output_file}" "${download_url}"
  else
    print_error "Unable to find a HTTP download program"
    return 1
  fi
}

extract_files() {
  archive_file=$1
  target_dir=$2

  if [ ! -d $target_dir ]; then
    mkdir -p $target_dir
  fi

  case "${archive_file}" in
    *.tar.gz | *.tgz) tar --no-same-owner -xzf "${archive_file}" -C "${target_dir}" ;;
    *.tar) tar --no-same-owner -xf "${archive_file}" -C "${target_dir}" ;;
    *.zip) unzip "${archive_file}" -d "${target_dir}" ;;
    *)
      print_error "extract_files unknown archive format for ${archive_file}"
      return 1
      ;;
  esac
}

main() {
  work_dir=$(mktemp -d)

  archive_name="shardcloud_${SYSTEM_TYPE}_${CPU_ARCH}.tar.gz"

  print_info "Installing Shard Cloud CLI, please wait..."
  fetch_release "${work_dir}/${archive_name}" $archive_name

  print_info "Unpacking ${archive_name}"
  extract_files "${work_dir}/${archive_name}" "${work_dir}/bin"

  sudo cp -f "${work_dir}/bin/shardcloud" "/usr/local/bin"
  print_completed "Successfuly installed Shard Cloud CLI $LATEST_VERSION!"

  rm -rf "${work_dir}"
}

main