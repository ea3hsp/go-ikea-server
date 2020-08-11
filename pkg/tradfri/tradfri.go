package tradfri

// ITradfri tradfri interface
type ITradfri interface {
	AuthExchange(string) (interface{}, error)
	DevicePower(int, int) (interface{}, error)
}
