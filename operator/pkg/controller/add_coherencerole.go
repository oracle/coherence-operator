package controller

import (
	"github.com/oracle/coherence-operator/pkg/controller/coherencerole"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, coherencerole.Add)
}
