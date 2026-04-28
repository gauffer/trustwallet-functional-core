# Wallet Service

HTTP-обёртка над [Trust Wallet Core](https://github.com/trustwallet/wallet-core) v4.6.6 (CGo).
Деривация адресов, валидация, подпись EIP-1559 офлайн. **Не бродкастит** (это делает [Laravel](https://github.com/gauffer/simple-ledger)).

## E2E на Sepolia проверено

ETH и USDC, входящие и исходящие, всё проверено вручную на живой сети. Транзакции реально попали в блокчейн.

| | Результат | Блок | Tx |
|---|---|---|---|
| ETH входящий | 0.05 ETH от Google Cloud Faucet, проиндексирован | 10748539 | |
| USDC входящий | 20 USDC от Circle Faucet, проиндексирован | 10782891 | [`0x34bc36...`](https://sepolia.etherscan.io/tx/0x34bc36675ba7d1c7715d285641dd1651f6b36aa6fd189babdcffaae136ebdd0c) |
| ETH исходящий | 0.001 ETH, BROADCASTED, подтверждён | 10782874 | [`0xc67165...`](https://sepolia.etherscan.io/tx/0xc671651942320889f7602d717d8e63166479317850f2114b760398e693c08068) |
| USDC исходящий | 1 USDC, BROADCASTED, подтверждён | 10782900 | [`0xa67706...`](https://sepolia.etherscan.io/tx/0xa6770690caad7ff3ab5ff3e1862d2f33e4b0fc237595ae313d8299106d445533) |

Офлайн-подпись через этот сервис, `signed_tx` всегда `0x02f8...` (EIP-2718 type 2), бродкаст через Laravel RPC. Ни разу не дал невалидной подписи.

## Архитектура

```
cmd/wallet/main.go       точка входа, только wiring (~25 строк)
internal/handler/        HTTP, JSON, маппинг ошибок (closures, не методы)
internal/wallet/         чистые функции: DeriveAddress, ValidateAddress, SignTransaction
internal/twcore/         единственный пакет с CGo; wallet.go, address.go, signer.go
wallet-core 4.6.6        статические архивы (.a): крипто, BIP-39, EIP-1559
```

Паттерн: functional core with imperative shell.

Типы полностью анемичные. 

`internal/wallet` чистое ядро, никакого IO. 

`internal/twcore` это единственное место с `import "C"`.

## Запуск

```bash
# Один раз: собрать wallet-core и положить в локальный Docker-реестр (~10-15 мин)
./scripts/bootstrap-wallet-core.sh

# Запуск
export WALLET_MNEMONIC_ETHEREUM="twelve words ..."
docker compose -f build/package/docker-compose.yml up -d    # порт 8000

# Быстрая итерация (бинарник на хосте, без пересборки Docker)
./scripts/build.sh && ./wallet-service
```

`config.json` закоммичен, мнемоника в нём всегда перезаписывается из `WALLET_MNEMONIC_<GATE_NAME_UPPERCASED>`.

## Эндпоинты

```
POST /api/v1/createaddress     {"gate":"ethereum","account":0,"change":0,"address_index":0}
POST /api/v1/validateaddress   {"gate":"ethereum","address":"0x..."}
POST /api/v1/tx                {"gate":"ethereum","account":0,"change":0,"address_index":0,"tx_params":{...}}
```

`signed_tx` всегда начинается с `0x02f8` (EIP-2718 type 2).
