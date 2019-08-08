// matchers contains custom implementations of Gomega matchers.
package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	coreV1 "k8s.io/api/core/v1"
	"reflect"
)

var emptyEnvVars []coreV1.EnvVar

//HaveEnvVar asserts that a []EnvVar contains a specified EnvVar.
func HaveEnvVar(expected coreV1.EnvVar) types.GomegaMatcher {
	return &HaveEnvVarMatcher{EnvVar: expected}
}

// HaveEnvVarMatcher is a Gomega matcher that matches an EnvVar in an EnvVar slice
type HaveEnvVarMatcher struct {
	EnvVar coreV1.EnvVar
}

// Match asserts that a value is an EnvVar slice containing the expected EnvVar.
func (h *HaveEnvVarMatcher) Match(actual interface{}) (success bool, err error) {
	if !isEnvVarSlice(actual) {
		return false, fmt.Errorf("HaveEnvVar matcher expects a []EnvVar. Got:%s", format.Object(actual, 1))
	}

	v := reflect.ValueOf(actual)
	l := v.Len()
	for i := 0; i < l; i++ {
		v := reflect.ValueOf(actual).Index(i)
		ev := v.Interface().(coreV1.EnvVar)
		if reflect.DeepEqual(h.EnvVar, ev) {
			return true, nil
		}
	}

	return false, nil
}

func (h *HaveEnvVarMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to contain EnvVar matching %s", h.EnvVar)
}

func (h *HaveEnvVarMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to not contain EnvVar matching %s", h.EnvVar)
}

// isEnvVarSlice determines whether the specfied interface is an EnvVar slice.
func isEnvVarSlice(a interface{}) bool {
	if a == nil {
		return false
	}
	return reflect.TypeOf(a) == reflect.TypeOf(emptyEnvVars)
}
