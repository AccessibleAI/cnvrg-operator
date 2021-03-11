package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("CnvrgApp controller", func() {

	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("Test PG", func() {
		//It("No tenancy is set", func() {
		//	ns := "test-pg-no-tenancy-set"
		//	ctx := context.Background()
		//	createNs(ns, k8sClient)
		//	testSpec := mlopsv1.DefaultSpec
		//	testSpec.CnvrgNs = ns
		//	testSpec.Networking.Enabled = "false"
		//	cnvrgApp := &mlopsv1.CnvrgApp{
		//		TypeMeta: metav1.TypeMeta{
		//			Kind:       "CnvrgApp",
		//			APIVersion: "mlops.cnvrg.io/v1",
		//		},
		//		ObjectMeta: metav1.ObjectMeta{
		//			Name:      "cnvrgapp",
		//			Namespace: "cnvrg",
		//		},
		//		Spec:   testSpec,
		//		Status: mlopsv1.CnvrgAppStatus{},
		//	}
		//	deployment := v1.Deployment{}
		//	Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())
		//	Eventually(func() bool {
		//		err := k8sClient.Get(ctx, types.NamespacedName{Name: "postgres", Namespace: ns}, &deployment)
		//		if err != nil {
		//			return false
		//		}
		//		return true
		//	}, timeout, interval).Should(BeTrue())
		//	var shouldBe map[string]string
		//	Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
		//
		//	//Eventually(func() map[string]string {
		//	//	return deployment.Spec.Template.Spec.NodeSelector
		//	//}, timeout, interval).Should(Equal(map[string]string{cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Value}))
		//
		//})
		//It("Tenancy is set", func() {
		//	ctx := context.Background()
		//	ns := "tenancy-set-pg"
		//	testSpec := mlopsv1.DefaultSpec
		//	testSpec.ControlPlan.Conf.Tenancy.Enabled = "true"
		//	testSpec.CnvrgNs = ns
		//	testSpec.Networking.Enabled = "false"
		//	cnvrgApp := &mlopsv1.CnvrgApp{
		//		TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
		//		ObjectMeta: metav1.ObjectMeta{Name: ns},
		//		Spec:       testSpec,
		//		Status:     mlopsv1.CnvrgAppStatus{},
		//	}
		//	deployment := v1.Deployment{}
		//	Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())
		//	Eventually(func() bool {
		//		err := k8sClient.Get(ctx, types.NamespacedName{Name: "postgres", Namespace: ns}, &deployment)
		//		if err != nil {
		//			return false
		//		}
		//		return true
		//	}, timeout, interval).Should(BeTrue())
		//	shouldBe := map[string]string{cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Value}
		//	Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
		//})
		//It("Tenancy is set, HostPath storage is set", func() {
		//	ctx := context.Background()
		//	ns := "tenancy-set-pg-hostpath"
		//	testSpec := mlopsv1.DefaultSpec
		//	testSpec.ControlPlan.Conf.Tenancy.Enabled = "true"
		//	testSpec.Storage.Enabled = "true"
		//	testSpec.Storage.Hostpath.Enabled = "true"
		//	testSpec.Storage.Hostpath.NodeName = "node-1"
		//	testSpec.CnvrgNs = ns
		//	testSpec.Networking.Enabled = "false"
		//	cnvrgApp := &mlopsv1.CnvrgApp{
		//		TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
		//		ObjectMeta: metav1.ObjectMeta{Name: ns},
		//		Spec:       testSpec,
		//		Status:     mlopsv1.CnvrgAppStatus{},
		//	}
		//	deployment := v1.Deployment{}
		//	Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())
		//	Eventually(func() bool {
		//		err := k8sClient.Get(ctx, types.NamespacedName{Name: "postgres", Namespace: ns}, &deployment)
		//		if err != nil {
		//			return false
		//		}
		//		return true
		//	}, timeout, interval).Should(BeTrue())
		//	shouldBe := map[string]string{
		//		cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Value,
		//		"kubernetes.io/hostname":                   testSpec.Storage.Hostpath.NodeName,
		//	}
		//	Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
		//})
		//It("Tenancy true, DedicatedNode true", func() {
		//	ctx := context.Background()
		//	ns := "tenancy-true-dedicated-true"
		//	testSpec := mlopsv1.DefaultSpec
		//	testSpec.CnvrgNs = ns
		//	testSpec.ControlPlan.Conf.Tenancy.Enabled = "true"
		//	testSpec.ControlPlan.Conf.Tenancy.DedicatedNodes = "true"
		//	testSpec.Networking.Enabled = "false"
		//	cnvrgApp := &mlopsv1.CnvrgApp{
		//		TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
		//		ObjectMeta: metav1.ObjectMeta{Name: ns},
		//		Spec:       testSpec,
		//		Status:     mlopsv1.CnvrgAppStatus{},
		//	}
		//	deployment := v1.Deployment{}
		//	Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())
		//	Eventually(func() bool {
		//		err := k8sClient.Get(ctx, types.NamespacedName{Name: "postgres", Namespace: ns}, &deployment)
		//		if err != nil {
		//			return false
		//		}
		//		return true
		//	}, timeout, interval).Should(BeTrue())
		//	shouldBe := map[string]string{cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Value}
		//	Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
		//	shouldBeToleration := []corev1.Toleration{
		//		{
		//			Key:      cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Key,
		//			Operator: "Equal",
		//			Value:    cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Value,
		//			Effect:   "NoSchedule",
		//		},
		//	}
		//	Expect(deployment.Spec.Template.Spec.Tolerations).Should(Equal(shouldBeToleration))
		//})
	})

	Context("Test Istio", func() {
		It("Testing istio", func() {
			ctx := context.Background()
			ns := "istio-test"
			//createNs(ns)
			testSpec := mlopsv1.DefaultSpec
			testSpec.CnvrgNs = ns
			testSpec.Networking.Enabled = "true"
			testSpec.Networking.Istio.Enabled = "true"
			testSpec.Networking.Istio.ExternalIP = "1.1.1.1;2.2.2.2"
			testSpec.Networking.Istio.IngressSvcAnnotations = "foo: bla    ;bar : bla-bla"
			testSpec.Networking.Istio.IngressSvcExtraPorts = "123;2323;2321"
			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
				ObjectMeta: metav1.ObjectMeta{Name: ns},
				Spec:       testSpec,
				Status:     mlopsv1.CnvrgAppStatus{},
			}
			Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVR])

			ns1 := &corev1.NamespaceList{}
			err := k8sClient.List(ctx, ns1)
			if err != nil {
				fmt.Println("x")
			}
			fmt.Println("x")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					fmt.Println(err)
					return false
				}

				return true
			}, timeout, interval).Should(BeTrue())

		})
	})
})

func createNs(ns string) {
	ctx := context.Background()
	ns1 := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
	err := k8sClient.Create(ctx, ns1)
	if err != nil {
		panic(err)
	}

}
