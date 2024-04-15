package bindings_test

import (
	"encoding/json"
	"fmt"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/controllers/services/bindings"
	. "code.cloudfoundry.org/korifi/controllers/controllers/workloads/testutils"
	. "code.cloudfoundry.org/korifi/tests/matchers"
	"code.cloudfoundry.org/korifi/tools"
	"code.cloudfoundry.org/korifi/tools/k8s"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	servicebindingv1beta1 "github.com/servicebinding/runtime/apis/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("CFServiceBinding", func() {
	var (
		testNamespace             string
		cfApp                     *korifiv1alpha1.CFApp
		instance                  *korifiv1alpha1.CFServiceInstance
		binding                   *korifiv1alpha1.CFServiceBinding
		instanceCredentialsSecret *corev1.Secret
	)

	BeforeEach(func() {
		testNamespace = uuid.NewString()
		Expect(adminClient.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: testNamespace,
			},
		})).To(Succeed())

		cfApp = BuildCFAppCRObject(uuid.NewString(), testNamespace)
		Expect(
			adminClient.Create(ctx, cfApp),
		).To(Succeed())

		Expect(k8s.Patch(ctx, adminClient, cfApp, func() {
			cfApp.Status = korifiv1alpha1.CFAppStatus{
				VCAPServicesSecretName: "foo",
			}
		})).To(Succeed())

		credentialsBytes, err := json.Marshal(map[string]any{
			"type":     "my-type",
			"provider": "my-provider",
			"obj": map[string]any{
				"foo": "bar",
			},
		})
		Expect(err).NotTo(HaveOccurred())
		instanceCredentialsSecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: testNamespace,
			},
			Data: map[string][]byte{
				korifiv1alpha1.CredentialsSecretKey: credentialsBytes,
			},
		}

		Expect(adminClient.Create(ctx, instanceCredentialsSecret)).To(Succeed())

		instance = &korifiv1alpha1.CFServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: testNamespace,
			},
			Spec: korifiv1alpha1.CFServiceInstanceSpec{
				DisplayName: "mongodb-service-instance-name",
				Type:        "user-provided",
				Tags:        []string{},
			},
		}
		Expect(adminClient.Create(ctx, instance)).To(Succeed())
		Expect(k8s.Patch(ctx, adminClient, instance, func() {
			instance.Status.Credentials.Name = instanceCredentialsSecret.Name
		})).To(Succeed())

		binding = &korifiv1alpha1.CFServiceBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: testNamespace,
			},
			Spec: korifiv1alpha1.CFServiceBindingSpec{
				Service: corev1.ObjectReference{
					Kind:       "ServiceInstance",
					Name:       instance.Name,
					APIVersion: "korifi.cloudfoundry.org/v1alpha1",
				},
				AppRef: corev1.LocalObjectReference{
					Name: cfApp.Name,
				},
			},
		}
		Expect(adminClient.Create(ctx, binding)).To(Succeed())
	})

	It("sets the ObservedGeneration status field", func() {
		Eventually(func(g Gomega) {
			g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
			g.Expect(binding.Status.ObservedGeneration).To(Equal(binding.Generation))
		}).Should(Succeed())
	})

	It("sets an owner reference from the instance to the binding", func() {
		Eventually(func(g Gomega) {
			g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
			g.Expect(binding.OwnerReferences).To(ConsistOf(MatchFields(IgnoreExtras, Fields{
				"Name": Equal(instance.Name),
			})))
		}).Should(Succeed())
	})

	It("sets the BindingSecretAvailable condition to true", func() {
		Eventually(func(g Gomega) {
			g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
			g.Expect(meta.IsStatusConditionTrue(binding.Status.Conditions, bindings.BindingSecretAvailableCondition)).To(BeTrue())
		}).Should(Succeed())
	})

	It("sets the binding status credentials name to the instance credentials secret", func() {
		Eventually(func(g Gomega) {
			g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
			g.Expect(binding.Status.Credentials.Name).To(Equal(instanceCredentialsSecret.Name))
		}).Should(Succeed())
	})

	It("creates the binding secret", func() {
		Eventually(func(g Gomega) {
			g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
			g.Expect(binding.Status.Binding.Name).To(Equal(binding.Name))

			bindingSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: binding.Namespace,
					Name:      binding.Status.Binding.Name,
				},
			}
			g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(bindingSecret), bindingSecret)).To(Succeed())
			g.Expect(bindingSecret.Type).To(BeEquivalentTo("servicebinding.io/my-type"))
			g.Expect(bindingSecret.Data).To(MatchAllKeys(Keys{
				"type":     Equal([]byte("my-type")),
				"provider": Equal([]byte("my-provider")),
				"obj":      Equal([]byte(`{"foo":"bar"}`)),
			}))
		}).Should(Succeed())
	})

	It("sets the binding status binding name to the binding secret name", func() {
		Eventually(func(g Gomega) {
			g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
			g.Expect(binding.Status.Binding.Name).To(Equal(binding.Name))
		}).Should(Succeed())
	})

	It("creates a servicebinding.io ServiceBinding", func() {
		Eventually(func(g Gomega) {
			sbServiceBinding := &servicebindingv1beta1.ServiceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: testNamespace,
					Name:      fmt.Sprintf("cf-binding-%s", binding.Name),
				},
			}
			g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(sbServiceBinding), sbServiceBinding)).To(Succeed())

			g.Expect(sbServiceBinding.Spec.Name).To(Equal(binding.Name))
			g.Expect(sbServiceBinding.Spec.Type).To(Equal("my-type"))

			g.Expect(sbServiceBinding.Labels).To(SatisfyAll(
				HaveKeyWithValue(bindings.ServiceBindingGUIDLabel, binding.Name),
				HaveKeyWithValue(korifiv1alpha1.CFAppGUIDLabelKey, cfApp.Name),
				HaveKeyWithValue(bindings.ServiceCredentialBindingTypeLabel, "app"),
			))

			g.Expect(sbServiceBinding.OwnerReferences).To(ConsistOf(MatchFields(IgnoreExtras, Fields{
				"Kind": Equal("CFServiceBinding"),
				"Name": Equal(binding.Name),
			})))

			g.Expect(sbServiceBinding.Spec.Workload).To(MatchFields(IgnoreExtras, Fields{
				"APIVersion": Equal("apps/v1"),
				"Kind":       Equal("StatefulSet"),
				"Selector": PointTo(Equal(metav1.LabelSelector{
					MatchLabels: map[string]string{
						korifiv1alpha1.CFAppGUIDLabelKey: cfApp.Name,
					},
				})),
			}))

			g.Expect(sbServiceBinding.Spec.Service).To(MatchFields(IgnoreExtras, Fields{
				"APIVersion": Equal("korifi.cloudfoundry.org/v1alpha1"),
				"Kind":       Equal("CFServiceBinding"),
				"Name":       Equal(binding.Name),
			}))
		}).Should(Succeed())
	})

	When("the credentials secret is not available", func() {
		BeforeEach(func() {
			Expect(k8s.Patch(ctx, adminClient, instance, func() {
				instance.Status.Credentials.Name = ""
			})).To(Succeed())
		})

		It("sets the BindingSecretAvailable condition to false", func() {
			Eventually(func(g Gomega) {
				g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
				g.Expect(binding.Status.Conditions).To(ContainElement(SatisfyAll(
					HasType(Equal(bindings.BindingSecretAvailableCondition)),
					HasStatus(Equal(metav1.ConditionFalse)),
					HasReason(Equal("CredentialsSecretNotAvailable")),
				)))
			}).Should(Succeed())
		})
	})

	When("the CFServiceBinding has a displayName set", func() {
		BeforeEach(func() {
			Expect(k8s.PatchResource(ctx, adminClient, binding, func() {
				binding.Spec.DisplayName = tools.PtrTo("a-custom-binding-name")
			})).To(Succeed())
		})

		It("sets the displayName as the name on the servicebinding.io ServiceBinding", func() {
			Eventually(func(g Gomega) {
				sbServiceBinding := &servicebindingv1beta1.ServiceBinding{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testNamespace,
						Name:      fmt.Sprintf("cf-binding-%s", binding.Name),
					},
				}
				g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(sbServiceBinding), sbServiceBinding)).To(Succeed())
				g.Expect(sbServiceBinding.Spec.Name).To(Equal("a-custom-binding-name"))
			}).Should(Succeed())
		})
	})

	When("the binding references a 'legacy' instance credentials secret", func() {
		JustBeforeEach(func() {
			Expect(k8s.Patch(ctx, adminClient, instance, func() {
				instance.Spec.SecretName = instance.Name
				instance.Status.Credentials.Name = instance.Name
			})).To(Succeed())

			Expect(k8s.Patch(ctx, adminClient, binding, func() {
				binding.Status.Binding.Name = instance.Name
			})).To(Succeed())
		})

		It("sets credentials secret not available condition", func() {
			Eventually(func(g Gomega) {
				g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
				g.Expect(binding.Status.Conditions).To(ContainElement(SatisfyAll(
					HasType(Equal(bindings.BindingSecretAvailableCondition)),
					HasStatus(Equal(metav1.ConditionFalse)),
					HasReason(Equal("FailedReconcilingCredentialsSecret")),
				)))
			}).Should(Succeed())
		})

		When("the referenced legacy binding secret exists", func() {
			BeforeEach(func() {
				Expect(adminClient.Create(ctx, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      instance.Name,
						Namespace: testNamespace,
					},
				})).To(Succeed())
			})

			It("does not update the binding status", func() {
				Eventually(func(g Gomega) {
					g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
					g.Expect(binding.Status.Binding.Name).To(Equal(instance.Name))
				}).Should(Succeed())
				Consistently(func(g Gomega) {
					g.Expect(adminClient.Get(ctx, client.ObjectKeyFromObject(binding), binding)).To(Succeed())
					g.Expect(binding.Status.Binding.Name).To(Equal(instance.Name))
				}).Should(Succeed())
			})
		})
	})
})
