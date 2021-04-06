package extra

import (
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"math"
)

// Must be a power of two between 4 and 65536 (inclusive)
const MortalEraPeriod = uint64(64)

func GetEra(currentBlockNumber uint64) types.ExtrinsicEra {
	// Adapted from https://substrate.dev/rustdocs/v2.0.1/src/sp_runtime/generic/era.rs.html#66
	phase := currentBlockNumber % MortalEraPeriod

	quantizeFactor := MortalEraPeriod >> 12
	if quantizeFactor < 1 {
		quantizeFactor = 1
	}
	quantizedPhase := phase / quantizeFactor * quantizeFactor

	encoded := uint16(math.Log2(float64(MortalEraPeriod))-1) | uint16((quantizedPhase/quantizeFactor)<<4)

	return types.ExtrinsicEra{
		IsMortalEra:   true,
		IsImmortalEra: false,
		AsMortalEra: types.MortalEra{
			First:  byte(encoded),
			Second: byte(encoded >> 8),
		},
	}
}

func GetOtherEra(blockNum uint64) *types.ExtrinsicEra {
	if blockNum == 0 {
		return nil
	}
	phase := blockNum % MortalEraPeriod
	index := uint64(6)
	trailingZero := index - 1

	var encoded uint64
	if trailingZero > 1 {
		encoded = trailingZero
	} else {
		encoded = 1
	}

	if trailingZero < 15 {
		encoded = trailingZero
	} else {
		encoded = 15
	}
	encoded += phase / 1 << 4
	first := byte(encoded >> 8)
	second := byte(encoded & 0xff)
	era := new(types.ExtrinsicEra)
	era.IsMortalEra = true
	era.AsMortalEra.First = first
	era.AsMortalEra.Second = second
	return era
}