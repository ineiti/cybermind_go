package cymidb

type PFHOperations uint32

const (
	PFHOUpdateFS = PFHOperations(1 << iota)
	PFHOUpdateDB
)

type PluginFileHook struct {
	Hook       Hook
	Root       string
	Operations PFHOperations
}

func NewPluginFileHook(h Hook, dir string, ops PFHOperations) (pfh PluginFileHook, err error) {
	return
}
