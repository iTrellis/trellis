package registry

import "github.com/go-trellis/trellis/configure"

type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}

const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// Result is returned by a call to Next on
// the watcher. Actions can be create, update, delete
type Result struct {
	Action  string
	Service configure.RegistServices
}

type watcher struct {
	opts WatchOption
}
