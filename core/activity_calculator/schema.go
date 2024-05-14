package activity_calculator

type ActiveAddressRes struct {
	Address   string `json:"address"`
	Transfers int    `json:"transfers"`
}
type TopActiveAddressesRes struct {
	TopActiveAddresses []ActiveAddressRes `json:"top_active_addresses"`
}
