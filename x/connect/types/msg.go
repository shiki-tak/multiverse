package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgTransfer struct {
	SrcPort    string         `json:"src_port"`
	SrcChannel string         `json:"src_channel"`
	DestHeight uint64         `json:"dest_height"`
	Receiver   sdk.AccAddress `json:"receiver"`
	Sender     sdk.AccAddress `json:"sender"`
	Denom      string         `json:"denom"`
	TokenID    string         `json:"token_id"`
}

// [src-port] [src-channel] [dest-height] [receiver] [token_id]

func NewMsgTransfer(srcPort string, srcChannel string, destHeight uint64, receiver sdk.AccAddress, sender sdk.AccAddress, denom, tokenID string) MsgTransfer {
	return MsgTransfer{
		SrcPort:    srcPort,
		SrcChannel: srcChannel,
		DestHeight: destHeight,
		Receiver:   receiver,
		Sender:     sender,
		Denom:      denom,
		TokenID:    tokenID,
	}
}

func (msg MsgTransfer) Route() string {
	return RouterKey
}

func (msg MsgTransfer) Type() string {
	return TypeTransfer
}

func (msg MsgTransfer) ValidateBasic() error {
	return nil
}

func (msg MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
