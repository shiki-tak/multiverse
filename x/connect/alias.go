package connect

import (
	"github.com/shiki-tak/connect/x/connect/keeper"
	"github.com/shiki-tak/connect/x/connect/types"
)

const (
	DefaultPacketTimeoutHeight    = keeper.DefaultPacketTimeoutHeight
	DefaultPacketTimeoutTimestamp = keeper.DefaultPacketTimeoutTimestamp
	EventTypeTimeout              = types.EventTypeTimeout
	EventTypePacket               = types.EventTypePacket
	EventTypeChannelClose         = types.EventTypeChannelClose
	AttributeKeyReceiver          = types.AttributeKeyReceiver
	AttributeKeyRefundReceiver    = types.AttributeKeyRefundReceiver
	AttributeKeyRefundValue       = types.AttributeKeyRefundValue
	AttributeKeyAckSuccess        = types.AttributeKeyAckSuccess
	AttributeKeyAckError          = types.AttributeKeyAckError

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	// QueryParams       = types.QueryParams
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewMsgTransfer      = types.NewMsgTransfer

	// variable aliases
	ModuleCdc              = types.ModuleCdc
	AttributeValueCategory = types.AttributeValueCategory
)

type (
	Keeper           = keeper.Keeper
	ChannelKeeper    = types.ChannelKeeper
	ClientKeeper     = types.ClientKeeper
	ConnectionKeeper = types.ConnectionKeeper
	NFTKeeper        = types.NFTKeeper

	GenesisState = types.GenesisState

	NonFungibleTokenPacketData            = types.NonFungibleTokenPacketData
	NonFungibleTokenPacketAcknowledgement = types.NonFungibleTokenPacketAcknowledgement
	MsgTransfer                           = types.MsgTransfer

	UnacknowledgedPacket               = types.UnacknowledgedPacket
	QueryUnacknowledgedPacketsRequest  = types.QueryUnacknowledgedPacketsRequest
	QueryUnacknowledgedPacketsResponse = types.QueryUnacknowledgedPacketsResponse
)
