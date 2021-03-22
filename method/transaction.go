package method

import (
	"errors"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v2/extra"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/vedhavyas/go-subkey"
	"math/big"
)

func (this *DotSdk) GetBlockTransaction(blockNum uint64, network uint8) ([]*TxItem, error) {
	txItems := make([]*TxItem, 0)

	types.SetSerDeOptions(types.SerDeOptions{NoPalletIndices: true})

	metadata, err := this.substrateApi.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	blockHash, err := this.substrateApi.RPC.Chain.GetBlockHash(blockNum)
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(metadata, "System", "Events", nil, nil)
	if err != nil {
		return nil, err
	}
	raw, err := this.substrateApi.RPC.State.GetStorageRaw(key, blockHash)
	if err != nil {
		return nil, err
	}
	events := types.EventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(metadata, &events)
	if err != nil {
		return nil, err
	}

	block, err := this.substrateApi.RPC.Chain.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}

	for _, event := range events.Balances_Transfer {
		send, _ := subkey.SS58Address(event.From[:], network)
		to, _ := subkey.SS58Address(event.To[:], network)
		ext := block.Block.Extrinsics[int(event.Phase.AsApplyExtrinsic)]

		resInter := DispatchInfo{}

		err := this.substrateApi.Client.Call(&resInter, "payment_queryInfo", ext, blockHash.Hex())
		if err != nil {
			return nil, err
		}

		partialFee := new(big.Int)
		partialFee, ok := partialFee.SetString(resInter.PartialFee, 10)
		if !ok {
			return nil, errors.New("failed: unable to set amount string")
		}

		item := new(TxItem)
		item.BlockHeight = int64(blockNum)
		item.From = send
		item.To = to
		item.EventId = fmt.Sprintf("%v-%v", blockNum, event.Phase.AsApplyExtrinsic)
		item.Amount = big.NewInt(event.Value.Int64())

		d, err := ext.MarshalJSON()
		if err != nil {
			return nil, err
		}
		tx := extra.CreateTxHash(string(d))
		item.Tx = tx

		item.PartialFee = partialFee
		item.Nonce = big.Int(ext.Signature.Nonce)
		item.Tip = big.Int(ext.Signature.Tip)
		item.Status = 0

		txItems = append(txItems, item)
	}

	//获取执行错误的tx
	for _, event := range events.System_ExtrinsicFailed {
		if event.Phase.IsApplyExtrinsic {
			ext := block.Block.Extrinsics[int(event.Phase.AsApplyExtrinsic)]
			d, err := ext.MarshalJSON()
			if err != nil {
				return nil, err
			}
			tx := extra.CreateTxHash(string(d))
			item := new(TxItem)
			item.BlockHeight = int64(blockNum)
			item.EventId = fmt.Sprintf("%v-%v", blockNum, event.Phase.AsApplyExtrinsic)
			item.Tx = tx
			item.Nonce = big.Int(ext.Signature.Nonce)
			item.Tip = big.Int(ext.Signature.Tip)
			item.Status = 1
			item.PartialFee = big.NewInt(0)

			txItems = append(txItems, item)
		}
	}

	return txItems, nil
}

type TxIn struct {
	Count string `json:"count"`
	// Chain                common.Chain `json:"chain"`
	TxArray              []TxInItem `json:"txArray"`
	Filtered             bool       `json:"filtered"`
	MemPool              bool       `json:"mem_pool"`          // indicate whether this item is in the mempool or not
	SentUnFinalised      bool       `json:"sent_un_finalised"` // indicate whehter unfinalised tx had been sent to THORChain
	Finalised            bool       `json:"finalised"`
	ConfirmationRequired int64      `json:"confirmation_required"`
}

// TxInItem Transaction Item
type TxInItem struct {
	BlockHeight int64  `json:"block_height"`
	Tx          string `json:"tx"`     // Block Hash
	Memo        string `json:"memo"`   // Remarks Text
	Sender      string `json:"sender"` // From Address
	To          string `json:"to"`     // To Address
	Coins       []Coin `json:"coins"`
	Gas         Coin   `json:"gas"` // Gas price
}

type TxItem struct {
	BlockHeight int64    `json:"block_height"`
	Tx          string   `json:"tx"` //Block Hash
	From        string   `json:"from"`
	To          string   `json:"to"`
	PartialFee  *big.Int `json:"partial_fee"` //旷工费
	Nonce       big.Int  `json:"nonce"`       //from的nonce
	Tip         big.Int  `json:"tip"`         //小费
	EventId     string   `json:"event_id"`    //事件id，交易的唯一识别标志
	Status      int      `json:"status"`      //交易状态,0成功，1失败
	Amount      *big.Int `json:"amount"`
}

// Coin struct
type Coin struct {
	Asset  Asset    `json:"asset"`
	Amount *big.Int `json:"amount"`
}

// Asset Struct
type Asset struct {
	Chain  Chain  `json:"chain"`
	Symbol Symbol `json:"symbol"`
	Ticker Ticker `json:"ticker"`
}

// DOTAsset DOT
var DOTAsset = Asset{Chain: "DOT", Symbol: "DOT", Ticker: "DOT"}

// Chain is an alias of string , represent a block chain
type Chain string

// Symbol represent an asset
type Symbol string

// Ticker represent an asset
type Ticker string

// DispatchInfo
type DispatchInfo struct {
	// Weight of this transaction
	Weight float64 `json:"weight"`
	// Class of this transaction
	Class string `json:"class"`
	// PaysFee indicates whether this transaction pays fees
	PartialFee string `json:"partialFee"`
}
