package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgChannel struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
}

func (msg MsgChannel) Route() string {
	return RouterKey
}

func (msg MsgChannel) Type() string {
	return "init"
}

func (msg MsgChannel) ValidateBasic() error {
	return nil
}

func (msg MsgChannel) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgChannel) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
