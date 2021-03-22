package extra

import (
	"encoding/hex"
	"golang.org/x/crypto/blake2b"
	"strings"
)

func CreateTxHash(extrinsic string) string {
	extrinsic = Remove0X(extrinsic)
	data, _ := hex.DecodeString(extrinsic)
	d := blake2b.Sum256(data)
	return "0x" + hex.EncodeToString(d[:])
}

func Remove0X(hexData string) string {
	hexData = strings.ReplaceAll(hexData,"\"","")
	if strings.HasPrefix(hexData, "0x") {
		return hexData[2:]
	}
	return hexData
}
