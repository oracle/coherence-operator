package coherencerole

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	)

func TestCoherenceRoleControler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CoherenceRole Controller Suite")
}
