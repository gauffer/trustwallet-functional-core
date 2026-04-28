# Многоэтапная сборка: wallet-core (CGo) → бинарник Go → минимальный рантайм.
#   docker build -t wallet-service .
#   docker run --rm -p 8000:8000 -e WALLET_MNEMONIC_ETHEREUM="..." wallet-service

FROM localhost:5000/wallet-core:4.6.6 AS builder

ARG GO_VERSION=1.26.2
RUN apt-get update > /dev/null 2>&1 && \
    apt-get install -y --no-install-recommends wget ca-certificates > /dev/null 2>&1 && \
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 \
    CC=clang-14 \
    CXX=clang++-14 \
    CGO_CFLAGS="-I/wallet-core/include" \
    CGO_LDFLAGS="-L/wallet-core/build -L/wallet-core/build/local/lib -L/wallet-core/build/trezor-crypto -lTrustWalletCore -lwallet_core_rs -lprotobuf -lTrezorCrypto -lstdc++ -lm" \
    go build -o wallet-service ./cmd/wallet

# Минимальный рантайм: нужен libstdc++, т.к. wallet-core 4.6.6 скомпилирован GCC.
FROM ubuntu:22.04
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates libstdc++6 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/wallet-service /usr/local/bin/wallet-service
COPY config.example.json /config.json

EXPOSE 8000
ENTRYPOINT ["wallet-service"]
