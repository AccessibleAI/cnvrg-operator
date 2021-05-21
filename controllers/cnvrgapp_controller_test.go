package controllers

import (
	"context"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/teris-io/shortid"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
	"time"
)

var defaultTrue = true

// +kubebuilder:docs-gen:collapse=Imports
var enabledAppTests = map[string]bool{
	"pg":         true,
	"redis":      true,
	"minio":      true,
	"es":         true,
	"ccp":        true,
	"monitoring": true,
}

var _ = Describe("CnvrgApp controller", func() {

	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)
	if enabledAppTests["pg"] {
		Context("Test PG", func() {
			It("PG Labels", func() {
				ns := createNs()
				ctx := context.Background()
				labels := map[string]string{"foo": "bar"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Pg.Enabled = &defaultTrue
				testApp.Spec.Labels = labels

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Pg.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			})
			It("PG Annotations", func() {
				ns := createNs()
				ctx := context.Background()
				annotations := map[string]string{"foo1": "bar1"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Pg.Enabled = &defaultTrue
				testApp.Spec.Annotations = annotations

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Pg.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			})
			It("PG HugePages - defaults", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Pg.Enabled = &defaultTrue
				testApp.Spec.Dbs.Pg.HugePages.Enabled = &defaultTrue

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Pg.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				v := corev1.Volume{
					Name: "hugepage",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium:    "HugePages",
							SizeLimit: nil,
						},
					},
				}
				Expect(deployment.Spec.Template.Spec.Volumes).Should(ContainElement(v))

				vm := corev1.VolumeMount{Name: "hugepage", MountPath: "/hugepages"}
				Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts).Should(ContainElement(vm))

				l := corev1.ResourceList{"hugepages-2Mi": resource.MustParse("4Gi")}
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(BeEquivalentTo(l))

			})
			It("PG HugePages - custom", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Pg.Enabled = &defaultTrue
				testApp.Spec.Dbs.Pg.HugePages.Enabled = &defaultTrue
				testApp.Spec.Dbs.Pg.HugePages.Size = "1Gi"
				testApp.Spec.Dbs.Pg.HugePages.Memory = "2Gi"

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Pg.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				v := corev1.Volume{
					Name: "hugepage",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium:    "HugePages",
							SizeLimit: nil,
						},
					},
				}
				Expect(deployment.Spec.Template.Spec.Volumes).Should(ContainElement(v))

				vm := corev1.VolumeMount{Name: "hugepage", MountPath: "/hugepages"}
				Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts).Should(ContainElement(vm))

				l := corev1.ResourceList{"hugepages-1Gi": resource.MustParse("2Gi")}
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(BeEquivalentTo(l))

			})
			It("PG NodeSelector", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Pg.Enabled = &defaultTrue
				testApp.Spec.Dbs.Pg.HugePages.Enabled = &defaultTrue
				testApp.Spec.Dbs.Pg.NodeSelector = map[string]string{"foo": "bar"}

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Pg.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))

			})
			It("PG Tenancy & NodeSelector", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Pg.Enabled = &defaultTrue
				testApp.Spec.Dbs.Pg.HugePages.Enabled = &defaultTrue
				testApp.Spec.Dbs.Pg.NodeSelector = map[string]string{"foo": "bar"}
				testApp.Spec.Tenancy.Enabled = &defaultTrue

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Pg.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				t := corev1.Toleration{
					Key:      testApp.Spec.Tenancy.Key,
					Operator: "Equal",
					Value:    testApp.Spec.Tenancy.Value,
					Effect:   "NoSchedule",
				}

				Expect(deployment.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))

			})

		})
	}
	if enabledAppTests["redis"] {
		Context("Test Redis", func() {
			It("Redis Labels", func() {
				ns := createNs()
				ctx := context.Background()
				labels := map[string]string{"foo": "bar"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Redis.Enabled = &defaultTrue
				testApp.Spec.Labels = labels

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Redis.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			})
			It("Redis Annotations", func() {
				ns := createNs()
				ctx := context.Background()
				annotations := map[string]string{"foo1": "bar1"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Redis.Enabled = &defaultTrue
				testApp.Spec.Annotations = annotations

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Redis.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			})
			It("Redis NodeSelector", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Redis.Enabled = &defaultTrue
				testApp.Spec.Dbs.Redis.NodeSelector = map[string]string{"foo": "bar"}

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Redis.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))

			})
			It("Redis Tenancy & NodeSelector", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Redis.Enabled = &defaultTrue
				testApp.Spec.Dbs.Redis.NodeSelector = map[string]string{"foo": "bar"}
				testApp.Spec.Tenancy.Enabled = &defaultTrue

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Redis.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				t := corev1.Toleration{
					Key:      testApp.Spec.Tenancy.Key,
					Operator: "Equal",
					Value:    testApp.Spec.Tenancy.Value,
					Effect:   "NoSchedule",
				}

				Expect(deployment.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))

			})

		})
	}
	if enabledAppTests["minio"] {
		Context("Test Minio", func() {
			It("Minio Labels", func() {
				ns := createNs()
				ctx := context.Background()
				labels := map[string]string{"foo": "bar"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Minio.Enabled = &defaultTrue
				testApp.Spec.Labels = labels

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Minio.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			})
			It("Minio Annotations", func() {
				ns := createNs()
				ctx := context.Background()
				annotations := map[string]string{"foo1": "bar1"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Minio.Enabled = &defaultTrue
				testApp.Spec.Annotations = annotations

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Minio.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			})
			It("Minio NodeSelector", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Minio.Enabled = &defaultTrue
				testApp.Spec.Dbs.Minio.NodeSelector = map[string]string{"foo": "bar"}

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Minio.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))

			})
			It("Minio Tenancy & NodeSelector", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Minio.Enabled = &defaultTrue
				testApp.Spec.Dbs.Minio.NodeSelector = map[string]string{"foo": "bar"}
				testApp.Spec.Tenancy.Enabled = &defaultTrue

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Minio.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				t := corev1.Toleration{
					Key:      testApp.Spec.Tenancy.Key,
					Operator: "Equal",
					Value:    testApp.Spec.Tenancy.Value,
					Effect:   "NoSchedule",
				}

				Expect(deployment.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))

			})

		})
	}
	if enabledAppTests["es"] {
		Context("Test Es", func() {
			It("Es Labels", func() {
				ns := createNs()
				ctx := context.Background()
				labels := map[string]string{"foo": "bar"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Es.Enabled = &defaultTrue
				testApp.Spec.Labels = labels

				sts := v1.StatefulSet{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Es.SvcName, Namespace: ns}, &sts)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(sts.Labels).Should(HaveKeyWithValue("foo", "bar"))
			})
			It("Es Annotations", func() {
				ns := createNs()
				ctx := context.Background()
				annotations := map[string]string{"foo1": "bar1"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Es.Enabled = &defaultTrue
				testApp.Spec.Annotations = annotations

				sts := v1.StatefulSet{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Es.SvcName, Namespace: ns}, &sts)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(sts.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			})
			It("Es NodeSelector", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Es.Enabled = &defaultTrue
				testApp.Spec.Dbs.Es.NodeSelector = map[string]string{"foo": "bar"}

				sts := v1.StatefulSet{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Es.SvcName, Namespace: ns}, &sts)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				Expect(sts.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))

			})
			It("Es Tenancy & NodeSelector", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.Dbs.Es.Enabled = &defaultTrue
				testApp.Spec.Dbs.Es.NodeSelector = map[string]string{"foo": "bar"}
				testApp.Spec.Tenancy.Enabled = &defaultTrue

				sts := v1.StatefulSet{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Es.SvcName, Namespace: ns}, &sts)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				t := corev1.Toleration{
					Key:      testApp.Spec.Tenancy.Key,
					Operator: "Equal",
					Value:    testApp.Spec.Tenancy.Value,
					Effect:   "NoSchedule",
				}

				Expect(sts.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
				Expect(sts.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("foo", "bar"))
				Expect(sts.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))

			})

		})
	}
	if enabledAppTests["ccp"] {
		Context("Test Cnvrg Control Plane", func() {

			It("WebApp Labels", func() {
				ns := createNs()
				ctx := context.Background()
				labels := map[string]string{"foo": "bar"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.WebApp.Enabled = &defaultTrue
				testApp.Spec.Labels = labels
				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.WebApp.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			})
			It("WebApp Annotations", func() {
				ns := createNs()
				ctx := context.Background()
				annotations := map[string]string{"foo1": "bar1"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.WebApp.Enabled = &defaultTrue
				testApp.Spec.Annotations = annotations

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.WebApp.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			})
			It("WebApp Tenancy", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.WebApp.Enabled = &defaultTrue
				testApp.Spec.Tenancy.Enabled = &defaultTrue

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.WebApp.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				t := corev1.Toleration{
					Key:      testApp.Spec.Tenancy.Key,
					Operator: "Equal",
					Value:    testApp.Spec.Tenancy.Value,
					Effect:   "NoSchedule",
				}

				Expect(deployment.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))

			})

			It("Sidekiq Labels", func() {
				ns := createNs()
				ctx := context.Background()
				labels := map[string]string{"foo": "bar"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Sidekiq.Split = &defaultTrue
				testApp.Spec.ControlPlane.Sidekiq.Enabled = &defaultTrue
				testApp.Spec.Labels = labels
				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: "sidekiq", Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			})
			It("Sidekiq Annotations", func() {
				ns := createNs()
				ctx := context.Background()
				annotations := map[string]string{"foo1": "bar1"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Sidekiq.Split = &defaultTrue
				testApp.Spec.ControlPlane.Sidekiq.Enabled = &defaultTrue
				testApp.Spec.Annotations = annotations

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: "sidekiq", Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			})
			It("Sidekiq Tenancy", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Sidekiq.Split = &defaultTrue
				testApp.Spec.ControlPlane.Sidekiq.Enabled = &defaultTrue
				testApp.Spec.Tenancy.Enabled = &defaultTrue

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: "sidekiq", Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				t := corev1.Toleration{
					Key:      testApp.Spec.Tenancy.Key,
					Operator: "Equal",
					Value:    testApp.Spec.Tenancy.Value,
					Effect:   "NoSchedule",
				}

				Expect(deployment.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))

			})

			It("Hyper Labels", func() {
				ns := createNs()
				ctx := context.Background()
				labels := map[string]string{"foo": "bar"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Hyper.Enabled = &defaultTrue
				testApp.Spec.Labels = labels
				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.Hyper.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			})
			It("Hyper Annotations", func() {
				ns := createNs()
				ctx := context.Background()
				annotations := map[string]string{"foo1": "bar1"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Hyper.Enabled = &defaultTrue
				testApp.Spec.Annotations = annotations

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.Hyper.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			})
			It("Hyper Tenancy", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Hyper.Enabled = &defaultTrue
				testApp.Spec.Tenancy.Enabled = &defaultTrue

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.Hyper.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				t := corev1.Toleration{
					Key:      testApp.Spec.Tenancy.Key,
					Operator: "Equal",
					Value:    testApp.Spec.Tenancy.Value,
					Effect:   "NoSchedule",
				}

				Expect(deployment.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))

			})

			It("Mpi Operator Labels", func() {
				ns := createNs()
				ctx := context.Background()
				labels := map[string]string{"foo": "bar"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Hyper.Enabled = &defaultTrue
				testApp.Spec.Labels = labels
				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.Hyper.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Labels).Should(HaveKeyWithValue("foo", "bar"))
			})
			It("Mpi Operator Annotations", func() {
				ns := createNs()
				ctx := context.Background()
				annotations := map[string]string{"foo1": "bar1"}
				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Hyper.Enabled = &defaultTrue
				testApp.Spec.Annotations = annotations

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.Hyper.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(deployment.Annotations).Should(HaveKeyWithValue("foo1", "bar1"))
			})
			It("Mpi Operator Tenancy", func() {
				ns := createNs()
				ctx := context.Background()

				testApp := getDefaultTestAppSpec(ns)
				testApp.Spec.ControlPlane.Hyper.Enabled = &defaultTrue
				testApp.Spec.Tenancy.Enabled = &defaultTrue

				deployment := v1.Deployment{}
				Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.Hyper.SvcName, Namespace: ns}, &deployment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				t := corev1.Toleration{
					Key:      testApp.Spec.Tenancy.Key,
					Operator: "Equal",
					Value:    testApp.Spec.Tenancy.Value,
					Effect:   "NoSchedule",
				}

				Expect(deployment.Spec.Template.Spec.Tolerations).Should(ContainElement(t))
				Expect(deployment.Spec.Template.Spec.NodeSelector).Should(HaveKeyWithValue("purpose", "cnvrg-control-plane"))

			})

		})
	}

})

func createNs() string {
	ns, _ := shortid.Generate()
	ns = strings.ReplaceAll(strings.ToLower(ns), "-", "z")
	ns = strings.ReplaceAll(ns, "_", "z")
	testNs := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
	err := k8sClient.Create(context.Background(), testNs)
	if err != nil {
		panic(err)
	}
	return ns
}

func getDefaultTestAppSpec(ns string) *mlopsv1.CnvrgApp {
	testSpec := mlopsv1.DefaultCnvrgAppSpec()

	return &mlopsv1.CnvrgApp{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CnvrgApp",
			APIVersion: "mlops.cnvrg.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cnvrgapp",
			Namespace: ns,
		},
		Spec: testSpec,
	}
}
