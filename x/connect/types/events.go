package types

import (
	"fmt"

	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

const (
	EventTypeTimeout      = "timeout"
	EventTypePacket       = "non_fungible_token_packet"
	EventTypeChannelClose = "channel_closed"

	AttributeKeyReceiver       = "receiver"
	AttributeKeyRefundReceiver = "refund_receiver"
	AttributeKeyRefundValue    = "refund_value"
	AttributeKeyAckSuccess     = "success"
	AttributeKeyAckError       = "error"
)

// IBC transfer events vars
var (
	AttributeValueCategory = fmt.Sprintf("%s_%s", ibctypes.ModuleName, ModuleName)
)
