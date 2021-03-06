package ibc_testing

import (
	"fmt"
)

// TestConnection is a testing helper struct to keep track of the connectionID, source clientID,
// and counterparty clientID used in creating and interacting with a connection.
type TestConnection struct {
	ID                   string
	ClientID             string
	CounterpartyClientID string
	NextChannelVersion   string
	Channels             []TestChannel
}

// AddTestChannel appends a new TestChannel which contains references to the port and channel ID
// used for channel creation and interaction.
//
// channel ID format: connectionid-<channel-index>
// the port is set to "transfer" to be compatible with the ICS-transfer module, this should
// eventually be updated as described in the issue: https://github.com/cosmos/cosmos-sdk/issues/6509
func (conn *TestConnection) AddTestChannel(portID string) TestChannel {
	channel := conn.NextTestChannel(portID)
	conn.Channels = append(conn.Channels, channel)
	return channel
}

// NextTestChannel returns the next test channel to be created on this connection, but does not
// add it to the list of created channels. This function is expected to be used when the caller
// has not created the associated channel in app state, but would still like to refer to the
// non-existent channel usually to test for its non-existence.
func (conn *TestConnection) NextTestChannel(portID string) TestChannel {
	channelID := fmt.Sprintf("%s%d", conn.ID, len(conn.Channels))
	return TestChannel{
		Version:              conn.NextChannelVersion,
		PortID:               portID,
		ID:                   channelID,
		ClientID:             conn.ClientID,
		CounterpartyClientID: conn.CounterpartyClientID,
	}
}

// FirstOrNextTestChannel returns the first test channel if it exists, otherwise it
// returns the next test channel to be created. This function is expected to be used
// when the caller does not know if the channel has or has not been created in app
// state, but would still like to refer to it to test existence or non-existence.
func (conn *TestConnection) FirstOrNextTestChannel() TestChannel {
	if len(conn.Channels) > 0 {
		return conn.Channels[0]
	}
	return conn.NextTestChannel("transfer")
}

// TestChannel is a testing helper struct to keep track of the portID and channelID
// used in creating and interacting with a channel. The clientID and counterparty
// client ID are also tracked to cut down on querying and argument passing.
type TestChannel struct {
	PortID               string
	ID                   string
	ClientID             string
	CounterpartyClientID string
	Version              string
}
