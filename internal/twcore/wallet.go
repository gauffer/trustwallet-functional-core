package twcore

// #include "tw_compat.h"
// #include <TrustWalletCore/TWHDWallet.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWAnyAddress.h>
// #include <TrustWalletCore/TWCoinType.h>
// #include <TrustWalletCore/TWMnemonic.h>
import "C"

import (
	"fmt"
	"unsafe"
)

const CoinTypeEthereum = C.TWCoinTypeEthereum

func IsMnemonicValid(mnemonic string) bool {
	str := TWStringCreate(mnemonic)
	defer TWStringDelete(str)
	return bool(C.TWMnemonicIsValid(str))
}

// DeriveKey возвращает байты приватного ключа для пути BIP-44 m/44'/60'/account'/change/addressIndex.
func DeriveKey(mnemonic string, account, change, addressIndex uint32) ([]byte, error) {
	if !IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	mnStr := TWStringCreate(mnemonic)
	defer TWStringDelete(mnStr)
	empty := TWStringCreate("")
	defer TWStringDelete(empty)

	wallet := C.TWHDWalletCreateWithMnemonic(mnStr, empty)
	if wallet == nil {
		return nil, fmt.Errorf("failed to create HD wallet")
	}
	defer C.TWHDWalletDelete(wallet)

	priKey := C.TWHDWalletGetDerivedKey(
		wallet,
		C.enum_TWCoinType(CoinTypeEthereum),
		C.uint32_t(account),
		C.uint32_t(change),
		C.uint32_t(addressIndex),
	)
	defer C.TWPrivateKeyDelete(priKey)

	priKeyData := C.TWPrivateKeyData(priKey)
	defer TWDataDelete(unsafe.Pointer(priKeyData))

	return TWDataGoBytes(unsafe.Pointer(priKeyData)), nil
}

func DeriveAddress(mnemonic string, account, change, addressIndex uint32) (string, error) {
	if !IsMnemonicValid(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic")
	}

	mnStr := TWStringCreate(mnemonic)
	defer TWStringDelete(mnStr)
	empty := TWStringCreate("")
	defer TWStringDelete(empty)

	wallet := C.TWHDWalletCreateWithMnemonic(mnStr, empty)
	if wallet == nil {
		return "", fmt.Errorf("failed to create HD wallet")
	}
	defer C.TWHDWalletDelete(wallet)

	path := fmt.Sprintf("m/44'/60'/%d'/%d/%d", account, change, addressIndex)
	pathStr := TWStringCreate(path)
	defer TWStringDelete(pathStr)

	priKey := C.TWHDWalletGetKey(
		wallet,
		C.enum_TWCoinType(CoinTypeEthereum),
		pathStr,
	)
	defer C.TWPrivateKeyDelete(priKey)

	pubKey := C.TWPrivateKeyGetPublicKeySecp256k1(priKey, false)
	defer C.TWPublicKeyDelete(pubKey)

	addr := C.TWAnyAddressCreateWithPublicKey(pubKey, C.enum_TWCoinType(CoinTypeEthereum))
	defer C.TWAnyAddressDelete(addr)

	addrStr := C.TWAnyAddressDescription(addr)
	defer TWStringDelete(unsafe.Pointer(addrStr))

	return TWStringGoString(unsafe.Pointer(addrStr)), nil
}
