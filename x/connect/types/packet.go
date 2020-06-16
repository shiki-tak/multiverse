package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	nftexported "github.com/shiki-tak/connect/x/nft/exported"
)

type NonFungibleTokenPacketData struct {
	Denom    string          `json:"denom"`
	NFT      nftexported.NFT `json:"nft"`
	Sender   sdk.AccAddress  `json:"sender"`
	Receiver sdk.AccAddress  `json:"receiver"`
}

type NonFungibleTokenPacketAcknowledgement struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewNonFungibleTokenPacketData(denom string, nft nftexported.NFT, sender, receiver sdk.AccAddress) NonFungibleTokenPacketData {
	return NonFungibleTokenPacketData{
		NFT:      nft,
		Sender:   sender,
		Receiver: receiver,
	}
}

func (ftpd NonFungibleTokenPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(ftpd))
}

// GetBytes is a helper for serialising
func (ack NonFungibleTokenPacketAcknowledgement) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(ack))
}
