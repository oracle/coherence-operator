package coherencecluster

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCoherenceCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CoherenceCluster Controller Suite")
}
