#!/usr/bin/env bash
set -eo pipefail

source /dev/stdin <<< "$(curl -Lsm 2 https://get2.belt.sh)"

ARCHIVE_FILE="/tmp/toolshed.tar.gz"

main() {
  local host="$1"

  [[ -z "$host" ]] && echo "Error: missing argument host" && exit 1
  [[ ! -e "$ARCHIVE_FILE" ]] && echo "Error: missing archive $ARCHIVE_FILE" && exit 1

  belt_begin_session "root" "$host"

  echo "Uploading app archive"
  app_upload "$ARCHIVE_FILE"

  echo "Stopping process"
  systemd_unit_stop "toolshed"

  echo "Setting up user"
  user_add "toolshed"

  echo "Copying files"
  app_copy_file "dist/toolshed"

  echo "Setting up permissions"
  app_set_permissions "toolshed:toolshed"

  echo "Copying unit file"
  systemd_add_unit "toolshed.service"

  echo "Starting process"
  systemd_unit_start "toolshed"

  echo "Adding caddy vhost"
  app_add_caddy_vhost

  echo "Restarting caddy"
  caddy_restart
}

echo "Deploying toolshed"
main "$1"
