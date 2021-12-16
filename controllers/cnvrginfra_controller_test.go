package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/scheduling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sort"
	"strings"
	"time"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("CnvrgInfra controller", func() {

	const (
		timeout  = time.Second * 30
		interval = time.Millisecond * 250
	)

	Context("Test Monitoring", func() {
		It("Prometheus Labels", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Monitoring.Prometheus.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			prom := &unstructured.Unstructured{}
			prom.SetGroupVersionKind(desired.Kinds[desired.PrometheusGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-infra-prometheus", Namespace: ns}, prom)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			a := prom.Object["spec"].(map[string]interface{})["podMetadata"].(map[string]interface{})["annotations"]
			l := prom.Object["spec"].(map[string]interface{})["podMetadata"].(map[string]interface{})["labels"]
			Expect(a).Should(BeNil())
			Expect(l).Should(HaveKeyWithValue("foo", "bar"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Prometheus Operator Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Monitoring.PrometheusOperator.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			deployment := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-prometheus-operator", Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Prometheus NodeExporter Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Monitoring.NodeExporter.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			ds := v1.DaemonSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "node-exporter", Namespace: ns}, &ds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(ds.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(ds.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Kube State Metrics Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Monitoring.KubeStateMetrics.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			deployment := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "kube-state-metrics", Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Grafana Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Monitoring.Grafana.Enabled = true
			infra.Spec.Monitoring.Prometheus.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			deployment := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "grafana", Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Dcgm exporter Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Monitoring.DcgmExporter.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			ds := v1.DaemonSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "dcgm-exporter", Namespace: ns}, &ds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(ds.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(ds.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Prom creds secret generator", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Monitoring.Prometheus.Enabled = true
			// create infra
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			// get prom creds
			promCreds := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Monitoring.Prometheus.CredsRef, Namespace: ns}, &promCreds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// enforce reconcile loop - enable Prometheus Operator and make sure it was deployed
			infraRes := mlopsv1.CnvrgInfra{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: ns}, &infraRes)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			rvBeforeUpdate := infraRes.ObjectMeta.ResourceVersion
			infraRes.Spec.Monitoring.PrometheusOperator.Enabled = true
			Expect(k8sClient.Update(ctx, &infraRes)).Should(Succeed())
			dep := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-prometheus-operator", Namespace: ns}, &dep)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Make sure resource version has been updated
			Expect(rvBeforeUpdate).Should(Not(Equal(infraRes.ObjectMeta.ResourceVersion)))
			// get redis creds after reconcile
			promCredsAfterReconcile := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Monitoring.Prometheus.CredsRef, Namespace: ns}, &promCredsAfterReconcile)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Make sure redis creds wasn't mutated between reconciliation loops
			Expect(promCreds.Data["CNVRG_PROMETHEUS_PASS"]).Should(Equal(promCredsAfterReconcile.Data["CNVRG_PROMETHEUS_PASS"]))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
	})

	Context("Test Storage", func() {
		It("Hostpath provisioner Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Storage.Hostpath.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			ds := v1.DaemonSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "hostpath-provisioner", Namespace: ns}, &ds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(ds.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(ds.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Nfs provisioner Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Storage.Nfs.Enabled = true
			infra.Spec.Storage.Nfs.Path = "/nfs-path"
			infra.Spec.Storage.Nfs.Server = "nfs-server"
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			d := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "nfs-client-provisioner", Namespace: ns}, &d)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(d.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(d.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})

	})

	Context("Test networking", func() {
		It("Istio Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Istio.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(istio.GetAnnotations()).Should(BeNil())
			Expect(istio.GetLabels()).Should(HaveKeyWithValue("foo", "bar"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio Service Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Istio.Enabled = true
			infra.Spec.Networking.Istio.IngressSvcAnnotations = map[string]string{"foo": "bar"}
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			serviceAnnotations := istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["serviceAnnotations"]
			Expect(serviceAnnotations).Should(HaveKeyWithValue("foo", "bar"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio LoadBalancer Source Ranges", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Istio.Enabled = true
			infra.Spec.Networking.Istio.LBSourceRanges = []string{"1.1.1.1/32", "2.2.2.2/24"}
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			lbSources := []string{
				istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["service"].(map[string]interface{})["loadBalancerSourceRanges"].([]interface{})[0].(string),
				istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["service"].(map[string]interface{})["loadBalancerSourceRanges"].([]interface{})[1].(string),
			}
			Expect(lbSources).Should(Equal([]string{"1.1.1.1/32", "2.2.2.2/24"}))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio ExternalIps", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			exIps := []string{"1.1.1.1", "2.2.2.2"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Istio.Enabled = true
			infra.Spec.Networking.Istio.ExternalIP = exIps
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			externalIps := []string{
				istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["service"].(map[string]interface{})["externalIPs"].([]interface{})[0].(string),
				istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["service"].(map[string]interface{})["externalIPs"].([]interface{})[1].(string),
			}
			Expect(externalIps).Should(Equal(exIps))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio ExtraPorts", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			exPorts := []int{1111, 2222}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Istio.Enabled = true
			infra.Spec.Networking.Istio.IngressSvcExtraPorts = exPorts
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			port1111 := istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["service"].(map[string]interface{})["ports"].([]interface{})[2]
			port2222 := istio.Object["spec"].(map[string]interface{})["components"].(map[string]interface{})["ingressGateways"].([]interface{})[0].(map[string]interface{})["k8s"].(map[string]interface{})["service"].(map[string]interface{})["ports"].([]interface{})[3]
			var port1111int64 int64 = 1111
			var port2222int64 int64 = 2222
			Expect(port1111).Should(HaveKeyWithValue("name", "port1111"))
			Expect(port1111).Should(HaveKeyWithValue("port", port1111int64))
			Expect(port2222).Should(HaveKeyWithValue("name", "port2222"))
			Expect(port2222).Should(HaveKeyWithValue("port", port2222int64))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio Default Gateway Name", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Istio.Enabled = true
			infra.Spec.Networking.Ingress.IstioGwEnabled = true
			infra.Spec.InfraNamespace = ns
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			gw := &unstructured.Unstructured{}
			gw.SetGroupVersionKind(desired.Kinds[desired.IstioGwGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: fmt.Sprintf(mlopsv1.IstioGwName, ns), Namespace: ns}, gw)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio Custom Gateway Name", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Istio.Enabled = true
			infra.Spec.Networking.Ingress.IstioGwEnabled = true
			infra.Spec.Networking.Ingress.IstioGwName = "foo-bar"
			infra.Spec.InfraNamespace = ns
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			gw := &unstructured.Unstructured{}
			gw.SetGroupVersionKind(desired.Kinds[desired.IstioGwGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "foo-bar", Namespace: ns}, gw)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio Disabled ", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)

			infra.Spec.InfraNamespace = ns
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVK])
			time.Sleep(time.Second * 3)
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-istio", Namespace: ns}, istio)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeFalse())
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio Disabled - Custom GW name for Prometheus and Grafana VS", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Monitoring.Grafana.Enabled = true
			infra.Spec.Monitoring.Prometheus.Enabled = true
			infra.Spec.Networking.Ingress.IstioGwName = "foo-bar"
			infra.Spec.InfraNamespace = ns
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())

			vs := &unstructured.Unstructured{}
			vs.SetGroupVersionKind(desired.Kinds[desired.IstioVsGVK])

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Monitoring.Grafana.SvcName, Namespace: ns}, vs)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(vs.Object["spec"].(map[string]interface{})["gateways"]).Should(ContainElement("foo-bar"))

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Monitoring.Prometheus.SvcName, Namespace: ns}, vs)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(vs.Object["spec"].(map[string]interface{})["gateways"]).Should(ContainElement("foo-bar"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())

		})
		It("Istio Disabled - No Istio GW are created", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.InfraNamespace = ns
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			gw := &unstructured.Unstructured{}
			gw.SetGroupVersionKind(desired.Kinds[desired.IstioGwGVK])
			time.Sleep(time.Second * 3)
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: fmt.Sprintf(mlopsv1.IstioGwName, ns), Namespace: ns}, gw)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeFalse())
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Istio Operator Deployment Tenancy", func() {
			ns := createNs()
			ctx := context.Background()

			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Tenancy.Enabled = true
			infra.Spec.Networking.Istio.Enabled = true
			infra.Spec.InfraNamespace = ns

			deployment := v1.Deployment{}
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			istio := &unstructured.Unstructured{}
			istio.SetGroupVersionKind(desired.Kinds[desired.IstioGVK])
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "istio-operator", Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			t := corev1.Toleration{
				Key:      infra.Spec.Tenancy.Key,
				Operator: "Equal",
				Value:    infra.Spec.Tenancy.Value,
				Effect:   "NoSchedule",
			}

			Expect(deployment.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
			Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})

	})

	Context("Test Config Reloader", func() {
		It("Config Reloader Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.ConfigReloader.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			d := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "config-reloader", Namespace: ns}, &d)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(d.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(d.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})

	})

	Context("Test Logging", func() {
		It("Fluentbit Labels/Annotations", func() {
			ns := createNs()
			ctx := context.Background()
			labels := map[string]string{"foo": "bar"}
			annotations := map[string]string{"foo1": "bar1"}
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Logging.Fluentbit.Enabled = true
			infra.Spec.InfraNamespace = ns
			infra.Spec.Labels = labels
			infra.Spec.Annotations = annotations
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			ds := v1.DaemonSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-fluentbit", Namespace: ns}, &ds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(ds.Labels).Should(HaveKeyWithValue("foo", "bar"))
			Expect(ds.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Fluentbit Default Logs Volumes", func() {
			ns := createNs()
			ctx := context.Background()

			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Logging.Fluentbit.Enabled = true
			infra.Spec.InfraNamespace = ns
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			ds := v1.DaemonSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-fluentbit", Namespace: ns}, &ds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			hostPathUnset := corev1.HostPathUnset
			v := corev1.Volume{
				Name: "varlog",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/var/log",
						Type: &hostPathUnset,
					},
				},
			}
			v1 := corev1.Volume{
				Name: "varlibdockercontainers",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/var/lib/docker/containers",
						Type: &hostPathUnset,
					},
				},
			}
			Expect(ds.Spec.Template.Spec.Volumes).Should(ContainElement(v))
			Expect(ds.Spec.Template.Spec.Volumes).Should(ContainElement(v1))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())

		})

		It("Fluentbit Extra Logs Volumes", func() {
			ns := createNs()
			ctx := context.Background()

			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Logging.Fluentbit.Enabled = true
			infra.Spec.Logging.Fluentbit.LogsMounts = map[string]string{"foobar": "/foo/bar"}
			infra.Spec.InfraNamespace = ns
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			ds := v1.DaemonSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-fluentbit", Namespace: ns}, &ds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			hostPathUnset := corev1.HostPathUnset
			v := corev1.Volume{
				Name: "foobar",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/foo/bar",
						Type: &hostPathUnset,
					},
				},
			}
			Expect(ds.Spec.Template.Spec.Volumes).Should(ContainElement(v))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})

		It("Fluentbit Tolerations", func() {
			ns := createNs()
			ctx := context.Background()

			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Logging.Fluentbit.Enabled = true
			infra.Spec.Logging.Fluentbit.NodeSelector = map[string]string{"foo": "bar"}
			infra.Spec.Tenancy.Enabled = true
			infra.Spec.InfraNamespace = ns
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			ds := v1.DaemonSet{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-fluentbit", Namespace: ns}, &ds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(ds.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))
			Expect(ds.Spec.Template.Spec.NodeSelector).ShouldNot(HaveKeyWithValue("purpose", "cnvrg-control-plane"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})

	})

	Context("Test Proxy", func() {
		It("Proxy configmap test creation - default no_proxy", func() {

			ctx := context.Background()
			ns := createNs()
			expectedNoProxy := networking.DefaultNoProxy("cluster.local")
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Proxy.Enabled = true
			infra.Spec.Networking.Proxy.HttpProxy = []string{
				"http://proxy1.org.local",
				"http://proxy2.org.local",
			}
			infra.Spec.Networking.Proxy.HttpsProxy = []string{
				"https://proxy1.org.local",
				"https://proxy2.org.local",
			}

			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())

			Eventually(func() bool {
				infraRes := mlopsv1.CnvrgInfra{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: ns}, &infraRes)
				if err != nil {
					return false
				}
				sort.Strings(infraRes.Spec.Networking.Proxy.NoProxy)
				sort.Strings(expectedNoProxy)
				return reflect.DeepEqual(infraRes.Spec.Networking.Proxy.NoProxy, expectedNoProxy)
			}, timeout, interval).Should(BeTrue())

			cm := corev1.ConfigMap{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Networking.Proxy.ConfigRef, Namespace: ns}, &cm)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(cm.Data).Should(HaveKeyWithValue("HTTP_PROXY", "http://proxy1.org.local,http://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("http_proxy", "http://proxy1.org.local,http://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("https_proxy", "https://proxy1.org.local,https://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("HTTPS_PROXY", "https://proxy1.org.local,https://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("NO_PROXY", strings.Join(expectedNoProxy, ",")))
			Expect(cm.Data).Should(HaveKeyWithValue("no_proxy", strings.Join(expectedNoProxy, ",")))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Proxy configmap test creation - custom no_proxy", func() {

			ctx := context.Background()
			ns := createNs()
			noProxy := []string{".foo.bar"}
			expectedNoProxy := append(noProxy, networking.DefaultNoProxy("cluster.local")...)
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Proxy.Enabled = true
			infra.Spec.Networking.Proxy.HttpProxy = []string{
				"http://proxy1.org.local",
				"http://proxy2.org.local",
			}
			infra.Spec.Networking.Proxy.HttpsProxy = []string{
				"https://proxy1.org.local",
				"https://proxy2.org.local",
			}
			infra.Spec.Networking.Proxy.NoProxy = noProxy

			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			Eventually(func() bool {
				infraRes := mlopsv1.CnvrgInfra{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: ns}, &infraRes)
				if err != nil {
					return false
				}
				sort.Strings(infraRes.Spec.Networking.Proxy.NoProxy)
				sort.Strings(expectedNoProxy)
				return reflect.DeepEqual(infraRes.Spec.Networking.Proxy.NoProxy, expectedNoProxy)
			}, timeout, interval).Should(BeTrue())
			cm := corev1.ConfigMap{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Networking.Proxy.ConfigRef, Namespace: ns}, &cm)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(cm.Data).Should(HaveKeyWithValue("HTTP_PROXY", "http://proxy1.org.local,http://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("http_proxy", "http://proxy1.org.local,http://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("https_proxy", "https://proxy1.org.local,https://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("HTTPS_PROXY", "https://proxy1.org.local,https://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("NO_PROXY", strings.Join(expectedNoProxy, ",")))
			Expect(cm.Data).Should(HaveKeyWithValue("no_proxy", strings.Join(expectedNoProxy, ",")))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("Proxy configmap test creation - k8s api server", func() {
			noProxy := []string{".foo.bar"}
			expectedNoProxy := append(noProxy, networking.DefaultNoProxy("cluster.local")...)
			ctx := context.Background()
			ns := createNs()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Networking.Proxy.Enabled = true
			infra.Spec.Networking.Proxy.HttpProxy = []string{
				"http://proxy1.org.local",
				"http://proxy2.org.local",
			}
			infra.Spec.Networking.Proxy.HttpsProxy = []string{
				"https://proxy1.org.local",
				"https://proxy2.org.local",
			}

			infra.Spec.Networking.Proxy.NoProxy = noProxy

			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			Eventually(func() bool {
				infraRes := mlopsv1.CnvrgInfra{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: ns}, &infraRes)
				if err != nil {
					return false
				}
				sort.Strings(infraRes.Spec.Networking.Proxy.NoProxy)
				sort.Strings(expectedNoProxy)
				return reflect.DeepEqual(infraRes.Spec.Networking.Proxy.NoProxy, expectedNoProxy)
			}, timeout, interval).Should(BeTrue())

			cm := corev1.ConfigMap{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Networking.Proxy.ConfigRef, Namespace: ns}, &cm)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(cm.Data).Should(HaveKeyWithValue("HTTP_PROXY", "http://proxy1.org.local,http://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("http_proxy", "http://proxy1.org.local,http://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("https_proxy", "https://proxy1.org.local,https://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("HTTPS_PROXY", "https://proxy1.org.local,https://proxy2.org.local"))
			Expect(cm.Data).Should(HaveKeyWithValue("NO_PROXY", strings.Join(expectedNoProxy, ",")))
			Expect(cm.Data).Should(HaveKeyWithValue("no_proxy", strings.Join(expectedNoProxy, ",")))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())

		})
		It("Proxy configmap test creation - proxy disabled", func() {

			ctx := context.Background()
			ns := createNs()
			infra := getEmptyTestInfraSpec(ns)

			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			cm := corev1.ConfigMap{}

			Eventually(func() bool {
				infraRes := mlopsv1.CnvrgInfra{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: ns}, &infraRes)
				if err != nil {
					return false
				}
				return infraRes.Spec.Networking.Proxy.Enabled == false
			}, timeout, interval).Should(BeTrue())

			time.Sleep(3)
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Networking.Proxy.ConfigRef, Namespace: ns}, &cm)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeFalse())
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
	})

	Context("Test Redis", func() {
		It("Redis creds secret generator", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Dbs.Redis.Enabled = true
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			// get redis creds
			redisCreds := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Dbs.Redis.CredsRef, Namespace: ns}, &redisCreds)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// enforce reconcile loop - enable Prometheus Operator and make sure it was deployed
			infraRes := mlopsv1.CnvrgInfra{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: ns}, &infraRes)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			rvBeforeUpdate := infraRes.ObjectMeta.ResourceVersion
			infraRes.Spec.Monitoring.PrometheusOperator.Enabled = true
			Expect(k8sClient.Update(ctx, &infraRes)).Should(Succeed())
			dep := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cnvrg-prometheus-operator", Namespace: ns}, &dep)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Make sure resource version has been updated
			Expect(rvBeforeUpdate).Should(Not(Equal(infraRes.ObjectMeta.ResourceVersion)))
			// get redis creds after reconcile
			redisCredsAfterReconcile := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Dbs.Redis.CredsRef, Namespace: ns}, &redisCredsAfterReconcile)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Make sure redis creds wasn't mutated between reconciliation loops
			Expect(redisCreds.Data["redis.conf"]).Should(Equal(redisCredsAfterReconcile.Data["redis.conf"]))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
	})

	Context("Test Capsule", func() {

		It("Capsule deployment", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Capsule.Enabled = true
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			dep := v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Capsule.SvcName, Namespace: ns}, &dep)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(dep.Spec.Template.Spec.ServiceAccountName).To(Equal("cnvrg-capsule"))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})

		It("Capsule PVC - default storage size", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Capsule.Enabled = true
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())

			pvc := corev1.PersistentVolumeClaim{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Capsule.SvcName, Namespace: ns}, &pvc)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			storageSize := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
			Expect(storageSize.String()).To(Equal(infra.Spec.Capsule.StorageSize))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})

		It("Capsule PVC - custom storage size", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Capsule.Enabled = true
			infra.Spec.Capsule.StorageSize = "20Gi"
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())

			pvc := corev1.PersistentVolumeClaim{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Capsule.SvcName, Namespace: ns}, &pvc)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			storageSize := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
			Expect(storageSize.String()).To(Equal(infra.Spec.Capsule.StorageSize))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})

		It("Capsule PVC - custom storage class", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Capsule.Enabled = true
			infra.Spec.Capsule.StorageClass = "foo-bar"
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			pvc := corev1.PersistentVolumeClaim{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.Capsule.SvcName, Namespace: ns}, &pvc)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(*pvc.Spec.StorageClassName).To(Equal(infra.Spec.Capsule.StorageClass))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
	})

	Context("Test Priority class", func() {
		It("cnvrg app priority class  ", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Capsule.Enabled = true
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			pr := v12.PriorityClass{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.CnvrgAppPriorityClass.Name, Namespace: ns}, &pr)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			expectedPr := v12.PriorityClass{
				ObjectMeta:  metav1.ObjectMeta{Name: infra.Spec.CnvrgAppPriorityClass.Name},
				Value:       infra.Spec.CnvrgAppPriorityClass.Value,
				Description: infra.Spec.CnvrgAppPriorityClass.Description,
			}
			Expect(expectedPr.Value).To(Equal(pr.Value))
			Expect(expectedPr.Description).To(Equal(pr.Description))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
		It("cnvrg job priority class  ", func() {
			ns := createNs()
			ctx := context.Background()
			infra := getDefaultTestInfraSpec(ns)
			infra.Spec.Capsule.Enabled = true
			Expect(k8sClient.Create(ctx, infra)).Should(Succeed())
			pr := v12.PriorityClass{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: infra.Spec.CnvrgJobPriorityClass.Name, Namespace: ns}, &pr)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			expectedPr := v12.PriorityClass{
				ObjectMeta:  metav1.ObjectMeta{Name: infra.Spec.CnvrgJobPriorityClass.Name},
				Value:       infra.Spec.CnvrgJobPriorityClass.Value,
				Description: infra.Spec.CnvrgJobPriorityClass.Description,
			}
			Expect(expectedPr.Value).To(Equal(pr.Value))
			Expect(expectedPr.Description).To(Equal(pr.Description))
			Expect(k8sClient.Delete(ctx, infra)).Should(Succeed())
		})
	})

})

func getEmptyTestInfraSpec(ns string) *mlopsv1.CnvrgInfra {

	infraEmptySpec := mlopsv1.CnvrgInfraSpec{}
	infraEmptySpec.InfraNamespace = ns
	return &mlopsv1.CnvrgInfra{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CnvrgInfra",
			APIVersion: "mlops.cnvrg.io/v1"},

		ObjectMeta: metav1.ObjectMeta{Name: ns},
		Spec:       infraEmptySpec,
	}
}

func getDefaultTestInfraSpec(ns string) *mlopsv1.CnvrgInfra {
	testSpec := mlopsv1.DefaultCnvrgInfraSpec()
	testSpec.Cri = mlopsv1.CriTypeDocker
	testSpec.InfraNamespace = ns
	testSpec.ClusterDomain = "test.local"
	return &mlopsv1.CnvrgInfra{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CnvrgInfra",
			APIVersion: "mlops.cnvrg.io/v1"},

		ObjectMeta: metav1.ObjectMeta{Name: ns},
		Spec:       testSpec,
	}
}
