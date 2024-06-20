package repositories_test

import (
	"code.cloudfoundry.org/korifi/api/repositories"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PodRepository", func() {
	var (
		podRepo *repositories.PodRepo
		org     *korifiv1alpha1.CFOrg
		space   *korifiv1alpha1.CFSpace
		pod     *corev1.Pod
		appGUID string
		app     *repositories.AppRecord
	)
	BeforeEach(func() {

		podRepo = repositories.NewPodRepo(userClientFactory)
		org = createOrgWithCleanup(ctx, prefixedGUID("org"))
		space = createSpaceWithCleanup(ctx, org.Name, prefixedGUID("space"))

		appGUID = uuid.NewString()
		app = &repositories.AppRecord{
			GUID:      appGUID,
			Revision:  "1",
			SpaceGUID: space.Name,
		}
		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "podname-2",
				Namespace: space.Name,
				Labels: map[string]string{
					"korifi.cloudfoundry.org/app-guid":     appGUID,
					"korifi.cloudfoundry.org/version":      "1",
					"korifi.cloudfoundry.org/process-type": "web",
				}},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "web",
						Image: "nginx",
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, pod)).To(Succeed())

	})
	Describe("DeletePod", func() {
		When("the user is not a a SpaceDeveloper", func() {

			It("fails to delete the pod", func() {
				Expect(podRepo.DeletePod(ctx, authInfo, *app, "2", "web")).To(HaveOccurred())
			})
		})
		When("the user is a SpaceDeveloper", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, space.Name)
			})

			It("deletes the pod", func() {
				Expect(podRepo.DeletePod(ctx, authInfo, *app, "2", "web")).To(Succeed())
			})
		})
		When("the instance does not exist", func() {
			BeforeEach(func() {
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, space.Name)
			})

			It("fails to delete instance", func() {
				Expect(podRepo.DeletePod(ctx, authInfo, *app, "3", "web")).To(HaveOccurred())
			})
		})

	})

})
