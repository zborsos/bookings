package dbmodels

import (
	// this blank import is needed for the migration script functionality
	_ "github.com/golang-migrate/migrate/source/file"
)

// Interface is going to be passed to all methods that use
// db.Client so that mock client instances can be passed for testing purposes
type Interface interface {
}
