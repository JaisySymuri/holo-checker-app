#!/usr/bin/env bash
set -euo pipefail

APP_NAME="holo-checker-app-linux-arm64"
SERVICE_NAME="holo-checker"
INSTALL_DIR="/opt/holo-checker"
OWNER_USER="opc"
OWNER_GROUP="opc"

echo "Building ${APP_NAME}..."
GOOS=linux GOARCH=arm64 go build -o "${APP_NAME}" .

echo "Installing binary to ${INSTALL_DIR}..."
sudo mkdir -p "${INSTALL_DIR}"
sudo cp "${APP_NAME}" "${INSTALL_DIR}/"
sudo chown "${OWNER_USER}:${OWNER_GROUP}" "${INSTALL_DIR}/${APP_NAME}"
sudo chmod 755 "${INSTALL_DIR}/${APP_NAME}"

if [[ -f .env ]]; then
  echo "Updating .env..."
  sudo cp .env "${INSTALL_DIR}/.env"
  sudo chown "${OWNER_USER}:${OWNER_GROUP}" "${INSTALL_DIR}/.env"
  sudo chmod 600 "${INSTALL_DIR}/.env"
fi

echo "Restarting service ${SERVICE_NAME}..."
sudo systemctl restart "${SERVICE_NAME}"

echo "Service status:"
sudo systemctl --no-pager --full status "${SERVICE_NAME}"

echo "Done."

## Save it, then make it executable once:

# chmod +x deploy-linux.sh


## After that, deploy is just:

# ./deploy-linux.sh