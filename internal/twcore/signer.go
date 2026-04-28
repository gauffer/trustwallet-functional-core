package twcore

// #include "tw_compat.h"
// #include <TrustWalletCore/TWAnySigner.h>
// #include <TrustWalletCore/TWCoinType.h>
import "C"

import "unsafe"

func Sign(inputBytes []byte) ([]byte, error) {
	input := TWDataCreate(inputBytes)
	defer TWDataDelete(input)

	output := C.TWAnySignerSign(input, C.enum_TWCoinType(CoinTypeEthereum))
	defer TWDataDelete(unsafe.Pointer(output))

	return TWDataGoBytes(unsafe.Pointer(output)), nil
}
