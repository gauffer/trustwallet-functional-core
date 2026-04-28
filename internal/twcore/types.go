package twcore

// #cgo CFLAGS: -I/wallet-core/include
// #cgo LDFLAGS: -L/wallet-core/build -L/wallet-core/build/local/lib -L/wallet-core/build/trezor-crypto -lTrustWalletCore -lwallet_core_rs -lprotobuf -lTrezorCrypto -lstdc++ -lm
// #include "tw_compat.h"
// #include <TrustWalletCore/TWString.h>
// #include <TrustWalletCore/TWData.h>
import "C"

import (
	"encoding/hex"
	"unsafe"
)

func TWStringGoString(s unsafe.Pointer) string {
	return C.GoString(C.TWStringUTF8Bytes(s))
}

func TWStringCreate(s string) unsafe.Pointer {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	return C.TWStringCreateWithUTF8Bytes(cStr)
}

func TWStringDelete(s unsafe.Pointer) {
	C.TWStringDelete(s)
}

func TWDataGoBytes(d unsafe.Pointer) []byte {
	cBytes := C.TWDataBytes(d)
	cSize := C.TWDataSize(d)
	return C.GoBytes(unsafe.Pointer(cBytes), C.int(cSize))
}

func TWDataCreate(d []byte) unsafe.Pointer {
	cBytes := C.CBytes(d)
	defer C.free(cBytes)
	return C.TWDataCreateWithBytes((*C.uchar)(cBytes), C.ulong(len(d)))
}

func TWDataDelete(d unsafe.Pointer) {
	C.TWDataDelete(d)
}

func TWDataHexString(d unsafe.Pointer) string {
	return hex.EncodeToString(TWDataGoBytes(d))
}
