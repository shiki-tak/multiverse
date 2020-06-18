package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSendPacket struct {
	SrcPort    string         `json:"src_port"`
	SrcChannel string         `json:"src_channel"`
	DestHeight uint64         `json:"dest_height"`
	Receiver   sdk.AccAddress `json:"receiver"`
	Sender     sdk.AccAddress `json:"sender"`
	// TODO: add field
}

// [src-port] [src-channel] [dest-height] [receiver] [token_id]

func NewMsgSendPacket(srcPort string, srcChannel string, destHeight uint64, receiver sdk.AccAddress, sender sdk.AccAddress) MsgSendPacket {
	return MsgSendPacket{
		SrcPort:    srcPort,
		SrcChannel: srcChannel,
		DestHeight: destHeight,
		Receiver:   receiver,
		Sender:     sender,
	}
}

func (msg MsgSendPacket) Route() string {
	return RouterKey
}

func (msg MsgSendPacket) Type() string {
	return TypeRelayer
}

func (msg MsgSendPacket) ValidateBasic() error {
	return nil
}

func (msg MsgSendPacket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSendPacket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
