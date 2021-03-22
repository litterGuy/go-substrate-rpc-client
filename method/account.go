package method

import (
	"errors"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"reflect"
)

func (this *DotSdk) GetAccountInfo(address string) (*types.AccountInfo, error) {
	publickey := this.GetPublicKeyFromAddress(address)
	meta, err := this.substrateApi.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", publickey, nil)
	if err != nil {
		return nil, err
	}
	var accountInfo types.AccountInfo
	ok, err := this.substrateApi.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		if err == nil {
			err = errors.New("获取帐号信息为空")
		}
		return nil, err
	}
	if reflect.DeepEqual(accountInfo, types.AccountInfo{}) {
		return nil, errors.New("获取帐号信息为空")
	}
	return &accountInfo, nil
}

func (this *DotSdk) GetAccountNonce(address string) (uint32, error) {
	account, err := this.GetAccountInfo(address)
	if err != nil {
		return 0, err
	}
	return uint32(account.Nonce), nil
}
