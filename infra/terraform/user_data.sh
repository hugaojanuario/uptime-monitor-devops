#!/bin/bash
set -euo pipefail
exec > >(tee -a /var/log/user-data.log) 2>&1

export DEBIAN_FRONTEND=noninteractive
APT="apt-get -o DPkg::Lock::Timeout=180"

REPO_URL="https://github.com/hugaojanuario/uptime-monitor-devops.git"
APP_DIR="/home/ubuntu/uptime-monitor-devops"

# --- Docker via repositorio oficial ---
$APT update
$APT install -y ca-certificates curl git

install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" > /etc/apt/sources.list.d/docker.list

$APT update
$APT install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

usermod -aG docker ubuntu

# --- codigo ---
sudo -u ubuntu git clone "$REPO_URL" "$APP_DIR"

# --- credenciais geradas no boot, nunca versionadas ---
umask 077
cat > "$APP_DIR/.env" <<EOF
DB_NAME=uptime
DB_USER=uptime
DB_PASSWORD=$(openssl rand -base64 24 | tr -d '/+=' | head -c 24)
EOF
chown ubuntu:ubuntu "$APP_DIR/.env"
chmod 600 "$APP_DIR/.env"

# --- sobe a stack a partir da imagem pronta ---
cd "$APP_DIR"
docker compose -f docker-compose.prod.yml up -d
