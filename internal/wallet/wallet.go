package wallet

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"wallet-service/internal/twcore"
	"wallet-service/proto/common"
	pb "wallet-service/proto/ethereum"

	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/proto"
)

type TxParams struct {
	To                      string `json:"to"`
	ValueWei                string `json:"value_wei"`
	Data                    string `json:"data"`
	Nonce                   uint64 `json:"nonce"`
	ChainID                 uint64 `json:"chain_id"`
	GasLimit                uint64 `json:"gas_limit"`
	MaxFeePerGasWei         string `json:"max_fee_per_gas_wei"`
	MaxPriorityFeePerGasWei string `json:"max_priority_fee_per_gas_wei"`
}

func DeriveAddress(mnemonic string, account, change, addressIndex uint32) (string, error) {
	return twcore.DeriveAddress(mnemonic, account, change, addressIndex)
}

func ValidateAddress(address string) bool {
	return twcore.ValidateAddress(address)
}

func SignTransaction(mnemonic string, account, change, addressIndex uint32, params TxParams) (string, string, error) {
	priKeyBytes, err := twcore.DeriveKey(mnemonic, account, change, addressIndex)
	if err != nil {
		return "", "", err
	}

	chainID := new(big.Int).SetUint64(params.ChainID)
	nonce := new(big.Int).SetUint64(params.Nonce)
	gasLimit := new(big.Int).SetUint64(params.GasLimit)

	maxFeePerGas, ok := new(big.Int).SetString(params.MaxFeePerGasWei, 10)
	if !ok {
		return "", "", fmt.Errorf("invalid max_fee_per_gas_wei: %s", params.MaxFeePerGasWei)
	}
	maxPriorityFee, ok := new(big.Int).SetString(params.MaxPriorityFeePerGasWei, 10)
	if !ok {
		return "", "", fmt.Errorf("invalid max_priority_fee_per_gas_wei: %s", params.MaxPriorityFeePerGasWei)
	}

	value, ok := new(big.Int).SetString(params.ValueWei, 10)
	if !ok {
		return "", "", fmt.Errorf("invalid value_wei: %s", params.ValueWei)
	}

	data, err := hexDecode(params.Data)
	if err != nil {
		return "", "", fmt.Errorf("invalid data hex: %w", err)
	}

	input := &pb.SigningInput{
		ChainId:               bigIntBytes(chainID),
		Nonce:                 bigIntBytes(nonce),
		TxMode:                pb.TransactionMode_Enveloped,
		MaxFeePerGas:          bigIntBytes(maxFeePerGas),
		MaxInclusionFeePerGas: bigIntBytes(maxPriorityFee),
		GasLimit:              bigIntBytes(gasLimit),
		ToAddress:             params.To,
		PrivateKey:            priKeyBytes,
		Transaction: &pb.Transaction{
			TransactionOneof: &pb.Transaction_ContractGeneric_{
				ContractGeneric: &pb.Transaction_ContractGeneric{
					Amount: bigIntBytes(value),
					Data:   data,
				},
			},
		},
	}

	inputBytes, err := proto.Marshal(input)
	if err != nil {
		return "", "", fmt.Errorf("marshal signing input: %w", err)
	}

	outputBytes, err := twcore.Sign(inputBytes)
	if err != nil {
		return "", "", fmt.Errorf("sign tx: %w", err)
	}

	var output pb.SigningOutput
	if err := proto.Unmarshal(outputBytes, &output); err != nil {
		return "", "", fmt.Errorf("unmarshal signing output: %w", err)
	}

	if output.Error != common.SigningError_OK {
		return "", "", fmt.Errorf("signing error: %s (%s)", output.Error.String(), output.ErrorMessage)
	}

	encoded := output.GetEncoded()
	signedTxHex := "0x" + hex.EncodeToString(encoded)
	txHash := "0x" + hex.EncodeToString(keccak256(encoded))

	return signedTxHex, txHash, nil
}

func bigIntBytes(n *big.Int) []byte {
	if n.Sign() == 0 {
		return []byte{0}
	}
	return n.Bytes()
}

func hexDecode(s string) ([]byte, error) {
	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}
	if len(s) == 0 {
		return []byte{}, nil
	}
	return hex.DecodeString(s)
}

func keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}
