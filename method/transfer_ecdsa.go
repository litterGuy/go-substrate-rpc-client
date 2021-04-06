package method

import (
	"github.com/centrifuge/go-substrate-rpc-client/v2/extra"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/vedhavyas/go-subkey/ecdsa"
	"math/big"
)

func (this *DotSdk) TransferECDSA(phrase string, network uint8, to string, amount *big.Int, nonce uint32) (*types.Hash, *big.Int, error) {

	//转化获取私钥
	secert, err := types.HexDecodeString(phrase)
	if err != nil {
		return nil, nil, err
	}

	scheme := ecdsa.Scheme{}
	pari, err := scheme.FromSeed(secert)
	if err != nil {
		return nil, nil, err
	}

	meta, err := this.substrateApi.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, nil, err
	}

	// Create a call, transferring
	accountId := this.GetPublicKeyFromAddress(to)
	toaddress := types.NewAddressFromAccountID(accountId)

	//c, err := types.NewCall(meta, "Balances.transfer", toaddress, types.NewUCompact(amount))
	c, err := extra.BalanceTransferCall(meta, "Balances.transfer", toaddress, types.NewUCompact(amount))
	if err != nil {
		return nil, nil, err
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := this.substrateApi.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, nil, err
	}

	rv, err := this.substrateApi.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, nil, err
	}

	blockhash, err := this.substrateApi.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		return nil, nil, err
	}
	lastestBlock, err := this.substrateApi.RPC.Chain.GetBlockLatest()
	if err != nil {
		return nil, nil, err
	}
	blockNum := uint64(lastestBlock.Block.Header.Number)
	era := extra.GetEra(blockNum)

	o := types.SignatureOptions{
		BlockHash:          blockhash,
		Era:                era,
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	extSigned, err := extra.SignUsingEcdsa(ext, pari, o)
	if err != nil {
		return nil, nil, err
	}

	partialFee, err := this.getPartialFee(extSigned)
	if err != nil {
		return nil, nil, err
	}

	// Send the extrinsic
	hash, err := this.substrateApi.RPC.Author.SubmitExtrinsic(extSigned)
	if err != nil {
		return nil, nil, err
	}
	return &hash, partialFee, nil
}
