package repository

// Storage is an interface that can be implemented in order to have control over how
// the repository of feature toggles is persisted.
type Storage interface {
	// Get returns the data for the specified feature toggle.
	Get(string, interface{}) error

	Set(string, interface{}) error

	Remove(string) error

	// List returns a list of all feature toggles.
	List() []interface{}
}
