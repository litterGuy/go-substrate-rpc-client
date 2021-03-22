package types

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v2/scale"
)

type MultiAddress struct {
	AccountId AccountID
	Index     UCompact
	Raw       Bytes
	Address32 H256
	Address20 H160
	types     int
}

// NewMultiAddressFromAccountID creates an Address from the given AccountID (public key)
func NewMultiAddressFromAccountID(b []byte) MultiAddress {
	return MultiAddress{
		AccountId: NewAccountID(b),
	}
}

// NewMultiAddressFromHexAccountID creates an Address from the given hex string that contains an AccountID (public key)
func NewMultiAddressFromHexAccountID(str string) (MultiAddress, error) {
	b, err := HexDecodeString(str)
	if err != nil {
		return MultiAddress{}, err
	}
	return NewMultiAddressFromAccountID(b), nil
}

func (d *MultiAddress) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return fmt.Errorf("generic MultiAddress read on bytes error: %v", err)
	}
	switch int(b) {
	case 0:
		err = decoder.Decode(&d.AccountId)
	case 1:
		err = decoder.Decode(&d.Index)
	case 2:
		err = decoder.Decode(&d.Address32)
	case 3:
		err = decoder.Decode(&d.Address20)
	default:
		err = fmt.Errorf("generic MultiAddress unsupport type=%d ", b)
	}
	if err != nil {
		return err
	}
	d.types = int(b)
	return nil
}

func (d MultiAddress) Encode(encoder scale.Encoder) error {
	t := NewU8(uint8(d.types))
	err := encoder.Encode(t)
	if err != nil {
		return err
	}
	switch d.types {
	case 0:
		if &d.AccountId == nil {
			err = fmt.Errorf("generic MultiAddress id is null:%v", d.AccountId)
		}
		err = encoder.Encode(d.AccountId)
	case 1:
		if &d.Index == nil {
			err = fmt.Errorf("generic MultiAddress index is null:%v", d.Index)
		}
		err = encoder.Encode(d.Index)
	case 2:
		if &d.Address32 == nil {
			err = fmt.Errorf("generic MultiAddress address32 is null:%v", d.Address32)
		}
		err = encoder.Encode(d.Address32)
	case 3:
		if &d.Address20 == nil {
			err = fmt.Errorf("generic MultiAddress address20 is null:%v", d.Address20)
		}
		err = encoder.Encode(d.Address20)
	default:
		err = fmt.Errorf("generic MultiAddress unsupport this types: %d", d.types)
	}
	if err != nil {
		return err
	}
	return nil
}

/*
No way, the underlying parsing is done like this
it can only be written like this, although it is very unfriendly
*/
func (d *MultiAddress) GetTypes() int {
	return d.types
}
func (d *MultiAddress) SetTypes(types int) {
	d.types = types
}
func (d *MultiAddress) GetAccountId() AccountID {
	return d.AccountId
}
func (d *MultiAddress) GetIndex() UCompact {
	return d.Index
}
func (d *MultiAddress) GetAddress32() H256 {
	return d.Address32
}
func (d *MultiAddress) GetAddress20() H160 {
	return d.Address20
}

func (d *MultiAddress) ToAddress() Address {
	if d.types != 0 {
		return Address{}
	}
	var ai []byte
	ai = append([]byte{0x00}, d.AccountId[:]...)
	return NewAddressFromAccountID(ai)
}
