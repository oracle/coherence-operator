package helm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"testing"
)

var HelmHelper *helper.HelmHelper

func TestCoherenceRoleControler(t *testing.T) {
	RegisterFailHandler(Fail)

	// Create a helper.HelmHelper
	h, err := helper.NewOperatorChartHelper()
	if err != nil {
		t.Error(err)
	}

	HelmHelper = h

	RunSpecs(t, "Coherence Operator Helm Suite")
}
