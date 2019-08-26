package v1_test

import (
	"github.com/onsi/ginkgo/reporters"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApiV1(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("test-report.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "API V1 Suite", []Reporter{junitReporter})
}
