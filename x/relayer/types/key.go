package types

import "fmt"

const (
	ModuleName = "relayer"

	Version = "relayer-1"

	PortID = "relayer"

	StoreKey = ModuleName

	RouterKey = ModuleName

	PortKey = "portID"

	QuerierRoute = ModuleName

	TypeRelayer = "relayer"
)

func GenerateKey(prefix, tokenID string) string {
	return fmt.Sprintf("%s%s", prefix, tokenID)
}
