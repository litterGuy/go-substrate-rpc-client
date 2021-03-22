// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"github.com/centrifuge/go-substrate-rpc-client/v2/scale"
)

// ChainProperties contains the SS58 format, the token decimals and the token symbol
type ChainProperties struct {
	IsSS58Format    bool
	SS58Format    U8
	IsTokenDecimals bool
	TokenDecimals U32
	IsTokenSymbol   bool
	TokenSymbol   Text
}

func (a *ChainProperties) Decode(decoder scale.Decoder) error {
	if err := decoder.DecodeOption(&a.IsSS58Format, &a.SS58Format); err != nil {
		return err
	}
	if err := decoder.DecodeOption(&a.IsTokenDecimals, &a.TokenDecimals); err != nil {
		return err
	}
	if err := decoder.DecodeOption(&a.IsTokenSymbol, &a.TokenSymbol); err != nil {
		return err
	}

	return nil
}

func (a ChainProperties) Encode(encoder scale.Encoder) error {
	if err := encoder.EncodeOption(a.IsSS58Format, a.SS58Format); err != nil {
		return err
	}
	if err := encoder.EncodeOption(a.IsTokenDecimals, a.TokenDecimals); err != nil {
		return err
	}
	if err := encoder.EncodeOption(a.IsTokenSymbol, a.TokenSymbol); err != nil {
		return err
	}

	return nil
}
