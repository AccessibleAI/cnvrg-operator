package installer

import (
	"fmt"
	"github.com/go-logr/logr"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	helmcli "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"time"
)

const (
	helmDriver         = "configmaps"
	helmInstallTimeout = 5 * time.Minute
)

var helmSettings = helmcli.New()

type Helm struct {
	cfg         *action.Configuration
	envSettings *helmcli.EnvSettings
	log         logr.Logger
	chartConfig ChartConfig
}

type ChartConfig struct {
	IncludeCRDs bool
	Namespace   string
	Pass        string
	ReleaseName string
	Url         string
	User        string
	Values      map[string]interface{}
	Version     string
}

func NewHelm(config ChartConfig, l logr.Logger) (*Helm, error) {
	helmSettings.SetNamespace(config.Namespace)
	helmSettings.RepositoryCache = "/tmp/.helmcache"
	helmSettings.RepositoryConfig = "/tmp/.helmrepo"
	installer := &Helm{
		envSettings: helmSettings,
		log:         l.WithValues("installer", "helm"),
		chartConfig: config,
	}

	if err := installer.initHelm(); err != nil {
		return nil, err
	}

	return installer, nil
}

func (h *Helm) Install() error {
	if h.shouldUpgrade() {
		// if helm chart already installed
		// run upgrade to keep idempotent logic
		return h.upgrade()
	} else {
		// execute new install
		return h.install()
	}
}

func (h *Helm) initHelm() error {
	// setup helm configuration
	h.cfg = &action.Configuration{}
	if err := h.cfg.Init(
		h.envSettings.RESTClientGetter(),
		h.chartConfig.Namespace,
		helmDriver,
		h.debug); err != nil {
		return err
	}
	if err := h.configureRegistryClient(); err != nil {
		return err
	}

	return nil
}

func (h *Helm) configureRegistryClient() error {
	if h.chartConfig.User != "" && h.chartConfig.Pass != "" {
		registryClient, err := registry.NewClient()
		if err != nil {
			return err
		}
		err = registryClient.Login(h.chartConfig.Url, registry.LoginOptBasicAuth(h.chartConfig.User, h.chartConfig.Pass))
		if err != nil {
			return err
		}
		h.cfg.RegistryClient = registryClient

	}
	return nil
}

func (h *Helm) loadChart(cpo *action.ChartPathOptions) (*chart.Chart, error) {
	cpo.Version = h.chartConfig.Version
	cpo.RepoURL = h.chartConfig.Url
	cp, err := cpo.LocateChart(h.chartConfig.ReleaseName, h.envSettings)
	if err != nil {
		return nil, err
	}
	return loader.Load(cp)
}

func (h *Helm) install() error {
	installCmd := action.NewInstall(h.cfg)
	// load (in fact download from the remote registry) the helm chart package
	chartPkg, err := h.loadChart(&installCmd.ChartPathOptions)
	if err != nil {
		return err
	}

	installCmd.Namespace = h.chartConfig.Namespace
	installCmd.ReleaseName = h.chartConfig.ReleaseName
	installCmd.CreateNamespace = true
	installCmd.Wait = false
	installCmd.DisableOpenAPIValidation = true
	installCmd.Timeout = helmInstallTimeout
	installCmd.IncludeCRDs = h.chartConfig.IncludeCRDs
	// execute helm chart installation
	rel, err := installCmd.Run(chartPkg, h.chartConfig.Values)
	if err != nil {
		return err
	}
	h.log.Info("successfully installed", "release", rel.Name)
	return nil
}

func (h *Helm) upgrade() error {
	upgradeCmd := action.NewUpgrade(h.cfg)
	chartPkg, err := h.loadChart(&upgradeCmd.ChartPathOptions)
	if err != nil {
		return err
	}

	upgradeCmd.Namespace = h.chartConfig.Namespace
	upgradeCmd.MaxHistory = 5
	upgradeCmd.Timeout = helmInstallTimeout
	upgradeCmd.DisableOpenAPIValidation = true
	upgradeCmd.Wait = false
	rel, err := upgradeCmd.Run(h.chartConfig.ReleaseName, chartPkg, h.chartConfig.Values)
	if err != nil {
		return err
	}
	h.log.Info("successfully upgraded", "release", rel.Name)
	return nil
}

func (h *Helm) Delete() error {
	uninstallCmd := action.NewUninstall(h.cfg)
	uninstallCmd.Wait = false
	uninstallCmd.Timeout = helmInstallTimeout
	if _, err := uninstallCmd.Run(h.chartConfig.ReleaseName); err != nil {
		return err
	}

	return nil
}

func (h *Helm) shouldUpgrade() bool {
	lastRelease, err := h.cfg.Releases.Last(h.chartConfig.ReleaseName)
	if err != nil {
		return false
	}

	// Concurrent `helm upgrade`s will either fail here with `errPending` or when creating the release with "already exists".
	// This should act as a pessimistic lock.
	if lastRelease.Info.Status.IsPending() {
		lastRelease.SetStatus(release.StatusFailed, "insure locking")
		if err = h.cfg.Releases.Update(lastRelease); err != nil {
			h.log.Error(err, "failed to update status for locking")
		}
	}
	return true
}

func (h *Helm) debug(format string, v ...interface{}) {
	h.log.Info(fmt.Sprintf(format, v...))
}
