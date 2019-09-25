/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"fmt"
	"github.com/go-test/deep"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing CoherenceRoleSpec struct", func() {

	Context("Copying an CoherenceRoleSpec using DeepCopyWithDefaults", func() {
		var (
			roleSpecOne = LoadCoherenceRoleFromCoherenceClusterYamlFile("full_role_one.yaml")
			roleSpecTwo = LoadCoherenceRoleFromCoherenceClusterYamlFile("full_role_two.yaml")
			original    *coherence.CoherenceRoleSpec
			defaults    *coherence.CoherenceRoleSpec
			clone       *coherence.CoherenceRoleSpec
		)

		// just before every "It" this method is executed to actually do the cloning
		JustBeforeEach(func() {
			clone = original.DeepCopyWithDefaults(defaults)
		})

		When("original and defaults are nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = nil
			})

			It("the copy should be nil", func() {
				Expect(clone).Should(BeNil())
			})
		})

		When("defaults is nil", func() {
			BeforeEach(func() {
				original = roleSpecOne.DeepCopy()
				defaults = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = roleSpecTwo.DeepCopy()
			})

			It("clone should be equal to the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = roleSpecOne.DeepCopy()
				defaults = roleSpecTwo.DeepCopy()
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original), fmt.Sprintf("Diff is %s", deep.Equal(clone, original)))
			})
		})

		// ----- Affinity -----------------------------------------------------------

		When("the original Affinity is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo.DeepCopy()
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Affinity = nil
			})

			It("clone should be equal to the original with the Affinity field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Affinity = defaults.Affinity

				Expect(clone).To(Equal(expected))
			})
		})

		// ----- Annotations --------------------------------------------------------

		Context("Annotations are merged", func() {
			var annOne = map[string]string{"One": "1", "Two": "2"}
			var annTwo = map[string]string{"Three": "3", "Four": "4"}

			When("the original Annotations is nil and default Annotations is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = nil
					original = roleSpecOne.DeepCopy()
					original.Annotations = nil
				})

				It("clone should have nil Affinity", func() {
					expected := original.DeepCopy()
					expected.Annotations = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Annotations is nil and default Annotations is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.Annotations = nil
				})

				It("clone should have empty Annotations", func() {
					expected := original.DeepCopy()
					expected.Annotations = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Annotations is nil and default Annotations is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = annOne
					original = roleSpecOne.DeepCopy()
					original.Annotations = nil
				})

				It("clone should have default Annotations", func() {
					expected := original.DeepCopy()
					expected.Annotations = defaults.Annotations
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Annotations is empty map and default Annotations is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = nil
					original = roleSpecOne.DeepCopy()
					original.Annotations = map[string]string{}
				})

				It("clone should have empty Annotations", func() {
					expected := original.DeepCopy()
					expected.Annotations = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Annotations is empty map and default Annotations is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.Annotations = map[string]string{}
				})

				It("clone should have empty Annotations", func() {
					expected := original.DeepCopy()
					expected.Annotations = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Annotations is empty map and default Annotations is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = annOne
					original = roleSpecOne.DeepCopy()
					original.Annotations = map[string]string{}
				})

				It("clone should have default Annotations", func() {
					expected := original.DeepCopy()
					expected.Annotations = defaults.Annotations
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Annotations is set map and default Annotations is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = nil
					original = roleSpecOne.DeepCopy()
					original.Annotations = annTwo
				})

				It("clone should have original Annotations", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Annotations is set and default Annotations is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.Annotations = annTwo
				})

				It("clone should have original Annotations", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Annotations is set and default Annotations is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Annotations = annOne
					original = roleSpecOne.DeepCopy()
					original.Annotations = annTwo
				})

				It("clone should have merged Annotations", func() {
					expected := original.DeepCopy()
					expected.Annotations = map[string]string{"One": "1", "Two": "2", "Three": "3", "Four": "4"}
					Expect(clone).To(Equal(expected))
				})

				When("the original Annotations is set and default Annotations is set with overlapping keys", func() {
					BeforeEach(func() {
						defaults = roleSpecTwo.DeepCopy()
						defaults.Annotations = annOne
						original = roleSpecOne.DeepCopy()
						original.Annotations = map[string]string{"Two": "22", "Four": "4"}
					})

					It("clone should have merged Annotations with original key taking precedence", func() {
						expected := original.DeepCopy()
						expected.Annotations = map[string]string{"One": "1", "Two": "22", "Four": "4"}
						Expect(clone).To(Equal(expected))
					})
				})
			})
		})

		// ----- Application --------------------------------------------------------

		Context("Application is merged", func() {
			var appOne = coherence.ApplicationSpec{
				Type:   stringPtr("java"),
				Main:   stringPtr("main"),
				LibDir: stringPtr("/lib"),
			}

			var appTwo = coherence.ApplicationSpec{
				ImageSpec: coherence.ImageSpec{
					Image: stringPtr("Foo:1.0"),
				},
				ConfigDir: stringPtr("/cfg"),
			}

			When("the original Application is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					original = roleSpecOne.DeepCopy()
					original.Application = nil
				})

				It("clone should be equal to the original with the Application field from the defaults", func() {
					expected := original.DeepCopy()
					expected.Application = defaults.Application
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Application is set and defaults Application is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Application = nil
					original = roleSpecOne.DeepCopy()
				})

				It("clone should be equal to the original with the Application field from the original", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Application is set and defaults Application is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Application = &appOne
					original = roleSpecOne.DeepCopy()
					original.Application = &appTwo
				})

				It("clone should be equal to the merged original and defaults Applications", func() {
					expected := original.DeepCopy()
					expected.Application = original.Application.DeepCopyWithDefaults(defaults.Application)
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- Coherence ----------------------------------------------------------

		Context("Coherence is merged", func() {
			var cohOne = coherence.CoherenceSpec{
				CacheConfig:    stringPtr("one.xml"),
				StorageEnabled: boolPtr(true),
			}

			var cohTwo = coherence.CoherenceSpec{
				OverrideConfig: stringPtr("override.xml"),
			}

			When("the original Coherence is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					original = roleSpecOne.DeepCopy()
					original.Coherence = nil
				})

				It("clone should be equal to the original with the Coherence field from the defaults", func() {
					expected := original.DeepCopy()
					expected.Coherence = defaults.Coherence
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Coherence is set and defaults Coherence is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Coherence = nil
					original = roleSpecOne.DeepCopy()
				})

				It("clone should be equal to the original with the Coherence field from the original", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Coherence is set and defaults Coherence is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Coherence = &cohOne
					original = roleSpecOne.DeepCopy()
					original.Coherence = &cohTwo
				})

				It("clone should be equal to the merged original and defaults Coherences", func() {
					expected := original.DeepCopy()
					expected.Coherence = original.Coherence.DeepCopyWithDefaults(defaults.Coherence)
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- Env ----------------------------------------------------------------

		Context("Environment Variables are merged", func() {
			var envOne = []corev1.EnvVar{{Name: "One", Value: "1"}, {Name: "Two", Value: "2"}}
			var envTwo = []corev1.EnvVar{{Name: "Three", Value: "3"}, {Name: "Four", Value: "4"}}

			When("the original Env is nil and default Env is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = nil
					original = roleSpecOne.DeepCopy()
					original.Env = nil
				})

				It("clone should have nil Env", func() {
					expected := original.DeepCopy()
					expected.Env = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Env is nil and default Env is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = []corev1.EnvVar{}
					original = roleSpecOne.DeepCopy()
					original.Env = nil
				})

				It("clone should have empty Env", func() {
					expected := original.DeepCopy()
					expected.Env = []corev1.EnvVar{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Env is nil and default Env is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = envOne
					original = roleSpecOne.DeepCopy()
					original.Env = nil
				})

				It("clone should have default Env", func() {
					expected := original.DeepCopy()
					expected.Env = defaults.Env
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Env is empty map and default Env is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = nil
					original = roleSpecOne.DeepCopy()
					original.Env = []corev1.EnvVar{}
				})

				It("clone should have empty Env", func() {
					expected := original.DeepCopy()
					expected.Env = []corev1.EnvVar{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Env is empty map and default Env is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = []corev1.EnvVar{}
					original = roleSpecOne.DeepCopy()
					original.Env = []corev1.EnvVar{}
				})

				It("clone should have empty Env", func() {
					expected := original.DeepCopy()
					expected.Env = []corev1.EnvVar{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Env is empty map and default Env is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = envOne
					original = roleSpecOne.DeepCopy()
					original.Env = []corev1.EnvVar{}
				})

				It("clone should have default Env", func() {
					expected := original.DeepCopy()
					expected.Env = defaults.Env
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Env is set map and default Env is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = nil
					original = roleSpecOne.DeepCopy()
					original.Env = envTwo
				})

				It("clone should have original Env", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Env is set and default Env is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = []corev1.EnvVar{}
					original = roleSpecOne.DeepCopy()
					original.Env = envTwo
				})

				It("clone should have original Env", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Env is set and default Env is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = envOne
					original = roleSpecOne.DeepCopy()
					original.Env = envTwo
				})

				It("clone should have merged Env", func() {
					expected := original.DeepCopy()
					expected.Env = []corev1.EnvVar{{Name: "Three", Value: "3"}, {Name: "Four", Value: "4"}, {Name: "One", Value: "1"}, {Name: "Two", Value: "2"}}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Env is set and default Env is set with overlapping keys", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Env = envOne
					original = roleSpecOne.DeepCopy()
					original.Env = []corev1.EnvVar{{Name: "Two", Value: "22"}, {Name: "Four", Value: "4"}}
				})

				It("clone should have merged Env with the original key taking precedence", func() {
					expected := original.DeepCopy()
					expected.Env = []corev1.EnvVar{{Name: "Two", Value: "22"}, {Name: "Four", Value: "4"}, {Name: "One", Value: "1"}}
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- JVM ----------------------------------------------------------------

		Context("JVM is merged", func() {
			var jvmOne = coherence.JVMSpec{
				HeapSize: stringPtr("1G"),
			}

			var jvmTwo = coherence.JVMSpec{
				GC: stringPtr("-XX:+UseG1GC"),
			}

			When("the original JVM is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					original = roleSpecOne.DeepCopy()
					original.JVM = nil
				})

				It("clone should be equal to the original with the JVM field from the defaults", func() {
					expected := original.DeepCopy()
					expected.JVM = defaults.JVM
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original JVM is set and defaults JVM is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.JVM = nil
					original = roleSpecOne.DeepCopy()
				})

				It("clone should be equal to the original with the JVM field from the original", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original JVM is set and defaults JVM is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.JVM = &jvmOne
					original = roleSpecOne.DeepCopy()
					original.JVM = &jvmTwo
				})

				It("clone should be equal to the merged original and defaults JVMs", func() {
					expected := original.DeepCopy()
					expected.JVM = original.JVM.DeepCopyWithDefaults(defaults.JVM)
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- Labels -------------------------------------------------------------

		Context("Labels are merged", func() {
			var labelsOne = map[string]string{"One": "1", "Two": "2"}
			var labelsTwo = map[string]string{"Three": "3", "Four": "4"}

			When("the original Labels is nil and default Labels is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = nil
					original = roleSpecOne.DeepCopy()
					original.Labels = nil
				})

				It("clone should have nil Affinity", func() {
					expected := original.DeepCopy()
					expected.Labels = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Labels is nil and default Labels is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.Labels = nil
				})

				It("clone should have empty Labels", func() {
					expected := original.DeepCopy()
					expected.Labels = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Labels is nil and default Labels is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = labelsOne
					original = roleSpecOne.DeepCopy()
					original.Labels = nil
				})

				It("clone should have default Labels", func() {
					expected := original.DeepCopy()
					expected.Labels = defaults.Labels
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Labels is empty map and default Labels is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = nil
					original = roleSpecOne.DeepCopy()
					original.Labels = map[string]string{}
				})

				It("clone should have empty Labels", func() {
					expected := original.DeepCopy()
					expected.Labels = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Labels is empty map and default Labels is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.Labels = map[string]string{}
				})

				It("clone should have empty Labels", func() {
					expected := original.DeepCopy()
					expected.Labels = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Labels is empty map and default Labels is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = labelsOne
					original = roleSpecOne.DeepCopy()
					original.Labels = map[string]string{}
				})

				It("clone should have default Labels", func() {
					expected := original.DeepCopy()
					expected.Labels = defaults.Labels
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Labels is set map and default Labels is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = nil
					original = roleSpecOne.DeepCopy()
					original.Labels = labelsTwo
				})

				It("clone should have original Labels", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Labels is set and default Labels is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.Labels = labelsTwo
				})

				It("clone should have original Labels", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Labels is set and default Labels is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Labels = labelsOne
					original = roleSpecOne.DeepCopy()
					original.Labels = labelsTwo
				})

				It("clone should have merged Labels", func() {
					expected := original.DeepCopy()
					expected.Labels = map[string]string{"One": "1", "Two": "2", "Three": "3", "Four": "4"}
					Expect(clone).To(Equal(expected))
				})

				When("the original Labels is set and default Labels is set with overlapping keys", func() {
					BeforeEach(func() {
						defaults = roleSpecTwo.DeepCopy()
						defaults.Labels = labelsOne
						original = roleSpecOne.DeepCopy()
						original.Labels = map[string]string{"Two": "22", "Four": "4"}
					})

					It("clone should have merged Labels with original key taking precedence", func() {
						expected := original.DeepCopy()
						expected.Labels = map[string]string{"One": "1", "Two": "22", "Four": "4"}
						Expect(clone).To(Equal(expected))
					})
				})
			})
		})

		// ----- Logging ------------------------------------------------------------

		Context("Logging is merged", func() {
			var logOne = coherence.LoggingSpec{
				ConfigMapName: stringPtr("cm-log"),
			}

			var logTwo = coherence.LoggingSpec{
				ConfigFile: stringPtr("logging.properties"),
			}

			When("the original Logging is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					original = roleSpecOne.DeepCopy()
					original.Logging = nil
				})

				It("clone should be equal to the original with the Logging field from the defaults", func() {
					expected := original.DeepCopy()
					expected.Logging = defaults.Logging
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Logging is set and defaults Logging is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Logging = nil
					original = roleSpecOne.DeepCopy()
				})

				It("clone should be equal to the original with the Logging field from the original", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Logging is set and defaults Logging is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Logging = &logOne
					original = roleSpecOne.DeepCopy()
					original.Logging = &logTwo
				})

				It("clone should be equal to the merged original and defaults Logging", func() {
					expected := original.DeepCopy()
					expected.Logging = original.Logging.DeepCopyWithDefaults(defaults.Logging)
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- NodeSelector -------------------------------------------------------

		Context("NodeSelector is not merged", func() {
			var selectorOne = map[string]string{"One": "1", "Two": "2"}
			var selectorTwo = map[string]string{"Three": "3", "Four": "4"}

			When("the original NodeSelector is nil and default NodeSelector is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = nil
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = nil
				})

				It("clone should have nil Affinity", func() {
					expected := original.DeepCopy()
					expected.NodeSelector = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original NodeSelector is nil and default NodeSelector is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = nil
				})

				It("clone should have empty NodeSelector", func() {
					expected := original.DeepCopy()
					expected.NodeSelector = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original NodeSelector is nil and default NodeSelector is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = selectorOne
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = nil
				})

				It("clone should have default NodeSelector", func() {
					expected := original.DeepCopy()
					expected.NodeSelector = defaults.NodeSelector
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original NodeSelector is empty map and default NodeSelector is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = nil
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = map[string]string{}
				})

				It("clone should have empty NodeSelector", func() {
					expected := original.DeepCopy()
					expected.NodeSelector = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original NodeSelector is empty map and default NodeSelector is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = map[string]string{}
				})

				It("clone should have empty NodeSelector", func() {
					expected := original.DeepCopy()
					expected.NodeSelector = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original NodeSelector is empty map and default NodeSelector is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = selectorOne
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = map[string]string{}
				})

				It("clone should have empty NodeSelector", func() {
					expected := original.DeepCopy()
					expected.NodeSelector = map[string]string{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original NodeSelector is set map and default NodeSelector is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = nil
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = selectorTwo
				})

				It("clone should have original NodeSelector", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original NodeSelector is set and default NodeSelector is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = map[string]string{}
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = selectorTwo
				})

				It("clone should have original NodeSelector", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original NodeSelector is set and default NodeSelector is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.NodeSelector = selectorOne
					original = roleSpecOne.DeepCopy()
					original.NodeSelector = selectorTwo
				})

				It("clone should have the original NodeSelector", func() {
					Expect(clone).To(Equal(original))
				})
			})
		})

		// ----- Ports --------------------------------------------------------------

		Context("Ports are merged", func() {
			var portOne = coherence.NamedPortSpec{
				Name:     "One",
				PortSpec: coherence.PortSpec{Port: 100},
			}
			var portTwo = coherence.NamedPortSpec{
				Name:     "Two",
				PortSpec: coherence.PortSpec{Port: 200},
			}
			var portTwoToo = coherence.NamedPortSpec{
				Name:     "Two",
				PortSpec: coherence.PortSpec{Port: 222},
			}
			var portThree = coherence.NamedPortSpec{
				Name:     "Three",
				PortSpec: coherence.PortSpec{Port: 300},
			}
			var portFour = coherence.NamedPortSpec{
				Name:     "Four",
				PortSpec: coherence.PortSpec{Port: 400},
			}

			var portsOne = []coherence.NamedPortSpec{portOne, portTwo}
			var portsTwo = []coherence.NamedPortSpec{portThree, portFour}

			When("the original Ports is nil and default Ports is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = nil
					original = roleSpecOne.DeepCopy()
					original.Ports = nil
				})

				It("clone should have nil Affinity", func() {
					expected := original.DeepCopy()
					expected.Ports = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Ports is nil and default Ports is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = []coherence.NamedPortSpec{}
					original = roleSpecOne.DeepCopy()
					original.Ports = nil
				})

				It("clone should have empty Ports", func() {
					expected := original.DeepCopy()
					expected.Ports = []coherence.NamedPortSpec{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Ports is nil and default Ports is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = portsOne
					original = roleSpecOne.DeepCopy()
					original.Ports = nil
				})

				It("clone should have default Ports", func() {
					expected := original.DeepCopy()
					expected.Ports = defaults.Ports
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Ports is empty map and default Ports is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = nil
					original = roleSpecOne.DeepCopy()
					original.Ports = []coherence.NamedPortSpec{}
				})

				It("clone should have empty Ports", func() {
					expected := original.DeepCopy()
					expected.Ports = []coherence.NamedPortSpec{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Ports is empty map and default Ports is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = []coherence.NamedPortSpec{}
					original = roleSpecOne.DeepCopy()
					original.Ports = []coherence.NamedPortSpec{}
				})

				It("clone should have empty Ports", func() {
					expected := original.DeepCopy()
					expected.Ports = []coherence.NamedPortSpec{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Ports is empty map and default Ports is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = portsOne
					original = roleSpecOne.DeepCopy()
					original.Ports = []coherence.NamedPortSpec{}
				})

				It("clone should have default Ports", func() {
					expected := original.DeepCopy()
					expected.Ports = defaults.Ports
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Ports is set map and default Ports is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = nil
					original = roleSpecOne.DeepCopy()
					original.Ports = portsTwo
				})

				It("clone should have original Ports", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Ports is set and default Ports is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = []coherence.NamedPortSpec{}
					original = roleSpecOne.DeepCopy()
					original.Ports = portsTwo
				})

				It("clone should have original Ports", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Ports is set and default Ports is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = portsOne
					original = roleSpecOne.DeepCopy()
					original.Ports = portsTwo
				})

				It("clone should have merged Ports", func() {
					expected := original.DeepCopy()
					expected.Ports = []coherence.NamedPortSpec{portThree, portFour, portOne, portTwo}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Ports is set and default Ports is set with overlapping keys", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Ports = portsOne
					original = roleSpecOne.DeepCopy()
					original.Ports = []coherence.NamedPortSpec{portTwoToo, portFour}
				})

				It("clone should have merged Ports with the original key taking precedence", func() {
					expected := original.DeepCopy()
					expected.Ports = []coherence.NamedPortSpec{portTwoToo, portFour, portOne}
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- ReadinessProbe -----------------------------------------------------

		Context("ReadinessProbe is merged", func() {
			var readyOne = coherence.ReadinessProbeSpec{
				InitialDelaySeconds: int32Ptr(10),
			}

			var readyTwo = coherence.ReadinessProbeSpec{
				TimeoutSeconds: int32Ptr(99),
			}

			When("the original ReadinessProbe is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					original = roleSpecOne.DeepCopy()
					original.ReadinessProbe = nil
				})

				It("clone should be equal to the original with the ReadinessProbe field from the defaults", func() {
					expected := original.DeepCopy()
					expected.ReadinessProbe = defaults.ReadinessProbe
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original ReadinessProbe is set and defaults ReadinessProbe is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.ReadinessProbe = nil
					original = roleSpecOne.DeepCopy()
				})

				It("clone should be equal to the original with the ReadinessProbe field from the original", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original ReadinessProbe is set and defaults ReadinessProbe is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.ReadinessProbe = &readyOne
					original = roleSpecOne.DeepCopy()
					original.ReadinessProbe = &readyTwo
				})

				It("clone should be equal to the merged original and defaults ReadinessProbes", func() {
					expected := original.DeepCopy()
					expected.ReadinessProbe = original.ReadinessProbe.DeepCopyWithDefaults(defaults.ReadinessProbe)
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- Replicas -----------------------------------------------------------

		When("the original Replicas is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo.DeepCopy()
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Replicas = nil
			})

			It("clone should be equal to the original with the Replicas field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Replicas = defaults.Replicas

				Expect(clone).To(Equal(expected))
			})
		})

		// ----- Resources ----------------------------------------------------------

		Context("Resources is not merged", func() {
			var resourceOne = corev1.ResourceRequirements{
				Limits: corev1.ResourceList{corev1.ResourceCPU: resource.Quantity{
					Format: "10Gi",
				}},
			}

			var resourceTwo = corev1.ResourceRequirements{
				Limits: corev1.ResourceList{corev1.ResourceCPU: resource.Quantity{
					Format: "90Gi",
				}},
			}

			When("the original Resources is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					original = roleSpecOne.DeepCopy()
					original.Resources = nil
				})

				It("clone should be equal to the original with the Resources field from the defaults", func() {
					expected := original.DeepCopy()
					expected.Resources = defaults.Resources
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original ReadinessProbe is set and defaults Resources is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Resources = nil
					original = roleSpecOne.DeepCopy()
				})

				It("clone should be equal to the original with the Resources field from the original", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original ReadinessProbe is set and defaults Resources is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Resources = &resourceOne
					original = roleSpecOne.DeepCopy()
					original.Resources = &resourceTwo
				})

				It("clone should be equal to the original Resources", func() {
					Expect(clone).To(Equal(original))
				})
			})
		})

		// ----- Role ---------------------------------------------------------------

		When("the original Role is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo.DeepCopy()
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Role = ""
			})

			It("clone should be equal to the original with the Role field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Role = defaults.Role

				Expect(clone).To(Equal(expected))
			})
		})

		// ----- Tolerations --------------------------------------------------------

		Context("Tolerations are not merged", func() {
			var tolerationsOne = []corev1.Toleration{{Key: "One"}, {Key: "Two"}}
			var tolerationsTwo = []corev1.Toleration{{Key: "Three"}, {Key: "Four"}}

			When("the original Tolerations is nil and defaults is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Tolerations = nil
					original = roleSpecOne.DeepCopy()
					original.Tolerations = nil
				})

				It("clone should have nil Tolerations", func() {
					expected := original.DeepCopy()
					expected.Tolerations = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Tolerations is nil and defaults is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Tolerations = tolerationsOne
					original = roleSpecOne.DeepCopy()
					original.Tolerations = nil
				})

				It("clone should be equal to the original with the Tolerations field from the defaults", func() {
					expected := original.DeepCopy()
					expected.Tolerations = defaults.Tolerations
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Tolerations is empty and defaults is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Tolerations = tolerationsOne
					original = roleSpecOne.DeepCopy()
					original.Tolerations = []corev1.Toleration{}
				})

				It("clone should be equal to the original with the Tolerations field from the defaults", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Tolerations is set and defaults is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Tolerations = tolerationsOne
					original = roleSpecOne.DeepCopy()
					original.Tolerations = tolerationsTwo
				})

				It("clone should be equal to the original with the Tolerations field from the defaults", func() {
					Expect(clone).To(Equal(original))
				})
			})
		})

		// ----- VolumeClaimTemplates -----------------------------------------------

		Context("VolumeClaimTemplates are merged", func() {
			var vcOne = corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{Name: "One"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: nil,
					VolumeName:  "1",
				},
			}

			var vcTwo = corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{Name: "Two"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: nil,
					VolumeName:  "2",
				},
			}

			var vcTwoToo = corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{Name: "Two"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: nil,
					VolumeName:  "2Too",
				},
			}

			var vcThree = corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{Name: "Three"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: nil,
					VolumeName:  "3",
				},
			}

			var vcFour = corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{Name: "Four"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: nil,
					VolumeName:  "4",
				},
			}

			var volsOne = []corev1.PersistentVolumeClaim{vcOne, vcTwo}
			var volsTwo = []corev1.PersistentVolumeClaim{vcThree, vcFour}

			When("the original VolumeClaimTemplates is nil and default VolumeClaimTemplates is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = nil
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = nil
				})

				It("clone should have nil VolumeClaimTemplates", func() {
					expected := original.DeepCopy()
					expected.VolumeClaimTemplates = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeClaimTemplates is nil and default VolumeClaimTemplates is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = nil
				})

				It("clone should have empty VolumeClaimTemplates", func() {
					expected := original.DeepCopy()
					expected.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeClaimTemplates is nil and default VolumeClaimTemplates is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = volsOne
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = nil
				})

				It("clone should have default VolumeClaimTemplates", func() {
					expected := original.DeepCopy()
					expected.VolumeClaimTemplates = defaults.VolumeClaimTemplates
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeClaimTemplates is empty map and default VolumeClaimTemplates is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = nil
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
				})

				It("clone should have empty VolumeClaimTemplates", func() {
					expected := original.DeepCopy()
					expected.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeClaimTemplates is empty map and default VolumeClaimTemplates is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
				})

				It("clone should have empty VolumeClaimTemplates", func() {
					expected := original.DeepCopy()
					expected.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeClaimTemplates is empty map and default VolumeClaimTemplates is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = volsOne
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
				})

				It("clone should have default VolumeClaimTemplates", func() {
					expected := original.DeepCopy()
					expected.VolumeClaimTemplates = defaults.VolumeClaimTemplates
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeClaimTemplates is set map and default VolumeClaimTemplates is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = nil
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = volsTwo
				})

				It("clone should have original VolumeClaimTemplates", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original VolumeClaimTemplates is set and default VolumeClaimTemplates is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{}
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = volsTwo
				})

				It("clone should have original VolumeClaimTemplates", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original VolumeClaimTemplates is set and default VolumeClaimTemplates is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = volsOne
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = volsTwo
				})

				It("clone should have merged VolumeClaimTemplates", func() {
					expected := original.DeepCopy()
					expected.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{vcThree, vcFour, vcOne, vcTwo}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeClaimTemplates is set and default VolumeClaimTemplates is set with overlapping keys", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeClaimTemplates = volsOne
					original = roleSpecOne.DeepCopy()
					original.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{vcTwoToo, vcFour}
				})

				It("clone should have merged VolumeClaimTemplates with the original key taking precedence", func() {
					expected := original.DeepCopy()
					expected.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{vcTwoToo, vcFour, vcOne}
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- VolumeMounts ------------------------------------------------------------

		Context("VolumeMounts are merged", func() {
			var vOne = corev1.VolumeMount{
				Name:      "One",
				MountPath: "/one",
			}

			var vTwo = corev1.VolumeMount{
				Name:      "Two",
				MountPath: "/two",
			}

			var vTwoToo = corev1.VolumeMount{
				Name:      "Two",
				MountPath: "/two-too",
			}

			var vThree = corev1.VolumeMount{
				Name:      "Three",
				MountPath: "/three",
			}

			var vFour = corev1.VolumeMount{
				Name:      "Four",
				MountPath: "/four",
			}

			var volsOne = []corev1.VolumeMount{vOne, vTwo}
			var volsTwo = []corev1.VolumeMount{vThree, vFour}

			When("the original VolumeMounts is nil and default VolumeMounts is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = nil
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = nil
				})

				It("clone should have nil VolumeMounts", func() {
					expected := original.DeepCopy()
					expected.VolumeMounts = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeMounts is nil and default VolumeMounts is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = []corev1.VolumeMount{}
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = nil
				})

				It("clone should have empty VolumeMounts", func() {
					expected := original.DeepCopy()
					expected.VolumeMounts = []corev1.VolumeMount{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeMounts is nil and default VolumeMounts is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = volsOne
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = nil
				})

				It("clone should have default VolumeMounts", func() {
					expected := original.DeepCopy()
					expected.VolumeMounts = defaults.VolumeMounts
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeMounts is empty map and default VolumeMounts is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = nil
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = []corev1.VolumeMount{}
				})

				It("clone should have empty VolumeMounts", func() {
					expected := original.DeepCopy()
					expected.VolumeMounts = []corev1.VolumeMount{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeMounts is empty map and default VolumeMounts is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = []corev1.VolumeMount{}
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = []corev1.VolumeMount{}
				})

				It("clone should have empty VolumeMounts", func() {
					expected := original.DeepCopy()
					expected.VolumeMounts = []corev1.VolumeMount{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeMounts is empty map and default VolumeMounts is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = volsOne
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = []corev1.VolumeMount{}
				})

				It("clone should have default VolumeMounts", func() {
					expected := original.DeepCopy()
					expected.VolumeMounts = defaults.VolumeMounts
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeMounts is set map and default VolumeMounts is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = nil
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = volsTwo
				})

				It("clone should have original VolumeMounts", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original VolumeMounts is set and default VolumeMounts is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = []corev1.VolumeMount{}
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = volsTwo
				})

				It("clone should have original VolumeMounts", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original VolumeMounts is set and default VolumeMounts is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = volsOne
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = volsTwo
				})

				It("clone should have merged VolumeMounts", func() {
					expected := original.DeepCopy()
					expected.VolumeMounts = []corev1.VolumeMount{vThree, vFour, vOne, vTwo}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original VolumeMounts is set and default VolumeMounts is set with overlapping keys", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.VolumeMounts = volsOne
					original = roleSpecOne.DeepCopy()
					original.VolumeMounts = []corev1.VolumeMount{vTwoToo, vFour}
				})

				It("clone should have merged VolumeMounts with the original key taking precedence", func() {
					expected := original.DeepCopy()
					expected.VolumeMounts = []corev1.VolumeMount{vTwoToo, vFour, vOne}
					Expect(clone).To(Equal(expected))
				})
			})
		})

		// ----- Volumes ------------------------------------------------------------

		Context("Volumes are merged", func() {
			var vOne = corev1.Volume{
				Name: "One",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "one",
					},
				},
			}

			var vTwo = corev1.Volume{
				Name: "Two",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "two",
					},
				},
			}

			var vTwoToo = corev1.Volume{
				Name: "Two",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "two-too",
					},
				},
			}

			var vThree = corev1.Volume{
				Name: "Three",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "three",
					},
				},
			}

			var vFour = corev1.Volume{
				Name: "FOur",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "four",
					},
				},
			}

			var volsOne = []corev1.Volume{vOne, vTwo}
			var volsTwo = []corev1.Volume{vThree, vFour}

			When("the original Volumes is nil and default Volumes is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = nil
					original = roleSpecOne.DeepCopy()
					original.Volumes = nil
				})

				It("clone should have nil Volumes", func() {
					expected := original.DeepCopy()
					expected.Volumes = nil
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Volumes is nil and default Volumes is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = []corev1.Volume{}
					original = roleSpecOne.DeepCopy()
					original.Volumes = nil
				})

				It("clone should have empty Volumes", func() {
					expected := original.DeepCopy()
					expected.Volumes = []corev1.Volume{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Volumes is nil and default Volumes is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = volsOne
					original = roleSpecOne.DeepCopy()
					original.Volumes = nil
				})

				It("clone should have default Volumes", func() {
					expected := original.DeepCopy()
					expected.Volumes = defaults.Volumes
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Volumes is empty map and default Volumes is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = nil
					original = roleSpecOne.DeepCopy()
					original.Volumes = []corev1.Volume{}
				})

				It("clone should have empty Volumes", func() {
					expected := original.DeepCopy()
					expected.Volumes = []corev1.Volume{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Volumes is empty map and default Volumes is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = []corev1.Volume{}
					original = roleSpecOne.DeepCopy()
					original.Volumes = []corev1.Volume{}
				})

				It("clone should have empty Volumes", func() {
					expected := original.DeepCopy()
					expected.Volumes = []corev1.Volume{}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Volumes is empty map and default Volumes is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = volsOne
					original = roleSpecOne.DeepCopy()
					original.Volumes = []corev1.Volume{}
				})

				It("clone should have default Volumes", func() {
					expected := original.DeepCopy()
					expected.Volumes = defaults.Volumes
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Volumes is set map and default Volumes is nil", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = nil
					original = roleSpecOne.DeepCopy()
					original.Volumes = volsTwo
				})

				It("clone should have original Volumes", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Volumes is set and default Volumes is empty", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = []corev1.Volume{}
					original = roleSpecOne.DeepCopy()
					original.Volumes = volsTwo
				})

				It("clone should have original Volumes", func() {
					Expect(clone).To(Equal(original))
				})
			})

			When("the original Volumes is set and default Volumes is set", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = volsOne
					original = roleSpecOne.DeepCopy()
					original.Volumes = volsTwo
				})

				It("clone should have merged Volumes", func() {
					expected := original.DeepCopy()
					expected.Volumes = []corev1.Volume{vThree, vFour, vOne, vTwo}
					Expect(clone).To(Equal(expected))
				})
			})

			When("the original Volumes is set and default Volumes is set with overlapping keys", func() {
				BeforeEach(func() {
					defaults = roleSpecTwo.DeepCopy()
					defaults.Volumes = volsOne
					original = roleSpecOne.DeepCopy()
					original.Volumes = []corev1.Volume{vTwoToo, vFour}
				})

				It("clone should have merged Volumes with the original key taking precedence", func() {
					expected := original.DeepCopy()
					expected.Volumes = []corev1.Volume{vTwoToo, vFour, vOne}
					Expect(clone).To(Equal(expected))
				})
			})
		})

	})

	// ----- Methods ------------------------------------------------------------

	Context("Getting Replica count", func() {
		var role coherence.CoherenceRoleSpec
		var replicas *int32

		JustBeforeEach(func() {
			role = coherence.CoherenceRoleSpec{Replicas: replicas}
		})

		When("Replicas is not set", func() {
			BeforeEach(func() {
				replicas = nil
			})

			It("should return the default replica count", func() {
				Expect(role.GetReplicas()).To(Equal(coherence.DefaultReplicas))
			})
		})

		When("Replicas is set", func() {
			BeforeEach(func() {
				replicas = int32Ptr(100)
			})

			It("should return the specified replica count", func() {
				Expect(role.GetReplicas()).To(Equal(*replicas))
			})
		})
	})

	When("Getting the full role name", func() {
		var cluster coherence.CoherenceCluster
		var role coherence.CoherenceRoleSpec

		BeforeEach(func() {
			cluster = coherence.CoherenceCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
			}

			role = coherence.CoherenceRoleSpec{
				Role: "storage",
			}
		})

		It("should return the specified replica count", func() {
			Expect(role.GetFullRoleName(&cluster)).To(Equal("test-cluster-storage"))
		})
	})

	Context("Getting the role name", func() {
		var role *coherence.CoherenceRoleSpec
		var name string

		JustBeforeEach(func() {
			name = role.GetRoleName()
		})

		When("role is not set", func() {
			BeforeEach(func() {
				role = &coherence.CoherenceRoleSpec{Role: ""}
			})

			It("should use the default name", func() {
				Expect(name).To(Equal(coherence.DefaultRoleName))
			})
		})

		When("role is set", func() {
			BeforeEach(func() {
				role = &coherence.CoherenceRoleSpec{Role: "test-role"}
			})

			It("should use the default name", func() {
				Expect(name).To(Equal("test-role"))
			})
		})

		When("role is nil", func() {
			BeforeEach(func() {
				role = nil
			})

			It("should use the default name", func() {
				Expect(name).To(Equal(coherence.DefaultRoleName))
			})
		})

	})

})
