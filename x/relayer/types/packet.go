package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NonFungibleTokenPacketData struct {
	Denom    string `json:"denom"`
	ID       string `json:"id"`
	TokenURI string `json:"token_uri"`
	// NFT      nftexported.NFT `json:"nft"`	// TODO:
	Sender   sdk.AccAddress `json:"sender"`
	Receiver sdk.AccAddress `json:"receiver"`
}

type NonFungibleTokenPacketAcknowledgement struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewNonFungibleTokenPacketData(denom, id, tokenURI string, sender, receiver sdk.AccAddress) NonFungibleTokenPacketData {
	return NonFungibleTokenPacketData{
		Denom:    denom,
		ID:       id,
		TokenURI: tokenURI,
		// NFT:      nft,
		Sender:   sender,
		Receiver: receiver,
	}
}

func (nftpd NonFungibleTokenPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(nftpd))
}

// GetBytes is a helper for serialising
func (ack NonFungibleTokenPacketAcknowledgement) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(ack))
}
