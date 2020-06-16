package types

import "fmt"

const (
	ModuleName = "connect"

	Version = "connect-1"

	PortID = "connect"

	StoreKey = ModuleName

	RouterKey = ModuleName

	PortKey = "portID"

	QuerierRoute = ModuleName

	TypeTransfer = "transfer"
)

func GenerateKey(prefix, tokenID string) string {
	return fmt.Sprintf("%s%s", prefix, tokenID)
}
