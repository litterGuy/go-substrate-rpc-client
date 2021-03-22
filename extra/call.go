package extra

import (
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
)

func BalanceTransferCall(m *types.Metadata, call string, to types.Address, amount types.UCompact) (types.Call, error) {
	c, err := m.FindCallIndex(call)
	if err != nil {
		return types.Call{}, err
	}

	var ma types.MultiAddress
	ma.SetTypes(0)
	ma.AccountId = to.AsAccountID

	return newCall(c, ma, amount)
}

/*
扩展： 创建一个新的Call方法
*/
func newCall(c types.CallIndex, args ...interface{}) (types.Call, error) {
	var a []byte
	for _, arg := range args {
		e, err := types.EncodeToBytes(arg)
		if err != nil {
			return types.Call{}, err
		}
		a = append(a, e...)
	}
	return types.Call{c, a}, nil
}
