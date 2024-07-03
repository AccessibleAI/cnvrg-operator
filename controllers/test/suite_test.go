package test

import (
	"context"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers/app"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"math/rand"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"testing"
	"time"
)

var tc TestCtx

type TestCtx struct {
	ctx context.Context

	env    *envtest.Environment
	mgr    ctrl.Manager
	client client.Client
	scheme *runtime.Scheme
	rnd    *rand.Rand
}

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	setExtBinsEnv()
	err := os.Setenv("ACK_GINKGO_DEPRECATIONS", "1.16.4")
	Expect(err).To(BeNil())

	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(time.Second)
	SetDefaultConsistentlyDuration(2 * time.Second)
	SetDefaultConsistentlyPollingInterval(250 * time.Millisecond)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	done := make(chan struct{})

	go func() {
		projectBaseDir := "../../"
		crdPaths := []string{
			filepath.Join(projectBaseDir, "charts/mlops/crds/"),
		}

		env := &envtest.Environment{
			CRDDirectoryPaths: crdPaths,
		}

		cfg, err := env.Start()
		Expect(err).ToNot(HaveOccurred())

		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme))
		Expect(mlopsv1.AddToScheme(scheme)).To(Succeed())

		mgr, err := ctrl.NewManager(cfg, ctrl.Options{
			Scheme:         scheme,
			LeaderElection: false,
		})
		Expect(err).To(BeNil())

		ctx := signals.SetupSignalHandler()
		tc = TestCtx{
			ctx:    ctx,
			env:    env,
			mgr:    mgr,
			scheme: mgr.GetScheme(),
			client: mgr.GetClient(),
			rnd:    rand.New(rand.NewSource(GinkgoRandomSeed())),
		}
		err = (&app.CnvrgAppReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr)

		Expect(err).ToNot(HaveOccurred())

		go func() {
			defer GinkgoRecover()
			err := mgr.Start(tc.ctx)
			Expect(err).ToNot(HaveOccurred())

			gexec.KillAndWait(4 * time.Second)

			// Teardown the test environment once controller is fnished.
			// Otherwise from Kubernetes 1.21+, teardon timeouts waiting on
			// kube-apiserver to return
			err = env.Stop()
			Expect(err).ToNot(HaveOccurred())
		}()

		close(done)
	}()

	Eventually(done, 60).Should(BeClosed())
})

func setExtBinsEnv() {
	const extBinsLocation = "/tmp/kubebuilder/bin/" // scripts/fetch_ext_bins.sh should run before
	Expect(os.Setenv("TEST_ASSET_KUBE_APISERVER", extBinsLocation+"kube-apiserver")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_ETCD", extBinsLocation+"etcd")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_KUBECTL", extBinsLocation+"kubectl")).To(Succeed())
}
