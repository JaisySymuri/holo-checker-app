#!/usr/bin/env bash
set -euo pipefail

APP_NAME="holo-checker-app-linux-arm64"
SERVICE_NAME="holo-checker"
INSTALL_DIR="/opt/holo-checker"
OWNER_USER="opc"
OWNER_GROUP="opc"

echo "Building ${APP_NAME}..."
GOOS=linux GOARCH=arm64 go build -o "${APP_NAME}" .

echo "Stopping service ${SERVICE_NAME} (if running)..."
if systemctl is-active --quiet "${SERVICE_NAME}"; then
  sudo systemctl stop "${SERVICE_NAME}"
else
  echo "Service not running, skip stop."
fi

echo "Installing binary to ${INSTALL_DIR}..."
sudo mkdir -p "${INSTALL_DIR}"

# Copy as temp file (atomic strategy)
sudo cp "${APP_NAME}" "${INSTALL_DIR}/${APP_NAME}.new"
sudo chown "${OWNER_USER}:${OWNER_GROUP}" "${INSTALL_DIR}/${APP_NAME}.new"
sudo chmod 755 "${INSTALL_DIR}/${APP_NAME}.new"

# Replace atomically
sudo mv -f "${INSTALL_DIR}/${APP_NAME}.new" "${INSTALL_DIR}/${APP_NAME}"

if [[ -f .env ]]; then
  echo "Updating .env..."
  sudo cp .env "${INSTALL_DIR}/.env"
  sudo chown "${OWNER_USER}:${OWNER_GROUP}" "${INSTALL_DIR}/.env"
  sudo chmod 600 "${INSTALL_DIR}/.env"
fi

echo "Reloading systemd daemon (just in case)..."
sudo systemctl daemon-reload

echo "Starting service ${SERVICE_NAME}..."
sudo systemctl start "${SERVICE_NAME}"

echo "Checking service health..."
if systemctl is-active --quiet "${SERVICE_NAME}"; then
  echo "✅ Service is running"
else
  echo "❌ Service failed to start"
  sudo journalctl -u "${SERVICE_NAME}" -n 50 --no-pager
  exit 1
fi

echo "Service status:"
sudo systemctl --no-pager --full status "${SERVICE_NAME}"

echo "Done."

## Save it, then make it executable once:

# chmod +x deploy-linux.sh


## After that, deploy is just:

# ./deploy-linux.sh