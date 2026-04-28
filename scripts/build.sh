#!/usr/bin/env bash
set -euo pipefail

# Собирает бинарник wallet-service внутри контейнера wallet-core, используя Go с хоста.
# Быстрее чем `docker build`: переиспользует кэш модулей хоста.

IMAGE="localhost:5000/wallet-core:4.6.6"
GO_ROOT="$(go env GOROOT)"
GO_MOD_CACHE="$(go env GOMODCACHE)"
CALLER_UID="$(id -u)"
CALLER_GID="$(id -g)"

echo "Building wallet-service inside $IMAGE..."
echo "Using host Go: $(go version)"

docker run --rm \
  -v "$(pwd)":/app \
  -v "$GO_ROOT":/usr/local/go:ro \
  -v "$GO_MOD_CACHE":/go/pkg/mod:ro \
  -w /app \
  -e PATH="/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
  -e GOPATH="/go" \
  -e GOMODCACHE="/go/pkg/mod" \
  "$IMAGE" \
  bash -c "
    CGO_ENABLED=1 \
    CC=clang-14 \
    CXX=clang++-14 \
    CGO_CFLAGS='-I/wallet-core/include' \
    CGO_LDFLAGS='-L/wallet-core/build -L/wallet-core/build/local/lib -L/wallet-core/build/trezor-crypto -lTrustWalletCore -lwallet_core_rs -lprotobuf -lTrezorCrypto -lstdc++ -lm' \
    go build -o wallet-service ./cmd/wallet && \
    chown $CALLER_UID:$CALLER_GID wallet-service
  "

echo "Done: ./wallet-service ($(du -h wallet-service | cut -f1))"
