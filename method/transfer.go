package method

import (
	"errors"
	"github.com/centrifuge/go-substrate-rpc-client/v2/extra"
	"github.com/centrifuge/go-substrate-rpc-client/v2/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"math/big"
)

/**
nonce: 从系统模块中查询的 nonce 不会考虑待处理的交易。如果您想要同时提交多个有效的交易，您必须跟踪并手动递增 nonce
*/
func (this *DotSdk) Transfer(phrase string, network uint8, to string, amount *big.Int, nonce uint32) (*types.Hash, *big.Int, error) {
	frompair, err := signature.KeyringPairFromSecret(phrase, network)
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
	err = ext.Sign(frompair, o)
	if err != nil {
		return nil, nil, err
	}

	partialFee, err := this.getPartialFee(ext)
	if err != nil {
		return nil, nil, err
	}

	// Send the extrinsic
	hash, err := this.substrateApi.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		return nil, nil, err
	}
	return &hash, partialFee, nil
}

func (this *DotSdk) getPartialFee(ext types.Extrinsic) (*big.Int, error) {
	blockHash, err := this.substrateApi.RPC.Chain.GetFinalizedHead()
	if err != nil {
		return nil, err
	}
	resInter := DispatchInfo{}
	err = this.substrateApi.Client.Call(&resInter, "payment_queryInfo", ext, blockHash.Hex())
	if err != nil {
		return nil, err
	}
	partialFee := new(big.Int)
	partialFee, ok := partialFee.SetString(resInter.PartialFee, 10)
	if !ok {
		return nil, errors.New("BXL: ERROR: unable to set amount string")
	}
	return partialFee, nil
}

func (this *DotSdk) GetPartialFee(phrase string, network uint8, to string, amount *big.Int, nonce uint32) (*big.Int, error) {
	frompair, err := signature.KeyringPairFromSecret(phrase, network)
	if err != nil {
		return nil, err
	}
	meta, err := this.substrateApi.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	// Create a call, transferring
	accountId := this.GetPublicKeyFromAddress(to)
	toaddress := types.NewAddressFromAccountID(accountId)

	c, err := extra.BalanceTransferCall(meta, "Balances.transfer", toaddress, types.NewUCompact(amount))
	if err != nil {
		return nil, err
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := this.substrateApi.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	rv, err := this.substrateApi.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(frompair, o)
	if err != nil {
		return nil, err
	}

	partialFee, err := this.getPartialFee(ext)
	if err != nil {
		return nil, err
	}
	return partialFee, nil
}
