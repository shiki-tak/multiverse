package relayer

import (
	"github.com/shiki-tak/connect/x/relayer/keeper"
	"github.com/shiki-tak/connect/x/relayer/types"
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

	// variable aliases
	ModuleCdc              = types.ModuleCdc
	AttributeValueCategory = types.AttributeValueCategory
)

type (
	Keeper        = keeper.Keeper
	ChannelKeeper = types.ChannelKeeper
	ClientKeeper  = types.ClientKeeper
	NFTKeeper     = types.NFTKeeper

	GenesisState = types.GenesisState

	MsgSendPacket = types.MsgSendPacket
)
