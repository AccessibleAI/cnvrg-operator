package controllers

import (
	"context"
	"fmt"
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

var _ = Describe("CnvrgApp controller", func() {

	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

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

			shouldLimits := map[corev1.ResourceName]resource.Quantity{
				"hugepages-2Mi":       resource.MustParse(testApp.Spec.Dbs.Pg.Limits.Memory),
				corev1.ResourceCPU:    resource.MustParse(testApp.Spec.Dbs.Pg.Limits.Cpu),
				corev1.ResourceMemory: resource.MustParse(testApp.Spec.Dbs.Pg.Limits.Memory),
			}

			Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(ContainElement(shouldLimits["hugepages-2Mi"]))
			Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(ContainElement(shouldLimits["cpu"]))
			Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(ContainElement(shouldLimits["memory"]))

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

			shouldLimits := map[corev1.ResourceName]resource.Quantity{
				"hugepages-2Mi":       resource.MustParse("2Gi"),
				corev1.ResourceCPU:    resource.MustParse(testApp.Spec.Dbs.Pg.Limits.Cpu),
				corev1.ResourceMemory: resource.MustParse(testApp.Spec.Dbs.Pg.Limits.Memory),
			}

			Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(ContainElement(shouldLimits["hugepages-2Mi"]))
			Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(ContainElement(shouldLimits["cpu"]))
			Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(ContainElement(shouldLimits["memory"]))

		})
		It("PG HugePages - disabled", func() {
			ns := createNs()
			ctx := context.Background()

			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.Dbs.Pg.Enabled = &defaultTrue

			deployment := v1.Deployment{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Pg.SvcName, Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			shouldLimits := map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse(testApp.Spec.Dbs.Pg.Limits.Cpu),
				corev1.ResourceMemory: resource.MustParse(testApp.Spec.Dbs.Pg.Limits.Memory),
			}

			v := corev1.Volume{
				Name: "hugepage",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						Medium:    "HugePages",
						SizeLimit: nil,
					},
				},
			}
			Expect(deployment.Spec.Template.Spec.Volumes).ShouldNot(ContainElement(v))
			vm := corev1.VolumeMount{Name: "hugepage", MountPath: "/hugepages"}
			Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts).ShouldNot(ContainElement(vm))
			Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits).Should(BeEquivalentTo(shouldLimits))

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
		FIt("Es default xms xmx", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.Dbs.Es.Requests.Memory = "4Gi"
			testApp.Spec.Dbs.Es.Enabled = &defaultTrue
			sts := v1.StatefulSet{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Es.SvcName, Namespace: ns}, &sts)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			javaOpts := corev1.EnvVar{Name: "ES_JAVA_OPTS", Value: "-Xms2g -Xmx2g"}
			Expect(sts.Spec.Template.Spec.Containers[0].Env).Should(ContainElement(javaOpts))
		})
		FIt("Es 5Gi requests  -  xms xmx", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.Dbs.Es.Requests.Memory = "5Gi"
			testApp.Spec.Dbs.Es.Enabled = &defaultTrue
			sts := v1.StatefulSet{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Es.SvcName, Namespace: ns}, &sts)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			javaOpts := corev1.EnvVar{Name: "ES_JAVA_OPTS", Value: "-Xms2g -Xmx2g"}
			Expect(sts.Spec.Template.Spec.Containers[0].Env).Should(ContainElement(javaOpts))
		})
		FIt("Es 6Gi requests -  xms xmx", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.Dbs.Es.Requests.Memory = "6Gi"
			testApp.Spec.Dbs.Es.Enabled = &defaultTrue
			sts := v1.StatefulSet{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Es.SvcName, Namespace: ns}, &sts)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			javaOpts := corev1.EnvVar{Name: "ES_JAVA_OPTS", Value: "-Xms3g -Xmx3g"}
			Expect(sts.Spec.Template.Spec.Containers[0].Env).Should(ContainElement(javaOpts))
		})
		FIt("Es 6000Mi requests - xms xmx", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.Dbs.Es.Requests.Memory = "6000Mi"
			testApp.Spec.Dbs.Es.Enabled = &defaultTrue
			sts := v1.StatefulSet{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.Dbs.Es.SvcName, Namespace: ns}, &sts)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			javaOpts := corev1.EnvVar{Name: "ES_JAVA_OPTS", Value: ""}
			Expect(sts.Spec.Template.Spec.Containers[0].Env).Should(ContainElement(javaOpts))
		})

	})

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

		It("ImageHub for WebApp - default ImageHub", func() {
			ns := createNs()
			ctx := context.Background()

			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.ControlPlane.WebApp.Enabled = &defaultTrue
			testApp.Spec.ControlPlane.Image = "app:1.2.3"

			dep := v1.Deployment{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.WebApp.SvcName, Namespace: ns}, &dep)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			shouldBe := fmt.Sprintf("%s/%s", testApp.Spec.ImageHub, testApp.Spec.ControlPlane.Image)
			Expect(dep.Spec.Template.Spec.Containers[0].Image).Should(Equal(shouldBe))

		})

		It("ImageHub for WebApp - custom ImageHub", func() {
			ns := createNs()
			ctx := context.Background()

			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.ControlPlane.WebApp.Enabled = &defaultTrue
			testApp.Spec.ImageHub = "foo/bar"
			testApp.Spec.ControlPlane.Image = "app:1.2.3"

			deployment := v1.Deployment{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.WebApp.SvcName, Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			shouldBe := fmt.Sprintf("%s/%s", testApp.Spec.ImageHub, testApp.Spec.ControlPlane.Image)
			Expect(deployment.Spec.Template.Spec.Containers[0].Image).Should(Equal(shouldBe))

		})

		It("ImageHub for Sidekiq - custom ImageHub", func() {
			ns := createNs()
			ctx := context.Background()

			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.ControlPlane.Sidekiq.Enabled = &defaultTrue
			testApp.Spec.ControlPlane.Sidekiq.Split = &defaultTrue
			testApp.Spec.ImageHub = "foo/bar"

			deployment := v1.Deployment{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "sidekiq", Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			shouldBe := fmt.Sprintf("%s/%s", testApp.Spec.ImageHub, testApp.Spec.ControlPlane.Image)
			Expect(deployment.Spec.Template.Spec.Containers[0].Image).Should(Equal(shouldBe))

		})

		It("Image for WebApp - disable  ImageHub", func() {
			ns := createNs()
			ctx := context.Background()

			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.ControlPlane.WebApp.Enabled = &defaultTrue
			testApp.Spec.ImageHub = "foo/bar"
			testApp.Spec.ControlPlane.Image = "foo/app:1.2.3"

			deployment := v1.Deployment{}
			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: testApp.Spec.ControlPlane.WebApp.SvcName, Namespace: ns}, &deployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			shouldBe := fmt.Sprintf("%s", testApp.Spec.ControlPlane.Image)
			Expect(deployment.Spec.Template.Spec.Containers[0].Image).Should(Equal(shouldBe))

		})

		It("Labels/Annotations CCP ConfigMap", func() {
			ns := createNs()
			ctx := context.Background()

			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.Labels = map[string]string{"foo": "bar", "foo1": "bar1"}
			testApp.Spec.Annotations = map[string]string{"foo1": "bar1"}
			testApp.Spec.ControlPlane.WebApp.Enabled = &defaultTrue

			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())

			cm := corev1.ConfigMap{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cp-annotation-label", Namespace: ns}, &cm)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(cm.Data["labels"]).Should(ContainSubstring("foo"))
			Expect(cm.Data["labels"]).Should(ContainSubstring("bar"))
			Expect(cm.Data["annotations"]).Should(ContainSubstring("foo1"))
			Expect(cm.Data["annotations"]).Should(ContainSubstring("bar1"))

		})

	})

	Context("Test Object Storage Secret", func() {
		It("Object Storage Secret - Minio enabled - random keys", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.Dbs.Minio.Enabled = &defaultTrue
			testApp.Spec.ControlPlane.ObjectStorage.Type = "minio"

			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			secret := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cp-object-storage", Namespace: ns}, &secret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(secret.Data["CNVRG_STORAGE_ACCESS_KEY"]).ShouldNot(BeEmpty())
			Expect(secret.Data["CNVRG_STORAGE_SECRET_KEY"]).ShouldNot(BeEmpty())
		})

		It("Object Storage Secret - Minio enabled - explicitly configured keys", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.Dbs.Minio.Enabled = &defaultTrue
			testApp.Spec.ControlPlane.ObjectStorage.Type = "minio"
			testApp.Spec.ControlPlane.ObjectStorage.AccessKey = "access-key"
			testApp.Spec.ControlPlane.ObjectStorage.SecretKey = "secret-key"

			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			secret := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cp-object-storage", Namespace: ns}, &secret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(secret.Data["CNVRG_STORAGE_ACCESS_KEY"]).Should(Equal([]byte("access-key")))
			Expect(secret.Data["CNVRG_STORAGE_SECRET_KEY"]).Should(Equal([]byte("secret-key")))
		})

		It("Object Storage Secret - Minio external - explicitly configured keys", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.ControlPlane.ObjectStorage.Type = "minio"
			testApp.Spec.ControlPlane.ObjectStorage.AccessKey = "access-key"
			testApp.Spec.ControlPlane.ObjectStorage.SecretKey = "secret-key"

			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			secret := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cp-object-storage", Namespace: ns}, &secret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(secret.Data["CNVRG_STORAGE_ACCESS_KEY"]).Should(Equal([]byte("access-key")))
			Expect(secret.Data["CNVRG_STORAGE_SECRET_KEY"]).Should(Equal([]byte("secret-key")))
		})

		It("Object Storage Secret - AWS S3 - explicitly configured keys", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.ControlPlane.ObjectStorage.Type = "aws"
			testApp.Spec.ControlPlane.ObjectStorage.AccessKey = "access-key"
			testApp.Spec.ControlPlane.ObjectStorage.SecretKey = "secret-key"

			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			secret := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cp-object-storage", Namespace: ns}, &secret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(secret.Data["CNVRG_STORAGE_ACCESS_KEY"]).Should(Equal([]byte("access-key")))
			Expect(secret.Data["CNVRG_STORAGE_SECRET_KEY"]).Should(Equal([]byte("secret-key")))
		})

		It("Object Storage Secret - AWS S3 - IAM S3 access", func() {
			ns := createNs()
			ctx := context.Background()
			testApp := getDefaultTestAppSpec(ns)
			testApp.Spec.ControlPlane.ObjectStorage.Type = "aws"

			Expect(k8sClient.Create(ctx, testApp)).Should(Succeed())
			secret := corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "cp-object-storage", Namespace: ns}, &secret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(secret.Data["CNVRG_STORAGE_ACCESS_KEY"]).Should(BeEmpty())
			Expect(secret.Data["CNVRG_STORAGE_SECRET_KEY"]).Should(BeEmpty())
		})

	})

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
