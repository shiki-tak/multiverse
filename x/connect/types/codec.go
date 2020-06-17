package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTransfer{}, "connect/MsgTransfer", nil)
	cdc.RegisterConcrete(NonFungibleTokenPacketData{}, "connect/NonFungibleTokenPacketData", nil)
}

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}
