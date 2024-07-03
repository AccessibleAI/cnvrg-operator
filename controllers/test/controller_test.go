package test

import (
	cnvrgv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
)

var _ = Describe("Cnvrg App controller", func() {
	labels := map[string]string{"foo": "bar"}
	annotations := map[string]string{"foo1": "bar1"}

	It("should pass crd sanity", func() {
		namespace := genNS(GenName("ns"))
		ExpectCreate(&namespace)

		By("dry-running a basic cnvg app")
		obj := genApp(namespace.Name, GenName("cnvrg"))
		ExpectDryRunCreate(&obj)
	})

	It("should create and destroy managed resources", func() {
		nsName := GenName("ns")
		namespace := genNS(nsName)
		ExpectCreate(&namespace)

		By("creating simple cnvrg app", func() {
			appName := GenName("cnvrg")
			obj := genApp(nsName, appName)
			ExpectCreate(&obj)
			EventuallyGet(Key(nsName, appName), &obj)
		})

		By("validating app with PG enabled", func() {
			appName := GenName("cnvrg")
			appKey := Key(nsName, appName)
			var obj cnvrgv1.CnvrgApp
			var managedDeployment v1.Deployment

			By("creating cnvrg app with PG", func() {
				obj = genApp(nsName, appName)
				obj.Spec.Dbs.Pg.Enabled = true
				obj.Spec.Labels = labels
				obj.Spec.Annotations = annotations
				ExpectCreate(&obj)
				EventuallyGet(appKey, &obj)
			})
			deployKey := Key(nsName, obj.Spec.Dbs.Pg.SvcName)
			By("validating managed deployment exists", func() {
				EventuallyGet(deployKey, &managedDeployment)

				// expected labels exists on managed deployment
				Expect(managedDeployment.Labels).To(SatisfyAll(
					MatchSubMap(labels),
				))

				// expected annotations exists on managed deployment
				Expect(managedDeployment.Annotations).To(SatisfyAll(
					MatchSubMap(annotations),
				))
			})

			By("deleting app", func() {
				ExpectDelete(appKey, &obj)
				EventuallyGone(appKey, &obj)
			})
		})

	})
},
)
