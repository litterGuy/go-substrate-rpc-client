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

// AccountInfo contains information of an account
type AccountInfo struct {
	Nonce     U32 `json:"nonce"`
	Consumers U32 `json:"consumers"`
	Providers U32 `json:"providers"`
	Data      struct {
		Free       U128 `json:"free"`
		Reserved   U128 `json:"reserved"`
		MiscFrozen U128 `json:"misc_frozen"`
		FreeFrozen U128 `json:"free_frozen"`
	} `json:"data"`
}
