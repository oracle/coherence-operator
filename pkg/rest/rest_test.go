/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package rest_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/pkg/fakes"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/rest"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
)

var _ = Describe("ReST Server", func() {
	const (
		siteLabel      = "test/site"
		rackLabel      = "test/rack"
		urlSitePattern = "http://127.0.0.1:%d/site/%s"
		urlRackPattern = "http://127.0.0.1:%d/rack/%s"
	)
	var (
		server       rest.Server
		mgr          *fakes.FakeManager
		err          error
		node         *corev1.Node
		clientErrors *fakes.ClientErrors
	)

	JustBeforeEach(func() {
		mgr = fakes.NewFakeManager()

		if node != nil {
			_ = mgr.GetClient().Create(context.TODO(), node)
		}

		if clientErrors != nil {
			mgr.Client.EnableErrors(*clientErrors)
		}

		cf := &flags.CoherenceOperatorFlags{
			RestHost:  "0.0.0.0",
			RestPort:  0,
			SiteLabel: siteLabel,
			RackLabel: rackLabel,
		}

		server, err = rest.StartRestServer(mgr, cf)

		if err != nil {
			fmt.Println(err)
		}
	})

	When("node exists with site and rack labels", func() {
		BeforeEach(func() {
			labels := map[string]string{
				siteLabel: "foo-site",
				rackLabel: "foo-rack",
			}

			node = &corev1.Node{}
			node.SetName("foo")
			node.SetLabels(labels)
		})

		It("the site should be blank", func() {
			url := fmt.Sprintf(urlSitePattern, server.GetPort(), "foo")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("foo-site"))
		})

		It("the rack should be blank", func() {
			url := fmt.Sprintf(urlRackPattern, server.GetPort(), "foo")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("foo-rack"))
		})
	})

	When("node exists without labels", func() {
		BeforeEach(func() {
			node = &corev1.Node{}
			node.SetName("foo")
		})

		It("the site should be blank", func() {
			url := fmt.Sprintf(urlSitePattern, server.GetPort(), "foo")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(""))
		})

		It("the rack should be blank", func() {
			url := fmt.Sprintf(urlRackPattern, server.GetPort(), "foo")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(""))
		})
	})

	When("node does not exist", func() {
		BeforeEach(func() {
			node = nil
		})

		It("the site should be blank", func() {
			url := fmt.Sprintf(urlSitePattern, server.GetPort(), "foo")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(""))
		})

		It("the rack should be blank", func() {
			url := fmt.Sprintf(urlRackPattern, server.GetPort(), "foo")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(""))
		})
	})

	When("k8s returns error getting node", func() {
		BeforeEach(func() {
			labels := map[string]string{
				siteLabel: "foo-site",
				rackLabel: "foo-rack",
			}

			node = &corev1.Node{}
			node.SetName("foo")
			node.SetLabels(labels)

			clientErrors = &fakes.ClientErrors{}
			clientErrors.AddGetError(fakes.ErrorIf{KeyIs: &types.NamespacedName{Name: "foo"}}, fmt.Errorf("oops"))
		})

		It("the site should be blank", func() {
			url := fmt.Sprintf(urlSitePattern, server.GetPort(), "foo")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(""))
		})

		It("the rack should be blank", func() {
			url := fmt.Sprintf(urlRackPattern, server.GetPort(), "foo")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(""))
		})
	})
})
