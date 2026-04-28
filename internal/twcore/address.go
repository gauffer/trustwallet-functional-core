package twcore

// #include "tw_compat.h"
// #include <TrustWalletCore/TWAnyAddress.h>
// #include <TrustWalletCore/TWCoinType.h>
import "C"

func ValidateAddress(address string) bool {
	str := TWStringCreate(address)
	defer TWStringDelete(str)
	return bool(C.TWAnyAddressIsValid(str, C.enum_TWCoinType(CoinTypeEthereum)))
}
