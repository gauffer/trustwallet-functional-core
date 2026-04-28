#!/usr/bin/env bash
# Однократная настройка: собирает wallet-core 4.6.6 из исходников и пушит в локальный Docker-реестр.
# Docker Hub не поддерживает trustwallet/wallet-core с 2021 года.
set -euo pipefail

WALLET_CORE_TAG="4.6.6"
WALLET_CORE_REPO="https://github.com/trustwallet/wallet-core.git"
CHECKOUT_DIR="${HOME}/c/third_party/trustwallet/wallet-core"
REGISTRY="localhost:5000"

echo "==> Starting local Docker registry on $REGISTRY"
if ! docker ps --format '{{.Names}}' | grep -q '^registry$'; then
  docker run -d --restart=always --name registry -p 5000:5000 registry:2
fi

echo "==> Cloning wallet-core (if not present)"
if [ ! -d "$CHECKOUT_DIR" ]; then
  mkdir -p "$(dirname "$CHECKOUT_DIR")"
  git clone "$WALLET_CORE_REPO" "$CHECKOUT_DIR"
fi

echo "==> Checking out $WALLET_CORE_TAG"
cd "$CHECKOUT_DIR"
git fetch --tags
git checkout "$WALLET_CORE_TAG"

echo "==> Building wallet-core image (this takes ~10-15 minutes)"
docker build -t "$REGISTRY/wallet-core:$WALLET_CORE_TAG" -t "$REGISTRY/wallet-core:latest" .

echo "==> Pushing to local registry"
docker push "$REGISTRY/wallet-core:$WALLET_CORE_TAG"
docker push "$REGISTRY/wallet-core:latest"

echo "==> Done. wallet-core $WALLET_CORE_TAG available at $REGISTRY/wallet-core:$WALLET_CORE_TAG"
