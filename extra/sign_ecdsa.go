package extra

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/ecdsa"
	"golang.org/x/crypto/blake2b"
)

func SignUsingEcdsa(e types.Extrinsic, signer subkey.KeyPair, o types.SignatureOptions) (types.Extrinsic, error) {
	if e.Type() != types.ExtrinsicVersion4 {
		return e, fmt.Errorf("unsupported extrinsic version: %v (isSigned: %v, type: %v)", e.Version, e.IsSigned(), e.Type())
	}

	mb, err := types.EncodeToBytes(e.Method)
	if err != nil {
		return e, err
	}

	era := o.Era
	if !o.Era.IsMortalEra {
		era = types.ExtrinsicEra{IsImmortalEra: true}
	}

	payload := types.ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: types.ExtrinsicPayloadV3{
			Method:      mb,
			Era:         era,
			Nonce:       o.Nonce,
			Tip:         o.Tip,
			SpecVersion: o.SpecVersion,
			GenesisHash: o.GenesisHash,
			BlockHash:   o.BlockHash,
		},
		TransactionVersion: o.TransactionVersion,
	}

	// You would use this if you are using Ecdsa/ Ed25519 since it needs to return bytes
	data, err := types.EncodeToBytes(payload)
	if err != nil {
		return e, err
	}

	sig, err := signECDSA(data, hex.EncodeToString(signer.Seed()))
	if err != nil {
		return e, err
	}

	//-----------------校验签名-------------------------

	digest := blake2b.Sum256(data)
	signature := sig[:64]
	ok := secp256k1.VerifySignature(signer.Public(), digest[:], signature)
	if !ok {
		return e, errors.New("verify signature error")
	}

	//-----------------校验签名-------------------------

	multiSig := types.MultiSignature{IsEcdsa: true, AsEcdsa: types.NewBytes(sig)}

	// multiSig := types.MultiSignature{IsEd25519: true, AsEd25519: sig}
	// You would use this if you are using Ecdsa since it needs to return bytes

	signerPubKey := types.NewMultiAddressFromAccountID(signer.Public())

	extSig := types.ExtrinsicSignatureV4{
		Signer:    signerPubKey,
		Signature: multiSig,
		Era:       era,
		Nonce:     o.Nonce,
		Tip:       o.Tip,
	}

	e.Signature = extSig

	// mark the extrinsic as signed
	e.Version |= types.ExtrinsicBitSigned

	return e, nil

}

func signECDSA(data []byte, privateKeyURI string) ([]byte, error) {
	// if data is longer than 256 bytes, hash it first
	if len(data) > 256 {
		h := blake2b.Sum256(data)
		data = h[:]
	}
	scheme := ecdsa.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, privateKeyURI)
	if err != nil {
		return nil, err
	}

	signature, err := kyr.Sign(data)
	if err != nil {
		return nil, err
	}
	return signature, nil
}
