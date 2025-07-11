package repositories_test

import (
	"encoding/json"
	"errors"
	"time"

	apierrors "code.cloudfoundry.org/korifi/api/errors"
	"code.cloudfoundry.org/korifi/api/repositories"
	"code.cloudfoundry.org/korifi/api/repositories/fake"
	"code.cloudfoundry.org/korifi/api/repositories/fakeawaiter"
	"code.cloudfoundry.org/korifi/api/repositories/k8sklient/descriptors"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/env"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/testutils"
	"code.cloudfoundry.org/korifi/tests/matchers"
	"code.cloudfoundry.org/korifi/tools"
	"code.cloudfoundry.org/korifi/tools/k8s"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	CFAppRevisionKey = "korifi.cloudfoundry.org/app-rev"

	CFAppDisplayNameLabel = "korifi.cloudfoundry.org/display-name"
	CFAppRevisionValue    = "1"
)

var _ = Describe("AppRepository", func() {
	var (
		appAwaiter *fakeawaiter.FakeAwaiter[
			*korifiv1alpha1.CFApp,
			korifiv1alpha1.CFAppList,
			*korifiv1alpha1.CFAppList,
		]
		appRepo *repositories.AppRepo
		cfOrg   *korifiv1alpha1.CFOrg
		cfSpace *korifiv1alpha1.CFSpace
		cfApp   *korifiv1alpha1.CFApp
	)

	BeforeEach(func() {
		appAwaiter = &fakeawaiter.FakeAwaiter[
			*korifiv1alpha1.CFApp,
			korifiv1alpha1.CFAppList,
			*korifiv1alpha1.CFAppList,
		]{}
		appRepo = repositories.NewAppRepo(spaceScopedKlient, appAwaiter)

		cfOrg = createOrgWithCleanup(ctx, prefixedGUID("org"))
		cfSpace = createSpaceWithCleanup(ctx, cfOrg.Name, prefixedGUID("space1"))

		cfApp = createAppWithGUID(cfSpace.Name, testutils.PrefixedGUID("cfapp1-"))
	})

	Describe("GetApp", func() {
		var (
			appGUID string
			app     repositories.AppRecord
			getErr  error
		)

		BeforeEach(func() {
			appGUID = cfApp.Name
		})

		JustBeforeEach(func() {
			app, getErr = appRepo.GetApp(ctx, authInfo, appGUID)
		})

		When("authorized in the space", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, orgUserRole.Name, cfOrg.Name)
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
			})

			It("can fetch the AppRecord CR we're looking for", func() {
				Expect(getErr).NotTo(HaveOccurred())

				Expect(app.GUID).To(Equal(cfApp.Name))
				Expect(app.EtcdUID).To(Equal(cfApp.GetUID()))
				Expect(app.Revision).To(Equal(CFAppRevisionValue))
				Expect(app.Name).To(Equal(cfApp.Spec.DisplayName))
				Expect(app.SpaceGUID).To(Equal(cfSpace.Name))
				Expect(app.State).To(Equal(repositories.DesiredState("STOPPED")))
				Expect(app.DropletGUID).To(Equal(cfApp.Spec.CurrentDropletRef.Name))
				Expect(app.Lifecycle).To(Equal(repositories.Lifecycle{
					Type: string(cfApp.Spec.Lifecycle.Type),
					Data: repositories.LifecycleData{
						Buildpacks: cfApp.Spec.Lifecycle.Data.Buildpacks,
						Stack:      cfApp.Spec.Lifecycle.Data.Stack,
					},
				}))
				Expect(app.IsStaged).To(BeTrue())
				Expect(app.DeletedAt).To(BeNil())

				Expect(app.Relationships()).To(Equal(map[string]string{
					"space": app.SpaceGUID,
				}))
			})

			When("the app has no current droplet set", func() {
				BeforeEach(func() {
					Expect(k8s.PatchResource(ctx, k8sClient, cfApp, func() {
						cfApp.Spec.CurrentDropletRef.Name = ""
					})).To(Succeed())
				})

				It("sets IsStaged to false", func() {
					Expect(getErr).ToNot(HaveOccurred())
					Expect(app.IsStaged).To(BeFalse())
				})
			})
		})

		When("the user is not authorized in the space", func() {
			It("returns a forbidden error", func() {
				Expect(getErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.ForbiddenError{}))
			})
		})

		When("duplicate Apps exist across namespaces with the same GUIDs", func() {
			BeforeEach(func() {
				space2 := createSpaceWithCleanup(ctx, cfOrg.Name, prefixedGUID("space2"))
				createAppWithGUID(space2.Name, appGUID)
			})

			It("returns an error", func() {
				Expect(getErr).To(HaveOccurred())
				Expect(getErr).To(MatchError(ContainSubstring("get-app duplicate records exist")))
			})
		})

		When("the app guid is not found", func() {
			BeforeEach(func() {
				appGUID = "does-not-exist"
			})

			It("returns an error", func() {
				Expect(getErr).To(HaveOccurred())
				Expect(getErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.NotFoundError{}))
			})
		})
	})

	Describe("ListApps", func() {
		var (
			message    repositories.ListAppsMessage
			listResult repositories.ListResult[repositories.AppRecord]
			cfApp2     *korifiv1alpha1.CFApp
			listErr    error
		)

		BeforeEach(func() {
			message = repositories.ListAppsMessage{}

			space2 := createSpaceWithCleanup(ctx, cfOrg.Name, prefixedGUID("space2"))
			createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
			createRoleBinding(ctx, userName, spaceDeveloperRole.Name, space2.Name)

			cfApp2 = createAppWithGUID(space2.Name, testutils.PrefixedGUID("cfapp2-"))
		})

		JustBeforeEach(func() {
			listResult, listErr = appRepo.ListApps(ctx, authInfo, message)
		})

		It("lists the apps", func() {
			Expect(listErr).NotTo(HaveOccurred())
			Expect(listResult.Records).To(ConsistOf(
				MatchFields(IgnoreExtras, Fields{"GUID": Equal(cfApp.Name)}),
				MatchFields(IgnoreExtras, Fields{"GUID": Equal(cfApp2.Name)}),
			))
			Expect(listResult.PageInfo).To(Equal(descriptors.PageInfo{
				TotalResults: 2,
				TotalPages:   1,
				PageNumber:   1,
				PageSize:     2,
			}))
		})

		Describe("list options", func() {
			var fakeKlient *fake.Klient

			BeforeEach(func() {
				fakeKlient = new(fake.Klient)
				appRepo = repositories.NewAppRepo(fakeKlient, appAwaiter)
			})

			Describe("parameters to list options", func() {
				BeforeEach(func() {
					message = repositories.ListAppsMessage{
						Names:         []string{"n1", "n2"},
						Guids:         []string{"g1", "g2"},
						SpaceGUIDs:    []string{"sg1", "sg2"},
						LabelSelector: "foo=bar",
						OrderBy:       "created_at",
						Pagination: repositories.Pagination{
							Page:    3,
							PerPage: 4,
						},
					}
				})

				It("translates parameters to klient list options", func() {
					Expect(listErr).NotTo(HaveOccurred())
					Expect(fakeKlient.ListCallCount()).To(Equal(1))
					_, _, listOptions := fakeKlient.ListArgsForCall(0)
					Expect(listOptions).To(ConsistOf(
						repositories.WithLabelIn(korifiv1alpha1.DisplayNameLabelKey, tools.EncodeValuesToSha224("n1", "n2")),
						repositories.WithLabelIn(korifiv1alpha1.GUIDLabelKey, []string{"g1", "g2"}),
						repositories.WithLabelIn(korifiv1alpha1.SpaceGUIDLabelKey, []string{"sg1", "sg2"}),
						repositories.WithLabelSelector("foo=bar"),
						repositories.WithOrdering("created_at"),
						repositories.WithPaging(repositories.Pagination{PerPage: 4, Page: 3}),
					))
				})
			})
		})

		DescribeTable("ordering",
			func(msg repositories.ListAppsMessage, match types.GomegaMatcher) {
				fakeKlient := new(fake.Klient)
				appRepo = repositories.NewAppRepo(fakeKlient, appAwaiter)

				_, err := appRepo.ListApps(ctx, authInfo, msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeKlient.ListCallCount()).To(Equal(1))
				_, _, listOptions := fakeKlient.ListArgsForCall(0)
				Expect(listOptions).To(match)
			},
			Entry("name", repositories.ListAppsMessage{OrderBy: "name"}, ContainElement(repositories.SortOpt{By: "Display Name", Desc: false})),
			Entry("-name", repositories.ListAppsMessage{OrderBy: "-name"}, ContainElement(repositories.SortOpt{By: "Display Name", Desc: true})),
			Entry("state", repositories.ListAppsMessage{OrderBy: "state"}, ContainElement(repositories.SortOpt{By: "State", Desc: false})),
			Entry("-state", repositories.ListAppsMessage{OrderBy: "-state"}, ContainElement(repositories.SortOpt{By: "State", Desc: true})),
		)
	})

	Describe("CreateApp", func() {
		const (
			testAppName = "test-app-name"
		)
		var (
			appCreateMessage repositories.CreateAppMessage
			createdAppRecord repositories.AppRecord
			createErr        error
		)

		BeforeEach(func() {
			appCreateMessage = repositories.CreateAppMessage{
				Name:      testAppName,
				SpaceGUID: cfSpace.Name,
				State:     "STOPPED",
				Lifecycle: repositories.Lifecycle{
					Type: "buildpack",
					Data: repositories.LifecycleData{
						Buildpacks: []string{},
						Stack:      "cflinuxfs3",
					},
				},
			}
		})

		JustBeforeEach(func() {
			createdAppRecord, createErr = appRepo.CreateApp(ctx, authInfo, appCreateMessage)
		})

		When("authorized in the space", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, orgUserRole.Name, cfOrg.Name)
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
			})

			It("creates a new app CR", func() {
				Expect(createErr).NotTo(HaveOccurred())
				Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())
			})

			It("returns an AppRecord with correct fields", func() {
				Expect(createErr).NotTo(HaveOccurred())
				Expect(createdAppRecord.GUID).To(MatchRegexp("^[-0-9a-f]{36}$"))
				Expect(createdAppRecord.SpaceGUID).To(Equal(cfSpace.Name))
				Expect(createdAppRecord.Name).To(Equal(testAppName))
				Expect(createdAppRecord.Lifecycle.Type).To(Equal("buildpack"))
				Expect(createdAppRecord.Lifecycle.Data.Buildpacks).To(BeEmpty())

				Expect(createdAppRecord.CreatedAt).To(BeTemporally("~", time.Now(), timeCheckThreshold))
			})

			It("creates an empty secret and sets the environment variable secret ref on the App", func() {
				Expect(createErr).NotTo(HaveOccurred())
				cfApp = &korifiv1alpha1.CFApp{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: createdAppRecord.SpaceGUID,
						Name:      createdAppRecord.GUID,
					},
				}
				Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())

				createdSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: cfSpace.Name,
						Name:      cfApp.Spec.EnvSecretName,
					},
				}
				Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(createdSecret), createdSecret)).To(Succeed())
				Expect(createdSecret.Data).To(BeEmpty())
			})

			When("environment variables are given", func() {
				BeforeEach(func() {
					appCreateMessage.EnvironmentVariables = map[string]string{
						"FOO": "foo",
						"BAR": "bar",
					}
				})

				It("creates an secret for the environment variables and sets the ref on the App", func() {
					Expect(createErr).NotTo(HaveOccurred())

					cfApp = &korifiv1alpha1.CFApp{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: createdAppRecord.SpaceGUID,
							Name:      createdAppRecord.GUID,
						},
					}
					Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())
					Expect(cfApp.Spec.EnvSecretName).NotTo(BeEmpty())

					createdSecret := &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: cfSpace.Name,
							Name:      cfApp.Spec.EnvSecretName,
						},
					}
					Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(createdSecret), createdSecret)).To(Succeed())
					Expect(createdSecret.Data).To(MatchAllKeys(Keys{
						"FOO": BeEquivalentTo("foo"),
						"BAR": BeEquivalentTo("bar"),
					}))
				})
			})

			When("buildpacks are given", func() {
				var buildpacks []string

				BeforeEach(func() {
					buildpacks = []string{"buildpack-1", "buildpack-2"}
					appCreateMessage.Lifecycle.Data.Buildpacks = buildpacks
				})

				It("creates a CFApp with the buildpacks set", func() {
					Expect(createErr).NotTo(HaveOccurred())
					cfApp = &korifiv1alpha1.CFApp{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: createdAppRecord.SpaceGUID,
							Name:      createdAppRecord.GUID,
						},
					}
					Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())
					Expect(cfApp.Spec.Lifecycle.Data.Buildpacks).To(Equal(buildpacks))
				})

				It("returns an AppRecord with the buildpacks set", func() {
					Expect(createdAppRecord.Lifecycle.Data.Buildpacks).To(Equal(buildpacks))
				})
			})

			When("the lifecycle is docker", func() {
				BeforeEach(func() {
					appCreateMessage.Lifecycle = repositories.Lifecycle{
						Type: "docker",
					}
				})

				It("creates an app with docker lifecycle", func() {
					Expect(createErr).NotTo(HaveOccurred())
					Expect(createdAppRecord.Lifecycle).To(Equal(repositories.Lifecycle{
						Type: "docker",
					}))

					cfApp = &korifiv1alpha1.CFApp{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: createdAppRecord.SpaceGUID,
							Name:      createdAppRecord.GUID,
						},
					}

					Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())
					Expect(cfApp.Spec.Lifecycle).To(Equal(korifiv1alpha1.Lifecycle{
						Type: "docker",
					}))
				})
			})
		})

		When("the user is not authorized in the space", func() {
			It("returns a forbidden error", func() {
				Expect(createErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.ForbiddenError{}))
			})
		})
	})

	Describe("PatchApp", func() {
		var (
			patchedAppRecord repositories.AppRecord
			patchErr         error

			appPatchMessage repositories.PatchAppMessage
		)

		BeforeEach(func() {
			Expect(k8sClient.Create(ctx, &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: cfApp.Namespace,
					Name:      cfApp.Spec.EnvSecretName,
				},
			})).To(Succeed())

			appPatchMessage = repositories.PatchAppMessage{
				Name:      cfApp.Spec.DisplayName,
				AppGUID:   cfApp.Name,
				SpaceGUID: cfSpace.Name,
				Lifecycle: &repositories.LifecyclePatch{
					Type: tools.PtrTo("docker"),
					Data: &repositories.LifecycleDataPatch{
						Buildpacks: &[]string{
							"some-buildpack",
						},
						Stack: "cflinuxfs3",
					},
				},
				MetadataPatch: repositories.MetadataPatch{
					Labels:      map[string]*string{"l": tools.PtrTo("lv")},
					Annotations: map[string]*string{"a": tools.PtrTo("av")},
				},
			}
		})

		JustBeforeEach(func() {
			patchedAppRecord, patchErr = appRepo.PatchApp(ctx, authInfo, appPatchMessage)
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())
		})

		When("authorized in the space", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, orgUserRole.Name, cfOrg.Name)
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
			})

			It("updates the app", func() {
				Expect(patchErr).NotTo(HaveOccurred())

				Expect(patchedAppRecord.GUID).To(Equal(cfApp.Name))
				Expect(patchedAppRecord.SpaceGUID).To(Equal(cfSpace.Name))
				Expect(patchedAppRecord.Name).To(Equal(appPatchMessage.Name))
				Expect(patchedAppRecord.Lifecycle).To(Equal(repositories.Lifecycle{
					Type: "docker",
					Data: repositories.LifecycleData{
						Buildpacks: []string{"some-buildpack"},
						Stack:      "cflinuxfs3",
					},
				}))

				Expect(cfApp.Spec.DisplayName).To(Equal(appPatchMessage.Name))
				Expect(cfApp.Spec.Lifecycle).To(Equal(korifiv1alpha1.Lifecycle{
					Type: "docker",
					Data: korifiv1alpha1.LifecycleData{
						Buildpacks: []string{"some-buildpack"},
						Stack:      "cflinuxfs3",
					},
				}))
				Expect(cfApp.Labels).To(HaveKeyWithValue("l", "lv"))
				Expect(cfApp.Annotations).To(HaveKeyWithValue("a", "av"))
			})

			Describe("partially patching the app", func() {
				var originalCFApp *korifiv1alpha1.CFApp

				BeforeEach(func() {
					originalCFApp = cfApp.DeepCopy()
				})

				When("name is empty", func() {
					BeforeEach(func() {
						appPatchMessage.Name = ""
					})

					It("does not change the display name", func() {
						Expect(patchedAppRecord.Name).To(Equal(originalCFApp.Spec.DisplayName))
						Expect(cfApp.Spec.DisplayName).To(Equal(originalCFApp.Spec.DisplayName))
					})
				})

				When("lifecycle is not specified", func() {
					BeforeEach(func() {
						appPatchMessage.Lifecycle = nil
					})

					It("does not change the app lifecyle", func() {
						Expect(patchedAppRecord.Lifecycle).To(Equal(repositories.Lifecycle{
							Type: string(originalCFApp.Spec.Lifecycle.Type),
							Data: repositories.LifecycleData{
								Buildpacks: originalCFApp.Spec.Lifecycle.Data.Buildpacks,
								Stack:      originalCFApp.Spec.Lifecycle.Data.Stack,
							},
						}))

						Expect(cfApp.Spec.Lifecycle).To(Equal(originalCFApp.Spec.Lifecycle))
					})
				})

				When("buildpacks are not specified", func() {
					BeforeEach(func() {
						appPatchMessage.Lifecycle.Data.Buildpacks = nil
					})

					It("does not change the app lifecyle buildpacks", func() {
						Expect(patchedAppRecord.Lifecycle.Data.Buildpacks).To(Equal(originalCFApp.Spec.Lifecycle.Data.Buildpacks))
						Expect(cfApp.Spec.Lifecycle.Data.Buildpacks).To(Equal(originalCFApp.Spec.Lifecycle.Data.Buildpacks))
					})
				})

				When("buildpacks are empty", func() {
					BeforeEach(func() {
						appPatchMessage.Lifecycle.Data.Buildpacks = &[]string{}
					})

					It("clears the app buildpacks", func() {
						Expect(patchedAppRecord.Lifecycle.Data.Buildpacks).To(BeEmpty())
						Expect(cfApp.Spec.Lifecycle.Data.Buildpacks).To(BeEmpty())
					})
				})

				When("stack is not specified", func() {
					BeforeEach(func() {
						appPatchMessage.Lifecycle.Data.Stack = ""
					})

					It("does not change the app lifecyle buildpacks", func() {
						Expect(patchedAppRecord.Lifecycle.Data.Stack).To(Equal(originalCFApp.Spec.Lifecycle.Data.Stack))
						Expect(cfApp.Spec.Lifecycle.Data.Stack).To(Equal(originalCFApp.Spec.Lifecycle.Data.Stack))
					})
				})
			})
		})

		When("the user is not authorized in the space", func() {
			It("returns a forbidden error", func() {
				Expect(patchErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.ForbiddenError{}))
			})
		})

		When("the user is a Space Manager (i.e. can view apps but not modify them)", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, orgUserRole.Name, cfOrg.Name)
				createRoleBinding(ctx, userName, spaceManagerRole.Name, cfSpace.Name)
			})

			It("returns a forbidden error", func() {
				Expect(patchErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.ForbiddenError{}))
			})
		})
	})

	Describe("PatchAppEnvVars", func() {
		const (
			key0 = "KEY0"
			key1 = "KEY1"
			key2 = "KEY2"
		)

		var (
			secretRecord repositories.AppEnvVarsRecord
			patchErr     error
		)

		BeforeEach(func() {
			envVars := map[string]string{
				key0: "VAL0",
				key1: "original-value",
			}
			secret := corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cfApp.Spec.EnvSecretName,
					Namespace: cfSpace.Name,
				},
				StringData: envVars,
			}
			Expect(k8sClient.Create(ctx, &secret)).To(Succeed())
		})

		JustBeforeEach(func() {
			var value1 *string
			value2 := "VAL2"

			newEnvVars := map[string]*string{
				key1: value1,
				key2: &value2,
			}
			patchEnvMsg := repositories.PatchAppEnvVarsMessage{
				AppGUID:              cfApp.Name,
				SpaceGUID:            cfSpace.Name,
				EnvironmentVariables: newEnvVars,
			}

			secretRecord, patchErr = appRepo.PatchAppEnvVars(ctx, authInfo, patchEnvMsg)
		})

		When("the user is authorized and an app exists with a secret", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
			})

			It("returns the updated secret record", func() {
				Expect(patchErr).NotTo(HaveOccurred())
				Expect(secretRecord.EnvironmentVariables).To(SatisfyAll(
					HaveLen(2),
					HaveKeyWithValue(key0, "VAL0"),
					HaveKeyWithValue(key2, "VAL2"),
				))
			})

			It("patches the underlying secret", func() {
				envSecret := corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      cfApp.Spec.EnvSecretName,
						Namespace: cfSpace.Name,
					},
				}
				err := k8sClient.Get(ctx, client.ObjectKeyFromObject(&envSecret), &envSecret)
				Expect(err).NotTo(HaveOccurred())

				Expect(asMapOfStrings(envSecret.Data)).To(SatisfyAll(
					HaveLen(2),
					HaveKeyWithValue(key0, "VAL0"),
					HaveKeyWithValue(key2, "VAL2"),
				))
			})
		})

		When("the user is not authorized", func() {
			It("return a forbidden error", func() {
				Expect(patchErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.ForbiddenError{}))
			})
		})
	})

	Describe("SetCurrentDroplet", func() {
		var (
			dropletGUID string
			appGUID     string

			currentDropletRecord repositories.CurrentDropletRecord
			setDropletErr        error
		)

		BeforeEach(func() {
			dropletGUID = uuid.NewString()
			appGUID = cfApp.Name

			Expect(k8sClient.Create(ctx, &korifiv1alpha1.CFBuild{
				ObjectMeta: metav1.ObjectMeta{
					Name:      dropletGUID,
					Namespace: cfSpace.Name,
				},
				Spec: korifiv1alpha1.CFBuildSpec{
					AppRef: corev1.LocalObjectReference{Name: appGUID},
					Lifecycle: korifiv1alpha1.Lifecycle{
						Type: "buildpack",
					},
				},
			})).To(Succeed())
		})

		JustBeforeEach(func() {
			currentDropletRecord, setDropletErr = appRepo.SetCurrentDroplet(ctx, authInfo, repositories.SetCurrentDropletMessage{
				AppGUID:     appGUID,
				DropletGUID: dropletGUID,
				SpaceGUID:   cfSpace.Name,
			})
		})

		When("user has the space developer role", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
			})

			It("awaits the ready condition", func() {
				Expect(appAwaiter.AwaitConditionCallCount()).To(Equal(1))
				obj, conditionType := appAwaiter.AwaitConditionArgsForCall(0)
				Expect(obj.GetName()).To(Equal(appGUID))
				Expect(obj.GetNamespace()).To(Equal(cfSpace.Name))
				Expect(conditionType).To(Equal(korifiv1alpha1.StatusConditionReady))
			})

			It("returns a CurrentDroplet record", func() {
				Expect(setDropletErr).NotTo(HaveOccurred())
				Expect(currentDropletRecord).To(Equal(repositories.CurrentDropletRecord{
					AppGUID:     cfApp.Name,
					DropletGUID: dropletGUID,
				}))
			})

			It("sets the spec.current_droplet_ref.name to the Droplet GUID", func() {
				lookupKey := client.ObjectKeyFromObject(cfApp)
				updatedApp := new(korifiv1alpha1.CFApp)
				Expect(k8sClient.Get(ctx, lookupKey, updatedApp)).To(Succeed())
				Expect(updatedApp.Spec.CurrentDropletRef.Name).To(Equal(dropletGUID))
			})

			When("the app never becomes ready", func() {
				BeforeEach(func() {
					appAwaiter.AwaitConditionReturns(&korifiv1alpha1.CFApp{}, errors.New("time-out-err"))
				})

				It("returns an error", func() {
					Expect(setDropletErr).To(MatchError(ContainSubstring("time-out-err")))
				})
			})

			When("the app doesn't exist", func() {
				BeforeEach(func() {
					appGUID = "no-such-app"
				})

				It("errors", func() {
					Expect(setDropletErr).To(MatchError(ContainSubstring("not found")))
				})
			})
		})

		When("the user is not authorized", func() {
			It("errors", func() {
				Expect(setDropletErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.ForbiddenError{}))
			})
		})
	})

	Describe("SetDesiredState", func() {
		const (
			appStartedValue = "STARTED"
			appStoppedValue = "STOPPED"
		)

		var (
			appGUID           string
			returnedAppRecord *repositories.AppRecord
			returnedErr       error
			initialAppState   string
			desiredAppState   string
		)

		BeforeEach(func() {
			initialAppState = appStartedValue
			desiredAppState = appStartedValue
		})

		JustBeforeEach(func() {
			appGUID = cfApp.Name
			Expect(k8s.PatchResource(ctx, k8sClient, cfApp, func() {
				cfApp.Spec.DesiredState = korifiv1alpha1.AppState(initialAppState)
			})).To(Succeed())
			appRecord, err := appRepo.SetAppDesiredState(ctx, authInfo, repositories.SetAppDesiredStateMessage{
				AppGUID:      appGUID,
				SpaceGUID:    cfSpace.Name,
				DesiredState: desiredAppState,
			})
			returnedAppRecord = &appRecord
			returnedErr = err
		})

		When("the user has permission to set the app state", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
			})

			When("starting an app", func() {
				BeforeEach(func() {
					initialAppState = appStoppedValue
				})

				It("doesn't return an error", func() {
					Expect(returnedErr).ToNot(HaveOccurred())
				})

				It("returns the updated app record", func() {
					Expect(returnedAppRecord.GUID).To(Equal(appGUID))
					Expect(returnedAppRecord.Name).To(Equal(cfApp.Spec.DisplayName))
					Expect(returnedAppRecord.SpaceGUID).To(Equal(cfSpace.Name))
				})

				It("waits for the started app state", func() {
					Expect(appAwaiter.AwaitStateCallCount()).To(Equal(1))
					Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())

					actualCFApp, actualStateCheck := appAwaiter.AwaitStateArgsForCall(0)
					Expect(actualCFApp.GetName()).To(Equal(cfApp.Name))
					Expect(actualCFApp.GetNamespace()).To(Equal(cfApp.Namespace))
					Expect(actualStateCheck(&korifiv1alpha1.CFApp{
						Spec: korifiv1alpha1.CFAppSpec{
							DesiredState: korifiv1alpha1.AppState(desiredAppState),
						},
						Status: korifiv1alpha1.CFAppStatus{
							Conditions: []metav1.Condition{{
								Type:               korifiv1alpha1.StatusConditionReady,
								Status:             metav1.ConditionTrue,
								ObservedGeneration: cfApp.Generation,
							}},
							ActualState: korifiv1alpha1.AppState(desiredAppState),
						},
					})).To(Succeed())
				})

				It("changes the desired state of the App", func() {
					Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())
					Expect(string(cfApp.Spec.DesiredState)).To(Equal(appStartedValue))
				})
			})

			When("stopping an app", func() {
				BeforeEach(func() {
					desiredAppState = appStoppedValue
				})

				It("doesn't return an error", func() {
					Expect(returnedErr).ToNot(HaveOccurred())
				})

				It("waits for the stopped app state", func() {
					Expect(appAwaiter.AwaitStateCallCount()).To(Equal(1))
					Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())

					actualCFApp, actualStateCheck := appAwaiter.AwaitStateArgsForCall(0)
					Expect(actualCFApp.GetName()).To(Equal(cfApp.Name))
					Expect(actualCFApp.GetNamespace()).To(Equal(cfApp.Namespace))
					Expect(actualStateCheck(&korifiv1alpha1.CFApp{
						Spec: korifiv1alpha1.CFAppSpec{
							DesiredState: korifiv1alpha1.AppState(desiredAppState),
						},
						Status: korifiv1alpha1.CFAppStatus{
							Conditions: []metav1.Condition{{
								Type:               korifiv1alpha1.StatusConditionReady,
								Status:             metav1.ConditionTrue,
								ObservedGeneration: cfApp.Generation,
							}},
							ActualState: korifiv1alpha1.AppState(desiredAppState),
						},
					})).To(Succeed())
				})

				It("changes the desired state of the App", func() {
					Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())
					Expect(string(cfApp.Spec.DesiredState)).To(Equal(appStoppedValue))
				})
			})

			When("the app doesn't exist", func() {
				It("returns an error", func() {
					_, err := appRepo.SetAppDesiredState(ctx, authInfo, repositories.SetAppDesiredStateMessage{
						AppGUID:      "fake-app-guid",
						SpaceGUID:    cfSpace.Name,
						DesiredState: appStartedValue,
					})

					Expect(err).To(MatchError(ContainSubstring("\"fake-app-guid\" not found")))
				})
			})
		})

		When("not allowed to set the application state", func() {
			It("returns a forbidden error", func() {
				Expect(returnedErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.ForbiddenError{}))
			})
		})
	})

	Describe("DeleteApp", func() {
		var (
			appGUID      string
			deleteAppErr error
		)

		BeforeEach(func() {
			appGUID = cfApp.Name
			createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
		})

		JustBeforeEach(func() {
			deleteAppErr = appRepo.DeleteApp(ctx, authInfo, repositories.DeleteAppMessage{
				AppGUID:   appGUID,
				SpaceGUID: cfSpace.Name,
			})
		})

		It("deletes the CFApp resource", func() {
			Expect(deleteAppErr).NotTo(HaveOccurred())
			app, err := appRepo.GetApp(ctx, authInfo, appGUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(app.DeletedAt).To(PointTo(BeTemporally("~", time.Now(), 5*time.Second)))
		})

		When("the app doesn't exist", func() {
			BeforeEach(func() {
				appGUID = "no-such-app"
			})

			It("errors", func() {
				Expect(deleteAppErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.NotFoundError{}))
			})
		})
	})

	Describe("GetAppEnv", func() {
		var (
			envVars      map[string]string
			secretName   string
			appGUID      string
			appEnvRecord repositories.AppEnvRecord
			getAppEnvErr error
		)

		BeforeEach(func() {
			appGUID = cfApp.Name
			secretName = "the-env-secret"

			envVars = map[string]string{
				"RAILS_ENV": "production",
				"LUNCHTIME": "12:00",
			}

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: cfSpace.Name,
				},
				StringData: envVars,
			}

			Expect(k8sClient.Create(ctx, secret)).To(Succeed())
		})

		JustBeforeEach(func() {
			appEnvRecord, getAppEnvErr = appRepo.GetAppEnv(ctx, authInfo, appGUID)
		})

		When("the user can read secrets in the space", func() {
			BeforeEach(func() {
				Expect(k8s.PatchResource(ctx, k8sClient, cfApp, func() {
					cfApp.Spec.EnvSecretName = secretName
				})).To(Succeed())

				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
			})

			It("returns the env vars stored on the secret", func() {
				Expect(getAppEnvErr).NotTo(HaveOccurred())
				Expect(appEnvRecord.AppGUID).To(Equal(cfApp.Name))
				Expect(appEnvRecord.SpaceGUID).To(Equal(cfApp.Namespace))
				Expect(appEnvRecord.EnvironmentVariables).To(Equal(envVars))
				Expect(appEnvRecord.SystemEnv).To(BeEmpty())
				Expect(appEnvRecord.AppEnv).To(BeEmpty())
			})

			When("the app has a service-binding secret", func() {
				var (
					vcapServiceSecretDataByte map[string][]byte
					vcapServiceSecretData     map[string]string
					vcapServiceDataPresenter  *env.VCAPServices
					err                       error
				)

				BeforeEach(func() {
					vcapServicesSecretName := prefixedGUID("vcap-secret")
					vcapServiceSecretDataByte, err = generateVcapServiceSecretDataByte()
					Expect(err).NotTo(HaveOccurred())
					vcapServiceSecretData = asMapOfStrings(vcapServiceSecretDataByte)
					vcapServiceDataPresenter = new(env.VCAPServices)
					err = json.Unmarshal(vcapServiceSecretDataByte["VCAP_SERVICES"], vcapServiceDataPresenter)
					Expect(err).NotTo(HaveOccurred())

					vcapSecret := &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      vcapServicesSecretName,
							Namespace: cfSpace.Name,
						},
						StringData: vcapServiceSecretData,
					}
					Expect(k8sClient.Create(ctx, vcapSecret)).To(Succeed())

					ogCFApp := cfApp.DeepCopy()
					cfApp.Status.VCAPServicesSecretName = vcapServicesSecretName
					Expect(k8sClient.Status().Patch(ctx, cfApp, client.MergeFrom(ogCFApp))).To(Succeed())
				})

				It("returns the env vars stored on the secret", func() {
					Expect(getAppEnvErr).NotTo(HaveOccurred())
					Expect(appEnvRecord.EnvironmentVariables).To(Equal(envVars))

					Expect(appEnvRecord.SystemEnv).NotTo(BeEmpty())
					Expect(appEnvRecord.SystemEnv["VCAP_SERVICES"]).To(Equal(vcapServiceDataPresenter))
				})
			})

			When("the app has a service-binding secret with empty VCAP_SERVICES data", func() {
				BeforeEach(func() {
					vcapServicesSecretName := prefixedGUID("vcap-secret")
					vcapSecret := &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      vcapServicesSecretName,
							Namespace: cfSpace.Name,
						},
						StringData: map[string]string{
							"VCAP_SERVICES": "{}",
						},
					}
					Expect(k8sClient.Create(ctx, vcapSecret)).To(Succeed())

					ogCFApp := cfApp.DeepCopy()
					cfApp.Status.VCAPServicesSecretName = vcapServicesSecretName
					Expect(k8sClient.Status().Patch(ctx, cfApp, client.MergeFrom(ogCFApp))).To(Succeed())
				})

				It("return an empty record for system env variables", func() {
					Expect(getAppEnvErr).NotTo(HaveOccurred())
					Expect(appEnvRecord.SystemEnv).To(BeEmpty())
				})
			})

			When("the app has a service-binding secret with missing VCAP_SERVICES data", func() {
				BeforeEach(func() {
					vcapServicesSecretName := prefixedGUID("vcap-secret")
					vcapSecret := &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      vcapServicesSecretName,
							Namespace: cfSpace.Name,
						},
					}
					Expect(k8sClient.Create(ctx, vcapSecret)).To(Succeed())

					ogCFApp := cfApp.DeepCopy()
					cfApp.Status.VCAPServicesSecretName = vcapServicesSecretName
					Expect(k8sClient.Status().Patch(ctx, cfApp, client.MergeFrom(ogCFApp))).To(Succeed())
				})

				It("return an empty record for system env variables", func() {
					Expect(getAppEnvErr).NotTo(HaveOccurred())
					Expect(appEnvRecord.SystemEnv).To(BeEmpty())
				})
			})

			When("the app has a vcap application secret", func() {
				BeforeEach(func() {
					vcapAppSecretName := prefixedGUID("vcap-app-secret")
					Expect(k8sClient.Create(ctx, &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      vcapAppSecretName,
							Namespace: cfApp.Namespace,
						},
						Data: map[string][]byte{
							"VCAP_APPLICATION": []byte(`{"foo":"bar"}`),
						},
					})).To(Succeed())
					Expect(k8s.Patch(ctx, k8sClient, cfApp, func() {
						cfApp.Status.VCAPApplicationSecretName = vcapAppSecretName
					})).To(Succeed())
				})

				It("returns the env vars stored on the secret", func() {
					Expect(getAppEnvErr).NotTo(HaveOccurred())
					Expect(appEnvRecord.EnvironmentVariables).To(Equal(envVars))

					Expect(appEnvRecord.AppEnv).To(HaveKey("VCAP_APPLICATION"))
				})
			})

			When("the EnvSecret doesn't exist", func() {
				BeforeEach(func() {
					secretName = "doIReallyExist"
					Expect(k8s.PatchResource(ctx, k8sClient, cfApp, func() {
						cfApp.Spec.EnvSecretName = secretName
					})).To(Succeed())
				})

				It("errors", func() {
					Expect(getAppEnvErr).To(MatchError(ContainSubstring("Secret")))
				})
			})

			When("the VCAPService secret doesn't exist", func() {
				BeforeEach(func() {
					vcapServicesSecretName := "doIReallyExist"

					ogCFApp := cfApp.DeepCopy()
					cfApp.Status.VCAPServicesSecretName = vcapServicesSecretName
					Expect(k8sClient.Status().Patch(ctx, cfApp, client.MergeFrom(ogCFApp))).To(Succeed())
				})

				It("errors", func() {
					Expect(getAppEnvErr).To(MatchError(ContainSubstring("Secret")))
				})
			})

			When("the VCAPApplication secret doesn't exist", func() {
				BeforeEach(func() {
					vcapApplicationSecretName := "doIReallyExist"

					ogCFApp := cfApp.DeepCopy()
					cfApp.Status.VCAPApplicationSecretName = vcapApplicationSecretName
					Expect(k8sClient.Status().Patch(ctx, cfApp, client.MergeFrom(ogCFApp))).To(Succeed())
				})

				It("errors", func() {
					Expect(getAppEnvErr).To(MatchError(ContainSubstring("Secret")))
				})
			})
		})

		When("EnvSecretName is blank", func() {
			BeforeEach(func() {
				secretName = ""
				Expect(k8s.PatchResource(ctx, k8sClient, cfApp, func() {
					cfApp.Spec.EnvSecretName = secretName
				})).To(Succeed())
			})

			It("returns an empty map", func() {
				Expect(appEnvRecord.EnvironmentVariables).To(BeEmpty())
			})
		})

		When("the user doesn't have permission to get secrets in the space", func() {
			It("errors", func() {
				Expect(getAppEnvErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.ForbiddenError{}))
			})
		})

		When("the app does not exist", func() {
			BeforeEach(func() {
				appGUID = "i don't exist"
			})
			It("returns an error", func() {
				Expect(getAppEnvErr).To(HaveOccurred())
				Expect(getAppEnvErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.NotFoundError{}))
			})
		})
	})

	Describe("GetDeletedAt", func() {
		var (
			deletionTime *time.Time
			getErr       error
		)

		BeforeEach(func() {
			createRoleBinding(ctx, userName, orgUserRole.Name, cfOrg.Name)
			createRoleBinding(ctx, userName, spaceDeveloperRole.Name, cfSpace.Name)
		})

		JustBeforeEach(func() {
			deletionTime, getErr = appRepo.GetDeletedAt(ctx, authInfo, cfApp.Name)
		})

		It("returns nil", func() {
			Expect(getErr).NotTo(HaveOccurred())
			Expect(deletionTime).To(BeNil())
		})

		When("the app is being deleted", func() {
			BeforeEach(func() {
				Expect(k8s.PatchResource(ctx, k8sClient, cfApp, func() {
					cfApp.Finalizers = append(cfApp.Finalizers, "foo")
				})).To(Succeed())

				Expect(k8sClient.Delete(ctx, cfApp)).To(Succeed())
			})

			It("returns the deletion time", func() {
				Expect(getErr).NotTo(HaveOccurred())
				Expect(deletionTime).To(PointTo(BeTemporally("~", time.Now(), time.Minute)))
			})
		})

		When("the app isn't found", func() {
			BeforeEach(func() {
				Expect(k8sClient.Delete(ctx, cfApp)).To(Succeed())
			})

			It("errors", func() {
				Expect(getErr).To(matchers.WrapErrorAssignableToTypeOf(apierrors.NotFoundError{}))
			})
		})
	})
})

func generateVcapServiceSecretDataByte() (map[string][]byte, error) {
	serviceDetails := env.ServiceDetails{
		Label:        "user-provided",
		Name:         "myupsi",
		Tags:         nil,
		InstanceGUID: "9779c01b-4b03-4a72-93c2-aae2ad4c75b2",
		InstanceName: "myupsi",
		BindingGUID:  "73f68d28-4602-47a3-8110-74ca991d5032",
		Credentials: map[string]any{
			"foo": "bar",
		},
		SyslogDrainURL: nil,
		VolumeMounts:   nil,
	}

	vcapServicesData, err := json.Marshal(env.VCAPServices{
		"user-provided": []env.ServiceDetails{
			serviceDetails,
		},
	})
	if err != nil {
		return nil, err
	}

	secretData := map[string][]byte{}
	secretData["VCAP_SERVICES"] = vcapServicesData

	return secretData, nil
}

func asMapOfStrings(data map[string][]byte) map[string]string {
	result := map[string]string{}

	for k, v := range data {
		result[k] = string(v)
	}

	return result
}
