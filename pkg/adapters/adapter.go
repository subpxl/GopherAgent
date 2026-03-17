package adapters

// Adapter is the interface any transport must implement
type Adapter interface {
	Start() error
	Name() string
}
