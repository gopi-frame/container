package container

import (
	"fmt"

	"github.com/gopi-frame/exception"
)

// EntryNotFoundException entry not found exception
type EntryNotFoundException struct {
	*exception.Exception
}

// NewEntryNotFoundException new entry not found exception
func NewEntryNotFoundException(name string) *EntryNotFoundException {
	return &EntryNotFoundException{
		Exception: exception.NewException(fmt.Sprintf("entry \"%s\" not found in container", name)),
	}
}
