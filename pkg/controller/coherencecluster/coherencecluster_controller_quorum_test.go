/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherencecluster

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("coherencecluster_controller start quorum tests", func() {
	var dataStatus coh.ClusterRoleStatus
	var proxyStatus coh.ClusterRoleStatus
	var cluster *coh.CoherenceCluster
	var desiredRole coh.CoherenceRoleSpec
	var canCreate bool
	var reason string

	BeforeEach(func() {
		dataStatus = coh.ClusterRoleStatus{
			Role:   "data",
			Ready:  false,
			Count:  0,
			Status: coh.RoleStatusCreated,
		}
		proxyStatus = coh.ClusterRoleStatus{
			Role:   "proxy",
			Ready:  false,
			Count:  0,
			Status: coh.RoleStatusCreated,
		}

		desiredRole = coh.CoherenceRoleSpec{Role: "test"}
	})

	JustBeforeEach(func() {
		cluster = &coh.CoherenceCluster{
			Spec: coh.CoherenceClusterSpec{
				Roles: []coh.CoherenceRoleSpec{
					{
						Role: "data",
					},
					{
						Role: "proxy",
					},
					desiredRole,
				},
			},
			Status: coh.CoherenceClusterStatus{
				Roles: 2,
				RoleStatus: []coh.ClusterRoleStatus{
					dataStatus,
					proxyStatus,
				},
			},
		}

		p := params{
			cluster:     cluster,
			desiredRole: desiredRole,
		}

		controller := &ReconcileCoherenceCluster{}
		// skip initialization for unit tests
		controller.SetInitialized(true)

		canCreate, reason = controller.canCreateRole(p)
	})

	When("a CoherenceRole does not specify a start quorum", func() {
		It("should be creatable", func() {
			Expect(canCreate).To(BeTrue())
		})

		It("should have empty reason", func() {
			Expect(reason).To(Equal(""))
		})
	})

	When("a CoherenceRole start quorum depending on single role being ready is met", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data"},
				},
			}

			dataStatus.Ready = true
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeTrue())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal(""))
		})
	})

	When("a CoherenceRole start quorum depending on single role being ready is not met", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data"},
				},
			}

			dataStatus.Ready = false
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeFalse())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal("Waiting for creation quorum to be met: \"role 'data' to be ready\""))
		})
	})

	When("a CoherenceRole start quorum depending on multiple being roles being ready is met", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data"},
					{Role: "proxy"},
				},
			}

			dataStatus.Ready = true
			proxyStatus.Ready = true
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeTrue())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal(""))
		})
	})

	When("a CoherenceRole start quorum depending on multiple being roles being ready is not met for one role", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data"},
					{Role: "proxy"},
				},
			}

			dataStatus.Ready = false
			proxyStatus.Ready = true
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeFalse())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal("Waiting for creation quorum to be met: \"role 'data' to be ready\""))
		})
	})

	When("a CoherenceRole start quorum depending on multiple being roles being ready is not met for both roles", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data"},
					{Role: "proxy"},
				},
			}

			dataStatus.Ready = false
			proxyStatus.Ready = false
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeFalse())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal("Waiting for creation quorum to be met: \"role 'data' to be ready\" and \"role 'proxy' to be ready\""))
		})
	})

	When("a CoherenceRole start quorum depending on single role having a pod count is met", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data", PodCount: 5},
				},
			}

			dataStatus.Ready = false
			dataStatus.Count = 5
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeTrue())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal(""))
		})
	})

	When("a CoherenceRole start quorum depending on single role having a pod count is not met", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data", PodCount: 5},
				},
			}

			dataStatus.Ready = false
			dataStatus.Count = 4
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeFalse())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal("Waiting for creation quorum to be met: \"role 'data' to have 5 ready Pods (ready=4)\""))
		})
	})

	When("a CoherenceRole start quorum depending on multiple being roles having a pod count is met", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data", PodCount: 5},
					{Role: "proxy", PodCount: 2},
				},
			}

			dataStatus.Ready = false
			dataStatus.Count = 5
			proxyStatus.Ready = false
			proxyStatus.Count = 5
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeTrue())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal(""))
		})
	})

	When("a CoherenceRole start quorum depending on multiple being roles having a pod count is not met for one role", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data", PodCount: 5},
					{Role: "proxy", PodCount: 2},
				},
			}

			dataStatus.Ready = false
			dataStatus.Count = 4
			proxyStatus.Ready = false
			proxyStatus.Count = 2
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeFalse())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal("Waiting for creation quorum to be met: \"role 'data' to have 5 ready Pods (ready=4)\""))
		})
	})

	When("a CoherenceRole start quorum depending on multiple being roles having a pod count is not met for both roles", func() {
		BeforeEach(func() {
			desiredRole = coh.CoherenceRoleSpec{
				Role: "test",
				StartQuorum: []coh.StartQuorum{
					{Role: "data", PodCount: 5},
					{Role: "proxy", PodCount: 2},
				},
			}

			dataStatus.Ready = false
			dataStatus.Count = 4
			proxyStatus.Ready = false
			proxyStatus.Count = 1
		})

		It("should be creatable", func() {
			Expect(canCreate).To(BeFalse())
		})
		It("should have empty reason", func() {
			Expect(reason).To(Equal("Waiting for creation quorum to be met: \"role 'data' to have 5 ready Pods (ready=4)\" and \"role 'proxy' to have 2 ready Pods (ready=1)\""))
		})
	})
})
