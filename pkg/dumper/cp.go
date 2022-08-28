package dumper

//
//import (
//	"fmt"
//	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
//	"github.com/AccessibleAI/cnvrg-operator/controllers"
//	"github.com/AccessibleAI/cnvrg-operator/pkg/controlplane"
//	"github.com/AccessibleAI/cnvrg-operator/pkg/dbs"
//	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
//	"github.com/AccessibleAI/cnvrg-operator/pkg/monitoring"
//	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
//	"github.com/AccessibleAI/cnvrg-operator/pkg/priorityclass"
//	"github.com/AccessibleAI/cnvrg-operator/pkg/registry"
//	"github.com/Dimss/crypt/apr1_crypt"
//	log "github.com/sirupsen/logrus"
//	"github.com/spf13/viper"
//	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"os"
//)
//
//type ControlPlane struct {
//	image          string
//	wildcardDomain string
//	cri            string
//	registryUser   string
//	registryPass   string
//	ingress        string
//	app            *mlopsv1.CnvrgApp
//	state          []*desired.State
//	https          bool
//	proxy          bool
//}
//
//func NewControlPlane(image, domain, cri, regUser, regPass, ingress, ns string, https, proxy bool) *ControlPlane {
//
//	app := &mlopsv1.CnvrgApp{
//		ObjectMeta: v1.ObjectMeta{
//			Name:      "cnvrg-app",
//			Namespace: ns,
//		},
//		Spec: mlopsv1.DefaultCnvrgAppSpec(),
//	}
//	infra := &mlopsv1.CnvrgInfra{
//		ObjectMeta: v1.ObjectMeta{
//			Name:      "cnvrg-infra",
//			Namespace: ns,
//		},
//		Spec: mlopsv1.DefaultCnvrgInfraSpec(),
//	}
//	if err := controllers.CalculateAndApplyAppDefaults(app, &app.Spec, infra, nil); err != nil {
//		log.Fatal(err)
//	}
//
//	return &ControlPlane{
//		image:          image,
//		wildcardDomain: domain,
//		cri:            cri,
//		registryUser:   regUser,
//		registryPass:   regPass,
//		ingress:        ingress,
//		https:          https,
//		app:            app,
//		proxy:          proxy,
//	}
//}
//
//func (p *ControlPlane) BuildState() error {
//	var err error
//	err = p.setControlPlaneState()
//	err = p.setDbsState()
//	p.setPriorityClassState()
//	p.setProxy()
//
//	p.state = controlplane.State(p.app)
//	p.state = append(p.state, dbs.RedisCreds(p.getRedisSecretData())...)
//	p.state = append(p.state, dbs.PgCreds(p.getPgSecretData())...)
//	p.state = append(p.state, dbs.EsCreds(p.getEsSecretData())...)
//	p.state = append(p.state, dbs.AppDbsState(p.app)...)
//	p.state = append(p.state, networking.CnvrgAppNetworkingState(p.app)...)
//	p.state = append(p.state, monitoring.PromCreds(p.getPromSecretData())...)
//	p.state = append(p.state, registry.State(p.getRegistryData())...)
//	p.state = append(p.state, priorityclass.State()...)
//
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (p *ControlPlane) Dump(preserveTmplDirs bool) error {
//
//	if err := os.RemoveAll(viper.GetString("templates-dump-dir")); err != nil {
//		log.Fatal(err)
//	}
//	for _, o := range p.state {
//
//		if o.TemplateData == nil {
//			o.TemplateData = p.app
//		}
//		if err := o.GenerateDeployable(); err != nil {
//			log.Fatal(err)
//		}
//		if err := o.DumpTemplateToFile(preserveTmplDirs); err != nil {
//			log.Fatal(err)
//		}
//	}
//
//	return nil
//}
//
//func (p *ControlPlane) setControlPlaneState() error {
//
//	p.app.Spec.Annotations = make(map[string]string)
//	p.app.Spec.Annotations["purpose"] = "cnvrg"
//	p.app.Spec.ControlPlane.WebApp.Enabled = true
//
//	p.app.Spec.ControlPlane.BaseConfig.FeatureFlags = make(map[string]string)
//	p.app.Spec.ControlPlane.BaseConfig.FeatureFlags["CNVRG_MOUNT_HOST_FOLDERS"] = "false"
//
//	p.app.Spec.ControlPlane.Image = viper.GetString("control-plane-image")
//	p.app.Spec.ControlPlane.Image = p.image
//	p.app.Spec.Cri = mlopsv1.CriType(p.cri)
//
//	p.app.Spec.Networking.Ingress.Type = mlopsv1.IngressType(p.ingress)
//	p.app.Spec.Networking.HTTPS.Enabled = p.https
//	if p.app.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress {
//		p.app.Spec.Networking.Ingress.IstioGwEnabled = true
//	}
//	p.app.Spec.ClusterDomain = p.wildcardDomain
//
//	p.app.Spec.ControlPlane.Sidekiq.Enabled = true
//	p.app.Spec.ControlPlane.Sidekiq.Split = true
//	p.app.Spec.ControlPlane.Hyper.Enabled = true
//	p.app.Spec.ControlPlane.CnvrgScheduler.Enabled = false
//	p.app.Spec.ControlPlane.Searchkiq.Enabled = true
//	p.app.Spec.ControlPlane.Systemkiq.Enabled = true
//
//	p.app.Spec.ControlPlane.ObjectStorage.Type = "minio"
//	p.app.Spec.ControlPlane.ObjectStorage.AzureContainer = "azure"
//	p.app.Spec.ControlPlane.ObjectStorage.AzureAccountName = "azure"
//	p.app.Spec.ControlPlane.ObjectStorage.AzureContainer = "azure"
//	p.app.Spec.ControlPlane.ObjectStorage.AccessKey = "AKIAIOSFODNN7EXAMPLE"
//	p.app.Spec.ControlPlane.ObjectStorage.SecretKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
//
//	p.app.Spec.ControlPlane.ObjectStorage.GcpProject = "gcp-project"
//
//	p.app.Spec.ControlPlane.SMTP.Port = 25
//	p.app.Spec.ControlPlane.SMTP.Server = "smtp"
//	p.app.Spec.ControlPlane.SMTP.Username = "username"
//	p.app.Spec.ControlPlane.SMTP.Password = "password"
//	p.app.Spec.ControlPlane.SMTP.Domain = "domain"
//	p.app.Spec.ControlPlane.SMTP.OpensslVerifyMode = "false"
//
//	p.app.Spec.ControlPlane.Ldap.Host = "host"
//	p.app.Spec.ControlPlane.Ldap.Port = "port"
//	p.app.Spec.ControlPlane.Ldap.Ssl = "false"
//	p.app.Spec.ControlPlane.Ldap.Base = "base"
//	p.app.Spec.ControlPlane.Ldap.AdminUser = "admin"
//	p.app.Spec.ControlPlane.Ldap.AdminPassword = "pass"
//
//	p.app.Spec.ControlPlane.BaseConfig.SentryURL = "sentry-url"
//
//	p.app.Spec.Dbs.Es.Enabled = true
//	p.app.Spec.Dbs.Pg.Enabled = true
//	p.app.Spec.Dbs.Minio.Enabled = true
//	p.app.Spec.Dbs.Redis.Enabled = true
//	p.app.Spec.CnvrgAppPriorityClass = mlopsv1.PriorityClass{Name: "cnvrg-apps", Value: 2000000, Description: "cnvrg control plane apps priority class"}
//	p.app.Spec.CnvrgJobPriorityClass = mlopsv1.PriorityClass{Name: "cnvrg-jobs", Value: 1000000, Description: "cnvrg jobs priority class"}
//
//	return nil
//}
//
//func (p *ControlPlane) setDbsState() error {
//	p.app.Spec.Dbs.Es.Enabled = true
//	p.app.Spec.Dbs.Pg.Enabled = true
//	p.app.Spec.Dbs.Minio.Enabled = true
//	p.app.Spec.Dbs.Redis.Enabled = true
//	return nil
//}
//
//func (p *ControlPlane) setPriorityClassState() {
//	p.app.Spec.CnvrgAppPriorityClass = mlopsv1.PriorityClass{Name: "cnvrg-apps", Value: 2000000, Description: "cnvrg control plane apps priority class"}
//	p.app.Spec.CnvrgJobPriorityClass = mlopsv1.PriorityClass{Name: "cnvrg-jobs", Value: 1000000, Description: "cnvrg jobs priority class"}
//}
//
//func (p *ControlPlane) getRedisSecretData() desired.TemplateData {
//	return desired.TemplateData{
//		Data: map[string]interface{}{
//			"Namespace":   p.app.Namespace,
//			"Annotations": p.app.Spec.Annotations,
//			"Labels":      p.app.Spec.Labels,
//			"CredsRef":    p.app.Spec.Dbs.Redis.CredsRef,
//			"SvcName":     p.app.Spec.Dbs.Redis.SvcName,
//		},
//	}
//}
//
//func (p *ControlPlane) getPgSecretData() desired.TemplateData {
//	return desired.TemplateData{
//		Data: map[string]interface{}{
//			"Namespace":          p.app.Namespace,
//			"CredsRef":           p.app.Spec.Dbs.Pg.CredsRef,
//			"Annotations":        p.app.Spec.Annotations,
//			"Labels":             p.app.Spec.Labels,
//			"MaxConnections":     p.app.Spec.Dbs.Pg.MaxConnections,
//			"SharedBuffers":      p.app.Spec.Dbs.Pg.SharedBuffers,
//			"EffectiveCacheSize": p.app.Spec.Dbs.Pg.EffectiveCacheSize,
//			"SvcName":            p.app.Spec.Dbs.Pg.SvcName,
//		},
//	}
//}
//
//func (p *ControlPlane) getPromSecretData() desired.TemplateData {
//	user := "cnvrg"
//	pass := desired.RandomString()
//	passHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return desired.TemplateData{
//		Data: map[string]interface{}{
//			"Namespace":   p.app.Namespace,
//			"Annotations": p.app.Spec.Annotations,
//			"Labels":      p.app.Spec.Labels,
//			"CredsRef":    p.app.Spec.Monitoring.Prometheus.CredsRef,
//			"User":        user,
//			"Pass":        pass,
//			"PassHash":    fmt.Sprintf("%s:%s", user, passHash),
//			"PromUrl": fmt.Sprintf("http://%s.%s.svc:%d",
//				p.app.Spec.Monitoring.Prometheus.SvcName,
//				p.app.Namespace,
//				p.app.Spec.Monitoring.Prometheus.Port,
//			),
//		},
//	}
//}
//
//func (p *ControlPlane) getEsSecretData() desired.TemplateData {
//	return desired.TemplateData{
//		Data: map[string]interface{}{
//			"Namespace": p.app.Namespace,
//			"CredsRef":  p.app.Spec.Dbs.Es.CredsRef,
//			"EsUrl": fmt.Sprintf("%s.%s.svc:%d",
//				p.app.Spec.Dbs.Es.SvcName,
//				p.app.Namespace,
//				p.app.Spec.Dbs.Es.Port),
//			"Annotations": p.app.Spec.Annotations,
//			"Labels":      p.app.Spec.Labels,
//		},
//	}
//}
//
//func (p *ControlPlane) getRegistryData() desired.TemplateData {
//	p.app.Spec.Registry.User = p.registryUser
//	p.app.Spec.Registry.Password = p.registryPass
//	return desired.TemplateData{
//		Namespace: p.app.Namespace,
//		Data: map[string]interface{}{
//			"Registry":    p.app.Spec.Registry,
//			"Annotations": p.app.Spec.Annotations,
//			"Labels":      p.app.Spec.Labels,
//		},
//	}
//}
//
//func (p *ControlPlane) setProxy() {
//	if p.proxy {
//		p.app.Spec.Networking.Proxy.Enabled = true
//		p.app.Spec.Networking.Proxy.HttpProxy = []string{"http://corp-proxy-example1"}
//		p.app.Spec.Networking.Proxy.HttpsProxy = []string{"http://corp-proxy-example1"}
//		p.app.Spec.Networking.Proxy.NoProxy = networking.DefaultNoProxy("cluster.local")
//	}
//}
