package method

import (
	"bytes"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v2"
	"github.com/centrifuge/go-substrate-rpc-client/v2/hash"
	"github.com/centrifuge/go-substrate-rpc-client/v2/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/cosmos/go-bip39"
)

const (
	Polkadot = 0
	Kusama   = 2
	Westend  = 42
)

var ChainMap = map[uint8]string{
	Polkadot: "Polkadot",
	Kusama:   "Kusama",
	Westend:  "Westend",
}

var NetworkMap = map[string]uint8{
	"Polkadot": Polkadot,
	"Kusama":   Kusama,
	"Westend":  Westend,
}

type DotSdk struct {
	substrateApi *gsrpc.SubstrateAPI
}

func NewDotSdk(rpcurl string) (*DotSdk, error) {
	api, err := gsrpc.NewSubstrateAPI(rpcurl)
	if err != nil {
		return nil, err
	}
	return &DotSdk{substrateApi: api}, nil
}

func (this *DotSdk) CheckChain(network uint8) error {
	net, ok := ChainMap[network]
	if !ok {
		return errors.New("配置节点环境不支持，请检查配置")
	}
	chain, err := this.substrateApi.RPC.System.Chain()
	if err != nil {
		return err
	}
	if chain != types.NewText(net) {
		return errors.New("连接节点和配置节点环境不一致，无法继续工作")
	}
	return nil
}

func (this *DotSdk) CheckSync() error {
	health, err := this.substrateApi.RPC.System.Health()
	if err != nil {
		return err
	}
	if health.IsSyncing {
		return errors.New("节点正在同步区块，无法提供服务")
	}
	return nil
}

func (this *DotSdk) Symbol() (*types.ChainProperties, error) {
	propreties, err := this.substrateApi.RPC.System.Properties()
	if err != nil {
		return nil, err
	}
	return &propreties, nil
}

func (this *DotSdk) GetHieght() (int64, error) {
	finalizedBlockHash, err := this.substrateApi.RPC.Chain.GetFinalizedHead()
	if err != nil {
		return 0, err
	}
	signedBlock, err := this.substrateApi.RPC.Chain.GetBlock(finalizedBlockHash)
	if err != nil {
		return 0, err
	}
	return int64(signedBlock.Block.Header.Number), nil
}

func (this *DotSdk) GenerateNewAddress(network uint8) (*signature.KeyringPair, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	keypair, err := signature.KeyringPairFromSecret(mnemonic, network)
	if err != nil {
		return nil, err
	}
	return &keypair, nil
}

func (this *DotSdk) ValidAddress(address string, network uint8) error {
	contextPrefix := []byte("SS58PRE")
	ss58d := base58.Decode(address)
	if ss58d[0] != network {
		return errors.New("network version check is not match")
	}
	noSum := ss58d[:len(ss58d)-2]
	all := append(contextPrefix, noSum...)
	checksum, err := hash.NewBlake2b512(nil)
	if err != nil {
		return err
	}
	checksum.Write(all)
	res := checksum.Sum(nil)
	// Verified checksum
	if !bytes.Equal(res[:2], ss58d[len(ss58d)-2:]) {
		return errors.New("invliad address")
	}
	return nil
}

func (this *DotSdk) GetPublicKeyFromAddress(address string) []byte {
	ss58d := base58.Decode(address)
	return ss58d[1 : len(ss58d)-2]
}
