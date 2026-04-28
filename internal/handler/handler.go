package handler

import (
	"encoding/json"
	"net/http"

	"wallet-service/internal/config"
	"wallet-service/internal/wallet"
)

func CreateAddress(cfg *config.Config) http.HandlerFunc {
	type req struct {
		Gate         string `json:"gate"`
		Account      uint32 `json:"account"`
		Change       uint32 `json:"change"`
		AddressIndex uint32 `json:"address_index"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body req
		if !decode(w, r, &body) {
			return
		}
		gate := config.FindGate(cfg, body.Gate)
		if gate == nil {
			writeError(w, http.StatusBadRequest, "gate not found")
			return
		}

		addr, err := wallet.DeriveAddress(gate.Mnemonic, body.Account, body.Change, body.AddressIndex)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "derivation failed")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"address": addr})
	}
}

func ValidateAddress(cfg *config.Config) http.HandlerFunc {
	type req struct {
		Gate    string `json:"gate"`
		Address string `json:"address"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body req
		if !decode(w, r, &body) {
			return
		}
		if config.FindGate(cfg, body.Gate) == nil {
			writeError(w, http.StatusBadRequest, "gate not found")
			return
		}

		writeJSON(w, http.StatusOK, map[string]bool{"valid": wallet.ValidateAddress(body.Address)})
	}
}

func SignTx(cfg *config.Config) http.HandlerFunc {
	type req struct {
		Gate         string          `json:"gate"`
		Account      uint32          `json:"account"`
		Change       uint32          `json:"change"`
		AddressIndex uint32          `json:"address_index"`
		TxParams     wallet.TxParams `json:"tx_params"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body req
		if !decode(w, r, &body) {
			return
		}
		gate := config.FindGate(cfg, body.Gate)
		if gate == nil {
			writeError(w, http.StatusBadRequest, "gate not found")
			return
		}

		signedTx, txHash, err := wallet.SignTransaction(
			gate.Mnemonic, body.Account, body.Change, body.AddressIndex, body.TxParams,
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "signing failed")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"tx_hash":   txHash,
			"signed_tx": signedTx,
		})
	}
}

func decode(w http.ResponseWriter, r *http.Request, body any) bool {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return false
	}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
