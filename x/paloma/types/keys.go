package types

const (
	// ModuleName defines the module name
	ModuleName = "paloma"

	// StoreKey defines the primary module store key
	StoreKey = "paloma-store"

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_paloma_store"
)

var (
	LightNodeClientLicenseKeyPrefix = []byte("light-node-client-license")
	LightNodeClientKeyPrefix        = []byte("light-node-client-store")
	LightNodeClientFeegranterKey    = []byte("light-node-client-feegranter")
	LightNodeClientFundersKey       = []byte("light-node-client-funders")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
