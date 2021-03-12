package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/desired"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
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
		It("No tenancy is set", func() {
			ns := "test-pg-no-tenancy-set"
			ctx := context.Background()
			createNs(ns)
			testSpec := mlopsv1.DefaultSpec
			testSpec.CnvrgNs = ns
			testSpec.Networking.Enabled = "false"
			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta: metav1.TypeMeta{
					Kind:       "CnvrgApp",
					APIVersion: "mlops.cnvrg.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cnvrgapp",
					Namespace: "cnvrg",
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
			var shouldBe map[string]string
			Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))

			//Eventually(func() map[string]string {
			//	return deployment.Spec.Template.Spec.NodeSelector
			//}, timeout, interval).Should(Equal(map[string]string{cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Conf.Tenancy.Value}))

		})
		It("Tenancy is set", func() {
			ctx := context.Background()
			ns := "tenancy-set-pg"
			testSpec := mlopsv1.DefaultSpec
			testSpec.ControlPlan.Tenancy.Enabled = "true"
			testSpec.CnvrgNs = ns
			testSpec.Networking.Enabled = "false"
			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
				ObjectMeta: metav1.ObjectMeta{Name: ns},
				Spec:       testSpec,
				Status:     mlopsv1.CnvrgAppStatus{},
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
			shouldBe := map[string]string{cnvrgApp.Spec.ControlPlan.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Tenancy.Value}
			Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
		})
		It("Tenancy is set, HostPath storage is set", func() {
			ctx := context.Background()
			ns := "tenancy-set-pg-hostpath"
			testSpec := mlopsv1.DefaultSpec
			testSpec.ControlPlan.Tenancy.Enabled = "true"
			testSpec.Storage.Enabled = "true"
			testSpec.Storage.Hostpath.Enabled = "true"
			testSpec.Storage.Hostpath.NodeName = "node-1"
			testSpec.CnvrgNs = ns
			testSpec.Networking.Enabled = "false"
			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
				ObjectMeta: metav1.ObjectMeta{Name: ns},
				Spec:       testSpec,
				Status:     mlopsv1.CnvrgAppStatus{},
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
			shouldBe := map[string]string{
				cnvrgApp.Spec.ControlPlan.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Tenancy.Value,
				"kubernetes.io/hostname":              testSpec.Storage.Hostpath.NodeName,
			}
			Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
		})
		It("Tenancy true, DedicatedNode true", func() {
			ctx := context.Background()
			ns := "tenancy-true-dedicated-true"
			testSpec := mlopsv1.DefaultSpec
			testSpec.CnvrgNs = ns
			testSpec.ControlPlan.Tenancy.Enabled = "true"
			testSpec.ControlPlan.Tenancy.DedicatedNodes = "true"
			testSpec.Networking.Enabled = "false"
			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
				ObjectMeta: metav1.ObjectMeta{Name: ns},
				Spec:       testSpec,
				Status:     mlopsv1.CnvrgAppStatus{},
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
			shouldBe := map[string]string{cnvrgApp.Spec.ControlPlan.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Tenancy.Value}
			Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))
			shouldBeToleration := []corev1.Toleration{
				{
					Key:      cnvrgApp.Spec.ControlPlan.Tenancy.Key,
					Operator: "Equal",
					Value:    cnvrgApp.Spec.ControlPlan.Tenancy.Value,
					Effect:   "NoSchedule",
				},
			}
			Expect(deployment.Spec.Template.Spec.Tolerations).Should(Equal(shouldBeToleration))
		})
	})

	Context("Test Istio", func() {
		It("Testing istio instance spec - extra params", func() {
			ctx := context.Background()
			ns := "istio-test"
			createNs(ns)
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

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					fmt.Println(err)
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			serviceAnnotations := istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["serviceAnnotations"]
			Expect(serviceAnnotations.(map[string]interface{})["bar"]).Should(Equal("bla-bla"))
			Expect(serviceAnnotations.(map[string]interface{})["foo"]).Should(Equal("bla"))

			externalIP := istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["service"].(map[string]interface{})["externalIPs"]
			Expect(externalIP.([]interface{})[0]).Should(Equal("1.1.1.1"))
			Expect(externalIP.([]interface{})[1]).Should(Equal("2.2.2.2"))

			ingressSvcExtraPorts := istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["service"].(map[string]interface{})["ports"]
			extraPorts := map[int64]bool{123: false, 2323: false, 2321: false}
			var expectedRes []bool
			for _, port := range ingressSvcExtraPorts.([]interface{}) {
				portNumber := port.(map[string]interface{})["port"].(int64)
				if _, ok := extraPorts[portNumber]; ok {
					expectedRes = append(expectedRes, true)
				}
			}
			Expect(expectedRes).Should(Equal([]bool{true, true, true}))
		})
		It("Testing istio instance - cnvrg tenancy enabled  ", func() {
			ctx := context.Background()
			ns := "istio-test-with-tenancy-enabled"
			createNs(ns)
			testSpec := mlopsv1.DefaultSpec
			testSpec.CnvrgNs = ns
			testSpec.ControlPlan.Tenancy.Enabled = "true"

			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
				ObjectMeta: metav1.ObjectMeta{Name: ns},
				Spec:       testSpec,
				Status:     mlopsv1.CnvrgAppStatus{},
			}
			Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVR])

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					fmt.Println(err)
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			nodeSelector := istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["nodeSelector"]
			Expect(nodeSelector.(map[string]interface{})["cnvrg-taint"]).Should(Equal("true"))
		})

		It("Testing istio Operator deployment - cnvrg tenancy enabled  ", func() {
			ctx := context.Background()
			ns := "istio-operator-deployment-with-tenancy-enabled"
			createNs(ns)
			testSpec := mlopsv1.DefaultSpec
			testSpec.CnvrgNs = ns
			testSpec.ControlPlan.Tenancy.Enabled = "true"
			cnvrgApp := &mlopsv1.CnvrgApp{
				TypeMeta:   metav1.TypeMeta{Kind: "CnvrgApp", APIVersion: "mlops.cnvrg.io/v1"},
				ObjectMeta: metav1.ObjectMeta{Name: ns},
				Spec:       testSpec,
				Status:     mlopsv1.CnvrgAppStatus{},
			}
			Expect(k8sClient.Create(ctx, cnvrgApp)).Should(Succeed())
			deployment := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "istio-operator", Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			shouldBe := map[string]string{cnvrgApp.Spec.ControlPlan.Tenancy.Key: cnvrgApp.Spec.ControlPlan.Tenancy.Value}
			Expect(deployment.Spec.Template.Spec.NodeSelector).Should(Equal(shouldBe))

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
