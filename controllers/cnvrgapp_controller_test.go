package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("CnvrgApp controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		CnvrgAppNamespace = "cnvrg"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("Test PG", func() {
		It("No tenancy is set", func() {
			ctx := context.Background()
			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta: metav1.TypeMeta{
					Kind:       "CnvrgApp",
					APIVersion: "mlops.cnvrg.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cnvrgapp",
					Namespace: "test-no-tenancy",
				},
				Spec:   mlopsv1.DefaultSpec,
				Status: mlopsv1.CnvrgAppStatus{},
			}
			deployment := v1.Deployment{}
			Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "postgres", Namespace: cnvrgApp.Spec.CnvrgNs}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			fmt.Println("------------------------------------")
			fmt.Println(deployment.Spec.Template.Spec.NodeSelector)
			fmt.Println("------------------------------------")
			var shouldBe map[string]string
			Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
			Expect(k8sClient.Delete(ctx, cnvrgApp)).Should(Succeed())
		})
	})
	Context("Test PG2", func() {
		It("No tenancy is set", func() {
			cnvrgapp := "test-with-tenancy"
			ns := "test-with-tenancy"
			ctx := context.Background()
			testSpec := mlopsv1.DefaultSpec
			testSpec.ControlPlan.Conf.Tenancy.Enabled = "true"
			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta: metav1.TypeMeta{
					Kind:       "CnvrgApp",
					APIVersion: "mlops.cnvrg.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      cnvrgapp,
					Namespace: ns,
				},
				Spec:   testSpec,
				Status: mlopsv1.CnvrgAppStatus{},
			}
			deployment := v1.Deployment{}
			Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "postgres", Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			fmt.Println("------------------------------------")
			fmt.Println(deployment.Spec.Template.Spec.NodeSelector)
			fmt.Println("------------------------------------")
			var shouldBe map[string]string
			Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
			Expect(k8sClient.Delete(ctx, cnvrgApp)).Should(Succeed())
		})
	})
})
